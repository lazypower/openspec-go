#!/usr/bin/env sh
set -eu

# audit-xref.sh — Verify test cross-reference matrix against actual code.
#
# Checks:
# 1. Every Test* function in Go source appears in the matrix
# 2. Every matrix entry points to a test function that exists
# 3. Every requirement in specs has at least one test in the matrix

REPO_ROOT="${REPO_ROOT:-.}"
SPEC_FILE="${REPO_ROOT}/openspec/specs/testing/spec.md"
INTERNAL="${REPO_ROOT}/internal"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'
ERRORS=0

info()  { printf "${GREEN}✓${NC} %s\n" "$1"; }
warn()  { printf "${YELLOW}⚠${NC} %s\n" "$1"; }
error() { printf "${RED}✗${NC} %s\n" "$1"; ERRORS=$((ERRORS + 1)); }

# --- 1. Extract actual Test* functions from Go source ---
echo "=== Extracting Test* functions from Go source ==="
ACTUAL_TESTS=$(grep -rh '^func Test' "${INTERNAL}" --include='*_test.go' | \
    sed 's/func \(Test[A-Za-z0-9_]*\).*/\1/' | sort -u)
ACTUAL_COUNT=$(echo "$ACTUAL_TESTS" | wc -l | tr -d ' ')
echo "Found ${ACTUAL_COUNT} test functions"
echo ""

# --- 2. Extract matrix entries from spec ---
echo "=== Extracting matrix entries from testing spec ==="
if [ ! -f "$SPEC_FILE" ]; then
    error "Testing spec not found: $SPEC_FILE"
    exit 1
fi

# Parse the markdown table: | TestName | file | requirement | impl |
# Skip header row (Test Function) and only match actual test names (Test*)
MATRIX_TESTS=$(grep -E '^\| Test[A-Z]' "$SPEC_FILE" | awk -F'|' '{print $2}' | sed 's/^ *//;s/ *$//' | sort -u)
MATRIX_COUNT=$(echo "$MATRIX_TESTS" | wc -l | tr -d ' ')
echo "Found ${MATRIX_COUNT} matrix entries"
echo ""

# --- 3. Check: every actual test is in the matrix ---
echo "=== Check: actual tests present in matrix ==="
MISSING_FROM_MATRIX=0
for test in $ACTUAL_TESTS; do
    if ! echo "$MATRIX_TESTS" | grep -qx "$test"; then
        error "Test function '${test}' exists in code but NOT in matrix"
        MISSING_FROM_MATRIX=$((MISSING_FROM_MATRIX + 1))
    fi
done
if [ "$MISSING_FROM_MATRIX" -eq 0 ]; then
    info "All actual test functions are listed in the matrix"
fi
echo ""

# --- 4. Check: every matrix entry points to a real test ---
echo "=== Check: matrix entries point to real tests ==="
PHANTOM_ENTRIES=0
for test in $MATRIX_TESTS; do
    if ! echo "$ACTUAL_TESTS" | grep -qx "$test"; then
        error "Matrix entry '${test}' not found in code"
        PHANTOM_ENTRIES=$((PHANTOM_ENTRIES + 1))
    fi
done
if [ "$PHANTOM_ENTRIES" -eq 0 ]; then
    info "All matrix entries reference existing test functions"
fi
echo ""

# --- 5. Check: every requirement in specs has at least one test ---
echo "=== Check: spec requirements have test coverage ==="
# Extract requirements from all spec files (not testing spec itself)
SPEC_DIR="${REPO_ROOT}/openspec/specs"
REQ_TMPFILE=$(mktemp)
for spec in "${SPEC_DIR}"/*/spec.md; do
    specname=$(basename "$(dirname "$spec")")
    [ "$specname" = "testing" ] && continue
    grep -E '^### Requirement:' "$spec" | sed 's/### Requirement: *//' | while IFS= read -r req; do
        echo "${specname}:${req}" >> "$REQ_TMPFILE"
    done
done

UNCOVERED=0
while IFS= read -r req; do
    [ -z "$req" ] && continue
    reqname=$(echo "$req" | cut -d: -f2-)
    if ! grep -q "$reqname" "$SPEC_FILE" 2>/dev/null; then
        warn "Requirement '${req}' has no test in the matrix"
        UNCOVERED=$((UNCOVERED + 1))
    fi
done < "$REQ_TMPFILE"
rm -f "$REQ_TMPFILE"
if [ "$UNCOVERED" -eq 0 ]; then
    info "All spec requirements have test coverage in the matrix"
fi
echo ""

# --- Summary ---
echo "═══════════════════════════════════════"
echo "  Cross-Reference Audit Summary"
echo "═══════════════════════════════════════"
echo "  Actual test functions:  ${ACTUAL_COUNT}"
echo "  Matrix entries:         ${MATRIX_COUNT}"
echo "  Missing from matrix:    ${MISSING_FROM_MATRIX}"
echo "  Phantom matrix entries: ${PHANTOM_ENTRIES}"
echo "  Uncovered requirements: ${UNCOVERED}"
echo "═══════════════════════════════════════"

if [ "$ERRORS" -gt 0 ]; then
    echo ""
    error "Audit failed with ${ERRORS} error(s)"
    exit 1
fi

info "Audit passed"
exit 0
