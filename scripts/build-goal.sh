#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
mkdir -p "$INSTALL_DIR"
(cd "$ROOT" && go build -o "$INSTALL_DIR/goal" .)
echo "Built: $INSTALL_DIR/goal"
echo "Ensure $INSTALL_DIR is on your PATH."
