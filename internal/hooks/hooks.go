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

type baseInput struct {
	ConversationID string `json:"conversation_id"`
	HookEventName  string `json:"hook_event_name"`
}

type preToolUseInput struct {
	baseInput
	ToolName  string `json:"tool_name"`
	ToolInput struct {
		Command string `json:"command"`
	} `json:"tool_input"`
}

type preToolUseOutput struct {
	UpdatedInput *struct {
		Command string `json:"command"`
	} `json:"updated_input,omitempty"`
}

type stopInput struct {
	baseInput
	Status    string `json:"status"`
	LoopCount int    `json:"loop_count"`
}

type stopOutput struct {
	FollowupMessage string `json:"followup_message,omitempty"`
}

var goalCommand = regexp.MustCompile(`^goal\s+`)

func injectSessionID(command, sessionID string) string {
	cmd := strings.TrimSpace(command)
	if sessionID == "" || !goalCommand.MatchString(cmd) || strings.Contains(cmd, "--session-id") {
		return command
	}
	return cmd + " --session-id " + sessionID
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

// PreToolUse injects --session-id into goal shell commands.
func PreToolUse() error {
	in, err := readInput[preToolUseInput]()
	if err != nil {
		return err
	}
	if in.ToolName != "Shell" {
		return writeJSON(preToolUseOutput{})
	}
	updated := injectSessionID(in.ToolInput.Command, in.ConversationID)
	if updated == in.ToolInput.Command {
		return writeJSON(preToolUseOutput{})
	}
	return writeJSON(preToolUseOutput{
		UpdatedInput: &struct {
			Command string `json:"command"`
		}{Command: updated},
	})
}

const followupTemplate = `The session goal is not complete. Continue working toward it, or run goal complete only when every requirement is satisfied.

Goal: %s

Do not stop until the goal is done and logged with goal complete.`

// Stop reprompts when an active goal exists and the agent did not complete it.
func Stop() error {
	in, err := readInput[stopInput]()
	if err != nil {
		return err
	}
	if in.Status == "aborted" {
		return writeJSON(stopOutput{})
	}
	g, err := store.Load(in.ConversationID)
	if err != nil {
		return writeJSON(stopOutput{})
	}
	if g.Status != store.StatusActive {
		return writeJSON(stopOutput{})
	}
	return writeJSON(stopOutput{
		FollowupMessage: fmt.Sprintf(followupTemplate, g.Text),
	})
}
