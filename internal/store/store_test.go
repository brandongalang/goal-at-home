package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetEditCompleteClear(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	id := "test-session-1"

	g, err := Set(id, "Build feature X")
	if err != nil {
		t.Fatal(err)
	}
	if g.Status != StatusActive || g.Text != "Build feature X" {
		t.Fatalf("unexpected goal: %+v", g)
	}

	g, err = Edit(id, "Build feature Y")
	if err != nil {
		t.Fatal(err)
	}
	if g.Text != "Build feature Y" {
		t.Fatalf("edit failed: %+v", g)
	}

	g, err = Complete(id)
	if err != nil {
		t.Fatal(err)
	}
	if g.Status != StatusComplete || g.CompletedAt == nil {
		t.Fatalf("complete failed: %+v", g)
	}

	if err := Clear(id); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(home, ".cursor", "goals", id+".json")
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Fatalf("expected file removed, stat err=%v", err)
	}
}
