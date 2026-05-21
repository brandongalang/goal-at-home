# Goal-at-home — agent install guide

Session goal enforcement requires **all three**: the `goal` CLI, lifecycle **hooks**, and the bundled skill. Instructions alone are not enough — without hooks, the agent can stop without running `goal complete`.

Repository: https://github.com/brandongalang/goal-at-home

## Install (Cursor)

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
- `~/.cursor/skills/goal-enforcement/SKILL.md` — skill (installed with the CLI, not separately)
- `~/.cursor/goals/` — per-session goal state

Do **not** install only `SKILL.md` via `npx skills add` or by copying the skill folder. That skips enforcement.

## Commands

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

## Other agents

| Agent | Supported? |
|-------|------------|
| Cursor | Yes — `./install.sh` |
| Claude Code, Codex, Gemini CLI | Planned (`goal install` adapters) |
| Antigravity, Windsurf, Copilot, Pi, OpenCode | No — no stop-hook enforcement yet |

See [docs/SUPPORTED_AGENTS.md](./docs/SUPPORTED_AGENTS.md).

## Verify install

```bash
which goal     # ~/.local/bin/goal
```

In an agent shell: `goal set "test"`, stop without completing — you should get a follow-up reprompt. Then `goal complete`.
