# upstream-tracking

## Purpose

## Requirements

### Requirement: Standalone Script
The upstream check logic SHALL be implemented as a standalone script (e.g., `scripts/upstream-check.sh`) rather than embedded in the openspec CLI binary.

#### Scenario: Script independence
- **WHEN** the upstream check runs
- **THEN** it requires only `curl`, `jq`, and `gh` CLI — no dependency on the openspec binary

#### Scenario: Local execution
- **WHEN** a maintainer runs `./scripts/upstream-check.sh` locally
- **THEN** it produces the same gap analysis output as the CI workflow
- **AND** supports `--dry-run` to preview without creating issues

