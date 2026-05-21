package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/brandongalang/goal-at-home/internal/hooks"
	"github.com/brandongalang/goal-at-home/internal/install"
	"github.com/brandongalang/goal-at-home/internal/store"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "hook":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: goal hook <pre-tool-use|stop>")
			os.Exit(2)
		}
		var err error
		switch os.Args[2] {
		case "pre-tool-use":
			err = hooks.PreToolUse()
		case "stop":
			err = hooks.Stop()
		default:
			fmt.Fprintf(os.Stderr, "unknown hook: %s\n", os.Args[2])
			os.Exit(2)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "hook error: %v\n", err)
			os.Exit(1)
		}
	case "install":
		if err := runInstall(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "set", "edit", "complete", "clear", "status":
		if err := runCommand(os.Args[1], os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `goal-at-home — session goal enforcement for coding agents

Usage:
  goal set    --session-id <id> <objective text>
  goal edit   --session-id <id> <objective text>
  goal complete --session-id <id>
  goal clear  --session-id <id>
  goal status --session-id <id>
  goal hook pre-tool-use   (pre-tool hook; Cursor/Claude/Codex/Gemini JSON on stdin)
  goal hook stop           (stop hook; same agents)

  goal install --agent <cursor|claude|codex|gemini|all>
  goal install list

Docs: https://github.com/brandongalang/goal-at-home
`)
}

func runInstall(args []string) error {
	if len(args) == 1 && args[0] == "list" {
		for _, a := range install.AllAgents() {
			fmt.Println(a)
		}
		return nil
	}
	var agents []install.Agent
	dryRun := false
	skipInstructions := false
	installDir := ""
	repo := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--agent", "-a":
			if i+1 >= len(args) {
				return fmt.Errorf("--agent requires a value")
			}
			i++
			parsed, err := install.ParseAgents([]string{args[i]})
			if err != nil {
				return err
			}
			agents = append(agents, parsed...)
		case "--bin-dir":
			if i+1 >= len(args) {
				return fmt.Errorf("--bin-dir requires a value")
			}
			i++
			installDir = args[i]
		case "--repo":
			if i+1 >= len(args) {
				return fmt.Errorf("--repo requires a value")
			}
			i++
			repo = args[i]
		case "--dry-run":
			dryRun = true
		case "--skip-instructions":
			skipInstructions = true
		case "help", "-h", "--help":
			printInstallUsage()
			return nil
		default:
			return fmt.Errorf("unknown install flag: %s", args[i])
		}
	}
	if len(agents) == 0 {
		printInstallUsage()
		return fmt.Errorf("specify --agent (e.g. goal install --agent cursor)")
	}
	return install.Run(install.Options{
		RepoRoot:         repo,
		InstallDir:       installDir,
		Agents:           agents,
		DryRun:           dryRun,
		SkipInstructions: skipInstructions,
	})
}

func printInstallUsage() {
	fmt.Fprintf(os.Stderr, `goal install — build CLI and wire hooks + skills

Usage:
  goal install --agent cursor     (~/.cursor: hooks, skills, AGENTS.md)
  goal install --agent claude     (~/.claude: hooks, skills, CLAUDE.md)
  goal install --agent codex      (~/.codex: hooks, skills, AGENTS.md)
  goal install --agent gemini     (~/.gemini: hooks, skills, GEMINI.md)
  goal install --agent all
  goal install list

First install from repo clone:  go run . install --agent cursor

Options:
  --bin-dir <path>        Install binary (default: ~/.local/bin/goal)
  --repo <path>           Repo root (default: find go.mod from cwd)
  --dry-run               Print actions without writing
  --skip-instructions     Do not append to global AGENTS.md / CLAUDE.md / GEMINI.md

`)
}

type parsedArgs struct {
	sessionID string
	text      string
}

func parseArgs(args []string) (parsedArgs, error) {
	var p parsedArgs
	var textParts []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--session-id":
			if i+1 >= len(args) {
				return p, fmt.Errorf("--session-id requires a value")
			}
			i++
			p.sessionID = args[i]
		default:
			textParts = append(textParts, args[i])
		}
	}
	p.text = strings.TrimSpace(strings.Join(textParts, " "))
	return p, nil
}

func runCommand(cmd string, args []string) error {
	p, err := parseArgs(args)
	if err != nil {
		return err
	}
	if p.sessionID == "" {
		return fmt.Errorf("--session-id is required")
	}
	switch cmd {
	case "set":
		if p.text == "" {
			return fmt.Errorf("objective text is required")
		}
		g, err := store.Set(p.sessionID, p.text)
		if err != nil {
			return err
		}
		fmt.Printf("goal set (active)\nsession=%s\n", g.SessionID)
		return nil
	case "edit":
		if p.text == "" {
			return fmt.Errorf("objective text is required")
		}
		g, err := store.Edit(p.sessionID, p.text)
		if err != nil {
			return err
		}
		fmt.Printf("goal updated (active)\nsession=%s\n", g.SessionID)
		return nil
	case "complete":
		g, err := store.Complete(p.sessionID)
		if err != nil {
			return err
		}
		fmt.Printf("goal complete\nsession=%s\n", g.SessionID)
		return nil
	case "clear":
		if err := store.Clear(p.sessionID); err != nil {
			return err
		}
		fmt.Printf("goal cleared\nsession=%s\n", p.sessionID)
		return nil
	case "status":
		g, err := store.Load(p.sessionID)
		if err != nil {
			return err
		}
		fmt.Print(store.FormatStatus(g))
		return nil
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}
