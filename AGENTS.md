# Goal-at-home — agent install guide

Session goal enforcement: agents run `goal set` at the start of sustained work and `goal complete` when done. While a goal is **active**, supported agents cannot end the turn without continuing (via hooks).

Repository: https://github.com/brandongalang/goal-at-home

## Quick install (Cursor — full enforcement)

Run from the repo root (requires `go`, `jq`, and `~/.local/bin` on PATH):

```bash
git clone https://github.com/brandongalang/goal-at-home.git
cd goal-at-home
./install.sh
```

Then **restart Cursor** so hooks reload.

What this installs:

- `~/.local/bin/goal` — CLI
- `~/.cursor/hooks.json` — merged `preToolUse` + `stop` hooks (existing hooks preserved)
- `~/.cursor/skills/goal-enforcement/SKILL.md` — skill for the agent
- `~/.cursor/goals/` — per-session goal state

## Skill-only install (any agent)

For agents without hook adapters (Antigravity, Pi, OpenCode, Copilot, etc.):

```bash
npx skills add brandongalang/goal-at-home --skill goal-enforcement -g -y
```

Or clone this repo and copy `skill/goal-enforcement/SKILL.md` into your agent’s skills folder (see [README.md](./README.md#install-by-agent)).

## Commands (after full install)

In Cursor, do **not** pass `--session-id` — hooks inject it on shell `goal` commands.

```bash
goal set "Concrete objective with clear done criteria"
goal status
goal edit "Updated objective if scope changed"
goal complete   # only when every requirement is satisfied
goal clear      # user cancelled or goal should be dropped
```

## When to use

- Multi-step implementation, migrations, refactors, or anything that must run to completion
- User says: set a goal, use goal enforcement, or do not stop until done

## Rules

1. Run `goal set "..."` at the start of substantive work.
2. Do not claim the task is finished without `goal complete` while a goal is active.
3. Do not run `goal complete` early — the stop hook will reprompt until the goal is real.
4. Casual Q&A does not need `goal set`.

## Agent-specific notes

| Agent | Full install | Skill only |
|-------|--------------|------------|
| Cursor | `./install.sh` | `npx skills add brandongalang/goal-at-home -g -y` |
| Claude Code | hooks adapter planned | `~/.claude/skills/goal-enforcement/` |
| Codex | hooks adapter planned | `~/.codex/skills/` or `.agents/skills/` |
| Gemini CLI | hooks adapter planned | `~/.gemini/skills/` |
| Antigravity | not supported (no stop hook) | `~/.gemini/antigravity/skills/` |
| Pi / OpenCode | extension/plugin TBD | `.agents/skills/` or agent-specific path |

See [docs/SUPPORTED_AGENTS.md](./docs/SUPPORTED_AGENTS.md) for the full matrix.

## Verify install (Cursor)

```bash
goal status    # may error outside a session — expected
which goal     # should be ~/.local/bin/goal
```

Set a test goal in an agent shell, stop without completing — you should get a follow-up reprompt. Then `goal complete`.
