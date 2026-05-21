# Supported agents

Goal-at-home is **CLI + hooks + skill**. Copying only `SKILL.md` does not enforce anything — the stop hook is what prevents the agent from ending while a goal is active.

## Requirements

1. **Pre-tool (shell)** — inject `--session-id` on `goal *` commands
2. **Stop / end-of-turn** — reprompt or block when a goal is still `active`
3. **`goal` binary** on PATH

## Tier 1 — Supported (CLI + hooks)

| Agent | Install | Pre-tool | Stop | Status |
|-------|---------|----------|------|--------|
| [Cursor](https://cursor.com) | `goal install --agent cursor` | `preToolUse` (Shell) | `stop` | **Available** |
| [Claude Code](https://code.claude.com) | `goal install --agent claude` | `PreToolUse` (Bash) | `Stop` | **Available** |
| [Codex CLI](https://developers.openai.com/codex) | `goal install --agent codex` | `PreToolUse` (Bash) | `Stop` | **Available** |
| [Gemini CLI](https://geminicli.com) | `goal install --agent gemini` | `BeforeTool` (`run_shell_command`) | `AfterAgent` | **Available** |

`goal hook` detects wire format from `hook_event_name` or `GOAL_HOOK_AGENT` (`cursor`, `claude`, `codex`, `gemini`). Agent setup prompt: [prompts/configure-agent-hooks.md](./prompts/configure-agent-hooks.md).

## Tier 2 — Possible via native extension (not shell hooks)

| Agent | Path forward |
|-------|----------------|
| [Pi](https://pi.dev) | TypeScript extension calling the same `goal` CLI + store |
| [OpenCode](https://opencode.ai) | Plugin on tool/stop events |

Still requires CLI + hard stop semantics — not skill-only.

## Not supported (yet)

These agents lack a Cursor-style stop reprompt we can wire today:

| Agent | Why |
|-------|-----|
| [Antigravity](https://antigravity.google) | Workflows/rules only; no shell stop hook ([forum](https://discuss.ai.google.dev/t/hooks-in-antigravity/120458)) |
| [Windsurf](https://windsurf.com) | Cascade hooks; no turn-level stop reprompt |
| GitHub Copilot, Continue, Aider, most assistants | No compatible hook surface |

Do not point users at goal-at-home for these until the agent exposes stop/reprompt hooks (or a native extension path).

## Compared to built-in `/goal`

Claude Code and Codex ship native goal/loop features. Goal-at-home targets:

- File-backed objectives (`goal set` / `goal complete`)
- Same CLI across hook-capable agents (see [CUSTOM_AGENT_SETUP.md](./CUSTOM_AGENT_SETUP.md))
- Deterministic enforcement via hooks
