# Custom agent setup (hooks + CLI)

Use this when your agent supports **pre-tool shell hooks** and a **turn-end stop hook** (or equivalent), but there is no `goal install <agent>` yet.

Goal-at-home needs three pieces:

1. **`goal` binary** on `PATH` (`goal install` or `go run . install` from the repo)
2. **Pre-tool hook** — inject `--session-id` on `goal *` shell commands
3. **Stop hook** — reprompt when `~/.cursor/goals/<session>.json` (or `$GOAL_HOME/goals/`) is still `active`

The bundled skills (**goalcraft-at-home**, **goal-enforcement**) are optional for crafting objectives but **not** sufficient alone.

## Hook contract (what `goal hook` implements)

| Hook command | When to run | stdin (minimum) | stdout when goal active |
|--------------|-------------|-----------------|-------------------------|
| `goal hook pre-tool-use` | Before shell/bash tool | `session_id` or `conversation_id`, `tool_name`, `tool_input.command` | Rewritten command with `--session-id` |
| `goal hook stop` | End of agent turn | `session_id` or `conversation_id` | Reprompt text (format depends on agent) |

`goal hook` auto-detects wire format from `hook_event_name` (or `GOAL_HOOK_AGENT`):

| Agent | Pre-tool event | Stop event | `GOAL_HOOK_AGENT` |
|-------|----------------|------------|-------------------|
| Cursor | `preToolUse` | `stop` | `cursor` |
| Claude Code | `PreToolUse` | `Stop` | `claude` |
| Codex CLI | `PreToolUse` | `Stop` | `codex` |
| Gemini CLI | `BeforeTool` | `AfterAgent` | `gemini` |

Shell tool names recognized: `Shell` (Cursor), `Bash` (Claude/Codex), `run_shell_command` (Gemini).

## Build the CLI

```bash
git clone https://github.com/brandongalang/goal-at-home.git
cd goal-at-home
go build -o ~/.local/bin/goal .
```

Ensure `~/.local/bin` is on your `PATH`.

Optional: set `GOAL_HOME` to store goals outside `~/.cursor` (e.g. `export GOAL_HOME=$HOME/.config/goal-at-home`).

## Example hook configs

Copy and merge into your agent’s hooks file (preserve existing hooks):

| Agent | Example file |
|-------|----------------|
| Cursor | [examples/hooks/cursor.hooks.json](../examples/hooks/cursor.hooks.json) |
| Claude Code | [examples/hooks/claude-code.settings.json](../examples/hooks/claude-code.settings.json) |
| Codex CLI | [examples/hooks/codex.hooks.json](../examples/hooks/codex.hooks.json) |
| Gemini CLI | [examples/hooks/gemini.settings.json](../examples/hooks/gemini.settings.json) |

After editing hooks, **restart the agent** or reload hook settings.

## Per-agent notes

### Cursor

```bash
goal install --agent cursor   # or claude | codex | gemini | all
```

Merges hooks, installs skills, and appends a marked block to global instructions:

| Agent | Global instructions file |
|-------|--------------------------|
| Cursor | `~/.cursor/AGENTS.md` |
| Claude Code | `~/.claude/CLAUDE.md` |
| Codex | `~/.codex/AGENTS.md` |
| Gemini CLI | `~/.gemini/GEMINI.md` |

Use `goal install --skip-instructions` to skip. Re-running install updates the block in place.

See [README](../README.md).

### Claude Code

- Hooks live in `~/.claude/settings.json` or project `.claude/settings.json` under `hooks`.
- Use matcher `Bash` for `PreToolUse`.
- Stop hook uses `decision: "block"` + `reason` (handled by `goal hook stop`).
- Docs: [Claude Code hooks](https://code.claude.com/docs/en/hooks)

### Codex CLI

- Hooks: `~/.codex/hooks.json` or repo `.codex/hooks.json`.
- Matcher `Bash` for `PreToolUse`; `Stop` has no matcher.
- Stop continuation uses `decision: "block"` + `reason`.
- Docs: [Codex hooks](https://developers.openai.com/codex/hooks)

### Gemini CLI

- Hooks in `settings.json` → `hooks.BeforeTool` (matcher `run_shell_command`) and `hooks.AfterAgent`.
- Map **AfterAgent** to `goal hook stop` (not a separate binary).
- Docs: [Gemini CLI hooks](https://geminicli.com/docs/hooks/)

## Install skills (any agent)

Copy into the agent’s skill directory if it supports skills:

```bash
SKILLS_DIR=~/.claude/skills   # or ~/.cursor/skills, etc.
mkdir -p "$SKILLS_DIR/goalcraft-at-home" "$SKILLS_DIR/goal-enforcement"
cp skill/goalcraft-at-home/SKILL.md "$SKILLS_DIR/goalcraft-at-home/"
cp skill/goal-enforcement/SKILL.md "$SKILLS_DIR/goal-enforcement/"
```

## Verify

```bash
export GOAL_HOME=/tmp/goal-test   # optional isolated store
goal set --session-id test-session "Finish the migration"
goal status --session-id test-session
```

In the agent: run `goal set "test"`, try to end the turn without `goal complete` — you should get a continuation reprompt. Then `goal complete`.

## Agent-assisted setup

Paste the prompt in [docs/prompts/configure-agent-hooks.md](./prompts/configure-agent-hooks.md) into your agent. It will inspect your agent’s hook format, build `goal`, merge hook entries, and run a smoke test.

## When this will not work

See [SUPPORTED_AGENTS.md](./SUPPORTED_AGENTS.md). Agents without turn-level stop/reprompt (Antigravity, Windsurf, most Copilot modes) cannot enforce goals yet — use Pi/OpenCode extensions instead of shell hooks.
