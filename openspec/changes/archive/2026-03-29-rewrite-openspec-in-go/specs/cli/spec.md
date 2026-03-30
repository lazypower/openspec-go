## ADDED Requirements

### Requirement: Root Command
The CLI SHALL provide a root `openspec` command that displays help text and version information when invoked without subcommands.

#### Scenario: Version flag
- **WHEN** user runs `openspec --version`
- **THEN** the CLI prints the version string and exits with code 0

#### Scenario: Help flag
- **WHEN** user runs `openspec --help`
- **THEN** the CLI prints usage information listing all subcommands

#### Scenario: No-color flag
- **WHEN** user runs any command with `--no-color`
- **THEN** all output MUST omit ANSI escape codes
- **AND** the `NO_COLOR` environment variable MUST also disable color

### Requirement: Init Command
The CLI SHALL provide `openspec init [path]` to initialize an OpenSpec directory structure in a project, supporting both interactive and non-interactive modes.

#### Scenario: Interactive mode (default)
- **WHEN** user runs `openspec init` without `--tools` flag and stdout is a TTY
- **THEN** the CLI launches a huh wizard with a multi-select checkbox form for tool selection
- **AND** the wizard lists all supported editors (claude-code, opencode, codex, goose)
- **AND** after selection, creates the directory structure and configures chosen editors

#### Scenario: Non-interactive with tool selection
- **WHEN** user runs `openspec init --tools claude-code,opencode,codex,goose`
- **THEN** the CLI configures slash commands for the specified editors without prompting
- **AND** generates editor-specific instruction files

#### Scenario: Non-interactive with all tools
- **WHEN** user runs `openspec init --tools all`
- **THEN** the CLI configures all supported editors (claude-code, opencode, codex, goose)

#### Scenario: Non-interactive with no tools
- **WHEN** user runs `openspec init --tools none`
- **THEN** the CLI creates the directory structure without any editor configuration

#### Scenario: Non-TTY fallback
- **WHEN** user runs `openspec init` without `--tools` flag and stdout is not a TTY
- **THEN** the CLI falls back to `--tools none` behavior and prints a message suggesting the `--tools` flag

#### Scenario: Directory structure created
- **WHEN** init completes regardless of mode
- **THEN** the CLI creates `openspec/` directory with `specs/`, `changes/`, `changes/archive/` subdirectories
- **AND** generates `openspec/project.md` from embedded template
- **AND** generates `openspec/AGENTS.md` from embedded template

#### Scenario: Initialize at custom path
- **WHEN** user runs `openspec init ./subdir`
- **THEN** the CLI initializes OpenSpec at the specified path

#### Scenario: Already initialized
- **WHEN** user runs `openspec init` in a directory that already has `openspec/`
- **THEN** the CLI warns and exits without overwriting existing files

### Requirement: Update Command
The CLI SHALL provide `openspec update [path]` to refresh managed instruction files.

#### Scenario: Update agents file
- **WHEN** user runs `openspec update`
- **THEN** the CLI replaces `openspec/AGENTS.md` with the latest embedded template
- **AND** updates slash command files for previously configured editors
- **AND** reports which files were updated

#### Scenario: Update preserves user files
- **WHEN** user runs `openspec update`
- **THEN** the CLI MUST NOT modify `openspec/project.md` or any spec/change content

### Requirement: List Command
The CLI SHALL provide `openspec list` to enumerate active changes or specs.

#### Scenario: List changes (default)
- **WHEN** user runs `openspec list`
- **THEN** the CLI lists all active changes with their task progress (completed/total)

#### Scenario: List specs
- **WHEN** user runs `openspec list --specs`
- **THEN** the CLI lists all specifications with their requirement counts

#### Scenario: List changes explicitly
- **WHEN** user runs `openspec list --changes`
- **THEN** the CLI lists active changes (same as default)

#### Scenario: Empty project
- **WHEN** no changes or specs exist
- **THEN** the CLI prints a message indicating none found

### Requirement: Show Command
The CLI SHALL provide `openspec show [item]` to display a change or spec.

#### Scenario: Show change in text mode
- **WHEN** user runs `openspec show my-change`
- **THEN** the CLI outputs the raw markdown content of `proposal.md`

