## ADDED Requirements

### Requirement: Delta Merge Engine
The archive system SHALL merge delta specs into main specs by applying ADDED, MODIFIED, REMOVED, and RENAMED operations.

#### Scenario: Apply ADDED deltas
- **WHEN** archiving a change with ADDED requirements
- **THEN** the new requirement blocks are appended to the end of the `## Requirements` section in the target spec

#### Scenario: Apply MODIFIED deltas
- **WHEN** archiving a change with MODIFIED requirements
- **THEN** the entire requirement block (header through scenarios) in the main spec is replaced with the delta content
- **AND** matching is performed by requirement name (whitespace-insensitive)

#### Scenario: Apply REMOVED deltas
- **WHEN** archiving a change with REMOVED requirements
- **THEN** the matching requirement block is deleted from the main spec

#### Scenario: Apply RENAMED deltas
- **WHEN** archiving a change with RENAMED requirements specifying FROM and TO
- **THEN** the requirement header in the main spec is updated from the old name to the new name

#### Scenario: Requirement not found for MODIFIED/REMOVED
- **WHEN** a MODIFIED or REMOVED delta references a requirement name that does not exist in the target spec
- **THEN** the archive reports an ERROR and aborts without modifying any files

#### Scenario: New spec creation
- **WHEN** a delta targets a spec that does not yet exist in `openspec/specs/`
- **THEN** the archive creates a new spec file with the ADDED requirements

### Requirement: Archive Workflow
The archive command SHALL follow a defined sequence of validation and file operations.

#### Scenario: Pre-archive validation
- **WHEN** user runs `openspec archive my-change`
- **THEN** the change and its delta specs are validated before any merge occurs
- **AND** validation errors abort the archive

#### Scenario: Skip validation
- **WHEN** user runs `openspec archive my-change --no-validate --yes`
- **THEN** validation is skipped entirely

#### Scenario: Post-merge validation
- **WHEN** deltas have been merged into main specs
- **THEN** the merged specs are validated before being written to disk
- **AND** if validation fails, no files are written and the operation is rolled back

#### Scenario: Directory move
- **WHEN** merge and validation succeed
- **THEN** the change directory is moved from `changes/<id>/` to `changes/archive/YYYY-MM-DD-<id>/`

#### Scenario: Skip spec updates
- **WHEN** user passes `--skip-specs`
- **THEN** the change is archived (moved to archive/) without merging any deltas into specs

### Requirement: Incomplete Task Handling
The archive system SHALL warn when tasks are incomplete.

#### Scenario: Incomplete tasks with interactive prompt
- **WHEN** tasks.md has unchecked items and `--yes` is not set
- **THEN** the CLI warns about N incomplete tasks and asks for confirmation to proceed

#### Scenario: Incomplete tasks with --yes
- **WHEN** tasks.md has unchecked items and `--yes` is set
- **THEN** the CLI prints a warning but proceeds without prompting

### Requirement: Atomic File Operations
The archive system SHALL ensure file modifications are atomic to prevent partial state.

#### Scenario: Write-then-move pattern
- **WHEN** merging deltas into a spec file
- **THEN** the merged content is written to a temporary file first, then renamed to the target path
- **AND** if any step fails, the original spec file remains unchanged

#### Scenario: Multi-spec consistency
- **WHEN** a change affects multiple specs
- **THEN** all spec merges are computed and validated before any files are written
- **AND** if any single merge fails, no specs are modified
