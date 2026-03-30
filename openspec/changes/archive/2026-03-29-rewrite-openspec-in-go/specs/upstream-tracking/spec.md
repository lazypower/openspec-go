## ADDED Requirements

### Requirement: Upstream Baseline Tracking
The project SHALL track the upstream TypeScript OpenSpec version it was last reconciled against.

#### Scenario: Baseline recorded in project metadata
- **WHEN** the Go rewrite is initialized
- **THEN** a file `UPSTREAM.md` at the repo root SHALL record the npm package name and version (e.g., `@fission-ai/openspec@0.17.2`) and the date of last reconciliation

#### Scenario: Baseline updated after reconciliation
- **WHEN** a gap analysis is completed and all relevant changes are either adopted or explicitly declined
- **THEN** the baseline version and date in `UPSTREAM.md` are updated to reflect the reconciled version

### Requirement: CI-Driven Upstream Check
The project SHALL provide a GitHub Actions workflow that periodically checks for new upstream releases.

#### Scenario: Scheduled workflow
- **WHEN** the cron schedule fires (e.g., weekly)
- **THEN** the workflow queries the npm registry for the latest version of `@fission-ai/openspec`
- **AND** compares it against the baseline version in `UPSTREAM.md`

#### Scenario: No new version
- **WHEN** the upstream version matches the baseline
- **THEN** the workflow exits successfully with no further action

#### Scenario: New version detected
- **WHEN** a newer upstream version exists
- **THEN** the workflow fetches the GitHub releases or changelog for versions between baseline and latest
- **AND** produces a gap summary comparing upstream changes against local spec capabilities

### Requirement: Gap Analysis as GitHub Issue
The CI workflow SHALL create a GitHub issue when upstream gaps are detected.

#### Scenario: Issue creation
- **WHEN** the upstream check detects a new version with changes
- **THEN** a GitHub issue is created with:
  - Title including the version delta (e.g., "Upstream sync: @fission-ai/openspec 0.17.2 → 0.18.0")
  - Body containing the changelog summary and gap analysis
  - Label `upstream-sync`

#### Scenario: Idempotent issue creation
- **WHEN** an open issue with the `upstream-sync` label already exists for the same version delta
- **THEN** the workflow skips issue creation

### Requirement: Decline Tracking
The project SHALL allow explicitly declining upstream features that are intentionally not adopted.

#### Scenario: Record a decline
- **WHEN** a maintainer adds an entry to `UPSTREAM.md` under a "Declined" section
- **THEN** the entry records the feature description, reason for declining, and the upstream version it appeared in

#### Scenario: Declined features acknowledged in gap reports
- **WHEN** a gap analysis runs and a detected gap matches a previously declined feature
- **THEN** the gap is listed separately as "declined" rather than flagged as an open gap

### Requirement: Standalone Script
The upstream check logic SHALL be implemented as a standalone script (e.g., `scripts/upstream-check.sh`) rather than embedded in the openspec CLI binary.

#### Scenario: Script independence
- **WHEN** the upstream check runs
- **THEN** it requires only `curl`, `jq`, and `gh` CLI — no dependency on the openspec binary

#### Scenario: Local execution
- **WHEN** a maintainer runs `./scripts/upstream-check.sh` locally
- **THEN** it produces the same gap analysis output as the CI workflow
- **AND** supports `--dry-run` to preview without creating issues
