---
name: goalcraft-at-home
description: >-
  Turn a rough task brief into a goal-at-home objective for enforced agent sessions.
  Use before goal set when the user wants sustained work with clear done criteria.
  Adapted from Goalcraft for goal CLI + stop hooks, not Codex /goal.
---

# Goalcraft at home

Shape messy intent into an objective that **goal-at-home** can enforce. After crafting, the agent runs `goal set "..."` (not `/goal`). Hooks keep the session alive until `goal complete`.

Upstream inspiration: [grp06/goalcraft](https://github.com/grp06/goalcraft). This skill is adapted for [goal-at-home](https://github.com/brandongalang/goal-at-home).

## Core contract

- Draft the **objective text only** — the string passed to `goal set`.
- Do **not** run `goal set`, `goal complete`, or `goal clear` unless the user explicitly asks to activate or finish a goal.
- Prefer evidence-based completion: "done" requires proof named in the objective, not intent or elapsed time.
- Enforcement is external: stop hooks reprompt while the goal is `active`; premature `goal complete` does not satisfy the user if verification failed.

## Workflow

1. Identify rough intent, workspace, and expected end state.
   - Preserve user intent when they supplied a draft.
   - Inspect the repo when paths, stack, or conventions matter.
   - Ask the smallest necessary question if a critical success criterion is missing; otherwise list assumptions.

2. Shape the goal around evidence, not effort. Use this checklist while thinking; compress before output.
   - Destination: what must be true at the end.
   - Starting point: branch, files, issue, or known state.
   - Objective/scope: concrete work and boundaries (files, systems, directories).
   - Preserve: must-not-regress, safety, approval boundaries.
   - Verification: commands, tests, screenshots, PR state, logs, or explicit user confirmation.
   - Done/stop: when to continue autonomously, when to ask, checkpoint rhythm, blockers.
   - Success metric: observable proof of completion.

3. Keep the objective usable by goal-at-home.
   - Output is plain text for `goal set "..."` — no `/goal` prefix.
   - Working target: **5,000 characters** or fewer (readability and hook follow-up size).
   - Hard limit: **20,000 characters** — compress if longer.
   - Make every requirement auditable (files, commands, behavior, artifacts).
   - Include exact commands only when known from the repo or user.
   - For vague or huge work, recommend narrowing scope or a planning pass before `goal set`.

4. Output shape (compact, default).
   - Single prose block with labeled sections: Destination, Starting point, Objective/scope, Preserve, Verification, Done/stop, Success metric.
   - Merge deliverables into Objective/scope; merge autonomy into Done/stop.
   - Do not paste example lists into the goal unless they are execution constraints.

5. Validate before returning.
   - Write the final objective to a temp file and run the bundled validator from this repo:
     `scripts/validate_goal_objective.py --strict-target /tmp/objective.txt`
   - Resolve the script relative to the **goal-at-home** install (this skill’s repo), not the user project.
   - Do not return until validation passes.

6. Activation.
   - Default: return **Assumptions** + **Ready for goal set** with the exact quoted string.
   - If the user asked to activate: run `goal set "<objective>"` only after validation (hooks inject `--session-id` on Cursor).

## Output format

```markdown
Assumptions:
- ...

Ready for goal set:

goal set "Destination: ...

Starting point: ...

Objective/scope: ...

Preserve: ...

Verification: ...

Done/stop: ...

Success metric: ..."
```

Report validated length, e.g. `Objective length: 2,914 characters`.

## Quality bar

- Another agent can resume after compaction or a new message and still know what "done" means.
- Premature completion is hard: unrelated green CI or "I tried" is not enough.
- User boundaries (no commits, no deploys, ask first) are explicit when relevant.
- The objective aligns with what hooks will keep reprompting until `goal complete` is justified.

## Pair with goal-enforcement

1. **goalcraft-at-home** (this skill) — craft the objective.
2. **goal-enforcement** — `goal set`, work, `goal complete`.

Requires goal-at-home installed (`./install.sh`): CLI + hooks + both skills.
