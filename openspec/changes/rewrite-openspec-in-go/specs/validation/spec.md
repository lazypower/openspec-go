## ADDED Requirements

### Requirement: Change Validation Rules
The validator SHALL check change proposals against a defined set of rules.

#### Scenario: Must have at least one delta
- **WHEN** a change has no spec files under `changes/<id>/specs/`
- **THEN** validation reports an ERROR: "Change must have at least one delta"

#### Scenario: ADDED/MODIFIED requirements must have scenarios
- **WHEN** a delta contains an ADDED or MODIFIED requirement without a `#### Scenario:` block
- **THEN** validation reports an ERROR: "Requirement must have at least one scenario"

#### Scenario: Requirements must use SHALL/MUST
- **WHEN** a requirement body does not contain "SHALL" or "MUST" (case-insensitive)
- **THEN** validation reports an ERROR: "Requirement must contain SHALL or MUST keyword"

#### Scenario: Proposal must have required sections
- **WHEN** proposal.md is missing a "Why" or "What Changes" section
- **THEN** validation reports an ERROR identifying the missing section

#### Scenario: Why section minimum length
- **WHEN** the "Why" section contains fewer than 50 characters
- **THEN** validation reports a WARNING: "Why section is too short"

#### Scenario: No duplicate requirements within sections
- **WHEN** two requirements in the same delta section share the same name
- **THEN** validation reports an ERROR: "Duplicate requirement name"

### Requirement: Spec Validation Rules
The validator SHALL check specifications against a defined set of rules.

#### Scenario: Must have Purpose and Requirements sections
- **WHEN** a spec.md is missing `## Purpose` or `## Requirements`
- **THEN** validation reports an ERROR identifying the missing section

#### Scenario: Purpose minimum length
- **WHEN** the Purpose section contains fewer than 50 characters
- **THEN** validation reports an ERROR: "Purpose section too short"

#### Scenario: Requirements must have scenarios
- **WHEN** a spec requirement lacks any `#### Scenario:` block
- **THEN** validation reports an ERROR: "Requirement must have at least one scenario"

#### Scenario: Requirements must use SHALL/MUST
- **WHEN** a spec requirement body does not contain "SHALL" or "MUST"
- **THEN** validation reports an ERROR

### Requirement: Strict Validation Mode
The validator SHALL perform additional checks when `--strict` is enabled.

#### Scenario: Cross-spec delta conflict detection
- **WHEN** strict mode is enabled and two active changes modify the same requirement in the same spec
- **THEN** validation reports a WARNING about the potential conflict

#### Scenario: Well-formed RENAMED pairs
- **WHEN** strict mode is enabled and a RENAMED delta lacks either FROM or TO
- **THEN** validation reports an ERROR: "RENAMED must specify both FROM and TO"

#### Scenario: MODIFIED completeness check
- **WHEN** strict mode is enabled and a MODIFIED requirement contains fewer scenarios than the original spec requirement
- **THEN** validation reports a WARNING suggesting the full requirement may not have been copied

### Requirement: Validation Report Format
The validator SHALL produce structured output in both text and JSON formats.

#### Scenario: Text output
- **WHEN** validation runs without `--json`
- **THEN** issues are printed to stderr grouped by file, with ERROR/WARNING/INFO prefixes and colors

#### Scenario: JSON output
- **WHEN** validation runs with `--json`
- **THEN** the output conforms to the report schema: `{ items: [{ id, type, valid, issues: [{ level, path, message, line }] }], summary: { totals: { items, passed, failed }, byType } }`

#### Scenario: Exit code reflects validation result
- **WHEN** any ERROR-level issue exists
- **THEN** the process exits with code 1
- **WHEN** only WARNING or INFO issues exist
- **THEN** the process exits with code 0

### Requirement: Concurrent Validation
The validator SHALL validate multiple items concurrently when validating more than one change or spec.

#### Scenario: Concurrent batch validation
- **WHEN** user runs `openspec validate --all` or `openspec validate --changes` or `openspec validate --specs`
- **THEN** the validator validates items concurrently using a bounded worker pool
- **AND** results are collected and reported in deterministic order (sorted by item ID)

#### Scenario: Concurrency flag
- **WHEN** user passes `--concurrency N`
- **THEN** the worker pool is limited to N goroutines
- **AND** the default is the value of `OPENSPEC_CONCURRENCY` env var, or `runtime.NumCPU()` if unset

#### Scenario: Single item validation
- **WHEN** user validates a single item (e.g., `openspec validate my-change`)
- **THEN** no worker pool is created; validation runs directly in the main goroutine

### Requirement: Validation Constants
The validator SHALL use configurable thresholds for validation checks.

#### Scenario: Threshold values
- **WHEN** validation checks are performed
- **THEN** the following thresholds apply:
  - MIN_WHY_SECTION_LENGTH: 50
  - MAX_WHY_SECTION_LENGTH: 1000
  - MIN_PURPOSE_LENGTH: 50
  - MAX_REQUIREMENT_TEXT_LENGTH: 500
  - MAX_DELTAS_PER_CHANGE: 10 (warning only)
