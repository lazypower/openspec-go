## 1. Project Scaffolding
- [ ] 1.1 Initialize Go module (`go mod init`)
- [ ] 1.2 Add cobra dependency
- [ ] 1.3 Add charmbracelet/huh dependency (interactive forms/wizards)
- [ ] 1.4 Add charmbracelet/lipgloss dependency (terminal styling)
- [ ] 1.5 Create package directory structure per design.md
- [ ] 1.6 Create `cmd/openspec/main.go` entry point with root cobra command
- [ ] 1.7 Create Makefile with build, test, lint targets

## 2. Domain Models
- [ ] 2.1 Define Spec, Requirement, Scenario structs in `internal/model/spec.go`
- [ ] 2.2 Define Change, Delta, TaskStatus structs in `internal/model/change.go`
- [ ] 2.3 Define Issue, ValidationReport structs in `internal/model/validation.go`

## 3. Markdown Parser
- [ ] 3.1 Implement line-scanner base in `internal/parser/`
- [ ] 3.2 Implement spec parser (`internal/parser/spec.go`)
- [ ] 3.3 Implement change/proposal parser (`internal/parser/change.go`)
- [ ] 3.4 Implement delta parser (`internal/parser/delta.go`)
- [ ] 3.5 Implement task progress parser (`internal/parser/task.go`)
- [ ] 3.6 Write parser unit tests with testdata fixtures
- [ ] 3.7 Create testdata fixtures (valid-spec, valid-change, invalid-spec, invalid-change)

## 4. Validator
- [ ] 4.1 Implement change validation rules (`internal/validator/change.go`)
- [ ] 4.2 Implement spec validation rules (`internal/validator/spec.go`)
- [ ] 4.3 Implement strict mode checks
- [ ] 4.4 Define validation constants (`internal/validator/constants.go`)
- [ ] 4.5 Write validator unit tests

## 5. Archive / Merge Engine
- [ ] 5.1 Implement delta merge operations (ADDED, MODIFIED, REMOVED, RENAMED)
- [ ] 5.2 Implement atomic file write (temp file + rename)
- [ ] 5.3 Implement archive workflow (validate → merge → move)
- [ ] 5.4 Write archive unit tests with before/after fixture pairs

## 6. Terminal Output
- [ ] 6.1 Implement lipgloss-based color/style helpers (`internal/output/color.go`)
- [ ] 6.2 Implement progress bar renderer
- [ ] 6.3 Implement JSON output helpers
- [ ] 6.4 Implement NO_COLOR and non-TTY detection

## 7. Editor Integrations
- [ ] 7.1 Define editor interface (`internal/editor/editor.go`)
- [ ] 7.2 Implement Claude Code configurator with slash commands
- [ ] 7.3 Implement OpenCode configurator
- [ ] 7.4 Implement Codex configurator
- [ ] 7.5 Implement Goose configurator
- [ ] 7.6 Create embedded templates (AGENTS.md, project.md, CLAUDE.md, slash commands)
- [ ] 7.7 Write editor integration tests

## 8. CLI Commands
- [ ] 8.1 Implement `init` command with huh wizard for tool selection
- [ ] 8.2 Implement `update` command
- [ ] 8.3 Implement `list` command (changes + specs modes)
- [ ] 8.4 Implement `show` command with type auto-detection and JSON output
- [ ] 8.5 Implement `validate` command with strict mode and JSON report
- [ ] 8.6 Implement `archive` command with confirmation prompts
- [ ] 8.7 Implement `view` command (terminal dashboard)

## 9. Upstream Tracking
- [ ] 9.1 Create `UPSTREAM.md` with baseline version `@fission-ai/openspec@0.17.2`
- [ ] 9.2 Write `scripts/upstream-check.sh` (curl + jq + gh, no Go dependency)
- [ ] 9.3 Create GitHub Actions workflow (`.github/workflows/upstream-check.yml`) on weekly cron
- [ ] 9.4 Implement idempotent issue creation with `upstream-sync` label
- [ ] 9.5 Add `--dry-run` flag to script for local preview
- [ ] 9.6 Document decline tracking format in `UPSTREAM.md`

## 10. Integration Tests
- [ ] 10.1 Write init integration test (full directory structure verification)
- [ ] 10.2 Write full lifecycle integration test (init → create → validate → archive)
- [ ] 10.3 Write view/dashboard integration test

## 11. Cross-Reference Verification
- [ ] 11.1 Verify test cross-reference matrix matches actual test functions
- [ ] 11.2 Verify every requirement in specs has at least one test in the matrix
- [ ] 11.3 Update matrix if implementation files diverge from design
