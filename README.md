# goal-at-home

Session goal enforcement for coding agents. Set a goal at the start of sustained work; the agent cannot end the session until it runs `goal complete` — enforced by the **CLI + hooks**, not prompts alone.

**Repository:** https://github.com/brandongalang/goal-at-home

Agents: read [AGENTS.md](./AGENTS.md) first.

## How it works

1. Agent runs `goal set "objective with clear done criteria"`.
2. A **pre-tool hook** injects `--session-id` so state is tied to the chat.
3. When the agent loop ends, a **stop hook** checks the session goal file.
4. If status is `active`, the agent gets a follow-up reprompt (up to `loop_limit`, default 10).
5. Agent runs `goal complete` when truly done; the next stop is allowed.

User **Stop** leaves the goal **active** so the next turn resumes enforcement.

Installing only a skill file (instructions without the binary and hooks) does **not** provide enforcement — the agent can ignore it. Use the full install below.

## Install (Cursor)

Requires `go`, `jq`, and `~/.local/bin` on your `PATH`.

```bash
git clone https://github.com/brandongalang/goal-at-home.git
cd goal-at-home
./install.sh
```

Restart Cursor after install.

This installs:

- `~/.local/bin/goal` — CLI
- `~/.cursor/hooks.json` — `preToolUse` + `stop` hooks (existing hooks preserved)
- `~/.cursor/skills/goal-enforcement/SKILL.md` — agent instructions (bundled with install)
- `~/.cursor/goals/` — per-session state

## Supported agents

| Agent | Status |
|-------|--------|
| **Cursor** | `./install.sh` today |
| **Claude Code, Codex, Gemini CLI** | Hook adapters planned |
| **Antigravity, Windsurf, Copilot, Pi, OpenCode, …** | Not supported until we ship an adapter with CLI + stop enforcement |

See [docs/SUPPORTED_AGENTS.md](./docs/SUPPORTED_AGENTS.md).

## Commands

| Command | Description |
|---------|-------------|
| `goal set "<text>"` | Create or replace the session goal (hooks add `--session-id`) |
| `goal edit "<text>"` | Update an active goal |
| `goal complete` | Mark complete (allows stop) |
| `goal clear` | Remove the goal without completing |
| `goal status` | Show current goal |

With hooks installed, agents should omit `--session-id` on shell commands.

## Hooks (Cursor)

| Hook | Command |
|------|---------|
| `preToolUse` (Shell) | `goal hook pre-tool-use` |
| `stop` | `goal hook stop` |

Example: [hooks.example.json](./hooks.example.json).

## Tell your agent

> Clone https://github.com/brandongalang/goal-at-home, run `./install.sh`, restart Cursor. Follow AGENTS.md. Do not use skill-only install — enforcement requires the CLI and hooks.

Optional Cursor user rule:

> For sustained agent work, run `goal set` at the start and `goal complete` only when every requirement is satisfied.

## License

MIT
