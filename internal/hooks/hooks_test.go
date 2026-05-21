package hooks

import (
	"strings"
	"testing"

	"github.com/brandongalang/goal-at-home/internal/store"
)

func TestDetectFormat(t *testing.T) {
	t.Setenv("GOAL_HOOK_AGENT", "")
	if detectFormat("preToolUse") != formatCursor {
		t.Fatal("expected cursor")
	}
	if detectFormat("PreToolUse") != formatClaudeCodex {
		t.Fatal("expected claude/codex")
	}
	if detectFormat("BeforeTool") != formatGemini {
		t.Fatal("expected gemini")
	}
	t.Setenv("GOAL_HOOK_AGENT", "codex")
	if detectFormat("anything") != formatClaudeCodex {
		t.Fatal("GOAL_HOOK_AGENT override")
	}
}

func TestInjectSessionID(t *testing.T) {
	got := injectSessionID(`goal set "Ship it"`, "conv-abc")
	if !strings.Contains(got, `--session-id conv-abc`) {
		t.Fatalf("got %q", got)
	}
	if injectSessionID("npm test", "conv-abc") != "npm test" {
		t.Fatal("should not modify non-goal commands")
	}
	if injectSessionID(`goal status --session-id x`, "conv-abc") != `goal status --session-id x` {
		t.Fatal("should not double-inject")
	}
}

func TestIsShellTool(t *testing.T) {
	for _, name := range []string{"Shell", "Bash", "run_shell_command"} {
		if !isShellTool(name) {
			t.Fatalf("%s should be shell", name)
		}
	}
	if isShellTool("Edit") {
		t.Fatal("Edit should not match")
	}
}

func TestStopFollowupWhenActive(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if _, err := store.Set("conv-stop", "Finish the migration"); err != nil {
		t.Fatal(err)
	}
	g, err := store.Load("conv-stop")
	if err != nil {
		t.Fatal(err)
	}
	msg := strings.TrimSpace(strings.Replace(followupTemplate, "%s", g.Text, 1))
	if !strings.Contains(msg, "Finish the migration") {
		t.Fatalf("unexpected message: %q", msg)
	}
}
