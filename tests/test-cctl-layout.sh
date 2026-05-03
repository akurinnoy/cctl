#!/usr/bin/env bash
set -euo pipefail

CCTL="${HOME}/bin/cctl"
FAIL=0

check_layout() {
    local label="$1"
    local output="$2"

    # Every main row must be on its own line (starts with a non-space char)
    # Every recap line must be on its own line (starts with "  └")
    # No two main rows should be on the same line
    local bad_lines
    bad_lines=$(echo "$output" | grep -cE '^\S.+\S{20,}  └' || true)
    if [[ "$bad_lines" -gt 0 ]]; then
        echo "FAIL [$label]: recap lines not on separate lines ($bad_lines violations)"
        FAIL=1
    else
        echo "OK   [$label]"
    fi
}

# Test default sort
out=$("$CCTL" 2>&1)
check_layout "default sort" "$out"

# Test all sort modes
for field in age project pid status session; do
    out=$("$CCTL" ls -s "$field" 2>&1)
    check_layout "sort=$field" "$out"
done

# Verify header exists
if echo "$out" | head -1 | grep -q "Active Sessions"; then
    echo "OK   [header present]"
else
    echo "FAIL [header missing]"
    FAIL=1
fi

# Verify column headers exist
if echo "$out" | grep -q "PROJECT.*PID.*STATUS.*AGE.*SESSION"; then
    echo "OK   [column headers]"
else
    echo "FAIL [column headers missing]"
    FAIL=1
fi

exit $FAIL