#### Scenario: Show spec in text mode
- **WHEN** user runs `openspec show my-spec`
- **THEN** the CLI outputs the raw markdown content of `spec.md`

#### Scenario: Auto-detect item type
- **WHEN** user runs `openspec show foo` and `foo` exists only in changes/
- **THEN** the CLI treats it as a change

#### Scenario: Disambiguate with type flag
- **WHEN** user runs `openspec show foo --type spec`
- **THEN** the CLI forces spec lookup regardless of auto-detection

#### Scenario: JSON output for change
- **WHEN** user runs `openspec show my-change --json`
- **THEN** the CLI outputs structured JSON with id, title, deltaCount, and deltas

#### Scenario: JSON output for spec
- **WHEN** user runs `openspec show my-spec --json`
- **THEN** the CLI outputs structured JSON with id, title, overview, requirementCount, and requirements

#### Scenario: Deltas only
- **WHEN** user runs `openspec show my-change --json --deltas-only`
- **THEN** the CLI outputs only the deltas array

#### Scenario: Requirements filtering
- **WHEN** user runs `openspec show my-spec --json --requirements`
- **THEN** the CLI outputs requirements without scenario content

#### Scenario: Single requirement
- **WHEN** user runs `openspec show my-spec --json -r 2`
- **THEN** the CLI outputs only the second requirement (1-based index)

#### Scenario: Item not found
- **WHEN** user runs `openspec show nonexistent`
- **THEN** the CLI prints an error with nearest-match suggestions

### Requirement: Validate Command
The CLI SHALL provide `openspec validate [item]` to check correctness of changes and specs.

#### Scenario: Validate single change
- **WHEN** user runs `openspec validate my-change`
- **THEN** the CLI validates the change and reports issues

#### Scenario: Validate single spec
- **WHEN** user runs `openspec validate my-spec --type spec`
- **THEN** the CLI validates the spec and reports issues

#### Scenario: Validate all
- **WHEN** user runs `openspec validate --all`
- **THEN** the CLI validates all changes and specs concurrently

#### Scenario: Concurrency limit
- **WHEN** user runs `openspec validate --all --concurrency 4`
- **THEN** validation runs with at most 4 concurrent workers

#### Scenario: Strict mode
- **WHEN** user runs `openspec validate my-change --strict`
- **THEN** the CLI performs comprehensive validation including cross-spec delta conflicts, requirement uniqueness, and well-formed RENAMED pairs

#### Scenario: JSON report
- **WHEN** user runs `openspec validate my-change --json`
- **THEN** the CLI outputs a structured validation report with items array and summary

#### Scenario: Exit codes
- **WHEN** validation finds errors
- **THEN** the CLI exits with code 1
- **WHEN** validation passes
- **THEN** the CLI exits with code 0

### Requirement: Archive Command
The CLI SHALL provide `openspec archive <change-id>` to archive completed changes and merge deltas.

#### Scenario: Archive with spec merge
- **WHEN** user runs `openspec archive my-change --yes`
- **THEN** the CLI merges delta specs into main specs
- **AND** moves the change directory to `changes/archive/YYYY-MM-DD-my-change/`

#### Scenario: Skip spec updates
- **WHEN** user runs `openspec archive my-change --skip-specs --yes`
- **THEN** the CLI archives without merging deltas into specs

#### Scenario: Confirmation prompt
- **WHEN** user runs `openspec archive my-change` without `--yes`
- **THEN** the CLI prompts for confirmation before proceeding

#### Scenario: Incomplete tasks warning
- **WHEN** tasks.md has unchecked items
- **THEN** the CLI warns about incomplete tasks before archiving

### Requirement: View Command
The CLI SHALL provide `openspec view` to display a terminal dashboard.

#### Scenario: Dashboard display
- **WHEN** user runs `openspec view`
- **THEN** the CLI renders a terminal dashboard showing spec count, requirement count, active changes with progress bars, completed changes, and specifications sorted by requirement count

#### Scenario: Empty project dashboard
- **WHEN** no specs or changes exist
- **THEN** the dashboard shows zero counts

### Requirement: Shell Completion
The CLI SHALL provide shell completions via cobra's built-in completion command.

#### Scenario: Generate completions
- **WHEN** user runs `openspec completion bash`
- **THEN** the CLI outputs a bash completion script to stdout
- **AND** supports bash, zsh, fish, and powershell
