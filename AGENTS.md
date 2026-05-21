# Goal-at-home — agent install guide

Session goal enforcement requires **all three**: the `goal` CLI, lifecycle **hooks**, and the bundled skill. Instructions alone are not enough — without hooks, the agent can stop without running `goal complete`.

Repository: https://github.com/brandongalang/goal-at-home

## Install (Cursor)

From the repo root (requires `go` and `~/.local/bin` on PATH):

```bash
git clone https://github.com/brandongalang/goal-at-home.git
cd goal-at-home
go run . install --agent cursor   # first time; later: goal install --agent cursor
```

Then **restart Cursor** so hooks reload.

What `goal install --agent cursor` installs:

- `~/.local/bin/goal` — CLI
- `~/.cursor/hooks.json` — merged `preToolUse` + `stop` hooks (existing hooks preserved)
- `~/.cursor/skills/goalcraft-at-home/SKILL.md` — craft objectives (adapted Goalcraft)
- `~/.cursor/skills/goal-enforcement/SKILL.md` — enforce until `goal complete`
- `~/.cursor/goals/` — per-session goal state
- `~/.cursor/AGENTS.md` — appends a marked **Session goals** section (idempotent; use `--skip-instructions` to omit)

Do **not** install only skills without `goal install` — that skips the CLI and hooks.

## Workflow

1. Use **goalcraft-at-home** to draft a validated objective from the user's brief.
2. Run `goal set "<objective>"`.
3. Work until verification criteria are met.
4. Run `goal complete`.

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
| Cursor | Yes — `goal install --agent cursor` |
| Claude Code, Codex, Gemini CLI | Yes — `goal install --agent claude` / `codex` / `gemini` |
| Antigravity, Windsurf, Copilot, Pi, OpenCode | No — no stop-hook enforcement yet |

**Configure a hook-capable agent:** follow [docs/CUSTOM_AGENT_SETUP.md](./docs/CUSTOM_AGENT_SETUP.md) or run the prompt in [docs/prompts/configure-agent-hooks.md](./docs/prompts/configure-agent-hooks.md).

See [docs/SUPPORTED_AGENTS.md](./docs/SUPPORTED_AGENTS.md).

## Verify install

```bash
which goal     # ~/.local/bin/goal
```

In an agent shell: `goal set "test"`, stop without completing — you should get a follow-up reprompt. Then `goal complete`.
