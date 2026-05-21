# Goalcraft attribution

The **goalcraft-at-home** skill and `scripts/validate_goal_objective.py` are adapted from [Goalcraft](https://github.com/grp06/goalcraft) by [@grp06](https://github.com/grp06).

## What we kept

- Evidence-based goal shaping (Destination, Verification, Done/stop, etc.)
- Compact objective discipline and quality bar
- Validator pattern (length checks before activation)

## What we changed

| Upstream (Codex) | goal-at-home |
|------------------|--------------|
| `/goal …` activation | `goal set "…"` |
| 3,400 / 3,999 char limits | 5,000 target / 20,000 max |
| `create_goal`, thread APIs | Session file + stop hooks |
| External skill install | Bundled in this repo |

We are not affiliated with upstream Goalcraft. Send improvements to the shaping workflow upstream when they apply to Codex; send goal-at-home-specific issues here.
