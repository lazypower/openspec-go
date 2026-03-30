## ADDED Requirements

### Requirement: Editor Interface
The editor integration system SHALL define a common interface that all editor configurators implement.

#### Scenario: Interface contract
- **WHEN** an editor configurator is registered
- **THEN** it implements: `Name() string`, `Configure(projectPath, openspecPath string) error`, `UpdateExisting(projectPath, openspecPath string) error`, and `IsConfigured(projectPath string) bool`

### Requirement: Claude Code Integration
The system SHALL configure Claude Code slash commands and instruction files.

#### Scenario: Configure slash commands
- **WHEN** Claude Code is selected during init
- **THEN** the configurator creates `.claude/commands/openspec/` directory with:
  - `proposal.md` — slash command for creating proposals
  - `apply.md` — slash command for implementing changes
  - `archive.md` — slash command for archiving changes

#### Scenario: Configure CLAUDE.md injection
- **WHEN** Claude Code is selected during init
- **THEN** the configurator adds an OpenSpec instruction block to the project CLAUDE.md
- **AND** the block is wrapped in `<!-- OPENSPEC:START -->` / `<!-- OPENSPEC:END -->` markers

#### Scenario: Update existing configuration
- **WHEN** `openspec update` is run and Claude Code was previously configured
- **THEN** the slash command files are replaced with latest templates
- **AND** the CLAUDE.md managed block is refreshed

### Requirement: OpenCode Integration
The system SHALL configure OpenCode with OpenSpec commands.

#### Scenario: Configure OpenCode
- **WHEN** OpenCode is selected during init
- **THEN** the configurator creates the appropriate command/prompt files for OpenCode's configuration format

#### Scenario: Update existing OpenCode config
- **WHEN** `openspec update` is run and OpenCode was previously configured
- **THEN** existing configuration files are refreshed with latest templates

### Requirement: Codex Integration
The system SHALL configure Codex with OpenSpec prompts.

#### Scenario: Configure Codex
- **WHEN** Codex is selected during init
- **THEN** the configurator creates prompt files in Codex's global config directory (`~/.codex/prompts/` or project-local equivalent)

#### Scenario: Update existing Codex config
- **WHEN** `openspec update` is run and Codex was previously configured
- **THEN** existing prompt files are refreshed with latest templates

### Requirement: Goose Integration
The system SHALL configure Goose with OpenSpec commands via `.goosehints` and recipe YAML files.

#### Scenario: Configure .goosehints
- **WHEN** Goose is selected during init
- **THEN** the configurator creates or appends to `.goosehints` at the project root with an OpenSpec instruction block wrapped in marker comments
- **AND** uses `@openspec/AGENTS.md` syntax to inline the full agent instructions

#### Scenario: Configure slash command recipes
- **WHEN** Goose is selected during init
- **THEN** the configurator creates recipe YAML files for proposal, apply, and archive operations under `.goose/recipes/openspec/`
- **AND** each recipe follows the Goose recipe schema (version, title, description, instructions)

#### Scenario: Update existing Goose config
- **WHEN** `openspec update` is run and Goose was previously configured
- **THEN** the `.goosehints` managed block is refreshed
- **AND** recipe YAML files are replaced with latest templates

### Requirement: AGENTS.md as Universal Fallback
The system SHALL generate AGENTS.md as an editor-agnostic instruction file.

#### Scenario: AGENTS.md always generated
- **WHEN** `openspec init` is run regardless of tool selection
- **THEN** `openspec/AGENTS.md` is created with comprehensive workflow instructions
- **AND** it serves as a fallback for any editor that reads markdown context files

### Requirement: Template Embedding
All editor templates SHALL be embedded in the binary at compile time.

#### Scenario: go:embed usage
- **WHEN** the binary is built
- **THEN** all `.md.tmpl` template files are embedded via `go:embed` directives
- **AND** no external template files are required at runtime

#### Scenario: Template variables
- **WHEN** templates are rendered
- **THEN** Go `text/template` substitutes project-specific values (paths, tool names) into the output
