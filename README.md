# goal-at-home

Session goal enforcement for coding agents. Set a goal at the start of sustained work; the agent cannot end the session until it runs `goal complete` (where hooks are supported).

**Repository:** https://github.com/brandongalang/goal-at-home

Agents: read [AGENTS.md](./AGENTS.md) first for install steps tailored to coding agents.

## How it works

1. Agent runs `goal set "objective with clear done criteria"`.
2. A **pre-tool hook** injects `--session-id` so state is tied to the chat.
3. When the agent loop ends, a **stop hook** checks the session goal file.
4. If status is `active`, the agent gets a follow-up reprompt (up to `loop_limit`, default 10).
5. Agent runs `goal complete` when truly done; the next stop is allowed.

User **Stop** leaves the goal **active** so the next turn resumes enforcement.

## Install

### Cursor (full — binary + hooks + skill)

Requires `go`, `jq`, and `~/.local/bin` on your `PATH`.

```bash
git clone https://github.com/brandongalang/goal-at-home.git
cd goal-at-home
./install.sh
```

Restart Cursor after install.

### Any agent (skill only — soft guidance)

Uses the [skills.sh](https://skills.sh) ecosystem ([Vercel skills CLI](https://github.com/vercel-labs/skills)):

```bash
npx skills add brandongalang/goal-at-home --skill goal-enforcement -g -y
```

This copies the skill into your agent’s skills directory. It does **not** install hooks or the `goal` binary.

### Install by agent

| Agent | Full enforcement | Skill path (global) |
|-------|-------------------|---------------------|
| **Cursor** | `./install.sh` | `~/.cursor/skills/goal-enforcement/` |
| **Claude Code** | Planned | `~/.claude/skills/goal-enforcement/` |
| **Codex** | Planned | `~/.codex/skills/` or `~/.agents/skills/` |
| **Gemini CLI** | Planned | `~/.gemini/skills/` |
| **Antigravity** | Not available (no stop hook) | `~/.gemini/antigravity/skills/` |
| **Windsurf** | Not available | `~/.codeium/windsurf/skills/` |
| **Pi / OpenCode** | Extension/plugin TBD | `.agents/skills/` |

See [docs/SUPPORTED_AGENTS.md](./docs/SUPPORTED_AGENTS.md) for tiers, Pi, OpenCode, and Antigravity details.

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

Point it at this repo and ask it to run install:

> Clone https://github.com/brandongalang/goal-at-home and run `./install.sh` for Cursor, or `npx skills add brandongalang/goal-at-home --skill goal-enforcement -g -y` for skill-only. Follow AGENTS.md.

Optional Cursor user rule:

> For sustained agent work, run `goal set` at the start and `goal complete` only when every requirement is satisfied.

## License

MIT
