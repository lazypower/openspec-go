#!/usr/bin/env sh
set -eu

# entrypoint.sh — Run verification suite inside the container.
# Usage:
#   docker run verify           # run both audits
#   docker run verify audit     # cross-reference audit only
#   docker run verify compat    # format compatibility only

MODE="${1:-all}"

case "$MODE" in
    audit)
        exec /work/scripts/audit-xref.sh
        ;;
    compat)
        exec /work/scripts/verify-compat.sh
        ;;
    all)
        echo "╔═══════════════════════════════════════╗"
        echo "║    OpenSpec Verification Suite        ║"
        echo "╚═══════════════════════════════════════╝"
        echo ""

        FAILED=0

        echo ">>> Cross-Reference Audit"
        echo ""
        /work/scripts/audit-xref.sh || FAILED=$((FAILED + 1))
        echo ""

        echo ">>> Format Compatibility Check"
        echo ""
        /work/scripts/verify-compat.sh || FAILED=$((FAILED + 1))
        echo ""

        if [ "$FAILED" -gt 0 ]; then
            echo "Verification failed: ${FAILED} suite(s) had errors"
            exit 1
        fi
        echo "All verification suites passed."
        ;;
    *)
        echo "Unknown mode: $MODE"
        echo "Usage: entrypoint.sh [all|audit|compat]"
        exit 1
        ;;
esac
