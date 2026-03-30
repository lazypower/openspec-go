## ADDED Requirements

### Requirement: Terminal Dashboard Rendering
The view command SHALL render a terminal dashboard showing project state using ASCII art and ANSI colors.

#### Scenario: Summary section
- **WHEN** user runs `openspec view`
- **THEN** the dashboard displays:
  - Total spec count and aggregate requirement count
  - Active change count
  - Completed (archived) change count
  - Overall task progress across all active changes (completed/total with percentage)

#### Scenario: Active changes with progress bars
- **WHEN** active changes exist with tasks.md files
- **THEN** each active change is listed with an ASCII progress bar showing completion percentage
- **AND** progress bars use `█` (filled) and `░` (empty) characters
- **AND** changes are sorted alphabetically

#### Scenario: Completed changes
- **WHEN** archived changes exist in `changes/archive/`
- **THEN** each is listed with a `✓` prefix, sorted alphabetically

#### Scenario: Specifications listing
- **WHEN** specs exist
- **THEN** each spec is listed with its requirement count, sorted by requirement count descending

### Requirement: Color and Formatting
The dashboard SHALL use ANSI colors for visual distinction, respecting terminal capabilities.

#### Scenario: Color scheme
- **WHEN** output is a TTY and NO_COLOR is not set
- **THEN** the dashboard uses:
  - Cyan for spec-related items
  - Yellow for active changes
  - Green for completed items and filled progress bars
  - Magenta for task progress metrics
  - Dim for secondary information

#### Scenario: No-color mode
- **WHEN** NO_COLOR environment variable is set or `--no-color` flag is passed
- **THEN** all output is plain text without ANSI escape codes

#### Scenario: Non-TTY output
- **WHEN** stdout is not a terminal (piped or redirected)
- **THEN** ANSI escape codes are omitted

### Requirement: Box Drawing
The dashboard SHALL use Unicode box-drawing characters for visual structure.

#### Scenario: Section separators
- **WHEN** rendering the dashboard
- **THEN** sections are separated by lines using `═` and `─` characters
- **AND** the header uses decorative box-drawing characters
