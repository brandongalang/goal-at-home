# Supported agents

Goal-at-home needs two capabilities:

1. **Pre-tool (shell)** — inject `--session-id` on `goal *` commands
2. **Stop / end-of-turn** — reprompt or block when a goal is still `active`

## Tier 1 — Full enforcement (installer supported today)

| Agent | Install | Pre-tool event | Stop event | Session ID |
|-------|---------|----------------|------------|------------|
| [Cursor](https://cursor.com) | `./install.sh` or `./scripts/install.sh` | `preToolUse` (Shell) | `stop` → `followup_message` | `conversation_id` |
| [Claude Code](https://code.claude.com) | Coming soon | `PreToolUse` (Bash) | `Stop` → `decision: block` | `session_id` |
| [Codex CLI](https://developers.openai.com/codex) | Coming soon | `PreToolUse` (Bash) | `Stop` | `session_id` |
| [Gemini CLI](https://geminicli.com) | Coming soon | `BeforeTool` | `AfterAgent` retry | `session_id` |

**Cursor** is the only agent with a working `install.sh` merge today. Others can use the skill (soft guidance) until adapters land.

## Tier 2 — Extension / plugin APIs (not shell hooks)

| Agent | Mechanism | Hard stop loop? |
|-------|-----------|-----------------|
| [Pi](https://pi.dev) | TypeScript extensions (`~/.pi/agent/extensions/`) | Possible via custom extension; no `hooks.json` |
| [OpenCode](https://opencode.ai) | TypeScript plugins (`.opencode/plugins/`) | Possible via plugin events; different config shape |

These need a native adapter package, not `goal hook` + `jq` merge.

## Tier 3 — Skill / rules only (no deterministic stop hook)

| Agent | What works |
|-------|------------|
| [Antigravity](https://antigravity.google) | `.agent/workflows/`, `.agent/rules/`, skills — no shell stop hook ([forum](https://discuss.ai.google.dev/t/hooks-in-antigravity/120458)) |
| [Windsurf](https://windsurf.com) | Cascade hooks for tools; no Cursor-style `stop` reprompt |
| GitHub Copilot, Continue, Aider, most IDE assistants | `npx skills add` for the skill only |

## Skill-only install (any agent on [skills.sh](https://skills.sh))

```bash
npx skills add brandongalang/goal-at-home --skill goal-enforcement -g -y
```

Installs `SKILL.md` to the agent’s skills directory. Does **not** install the `goal` binary or lifecycle hooks.

## Compared to built-in `/goal` commands

Claude Code and Codex ship native goal/loop features. Goal-at-home is for:

- A **file-backed** objective (`goal set` / `goal complete`)
- The same **CLI** across agents (as adapters ship)
- **Deterministic** stop enforcement via hooks (where supported)
