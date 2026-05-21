package hooks

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/brandongalang/goal-at-home/internal/store"
)

type agentFormat int

const (
	formatCursor agentFormat = iota
	formatClaudeCodex
	formatGemini
)

type baseInput struct {
	ConversationID string `json:"conversation_id"`
	SessionID      string `json:"session_id"`
	HookEventName  string `json:"hook_event_name"`
}

type preToolUseInput struct {
	baseInput
	ToolName  string          `json:"tool_name"`
	ToolInput json.RawMessage `json:"tool_input"`
}

type toolInputCommand struct {
	Command string `json:"command"`
}

type stopInput struct {
	baseInput
	Status         string `json:"status"`
	LoopCount      int    `json:"loop_count"`
	StopHookActive bool   `json:"stop_hook_active"`
}

type cursorPreOutput struct {
	UpdatedInput *struct {
		Command string `json:"command"`
	} `json:"updated_input,omitempty"`
}

type cursorStopOutput struct {
	FollowupMessage string `json:"followup_message,omitempty"`
}

type claudeCodexPreOutput struct {
	HookSpecificOutput *struct {
		HookEventName      string `json:"hookEventName"`
		PermissionDecision string `json:"permissionDecision"`
		UpdatedInput       *struct {
			Command string `json:"command"`
		} `json:"updatedInput,omitempty"`
	} `json:"hookSpecificOutput,omitempty"`
}

type claudeCodexStopOutput struct {
	Decision string `json:"decision,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

type geminiPreOutput struct {
	HookSpecificOutput *struct {
		ToolInput *struct {
			Command string `json:"command"`
		} `json:"tool_input,omitempty"`
	} `json:"hookSpecificOutput,omitempty"`
}

var goalCommand = regexp.MustCompile(`^goal\s+`)

var shellToolNames = map[string]bool{
	"Shell":             true,
	"Bash":              true,
	"run_shell_command": true,
}

func sessionID(in baseInput) string {
	if in.ConversationID != "" {
		return in.ConversationID
	}
	return in.SessionID
}

func detectFormat(hookEventName string) agentFormat {
	switch strings.ToLower(os.Getenv("GOAL_HOOK_AGENT")) {
	case "cursor":
		return formatCursor
	case "claude", "codex":
		return formatClaudeCodex
	case "gemini":
		return formatGemini
	}
	switch hookEventName {
	case "preToolUse", "stop":
		return formatCursor
	case "PreToolUse", "Stop":
		return formatClaudeCodex
	case "BeforeTool", "AfterAgent":
		return formatGemini
	default:
		return formatCursor
	}
}

func injectSessionID(command, sessionID string) string {
	cmd := strings.TrimSpace(command)
	if sessionID == "" || !goalCommand.MatchString(cmd) || strings.Contains(cmd, "--session-id") {
		return command
	}
	return cmd + " --session-id " + sessionID
}

func toolCommand(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var t toolInputCommand
	if err := json.Unmarshal(raw, &t); err != nil {
		return ""
	}
	return t.Command
}

func isShellTool(name string) bool {
	return shellToolNames[name]
}

func readInput[T any]() (T, error) {
	var out T
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return out, err
	}
	if len(data) == 0 {
		return out, fmt.Errorf("empty hook input")
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return out, err
	}
	return out, nil
}

func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(v)
}

func writePreOutput(format agentFormat, command string) error {
	switch format {
	case formatClaudeCodex:
		return writeJSON(claudeCodexPreOutput{
			HookSpecificOutput: &struct {
				HookEventName      string `json:"hookEventName"`
				PermissionDecision string `json:"permissionDecision"`
				UpdatedInput       *struct {
					Command string `json:"command"`
				} `json:"updatedInput,omitempty"`
			}{
				HookEventName:      "PreToolUse",
				PermissionDecision: "allow",
				UpdatedInput:       &struct{ Command string `json:"command"` }{Command: command},
			},
		})
	case formatGemini:
		return writeJSON(geminiPreOutput{
			HookSpecificOutput: &struct {
				ToolInput *struct {
					Command string `json:"command"`
				} `json:"tool_input,omitempty"`
			}{
				ToolInput: &struct{ Command string `json:"command"` }{Command: command},
			},
		})
	default:
		return writeJSON(cursorPreOutput{
			UpdatedInput: &struct {
				Command string `json:"command"`
			}{Command: command},
		})
	}
}

func writeEmptyPreOutput(format agentFormat) error {
	switch format {
	case formatClaudeCodex:
		return writeJSON(claudeCodexPreOutput{})
	case formatGemini:
		return writeJSON(geminiPreOutput{})
	default:
		return writeJSON(cursorPreOutput{})
	}
}

// PreToolUse injects --session-id into goal shell commands (Cursor, Claude, Codex, Gemini).
func PreToolUse() error {
	in, err := readInput[preToolUseInput]()
	if err != nil {
		return err
	}
	format := detectFormat(in.HookEventName)
	if !isShellTool(in.ToolName) {
		return writeEmptyPreOutput(format)
	}
	cmd := toolCommand(in.ToolInput)
	updated := injectSessionID(cmd, sessionID(in.baseInput))
	if updated == cmd {
		return writeEmptyPreOutput(format)
	}
	return writePreOutput(format, updated)
}

const followupTemplate = `The session goal is not complete. Continue working toward it, or run goal complete only when every requirement is satisfied.

Goal: %s

Do not stop until the goal is done and logged with goal complete.`

func writeStopOutput(format agentFormat, message string) error {
	switch format {
	case formatClaudeCodex, formatGemini:
		return writeJSON(claudeCodexStopOutput{
			Decision: "block",
			Reason:   message,
		})
	default:
		return writeJSON(cursorStopOutput{FollowupMessage: message})
	}
}

func writeEmptyStopOutput(format agentFormat) error {
	switch format {
	case formatClaudeCodex, formatGemini:
		return writeJSON(claudeCodexStopOutput{})
	default:
		return writeJSON(cursorStopOutput{})
	}
}

// Stop reprompts when an active goal exists and the agent did not complete it.
func Stop() error {
	in, err := readInput[stopInput]()
	if err != nil {
		return err
	}
	format := detectFormat(in.HookEventName)
	if in.Status == "aborted" {
		return writeEmptyStopOutput(format)
	}
	g, err := store.Load(sessionID(in.baseInput))
	if err != nil {
		return writeEmptyStopOutput(format)
	}
	if g.Status != store.StatusActive {
		return writeEmptyStopOutput(format)
	}
	return writeStopOutput(format, fmt.Sprintf(followupTemplate, g.Text))
}
