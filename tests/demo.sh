#!/usr/bin/env bash
# Generate screenshot-ready cctl output with fake data.
# Usage: bash tests/demo.sh
#        bash tests/demo.sh -s status
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(dirname "$SCRIPT_DIR")"

chmod +x "$SCRIPT_DIR/demo-cctl-data"

# Put fake cctl-data first in PATH so the real cctl wrapper finds it
export PATH="$SCRIPT_DIR:$REPO_DIR:$PATH"

# Rename trick: cctl wrapper calls "cctl-data", our fake is "demo-cctl-data"
# Create a temp symlink
TMPDIR=$(mktemp -d)
ln -s "$SCRIPT_DIR/demo-cctl-data" "$TMPDIR/cctl-data"
export PATH="$TMPDIR:$PATH"

# Run cctl — pass through all arguments (supports: ls, pick, recap, etc.)
"$REPO_DIR/cctl" "$@"

rm -rf "$TMPDIR"
