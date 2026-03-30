## 1. Project Scaffolding
- [x] 1.1 Initialize Go module (`go mod init`)
- [x] 1.2 Add cobra dependency
- [x] 1.3 Add charmbracelet/huh dependency (interactive forms/wizards)
- [x] 1.4 Add charmbracelet/lipgloss dependency (terminal styling)
- [x] 1.5 Create package directory structure per design.md
- [x] 1.6 Create `cmd/openspec/main.go` entry point with root cobra command
- [x] 1.7 Create Makefile with build, test, lint targets

## 2. Domain Models
- [x] 2.1 Define Spec, Requirement, Scenario structs in `internal/model/spec.go`
- [x] 2.2 Define Change, Delta, TaskStatus structs in `internal/model/change.go`
- [x] 2.3 Define Issue, ValidationReport structs in `internal/model/validation.go`

## 3. Markdown Parser
- [x] 3.1 Implement line-scanner base in `internal/parser/`
- [x] 3.2 Implement spec parser (`internal/parser/spec.go`)
- [x] 3.3 Implement change/proposal parser (`internal/parser/change.go`)
- [x] 3.4 Implement delta parser (`internal/parser/delta.go`)
- [x] 3.5 Implement task progress parser (`internal/parser/task.go`)
- [x] 3.6 Write parser unit tests with testdata fixtures
- [x] 3.7 Create testdata fixtures (valid-spec, valid-change, invalid-spec, invalid-change)

## 4. Validator
- [x] 4.1 Implement change validation rules (`internal/validator/change.go`)
- [x] 4.2 Implement spec validation rules (`internal/validator/spec.go`)
- [x] 4.3 Implement strict mode checks
- [x] 4.4 Define validation constants (`internal/validator/constants.go`)
- [x] 4.5 Write validator unit tests

## 5. Archive / Merge Engine
- [x] 5.1 Implement delta merge operations (ADDED, MODIFIED, REMOVED, RENAMED)
- [x] 5.2 Implement atomic file write (temp file + rename)
- [x] 5.3 Implement archive workflow (validate → merge → move)
- [x] 5.4 Write archive unit tests with before/after fixture pairs

## 6. Terminal Output
- [x] 6.1 Implement lipgloss-based color/style helpers (`internal/output/color.go`)
- [x] 6.2 Implement progress bar renderer
- [x] 6.3 Implement JSON output helpers
- [x] 6.4 Implement NO_COLOR and non-TTY detection

## 7. Editor Integrations
- [x] 7.1 Define editor interface (`internal/editor/editor.go`)
- [x] 7.2 Implement Claude Code configurator with slash commands
- [x] 7.3 Implement OpenCode configurator
- [x] 7.4 Implement Codex configurator
- [x] 7.5 Implement Goose configurator
- [x] 7.6 Create embedded templates (AGENTS.md, project.md, CLAUDE.md, slash commands)
- [x] 7.7 Write editor integration tests

## 8. CLI Commands
- [x] 8.1 Implement `init` command with huh wizard for tool selection
- [x] 8.2 Implement `update` command
- [x] 8.3 Implement `list` command (changes + specs modes)
- [x] 8.4 Implement `show` command with type auto-detection and JSON output
- [x] 8.5 Implement `validate` command with strict mode and JSON report
- [x] 8.6 Implement `archive` command with confirmation prompts
- [x] 8.7 Implement `view` command (terminal dashboard)

## 9. Upstream Tracking
- [x] 9.1 Create `UPSTREAM.md` with baseline version `@fission-ai/openspec@0.17.2`
- [x] 9.2 Write `scripts/upstream-check.sh` (curl + jq + gh, no Go dependency)
- [x] 9.3 Create GitHub Actions workflow (`.github/workflows/upstream-check.yml`) on weekly cron
- [x] 9.4 Implement idempotent issue creation with `upstream-sync` label
- [x] 9.5 Add `--dry-run` flag to script for local preview
- [x] 9.6 Document decline tracking format in `UPSTREAM.md`
- [x] 9.7 Containerize deps in Wolfi image (`containers/upstream-check/Containerfile`)

## 10. Integration Tests
- [x] 10.1 Write init integration test (full directory structure verification)
- [x] 10.2 Write full lifecycle integration test (init → create → validate → archive)
- [x] 10.3 Write view/dashboard integration test

## 11. Cross-Reference Verification
- [x] 11.1 Verify test cross-reference matrix matches actual test functions
- [x] 11.2 Verify every requirement in specs has at least one test in the matrix
- [x] 11.3 Update matrix if implementation files diverge from design
- [x] 11.4 Containerize audit + format compat in Wolfi image (`containers/verify/Containerfile`)
- [x] 11.5 CI workflow for automated verification (`.github/workflows/verify.yml`)
