#!/usr/bin/env sh
set -eu

# upstream-check.sh — Compare local baseline against upstream npm releases.
# Requires: curl, jq, gh
# Runs inside the upstream-check container or locally if deps are present.

UPSTREAM_FILE="${UPSTREAM_FILE:-UPSTREAM.md}"
PACKAGE="@fission-ai/openspec"
REGISTRY="https://registry.npmjs.org"
LABEL="upstream-sync"
DRY_RUN=false
REPO="${GITHUB_REPOSITORY:-}"

usage() {
    echo "Usage: upstream-check.sh [--dry-run] [--repo OWNER/REPO]"
    echo ""
    echo "Options:"
    echo "  --dry-run   Preview output without creating GitHub issues"
    echo "  --repo      GitHub repo (default: \$GITHUB_REPOSITORY)"
    exit 0
}

while [ $# -gt 0 ]; do
    case "$1" in
        --dry-run) DRY_RUN=true; shift ;;
        --repo) REPO="$2"; shift 2 ;;
        --help|-h) usage ;;
        *) echo "Unknown option: $1"; exit 1 ;;
    esac
done

# --- Parse baseline from UPSTREAM.md ---
parse_baseline() {
    if [ ! -f "$UPSTREAM_FILE" ]; then
        echo "ERROR: $UPSTREAM_FILE not found" >&2
        exit 1
    fi
    # Extract version from the markdown table row containing "Version"
    grep -E '^\| Version' "$UPSTREAM_FILE" | sed 's/.*`\(.*\)`.*/\1/' | head -1
}

# --- Parse declined features ---
parse_declined() {
    # Extract feature names from the declined table (skip header and separator rows)
    if grep -q "^| .*|.*|.*|.*|$" "$UPSTREAM_FILE" 2>/dev/null; then
        grep -E '^\| [^-]' "$UPSTREAM_FILE" | grep -v "^| Feature" | grep -v "^| Field" | grep -v "^| Package" | grep -v "^| Version" | grep -v "^| Reconciled" | grep -v "^<!--" | awk -F'|' '{print $2}' | sed 's/^ *//;s/ *$//' | grep -v '^$'
    fi
}

# --- Query npm registry ---
fetch_latest_version() {
    curl -sf "${REGISTRY}/${PACKAGE}/latest" | jq -r '.version'
}

fetch_versions_since() {
    local baseline="$1"
    curl -sf "${REGISTRY}/${PACKAGE}" | jq -r --arg base "$baseline" '
        .versions | keys[] | select(. > $base)
    ' | sort -V
}

fetch_changelog_entry() {
    local version="$1"
    # Try to get description or changelog from the version metadata
    curl -sf "${REGISTRY}/${PACKAGE}/${version}" | jq -r '
        "### " + .version + "\n" +
        (if .description then "- " + .description else "" end) + "\n" +
        (if .gitHead then "- Commit: " + .gitHead[:8] else "" end)
    '
}

# --- Compare versions ---
version_gt() {
    # Returns 0 (true) if $1 > $2 using sort -V
    [ "$(printf '%s\n%s' "$1" "$2" | sort -V | tail -1)" = "$1" ] && [ "$1" != "$2" ]
}

# --- Build gap summary ---
build_gap_summary() {
    local baseline="$1"
    local latest="$2"
    local declined_features="$3"

    echo "## Upstream Sync: ${PACKAGE} ${baseline} → ${latest}"
    echo ""
    echo "New versions detected since baseline \`${baseline}\`:"
    echo ""

    new_versions=$(fetch_versions_since "$baseline")
    for v in $new_versions; do
        fetch_changelog_entry "$v"
        echo ""
    done

    if [ -n "$declined_features" ]; then
        echo "## Previously Declined"
        echo ""
        echo "The following features were previously reviewed and declined:"
        echo ""
        echo "$declined_features" | while read -r feature; do
            [ -n "$feature" ] && echo "- ~~${feature}~~"
        done
        echo ""
    fi

    echo "---"
    echo ""
    echo "Review each change and either:"
    echo "1. Create an OpenSpec change proposal to adopt it"
    echo "2. Add to the Declined section of UPSTREAM.md with reason"
    echo "3. Update the baseline version once all items are resolved"
}

# --- Check for existing issue ---
issue_exists() {
    local title_prefix="$1"
    if [ -z "$REPO" ]; then
        echo "WARNING: No --repo or \$GITHUB_REPOSITORY set, skipping issue check" >&2
        return 1
    fi
    count=$(gh issue list --repo "$REPO" --label "$LABEL" --state open --search "$title_prefix" --json number --jq 'length')
    [ "$count" -gt 0 ]
}

# --- Create GitHub issue ---
create_issue() {
    local title="$1"
    local body="$2"
    if [ -z "$REPO" ]; then
        echo "WARNING: No --repo or \$GITHUB_REPOSITORY set, skipping issue creation" >&2
        return 0
    fi
    gh issue create --repo "$REPO" --title "$title" --body "$body" --label "$LABEL"
}

# --- Main ---
main() {
    baseline=$(parse_baseline)
    if [ -z "$baseline" ]; then
        echo "ERROR: Could not parse baseline version from $UPSTREAM_FILE" >&2
        exit 1
    fi
    echo "Baseline: ${PACKAGE}@${baseline}"

    latest=$(fetch_latest_version)
    if [ -z "$latest" ]; then
        echo "ERROR: Could not fetch latest version from npm" >&2
        exit 1
    fi
    echo "Latest:   ${PACKAGE}@${latest}"

    if ! version_gt "$latest" "$baseline"; then
        echo "Up to date. No action needed."
        exit 0
    fi

    echo "New version detected: ${baseline} → ${latest}"
    echo ""

    declined=$(parse_declined)
    summary=$(build_gap_summary "$baseline" "$latest" "$declined")

    title="Upstream sync: ${PACKAGE} ${baseline} → ${latest}"

    if [ "$DRY_RUN" = true ]; then
        echo "--- DRY RUN ---"
        echo ""
        echo "Title: ${title}"
        echo "Label: ${LABEL}"
        echo ""
        echo "$summary"
        exit 0
    fi

    # Idempotent: skip if issue already exists
    if issue_exists "Upstream sync: ${PACKAGE} ${baseline}"; then
        echo "Issue already exists for this version delta. Skipping."
        exit 0
    fi

    create_issue "$title" "$summary"
    echo "Issue created."
}

main
