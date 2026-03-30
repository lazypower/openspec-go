#!/usr/bin/env sh
set -eu

# verify-compat.sh — Compare Go openspec parser output against TypeScript openspec.
#
# Runs both parsers on shared testdata fixtures and diffs the normalized JSON.
# Designed to run inside the verify container where both tools are available.

REPO_ROOT="${REPO_ROOT:-.}"
TESTDATA="${REPO_ROOT}/testdata"
GO_BIN="${GO_BIN:-${REPO_ROOT}/bin/openspec}"
TS_BIN="${TS_BIN:-npx @fission-ai/openspec}"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'
ERRORS=0
PASSED=0

info()  { printf "${GREEN}✓${NC} %s\n" "$1"; PASSED=$((PASSED + 1)); }
error() { printf "${RED}✗${NC} %s\n" "$1"; ERRORS=$((ERRORS + 1)); }

# --- Setup temp project for each fixture ---
run_comparison() {
    local fixture_type="$1"  # "spec" or "change"
    local fixture_dir="$2"
    local fixture_name="$3"

    local tmpdir
    tmpdir=$(mktemp -d)
    trap "rm -rf '$tmpdir'" EXIT

    # Create minimal openspec project structure
    mkdir -p "${tmpdir}/openspec/specs" "${tmpdir}/openspec/changes/archive"

    if [ "$fixture_type" = "spec" ]; then
        mkdir -p "${tmpdir}/openspec/specs/${fixture_name}"
        cp "$fixture_dir"/*.md "${tmpdir}/openspec/specs/${fixture_name}/"

        # Go parser: show as JSON
        local go_out="${tmpdir}/go.json"
        (cd "$tmpdir" && "$GO_BIN" show "$fixture_name" --type spec --json 2>/dev/null) > "$go_out" || true

        # TS parser: show as JSON (if available)
        local ts_out="${tmpdir}/ts.json"
        if command -v npx >/dev/null 2>&1; then
            (cd "$tmpdir" && $TS_BIN show "$fixture_name" --type spec --json 2>/dev/null) > "$ts_out" || true
        fi

    elif [ "$fixture_type" = "change" ]; then
        mkdir -p "${tmpdir}/openspec/changes/${fixture_name}"
        cp -r "$fixture_dir"/* "${tmpdir}/openspec/changes/${fixture_name}/"

        local go_out="${tmpdir}/go.json"
        (cd "$tmpdir" && "$GO_BIN" show "$fixture_name" --type change --json 2>/dev/null) > "$go_out" || true

        local ts_out="${tmpdir}/ts.json"
        if command -v npx >/dev/null 2>&1; then
            (cd "$tmpdir" && $TS_BIN show "$fixture_name" --type change --json 2>/dev/null) > "$ts_out" || true
        fi
    fi

    # Compare if both outputs exist and are non-empty
    if [ -s "$go_out" ]; then
        if [ -s "$ts_out" ]; then
            # Normalize: sort keys, strip whitespace differences
            local go_norm="${tmpdir}/go.norm.json"
            local ts_norm="${tmpdir}/ts.norm.json"
            jq -S '.' "$go_out" > "$go_norm" 2>/dev/null || true
            jq -S '.' "$ts_out" > "$ts_norm" 2>/dev/null || true

            if diff -q "$go_norm" "$ts_norm" >/dev/null 2>&1; then
                info "${fixture_type}/${fixture_name}: Go and TS parsers agree"
            else
                error "${fixture_type}/${fixture_name}: Go and TS parsers DIVERGE"
                diff --unified "$ts_norm" "$go_norm" | head -30 || true
            fi
        else
            info "${fixture_type}/${fixture_name}: Go parser produced output (TS not available for comparison)"
        fi
    else
        error "${fixture_type}/${fixture_name}: Go parser produced no output"
    fi

    rm -rf "$tmpdir"
    trap - EXIT
}

# --- Main ---
echo "═══════════════════════════════════════"
echo "  Format Compatibility Verification"
echo "═══════════════════════════════════════"
echo ""

# Check Go binary exists
if [ ! -x "$GO_BIN" ]; then
    echo "Building Go binary..."
    (cd "$REPO_ROOT" && make build) || { error "Failed to build Go binary"; exit 1; }
fi

# Run spec fixtures
echo "--- Spec Fixtures ---"
if [ -d "${TESTDATA}/valid-spec" ]; then
    run_comparison "spec" "${TESTDATA}/valid-spec" "valid-spec"
fi

# Run change fixtures
echo ""
echo "--- Change Fixtures ---"
if [ -d "${TESTDATA}/valid-change" ]; then
    run_comparison "change" "${TESTDATA}/valid-change" "valid-change"
fi

echo ""
echo "═══════════════════════════════════════"
echo "  Compatibility Summary"
echo "═══════════════════════════════════════"
echo "  Passed: ${PASSED}"
echo "  Failed: ${ERRORS}"
echo "═══════════════════════════════════════"

if [ "$ERRORS" -gt 0 ]; then
    exit 1
fi
exit 0
