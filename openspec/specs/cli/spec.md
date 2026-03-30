# cli

## Purpose

## Requirements

### Requirement: Shell Completion
The CLI SHALL provide shell completions via cobra's built-in completion command.

#### Scenario: Generate completions
- **WHEN** user runs `openspec completion bash`
- **THEN** the CLI outputs a bash completion script to stdout
- **AND** supports bash, zsh, fish, and powershell

