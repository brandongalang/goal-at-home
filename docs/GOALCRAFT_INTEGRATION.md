# Goalcraft integration plan

How to bake [Goalcraft](https://github.com/grp06/goalcraft) into [goal-at-home](https://github.com/brandongalang/goal-at-home) so every enforced goal is written to a high-quality, evidence-based contract — not a one-line wish.

## Roles (do not merge them)

| Project | Job |
|---------|-----|
| **Goalcraft** | *Shape* the objective: interview, compact template, length validation, Codex `/goal` literacy |
| **goal-at-home** | *Enforce* the objective: persist per session, inject session ID, stop hook until `goal complete` |

Goalcraft is an **agent skill** (LLM workflow + bundled validator). It is not a Go library and should not be invoked as a black-box API from `goal` without a model in the loop.

```text
User intent (messy)
       │
       ▼
┌──────────────────┐
│ Goalcraft skill  │  ← agent follows SKILL.md, drafts compact objective
└────────┬─────────┘
         │ validated text (≤3,400 chars target for Codex-shaped goals)
         ▼
┌──────────────────┐
│ goal set "..."   │  ← goal-at-home stores + hooks enforce
└────────┬─────────┘
         │
         ▼
   work … goal complete
```

On **Cursor**, we store the crafted objective in `~/.cursor/goals/<session>.json` and enforce via hooks — we do not call Codex `create_goal` unless the user also wants Codex `/goal` parity later.

## What to reuse from Goalcraft

| Asset | Use in goal-at-home |
|-------|---------------------|
| `SKILL.md` workflow | Chained skill: “always goalcraft before `goal set`” |
| `scripts/validate_goal_length.py` | Vendored copy; run on `goal set` / `goal edit` (optional strictness per agent) |
| Compact output shape | Default template for `goal template` / docs |
| `references/codex-goal-contract.md` | Link from goal-enforcement skill for Codex users |

Do **not** fork Goalcraft’s LLM prompts into Go strings long-term — stay aligned with upstream via skill install or git submodule docs.

## Integration phases

### Phase 0 — Skill chaining (no CLI changes)

**Effort:** small · **Value:** immediate

1. Extend `skill/goal-enforcement/SKILL.md`:
   - Before `goal set`, invoke Goalcraft (user has `$goalcraft` / skill installed).
   - Paste only the **objective body** into `goal set` (strip leading `/goal ` if present).
   - Run Goalcraft’s validator on the final string before `goal set`.
2. Extend `install.sh` (optional flag `--with-goalcraft`):
   - Clone or symlink `grp06/goalcraft` → `~/.cursor/skills/goalcraft` (and later per-agent paths).
   - Print reminder: Goalcraft is required for the recommended workflow.
3. Document in README: recommended path is Goalcraft → `goal set`, not raw `goal set "fix the thing"`.

**Acceptance:** Agent given a vague task runs Goalcraft, then `goal set` with a compact goal; stop hook still fires.

### Phase 1 — Validator in the `goal` CLI

**Effort:** small · **Value:** deterministic guardrails

1. Vendor `validate_goal_length.py` under `goal-at-home/scripts/` (MIT-compatible; attribute Goalcraft in NOTICE).
2. New behavior on `goal set` and `goal edit`:
   ```bash
   goal set "Destination: …"     # runs validator by default
   goal set --no-validate "…"  # escape hatch (Cursor may exceed Codex limits)
   goal validate [file|-]       # explicit check; exit 0/1 for hooks/CI
   ```
3. Defaults:
   - `--target-chars 3400 --strict-target` when `GOALCRAFT_STRICT=1` or `--strict`
   - Hard max 3999 only (warn above 3400) for Cursor-first installs
4. Strip `/goal ` prefix before save (same as Goalcraft validator).

**Acceptance:** `goal set` rejects (or warns) objectives over limit; accepts Goalcraft-shaped compact goals.

### Phase 2 — `goal template` + structured metadata (optional)

**Effort:** medium · **Value:** consistency without an LLM in CLI

1. `goal template` prints the compact scaffold (Destination, Starting point, …) for the agent to fill.
2. Extend `store.Goal` optionally:
   ```json
   {
     "session_id": "...",
     "text": "full enforcement string",
     "status": "active",
     "crafted_by": "goalcraft",
     "schema": "compact-v1",
     "fields": { "destination": "...", "verification": "..." }
   }
   ```
   Enforcement and stop hook continue to use `text` only; structured fields are for UI, resume, and analytics.

**Acceptance:** Same enforcement behavior; richer persisted goals for debugging.

### Phase 3 — `goal install` bundles Goalcraft

**Effort:** medium · **Value:** one-shot setup

As part of `goal install` (multi-agent installer):

| Agent | goal-at-home | Goalcraft skill |
|-------|--------------|-----------------|
| Cursor | binary + hooks + goal-enforcement | `~/.cursor/skills/goalcraft` |
| Codex | binary + hooks + goal-enforcement | `~/.codex/skills/goalcraft` |
| Claude Code | TBD | `~/.claude/skills/goalcraft` |

Implementation options:

- **Submodule:** `third_party/goalcraft` at pinned tag; install copies skill + validator.
- **Install script:** `git clone --depth 1` into skills dir on demand.
- **Document only:** user runs `ln -s` per Goalcraft README (simplest; no coupling).

Prefer **pinned vendor of validator + symlink skill** so offline install works.

### Phase 4 — Codex `/goal` bridge (optional, Codex-only)

**Effort:** larger · **Value:** single source of truth on Codex

When `goal install -a codex` and user wants native Codex goals:

- After Goalcraft + `goal set`, optionally mirror to Codex `thread/goal/set` via app-server (only if user passes `--sync-codex`).
- Keep goal-at-home file store as backup for hooks that read session id from Codex hook payload.

Defer until Codex adapter lands; avoid dual enforcement without clear precedence rules.

## Recommended agent workflow (document this everywhere)

```markdown
1. User describes messy intent.
2. Agent applies Goalcraft skill → compact, validated objective (<3,400 chars).
3. Agent runs: goal set "<objective>"
4. Agent works until verification gates in the objective are satisfied.
5. Agent runs: goal complete
```

For **review-only** requests (“sharpen this goal”), Goalcraft alone — no `goal set` until user confirms activation.

## CLI surface (proposed)

| Command | Purpose |
|---------|---------|
| `goal craft` | *Deprecated alias* — print “use Goalcraft skill, then goal set” (avoid implying Go crafts goals) |
| `goal template` | Emit compact Goalcraft scaffold |
| `goal validate` | Run length/structure check on stdin or file |
| `goal set` | Save active goal; validate by default |
| `goal set --no-validate` | Save verbatim (power users / long Cursor goals) |

Avoid `goal set --craft` that calls an LLM from Go — that duplicates Goalcraft poorly.

## Hook interaction

- **preToolUse:** unchanged; still injects `--session-id` on `goal …` commands.
- **stop:** unchanged; still reads `text` from store.
- Optional future **preToolUse** matcher: block `goal set` if validator fails when `GOALCRAFT_ENFORCE=1` — heavy-handed; prefer CLI rejection first.

## Cursor vs Codex limits

| Surface | Character policy |
|---------|------------------|
| Codex `/goal` | Hard 3,999; Goalcraft target 3,400 |
| Cursor `goal set` | No platform limit today; still recommend Goalcraft shape for quality |

Default: validate with **warn** above 3,400 on Cursor, **strict** when env `GOAL_AT_HOME_PROFILE=codex` or flag `--profile codex`.

## Install / dependency story

```bash
# Recommended one-liner (future)
./install.sh --with-goalcraft

# Manual today
git clone https://github.com/grp06/goalcraft.git ~/.cursor/skills/goalcraft
git clone https://github.com/brandongalang/goal-at-home.git && cd goal-at-home && ./install.sh
```

Goal-at-home should **not** vendor the entire Goalcraft repo into `goal-at-home` git history without a submodule — pin version in `docs/GOALCRAFT_INTEGRATION.md` and install script.

## Open questions

1. **Upstream coupling:** Submodule vs copy validator only vs runtime `npx`/path to user’s existing `~/.codex/skills/goalcraft`?
2. **Strict by default?** Reject `goal set` over 3,400 on Cursor might frustrate users; suggest strict only with `--with-goalcraft` install profile.
3. **Single skill vs merged:** Merge Goalcraft checklist into `goal-enforcement` (drift risk) vs keep two skills (clear separation)?
4. **Attribution:** NOTICE file + link to grp06/goalcraft; confirm license compatibility (Goalcraft repo has LICENSE).

## Suggested implementation order

1. Phase 0 — skill chaining + README/AGENTS.md workflow (this week)
2. Phase 1 — vendor validator + `goal validate` + validate on set (next)
3. Phase 3 — `--with-goalcraft` on install.sh when multi-agent installer exists
4. Phase 2 / 4 — only if users need structured JSON or Codex sync

## Success criteria

- No `goal set` in docs without Goalcraft in the recommended path.
- Every enforced goal is auditable (Destination, Verification, Done/stop) per Goalcraft compact shape.
- Validator prevents accidental 4k+ Codex-style paste into `goal set`.
- Goal-at-home remains enforceable without Goalcraft installed (validator off / `--no-validate`), but docs call that suboptimal.
