# goal-at-home

Session goal enforcement for coding agents. Set a goal at the start of sustained work; the agent cannot end the session until it runs `goal complete` — enforced by the **CLI + hooks**, not prompts alone.

**Adapts the [Goalcraft](https://github.com/grp06/goalcraft) pattern for better goals** — evidence-based objectives (Destination, Verification, Done/stop, etc.) via bundled **goalcraft-at-home**, then enforced with `goal set` and stop hooks. See [integration](./docs/GOALCRAFT_INTEGRATION.md) · [attribution](./docs/GOALCRAFT_ATTRIBUTION.md).

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

## Install

Requires `go` and `~/.local/bin` on your `PATH`.

```bash
git clone https://github.com/brandongalang/goal-at-home.git
cd goal-at-home
go run . install --agent cursor   # first time from clone; then: goal install --agent cursor
```

Other agents:

```bash
goal install --agent claude    # Claude Code
goal install --agent codex     # Codex CLI
goal install --agent gemini    # Gemini CLI
goal install --agent all       # every supported agent
goal install list              # show agent names
```

`goal install` builds the binary, merges hooks, copies skills, and appends a short **Session goals** block to the agent’s global instructions file (`AGENTS.md`, `CLAUDE.md`, or `GEMINI.md`). No `jq` required. Restart the agent after install.

This installs:

- `~/.local/bin/goal` — CLI
- `~/.cursor/hooks.json` — `preToolUse` + `stop` hooks (existing hooks preserved)
- `~/.cursor/skills/goalcraft-at-home/SKILL.md` — craft objectives before `goal set`
- `~/.cursor/skills/goal-enforcement/SKILL.md` — enforce until `goal complete`
- `~/.cursor/goals/` — per-session state

## Supported agents

| Agent | Install command |
|-------|-----------------|
| **Cursor** | `goal install --agent cursor` |
| **Claude Code** | `goal install --agent claude` |
| **Codex CLI** | `goal install --agent codex` |
| **Gemini CLI** | `goal install --agent gemini` |
| **All of the above** | `goal install --agent all` |
| **Antigravity, Windsurf, Copilot, …** | Not supported — no turn-end stop hook to wire |

See [docs/SUPPORTED_AGENTS.md](./docs/SUPPORTED_AGENTS.md). For troubleshooting or custom paths, [docs/CUSTOM_AGENT_SETUP.md](./docs/CUSTOM_AGENT_SETUP.md).

## Goalcraft at home (bundled)

This repo includes **goalcraft-at-home** — an adapted [Goalcraft](https://github.com/grp06/goalcraft) workflow for `goal set` (not Codex `/goal`). `goal install` installs both skills plus hooks and global agent instructions.

Details: [docs/GOALCRAFT_INTEGRATION.md](./docs/GOALCRAFT_INTEGRATION.md) · [attribution](./docs/GOALCRAFT_ATTRIBUTION.md)

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

**Cursor:** Clone https://github.com/brandongalang/goal-at-home, run `goal install --agent cursor`, restart Cursor.

**Other agents:** `goal install --agent <name>`. If install fails on an unusual setup, use [docs/prompts/configure-agent-hooks.md](./docs/prompts/configure-agent-hooks.md).

Do not use skill-only install — enforcement requires the CLI and hooks.

Optional Cursor user rule:

> For sustained agent work, run `goal set` at the start and `goal complete` only when every requirement is satisfied.

## License

MIT
