# archive

## Purpose

## Requirements

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

