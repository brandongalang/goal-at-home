#!/usr/bin/env python3
"""Validate a goal-at-home objective (for goal set), adapted from Goalcraft."""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

# Match Codex /goal best practice (same as upstream Goalcraft).
DEFAULT_MAX_CHARS = 3_999
TARGET_CHARS = 3_400
STRICT_WARN_CHARS = 3_800


def objective_text(text: str) -> str:
    text = text.strip()
    if text.startswith("```"):
        lines = text.splitlines()
        if lines and lines[0].startswith("```"):
            lines = lines[1:]
        if lines and lines[-1].startswith("```"):
            lines = lines[:-1]
        text = "\n".join(lines).strip()
    # goal set "..." wrapper from agent output
    m = re.match(r'^goal\s+set\s+"(.*)"\s*$', text, re.DOTALL)
    if m:
        return m.group(1).strip()
    m = re.match(r"^goal\s+set\s+'(.*)'\s*$", text, re.DOTALL)
    if m:
        return m.group(1).strip()
    # Codex /goal paste (strip if user pasted legacy shape)
    if text.startswith("/goal"):
        rest = text[len("/goal") :]
        if rest.startswith((" ", "\n", "\t")):
            return rest.strip()
    return text


def read_input(path: str | None) -> str:
    if path:
        return Path(path).read_text(encoding="utf-8")
    return sys.stdin.read()


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Validate a goal-at-home objective for goal set."
    )
    parser.add_argument(
        "path",
        nargs="?",
        help="File containing the objective or goal set command.",
    )
    parser.add_argument(
        "--max-chars",
        type=int,
        default=DEFAULT_MAX_CHARS,
        help=f"Hard maximum. Default: {DEFAULT_MAX_CHARS}.",
    )
    parser.add_argument(
        "--target-chars",
        type=int,
        default=TARGET_CHARS,
        help=f"Recommended target. Default: {TARGET_CHARS}.",
    )
    parser.add_argument(
        "--strict-target",
        action="store_true",
        help="Exit non-zero when objective exceeds --target-chars.",
    )
    args = parser.parse_args()

    objective = objective_text(read_input(args.path))
    count = len(objective)
    print(f"objective_chars={count}")
    print(f"target_chars={args.target_chars}")
    print(f"max_chars={args.max_chars}")
    if count > args.max_chars:
        print("error=objective exceeds /goal character limit (3999)", file=sys.stderr)
        return 1
    if count >= STRICT_WARN_CHARS:
        print(
            "warning=objective at or above 3800 chars; compress per Goalcraft practice",
            file=sys.stderr,
        )
        if args.strict_target:
            return 1
    elif count > args.target_chars:
        print("warning=objective exceeds 3400 char target", file=sys.stderr)
        if args.strict_target:
            return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
