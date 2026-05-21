package hooks

import (
	"strings"
	"testing"

	"github.com/brandongalang/goal-at-home/internal/store"
)

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
