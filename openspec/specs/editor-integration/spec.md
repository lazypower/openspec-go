# editor-integration

## Purpose

## Requirements

### Requirement: Template Embedding
All editor templates SHALL be embedded in the binary at compile time.

#### Scenario: go:embed usage
- **WHEN** the binary is built
- **THEN** all `.md.tmpl` template files are embedded via `go:embed` directives
- **AND** no external template files are required at runtime

#### Scenario: Template variables
- **WHEN** templates are rendered
- **THEN** Go `text/template` substitutes project-specific values (paths, tool names) into the output

