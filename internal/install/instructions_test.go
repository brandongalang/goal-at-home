package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppendInstructionsIdempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "CLAUDE.md")
	paths := targetPaths{instructionsFile: path}

	if err := appendInstructions(paths); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(first), instructionsStart) {
		t.Fatal("missing marker")
	}
	if err := appendInstructions(paths); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(second), instructionsStart) != 1 {
		t.Fatalf("expected one block, got:\n%s", second)
	}
}

func TestRemoveMarkedBlock(t *testing.T) {
	in := "before\n\n" + instructionsStart + "\nold\n" + instructionsEnd + "\n\nafter"
	out := removeMarkedBlock(in)
	if strings.Contains(out, instructionsStart) || strings.Contains(out, "old") {
		t.Fatalf("got %q", out)
	}
	if !strings.Contains(out, "before") || !strings.Contains(out, "after") {
		t.Fatalf("got %q", out)
	}
}
