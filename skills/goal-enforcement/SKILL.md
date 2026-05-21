---
name: goal-enforcement
description: >-
  Session goal enforcement for coding agents. Use when the user asks for sustained
  work that must not stop until explicitly finished. Requires the goal CLI and
  hooks (goal-at-home) for hard enforcement; skill-only install is soft guidance.
---

# Session goal enforcement

This workspace uses **goal-at-home**: the agent should not end sustained work while a goal is **active** unless it runs `goal complete` (when hooks are installed).

## When to use

- User asks for multi-step work, refactors, migrations, or anything that must run to completion
- User says to set a goal, use goal enforcement, or not stop until done

## Install

- **Cursor (full):** clone https://github.com/brandongalang/goal-at-home and run `./install.sh`
- **Other agents (skill only):** `npx skills add brandongalang/goal-at-home --skill goal-enforcement -g -y`

See https://github.com/brandongalang/goal-at-home/blob/main/AGENTS.md

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

- Never say you are finished without running `goal complete` while a goal is active (when hooks are installed).
- Do not run `goal complete` to bypass work — the stop hook will keep reprompting until the goal is real.
- If the user hits Stop, the goal stays active; the next turn should continue toward it.
- Casual questions without sustained work do not need `goal set`.
