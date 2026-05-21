package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	instructionsStart = "<!-- goal-at-home:start -->"
	instructionsEnd   = "<!-- goal-at-home:end -->"
)

func appendInstructions(paths targetPaths) error {
	block := instructionsBlock(paths.instructionsFile)
	existing, err := os.ReadFile(paths.instructionsFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	content := string(existing)
	if strings.Contains(content, instructionsStart) {
		content = removeMarkedBlock(content)
	}
	content = strings.TrimRight(content, "\n")
	if content != "" {
		content += "\n\n"
	}
	content += block + "\n"
	if err := os.MkdirAll(filepath.Dir(paths.instructionsFile), 0o755); err != nil {
		return err
	}
	return os.WriteFile(paths.instructionsFile, []byte(content), 0o644)
}

func instructionsBlock(path string) string {
	intro := ""
	if line := instructionsDocName(path); line != "" {
		intro = line + "\n\n"
	}
	return fmt.Sprintf(`%s
%s## Session goals (goal-at-home)

Enforcement requires the **goal** CLI and hooks (`+"`goal install`"+`). Skills alone are not enough.

**Sustained work:** **goalcraft-at-home** → `+"`goal set \"…\"`"+` → work → `+"`goal complete`"+` when every requirement is met.

- Do not pass `+"`--session-id`"+`; hooks inject it.
- Do not stop or claim done while a goal is **active**.
- Casual Q&A does not need `+"`goal set`"+`.

https://github.com/brandongalang/goal-at-home
%s`, instructionsStart, intro, instructionsEnd)
}

func instructionsDocName(path string) string {
	base := path
	if i := strings.LastIndex(path, "/"); i >= 0 {
		base = path[i+1:]
	}
	switch base {
	case "CLAUDE.md":
		return "Applies in every Claude Code session (global " + base + ")."
	case "GEMINI.md":
		return "Applies in every Gemini CLI session (global " + base + ")."
	case "AGENTS.md":
		return "Applies in every session that reads global " + base + "."
	default:
		return ""
	}
}

func removeMarkedBlock(content string) string {
	start := strings.Index(content, instructionsStart)
	if start < 0 {
		return content
	}
	end := strings.Index(content[start:], instructionsEnd)
	if end < 0 {
		return strings.TrimSpace(content[:start])
	}
	end = start + end + len(instructionsEnd)
	before := strings.TrimRight(content[:start], "\n")
	after := strings.TrimLeft(content[end:], "\n")
	if before == "" {
		return after
	}
	if after == "" {
		return before
	}
	return before + "\n\n" + after
}
