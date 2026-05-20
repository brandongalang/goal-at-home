#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
CURSOR_DIR="${CURSOR_DIR:-$HOME/.cursor}"
GOALS_DIR="$CURSOR_DIR/goals"
HOOKS_JSON="$CURSOR_DIR/hooks.json"
SKILL_DIR="$CURSOR_DIR/skills/goal-enforcement"

mkdir -p "$INSTALL_DIR" "$GOALS_DIR" "$SKILL_DIR" "$(dirname "$HOOKS_JSON")"

echo "Building goal binary..."
(cd "$ROOT" && go build -o "$INSTALL_DIR/goal" .)

echo "Installing skill..."
cp "$ROOT/skill/SKILL.md" "$SKILL_DIR/SKILL.md"

merge_hooks() {
  local tmp
  tmp="$(mktemp)"
  if [[ -f "$HOOKS_JSON" ]]; then
    jq --arg goal_pre "goal hook pre-tool-use" --arg goal_stop "goal hook stop" '
      .version //= 1
      | .hooks //= {}
      | .hooks.preToolUse //= []
      | .hooks.stop //= []
      | .hooks.preToolUse |= (
          [{ "command": $goal_pre, "matcher": "Shell" }]
          + map(select(.command != $goal_pre))
        )
      | .hooks.stop |= (
          map(select(.command != $goal_stop))
          + [{ "command": $goal_stop, "loop_limit": 10 }]
        )
    ' "$HOOKS_JSON" > "$tmp"
  else
    jq -n --arg goal_pre "goal hook pre-tool-use" --arg goal_stop "goal hook stop" '
      {
        version: 1,
        hooks: {
          preToolUse: [{ command: $goal_pre, matcher: "Shell" }],
          stop: [{ command: $goal_stop, loop_limit: 10 }]
        }
      }
    ' > "$tmp"
  fi
  mv "$tmp" "$HOOKS_JSON"
}

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required to merge hooks.json. Install jq and re-run, or add manually:" >&2
  cat "$ROOT/hooks.example.json" >&2
  exit 1
fi

echo "Merging ~/.cursor/hooks.json..."
merge_hooks

echo ""
echo "Installed:"
echo "  binary:  $INSTALL_DIR/goal"
echo "  goals:   $GOALS_DIR/"
echo "  hooks:   $HOOKS_JSON"
echo "  skill:   $SKILL_DIR/SKILL.md"
echo ""
echo "Ensure $INSTALL_DIR is on your PATH."
echo "Restart Cursor or reload hooks after install."
