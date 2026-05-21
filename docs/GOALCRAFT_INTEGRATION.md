# Goalcraft at home (bundled, adapted)

[goal-at-home](https://github.com/brandongalang/goal-at-home) ships **goalcraft-at-home**: an adapted Goalcraft workflow for our CLI and hooks — not a dependency on Codex `/goal`.

Upstream reference: [grp06/goalcraft](https://github.com/grp06/goalcraft). See [GOALCRAFT_ATTRIBUTION.md](./GOALCRAFT_ATTRIBUTION.md).

## Two skills, one install

| Skill | Role |
|-------|------|
| **goalcraft-at-home** | Shape the objective (evidence, compact sections, validation) |
| **goal-enforcement** | Activate and enforce (`goal set` → work → `goal complete`) |

```text
messy brief
    → goalcraft-at-home (agent)
    → goal set "Destination: … Verification: …"
    → hooks enforce until goal complete
```

No external Goalcraft install required. No `/goal` prefix. No Codex `create_goal`.

## Bundled assets

| Path | Purpose |
|------|---------|
| `skill/goalcraft-at-home/SKILL.md` | Crafting workflow for agents |
| `skill/goal-enforcement/SKILL.md` | Enforcement workflow |
| `scripts/validate_goal_objective.py` | Length check before `goal set` |

`./install.sh` copies both skills to `~/.cursor/skills/`.

## Validator limits (same as Codex `/goal`)

| Limit | Chars | Rationale |
|-------|-------|-----------|
| Working target | 3,400 | Goalcraft default; keeps objectives operable |
| Fail draft | ≥ 3,800 | Compress even if under hard cap |
| Hard max | 3,999 | Matches Codex `/goal` TUI limit |

```bash
python3 scripts/validate_goal_objective.py --strict-target /tmp/objective.txt
```

## Agent workflow (canonical)

1. User describes work; agent uses **goalcraft-at-home** to draft and validate.
2. Agent runs `goal set "<objective>"` (hooks add `--session-id` on Cursor).
3. Agent works until **Verification** and **Success metric** are satisfied.
4. Agent runs `goal complete`.

Do not `goal set` a one-liner without crafting unless the user explicitly insists.

## Future CLI wiring (optional)

| Command | Behavior |
|---------|----------|
| `goal validate` | Shell out to `validate_goal_objective.py` |
| `goal set` | Validate by default; `--no-validate` escape hatch |

Crafting stays in the **skill** (LLM). The CLI only validates and persists.

## Codex users

If someone also uses Codex native `/goal`, they can keep upstream Goalcraft installed separately. goal-at-home does not conflict; **goalcraft-at-home** is the path for enforced Cursor (and future hook adapters).

## Success criteria

- Every recommended `goal set` follows the compact Goalcraft shape.
- Docs never tell users to install grp06/goalcraft for Cursor enforcement.
- Attribution and deltas are documented in GOALCRAFT_ATTRIBUTION.md.
