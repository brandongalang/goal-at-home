package install

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	goalPreCmd = "goal hook pre-tool-use"
	goalStopCmd = "goal hook stop"
)

// Agent identifies a supported install target.
type Agent string

const (
	AgentCursor Agent = "cursor"
	AgentClaude Agent = "claude"
	AgentCodex  Agent = "codex"
	AgentGemini Agent = "gemini"
	AgentAll    Agent = "all"
)

// Options configures an install run.
type Options struct {
	RepoRoot          string
	InstallDir        string
	Agents            []Agent
	DryRun            bool
	SkipInstructions  bool
}

// AllAgents returns install targets in stable order.
func AllAgents() []Agent {
	return []Agent{AgentCursor, AgentClaude, AgentCodex, AgentGemini}
}

// ParseAgents expands "all" and validates names.
func ParseAgents(names []string) ([]Agent, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("specify at least one agent: cursor, claude, codex, gemini, or all")
	}
	var out []Agent
	seen := map[Agent]bool{}
	for _, raw := range names {
		for _, part := range strings.Split(raw, ",") {
			part = strings.TrimSpace(strings.ToLower(part))
			if part == "" {
				continue
			}
			if part == string(AgentAll) {
				for _, a := range AllAgents() {
					if !seen[a] {
						seen[a] = true
						out = append(out, a)
					}
				}
				continue
			}
			a := Agent(part)
			switch a {
			case AgentCursor, AgentClaude, AgentCodex, AgentGemini:
				if !seen[a] {
					seen[a] = true
					out = append(out, a)
				}
			default:
				return nil, fmt.Errorf("unknown agent %q (use cursor, claude, codex, gemini, or all)", part)
			}
		}
	}
	return out, nil
}

// Run builds the CLI, installs skills, and merges hooks for each agent.
func Run(opts Options) error {
	if opts.RepoRoot == "" {
		root, err := FindRepoRoot()
		if err != nil {
			return err
		}
		opts.RepoRoot = root
	}
	if opts.InstallDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		opts.InstallDir = filepath.Join(home, ".local", "bin")
	}
	agents, err := ParseAgents(agentNames(opts.Agents))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(opts.InstallDir, 0o755); err != nil {
		return err
	}
	binPath := filepath.Join(opts.InstallDir, "goal")
	if opts.DryRun {
		fmt.Printf("would build: %s\n", binPath)
	} else if err := buildBinary(opts.RepoRoot, binPath); err != nil {
		return err
	}

	for _, agent := range agents {
		if err := installAgent(opts, agent); err != nil {
			return fmt.Errorf("%s: %w", agent, err)
		}
	}

	fmt.Println("")
	fmt.Println("Installed:")
	fmt.Printf("  binary: %s\n", binPath)
	for _, agent := range agents {
		for _, line := range agentSummary(agent) {
			fmt.Printf("  %s: %s\n", agent, line)
		}
	}
	fmt.Printf("\nEnsure %s is on your PATH.\n", opts.InstallDir)
	fmt.Println(restartHint(agents))
	return nil
}

func agentNames(agents []Agent) []string {
	if len(agents) > 0 {
		out := make([]string, len(agents))
		for i, a := range agents {
			out[i] = string(a)
		}
		return out
	}
	return nil
}

func buildBinary(repoRoot, dest string) error {
	// Re-exec go build from repo root; keeps install logic in one binary after first build.
	return runGoBuild(repoRoot, dest)
}

func installAgent(opts Options, agent Agent) error {
	paths, err := pathsForAgent(agent)
	if err != nil {
		return err
	}
	if opts.DryRun {
		fmt.Printf("would install %s hooks → %s\n", agent, paths.hooksFile)
		fmt.Printf("would install %s skills → %s\n", agent, paths.skillsBase)
		if !opts.SkipInstructions {
			fmt.Printf("would append %s instructions → %s\n", agent, paths.instructionsFile)
		}
		return nil
	}
	if err := os.MkdirAll(paths.goalsDir, 0o755); err != nil {
		return err
	}
	if err := installSkills(opts.RepoRoot, paths); err != nil {
		return err
	}
	if err := mergeHooks(agent, paths); err != nil {
		return err
	}
	if opts.SkipInstructions {
		return nil
	}
	return appendInstructions(paths)
}

type targetPaths struct {
	hooksFile         string
	goalsDir          string
	skillsBase        string
	instructionsFile  string
}

