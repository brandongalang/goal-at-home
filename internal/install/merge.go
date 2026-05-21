package install

func mergeCursorHooks(root map[string]any) {
	hooks, _ := root["hooks"].(map[string]any)
	if hooks == nil {
		hooks = map[string]any{}
		root["hooks"] = hooks
	}
	if _, ok := root["version"]; !ok {
		root["version"] = 1
	}
	pre, _ := hooks["preToolUse"].([]any)
	pre = prependIfMissing(pre, map[string]any{
		"command": goalPreCmd,
		"matcher": "Shell",
	}, goalPreCmd)
	hooks["preToolUse"] = pre

	stop, _ := hooks["stop"].([]any)
	stop = appendIfMissing(stop, map[string]any{
		"command":    goalStopCmd,
		"loop_limit": 10,
	}, goalStopCmd)
	hooks["stop"] = stop
}

func mergeCodexHooks(root map[string]any) {
	hooks, _ := root["hooks"].(map[string]any)
	if hooks == nil {
		hooks = map[string]any{}
		root["hooks"] = hooks
	}
	mergePascalShellHooks(hooks)
}

func mergeClaudeHooks(root map[string]any) {
	hooks, _ := root["hooks"].(map[string]any)
	if hooks == nil {
		hooks = map[string]any{}
		root["hooks"] = hooks
	}
	mergePascalShellHooks(hooks)
}

func mergePascalShellHooks(hooks map[string]any) {
	preEntry := map[string]any{
		"matcher": "Bash",
		"hooks": []any{
			map[string]any{
				"type":    "command",
				"command": goalPreCmd,
			},
		},
	}
	pre, _ := hooks["PreToolUse"].([]any)
	hooks["PreToolUse"] = prependPascalHookGroup(pre, preEntry, goalPreCmd)

	stopEntry := map[string]any{
		"hooks": []any{
			map[string]any{
				"type":    "command",
				"command": goalStopCmd,
			},
		},
	}
	stop, _ := hooks["Stop"].([]any)
	hooks["Stop"] = appendPascalHookGroup(stop, stopEntry, goalStopCmd)
}

func mergeGeminiHooks(root map[string]any) {
	hooks, _ := root["hooks"].(map[string]any)
	if hooks == nil {
		hooks = map[string]any{}
		root["hooks"] = hooks
	}

	beforeEntry := map[string]any{
		"matcher": "run_shell_command",
		"hooks": []any{
			map[string]any{
				"type":    "command",
				"command": goalPreCmd,
				"name":    "goal-at-home-pre-tool",
			},
		},
	}
	before, _ := hooks["BeforeTool"].([]any)
	hooks["BeforeTool"] = prependPascalHookGroup(before, beforeEntry, goalPreCmd)

	afterEntry := map[string]any{
		"hooks": []any{
			map[string]any{
				"type":    "command",
				"command": goalStopCmd,
				"name":    "goal-at-home-stop",
			},
		},
	}
	after, _ := hooks["AfterAgent"].([]any)
	hooks["AfterAgent"] = appendPascalHookGroup(after, afterEntry, goalStopCmd)
}

func prependIfMissing(list []any, item map[string]any, command string) []any {
	filtered := removeByCommand(list, command)
	return append([]any{item}, filtered...)
}

func appendIfMissing(list []any, item map[string]any, command string) []any {
	return append(removeByCommand(list, command), item)
}

func prependPascalHookGroup(list []any, item map[string]any, command string) []any {
	filtered := removePascalGroupsWithCommand(list, command)
	return append([]any{item}, filtered...)
}

func appendPascalHookGroup(list []any, item map[string]any, command string) []any {
	filtered := removePascalGroupsWithCommand(list, command)
	return append(filtered, item)
}

func removeByCommand(list []any, command string) []any {
	var out []any
	for _, raw := range list {
		m, ok := raw.(map[string]any)
		if !ok {
			out = append(out, raw)
			continue
		}
		if cmd, _ := m["command"].(string); cmd == command {
			continue
		}
		out = append(out, raw)
	}
	return out
}

func removePascalGroupsWithCommand(list []any, command string) []any {
	var out []any
	for _, raw := range list {
		group, ok := raw.(map[string]any)
		if !ok {
			out = append(out, raw)
			continue
		}
		inner, _ := group["hooks"].([]any)
		if pascalGroupHasCommand(inner, command) {
			continue
		}
		out = append(out, raw)
	}
	return out
}

func pascalGroupHasCommand(inner []any, command string) bool {
	for _, raw := range inner {
		m, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		if cmd, _ := m["command"].(string); cmd == command {
			return true
		}
	}
	return false
}
