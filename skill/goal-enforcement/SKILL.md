---
name: goal-enforcement
description: >-
  Session goal enforcement for coding agents. Use when the user asks for sustained
  work that must not stop until explicitly finished. Requires goal-at-home CLI and
  hooks (run goal install); this skill alone does not enforce stopping.
---

# Session goal enforcement

This workspace uses **goal-at-home**: the agent cannot end sustained work while a goal is **active** unless it runs `goal complete`. Enforcement comes from hooks, not this file alone.

If `goal` or hooks are missing, run `goal install --agent <cursor|claude|codex|gemini>` from https://github.com/brandongalang/goal-at-home

## Before goal set

Use the bundled **goalcraft-at-home** skill to turn the user's brief into a validated objective (Destination, Verification, Done/stop, etc.). Do not `goal set` a vague one-liner unless the user explicitly refuses crafting.

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

1. **goalcraft-at-home** — craft and validate the objective from the user's brief.
2. `goal set "..."` — activate enforcement for this session.
3. Do the work; if scope changes, `goal edit "..."`.
4. Verify against the objective (tests, files, behavior).
5. `goal complete` only when satisfied.
6. `goal clear` only when the user cancels or the goal should be dropped.

## Rules

- Never say you are finished without running `goal complete` while a goal is active.
- Do not run `goal complete` to bypass work — the stop hook will keep reprompting until the goal is real.
- If the user hits Stop, the goal stays active; the next turn should continue toward it.
- Casual questions without sustained work do not need `goal set`.
