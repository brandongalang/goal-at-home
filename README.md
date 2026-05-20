# cursor-goal

Session-scoped goal enforcement for [Cursor](https://cursor.com) agents. If the agent sets a goal and tries to stop without running `goal complete`, a `stop` hook auto-reprompts it to continue.

## How it works

1. Agent runs `goal set "objective"` (you can instruct Cursor to do this).
2. A `preToolUse` hook injects `--session-id` so state is tied to the chat.
3. When the agent loop ends, a `stop` hook checks `~/.cursor/goals/<session-id>.json`.
4. If status is `active`, Cursor submits a follow-up message (up to `loop_limit`, default 10).
5. Agent runs `goal complete` when truly done; the next stop is allowed.

User **Stop** leaves the goal **active** so the next message resumes enforcement.

## Install

```bash
git clone https://github.com/brandongalang/cursor-goal.git
cd cursor-goal
./scripts/install.sh
```

Requires `go`, `jq`, and `~/.local/bin` on your `PATH`.

The installer:

- Builds `goal` to `~/.local/bin/goal`
- Merges hooks into `~/.cursor/hooks.json` (preserves existing hooks like `rtk`)
- Installs `~/.cursor/skills/goal-enforcement/SKILL.md`

Restart Cursor after install.

## Commands

| Command | Description |
|---------|-------------|
| `goal set --session-id <id> "<text>"` | Create or replace the session goal |
| `goal edit --session-id <id> "<text>"` | Update text on an active goal |
| `goal complete --session-id <id>` | Mark complete (allows stop) |
| `goal clear --session-id <id>` | Remove the session goal |
| `goal status --session-id <id>` | Show current goal |

In Cursor agent shells, `--session-id` is injected by the hook; agents should run `goal set "..."` without the flag.

## Hooks

| Hook | Command |
|------|---------|
| `preToolUse` (Shell) | `goal hook pre-tool-use` |
| `stop` | `goal hook stop` |

See [hooks.example.json](./hooks.example.json).

## User rule (optional)

Add to Cursor user rules:

> For sustained agent work, tell the agent to `goal set` at the start and not finish until `goal complete`.

## License

MIT
