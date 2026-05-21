#!/usr/bin/env bash
# Goal-at-home installer — delegates to scripts/install.sh (Cursor today).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")" && pwd)"
exec "$ROOT/scripts/install.sh" "$@"
