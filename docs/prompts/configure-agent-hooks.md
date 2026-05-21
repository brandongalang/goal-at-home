# Prompt: configure goal-at-home for this agent

Copy everything below the line into your coding agent (Claude Code, Codex, Gemini CLI, or another hook-capable CLI).

---

You are configuring **goal-at-home** for the agent runtime I am using in this session.

Repository: https://github.com/brandongalang/goal-at-home

## Goal

Install full session goal enforcement: **CLI + pre-tool hook + stop hook + skills**. Skill-only install is invalid.

## Requirements

1. **Build** the `goal` binary and put it on my PATH (`~/.local/bin/goal` preferred).
2. **Wire hooks** so:
   - Before shell/bash commands, run `goal hook pre-tool-use` (injects `--session-id` on `goal *` commands).
   - At end of turn, run `goal hook stop` (reprompts while goal status is `active`).
3. **Install skills** `goalcraft-at-home` and `goal-enforcement` into this agent’s skill directory if supported.
4. **Do not** tell me to use Codex `/goal` or install external grp06/goalcraft for enforcement — use this repo’s bundled workflow.

## Steps you must perform

1. Clone or use an existing clone of goal-at-home. If missing: `git clone https://github.com/brandongalang/goal-at-home.git`.
2. Read `docs/CUSTOM_AGENT_SETUP.md` and `docs/SUPPORTED_AGENTS.md`.
3. Identify **this agent’s** hook config file path and JSON/TOML shape (Cursor `hooks.json`, Claude `settings.json`, Codex `hooks.json`, Gemini `settings.json`, etc.).
4. Compare with examples in `examples/hooks/` and merge goal hooks **without removing** my existing hooks.
5. Run `goal install --agent <name>` from the repo (or `go run . install --agent <name>` if `goal` is not on PATH yet).
6. Set `GOAL_HOME` only if goals should not live under `~/.cursor/goals` (explain if you change it).
7. If the agent uses different hook event names, set `GOAL_HOOK_AGENT` in hook commands only when auto-detect is insufficient:
   - `GOAL_HOOK_AGENT=cursor goal hook pre-tool-use`
   - `GOAL_HOOK_AGENT=claude goal hook stop`
   - Values: `cursor`, `claude`, `codex`, `gemini`
8. Restart or reload hooks per agent docs.
9. **Verify**: `goal set --session-id smoke-test "Say hello and complete"` then simulate or run stop hook behavior; run `goal complete --session-id smoke-test`.

## Hook commands (canonical)

| Purpose | Command |
|---------|---------|
| Pre-tool | `goal hook pre-tool-use` |
| Stop | `goal hook stop` |

Gemini: attach `BeforeTool` → pre-tool-use, `AfterAgent` → stop.

## Matcher cheat sheet

| Agent | Pre-tool matcher |
|-------|------------------|
| Cursor | `Shell` |
| Claude Code | `Bash` |
| Codex | `Bash` |
| Gemini CLI | `run_shell_command` |

## Deliverables

Reply with:

1. Agent name and hook file path(s) you modified
2. Exact diff or JSON you added
3. `which goal` output
4. Verification commands you ran and their output
5. Anything still manual (restart IDE, auth, etc.)

If this agent **cannot** block or continue the turn at end-of-response, say so clearly and point to `docs/SUPPORTED_AGENTS.md` — do not pretend skill text alone enforces goals.

Reference docs in-repo: `docs/CUSTOM_AGENT_SETUP.md`, `AGENTS.md`, `hooks.example.json`.