func pathsForAgent(agent Agent) (targetPaths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return targetPaths{}, err
	}
	switch agent {
	case AgentCursor:
		base := envOr("CURSOR_DIR", filepath.Join(home, ".cursor"))
		return targetPaths{
			hooksFile:        filepath.Join(base, "hooks.json"),
			goalsDir:         filepath.Join(base, "goals"),
			skillsBase:       filepath.Join(base, "skills"),
			instructionsFile: filepath.Join(base, "AGENTS.md"),
		}, nil
	case AgentClaude:
		base := envOr("CLAUDE_CONFIG_DIR", filepath.Join(home, ".claude"))
		return targetPaths{
			hooksFile:        filepath.Join(base, "settings.json"),
			goalsDir:         filepath.Join(home, ".cursor", "goals"),
			skillsBase:       filepath.Join(base, "skills"),
			instructionsFile: filepath.Join(base, "CLAUDE.md"),
		}, nil
	case AgentCodex:
		base := envOr("CODEX_HOME", filepath.Join(home, ".codex"))
		return targetPaths{
			hooksFile:        filepath.Join(base, "hooks.json"),
			goalsDir:         filepath.Join(home, ".cursor", "goals"),
			skillsBase:       filepath.Join(base, "skills"),
			instructionsFile: filepath.Join(base, "AGENTS.md"),
		}, nil
	case AgentGemini:
		base := envOr("GEMINI_CLI_HOME", home)
		geminiDir := filepath.Join(base, ".gemini")
		return targetPaths{
			hooksFile:        filepath.Join(geminiDir, "settings.json"),
			goalsDir:         filepath.Join(home, ".cursor", "goals"),
			skillsBase:       filepath.Join(geminiDir, "skills"),
			instructionsFile: filepath.Join(geminiDir, "GEMINI.md"),
		}, nil
	default:
		return targetPaths{}, fmt.Errorf("unsupported agent %q", agent)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func installSkills(repoRoot string, paths targetPaths) error {
	for _, name := range []string{"goal-enforcement", "goalcraft-at-home"} {
		destDir := filepath.Join(paths.skillsBase, name)
		if err := os.MkdirAll(destDir, 0o755); err != nil {
			return err
		}
		src := filepath.Join(repoRoot, "skill", name, "SKILL.md")
		data, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("read skill %s: %w", name, err)
		}
		if err := os.WriteFile(filepath.Join(destDir, "SKILL.md"), data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func mergeHooks(agent Agent, paths targetPaths) error {
	data, err := readJSONFile(paths.hooksFile)
	if err != nil {
		return err
	}
	switch agent {
	case AgentCursor:
		mergeCursorHooks(data)
	case AgentClaude:
		mergeClaudeHooks(data)
	case AgentCodex:
		mergeCodexHooks(data)
	case AgentGemini:
		mergeGeminiHooks(data)
	default:
		return fmt.Errorf("unsupported agent %q", agent)
	}
	return writeJSONFile(paths.hooksFile, data)
}

func agentSummary(agent Agent) []string {
	p, err := pathsForAgent(agent)
	if err != nil {
		return []string{err.Error()}
	}
	return []string{
		"hooks " + p.hooksFile,
		"skills " + p.skillsBase,
		"goals " + p.goalsDir,
		"instructions " + p.instructionsFile,
	}
}

func restartHint(agents []Agent) string {
	var hints []string
	for _, a := range agents {
		switch a {
		case AgentCursor:
			hints = append(hints, "Restart Cursor to reload hooks.")
		case AgentClaude:
			hints = append(hints, "Restart Claude Code to reload hooks.")
		case AgentCodex:
			hints = append(hints, "Restart Codex CLI to reload hooks.")
		case AgentGemini:
			hints = append(hints, "Restart Gemini CLI to reload hooks.")
		}
	}
	return strings.Join(hints, " ")
}

// FindRepoRoot locates the directory containing go.mod.
func FindRepoRoot() (string, error) {
	if v := os.Getenv("GOAL_REPO"); v != "" {
		return v, nil
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("could not find goal-at-home repo (go.mod); run from repo root or set GOAL_REPO")
}

func readJSONFile(path string) (map[string]any, error) {
	data := map[string]any{}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	if err := dec.Decode(&data); err != nil && err != io.EOF {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return data, nil
}

func writeJSONFile(path string, data map[string]any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0o644)
}
