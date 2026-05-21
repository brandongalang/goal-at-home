---
name: goal-enforcement
description: >-
  Session goal enforcement for coding agents. Use when the user asks for sustained
  work that must not stop until explicitly finished. Requires goal-at-home CLI and
  hooks (run ./install.sh); this skill alone does not enforce stopping.
---

# Session goal enforcement

This workspace uses **goal-at-home**: the agent cannot end sustained work while a goal is **active** unless it runs `goal complete`. Enforcement comes from hooks, not this file alone.

If `goal` or hooks are missing, run `./install.sh` from https://github.com/brandongalang/goal-at-home (see AGENTS.md).

## When to use

- User asks for multi-step work, refactors, migrations, or anything that must run to completion
- User says to set a goal, use goal enforcement, or not stop until done

## Commands

All commands are shell tools. On Cursor, the `preToolUse` hook injects `--session-id` automatically — do not add it yourself.

```bash
goal set "One-sentence objective with clear done criteria"
goal edit "Updated objective if scope changed"
goal status
goal complete    # only when every requirement is satisfied
goal clear       # remove the goal without completing (user asked to cancel, or goal is wrong)
```

## Workflow

1. At the start of substantive work, run `goal set "..."` with a concrete, verifiable objective.
2. Do the work. If scope changes, `goal edit "..."`.
3. Before claiming done, verify against the objective (tests, files, behavior).
4. Run `goal complete` only when satisfied.
5. Use `goal clear` only when the user cancels or the goal should be dropped.

## Rules

- Never say you are finished without running `goal complete` while a goal is active.
- Do not run `goal complete` to bypass work — the stop hook will keep reprompting until the goal is real.
- If the user hits Stop, the goal stays active; the next turn should continue toward it.
- Casual questions without sustained work do not need `goal set`.
