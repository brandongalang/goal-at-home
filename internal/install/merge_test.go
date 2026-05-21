package install

import (
	"encoding/json"
	"testing"
)

func TestMergeCursorHooks(t *testing.T) {
	root := map[string]any{}
	mergeCursorHooks(root)
	b, _ := json.Marshal(root)
	if !contains(string(b), goalPreCmd) || !contains(string(b), goalStopCmd) {
		t.Fatalf("missing goal hooks: %s", b)
	}
	mergeCursorHooks(root)
	b2, _ := json.Marshal(root)
	if countSubstring(string(b2), goalPreCmd) != 1 {
		t.Fatalf("duplicate pre hook: %s", b2)
	}
}

func TestMergeClaudeHooks(t *testing.T) {
	root := map[string]any{"other": true}
	mergeClaudeHooks(root)
	b, _ := json.Marshal(root)
	if !contains(string(b), "PreToolUse") || !contains(string(b), goalStopCmd) {
		t.Fatalf("missing hooks: %s", b)
	}
}

func TestMergeGeminiHooks(t *testing.T) {
	root := map[string]any{}
	mergeGeminiHooks(root)
	b, _ := json.Marshal(root)
	if !contains(string(b), "BeforeTool") || !contains(string(b), "AfterAgent") {
		t.Fatalf("missing hooks: %s", b)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func countSubstring(s, sub string) int {
	n := 0
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			n++
		}
	}
	return n
}
