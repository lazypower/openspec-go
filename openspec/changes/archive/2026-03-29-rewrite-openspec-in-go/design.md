# Design: OpenSpec Go Rewrite

## Context
OpenSpec is a spec-driven development tool that manages structured markdown specifications and change proposals. The TypeScript implementation (v0.17.2) works but carries Node.js runtime overhead and 20 editor integrations that add maintenance burden. This rewrite targets Go for distribution simplicity and runtime efficiency while preserving format compatibility.

## Goals
- Single static binary, zero runtime dependencies
- Format-compatible drop-in replacement for TypeScript CLI
- Sub-100ms cold start for all commands
- Test coverage with explicit test-to-requirement traceability
- Support OpenCode, Claude Code, Codex, and Goose editor integrations only

## Non-Goals
- Web UI or HTTP server (dashboard stays terminal-only)
- Plugin system or extensible editor registry
- Backward compatibility with deprecated noun-first commands
- Interactive wizard for init (simple prompts or flags instead)
- Config file management (`~/.openspec/config.json`) — defer until needed

## Package Layout

```
openspec-go/
├── cmd/
│   └── openspec/
│       └── main.go              # Entry point, cobra root command
├── internal/
│   ├── cmd/                     # Cobra command definitions
│   │   ├── root.go
│   │   ├── init.go
│   │   ├── update.go
│   │   ├── list.go
│   │   ├── show.go
│   │   ├── validate.go
│   │   ├── archive.go
│   │   └── view.go
│   ├── parser/                  # Markdown parsing
│   │   ├── spec.go              # Spec file parser
│   │   ├── change.go            # Change/proposal parser
│   │   ├── delta.go             # Delta section parser
│   │   ├── task.go              # Task progress counter
│   │   └── parser_test.go
│   ├── model/                   # Domain types
│   │   ├── spec.go              # Spec, Requirement, Scenario
│   │   ├── change.go            # Change, Delta, TaskStatus
│   │   └── validation.go        # Issue, ValidationReport
│   ├── validator/               # Validation rules
│   │   ├── spec.go              # Spec validation
│   │   ├── change.go            # Change validation
│   │   ├── constants.go         # Thresholds
│   │   └── validator_test.go
│   ├── archive/                 # Archive/merge logic
│   │   ├── archive.go           # Delta merge engine
│   │   └── archive_test.go
│   ├── editor/                  # Editor integrations
│   │   ├── editor.go            # Interface + registry
│   │   ├── claude.go            # Claude Code configurator
│   │   ├── opencode.go          # OpenCode configurator
│   │   ├── codex.go             # Codex configurator
│   │   ├── goose.go             # Goose configurator
│   │   └── editor_test.go
│   ├── template/                # Embedded templates
│   │   ├── templates.go         # go:embed for .md templates
│   │   ├── agents.md.tmpl
│   │   ├── project.md.tmpl
│   │   └── claude.md.tmpl
│   └── output/                  # Terminal rendering
│       ├── color.go             # ANSI color helpers
│       ├── table.go             # Tabular output
│       ├── progress.go          # Progress bars
│       └── json.go              # JSON output mode
├── testdata/                    # Fixture files for tests
│   ├── valid-spec/
│   ├── valid-change/
│   ├── invalid-spec/
│   └── invalid-change/
├── go.mod
├── go.sum
└── Makefile
```

## Decisions

### Markdown Parsing: strings + regexp, no AST
The spec format uses a small, fixed grammar: `##` section headers, `###` requirement headers, `####` scenario headers, bullet points with bold keywords. This is mechanical extraction, not general-purpose markdown rendering. A line-scanner with regexp matching handles all cases with ~200 lines of code vs. pulling in a full markdown parser dependency.

**Parsing strategy:**
- Line-by-line scanner
- Header detection: `^(#{1,6})\s+(.+)$`
- Requirement detection: `^###\s+Requirement:\s+(.+)$`
- Scenario detection: `^####\s+Scenario:\s+(.+)$`
- Delta section detection: `^##\s+(ADDED|MODIFIED|REMOVED|RENAMED)\s+Requirements$`
- Task detection: `^[-*]\s+\[([ x])\]` (case-insensitive x)
- Content accumulation between headers

### CLI Framework: cobra
Industry standard for Go CLIs. Provides subcommands, flags, help generation, and shell completions for free. Maps cleanly to the existing command surface.

### Editor Support: Four editors only
The TypeScript version supports 20 editors. Most are copy-paste configurators with trivial template differences. Supporting four (OpenCode, Claude Code, Codex, Goose) covers the user's actual toolchain. Adding more later is trivial — each editor is ~50 lines implementing a single interface.

### Interactive TUI: charmbracelet/huh
The Charm ecosystem's `huh` library is the Go equivalent of `@inquirer/prompts` — multi-step forms, selection lists, confirmations, and text input. This preserves the TypeScript version's init wizard experience (tool selection with checkboxes, multi-step flow) without rolling custom stdin handling.

### Terminal Styling: charmbracelet/lipgloss
`lipgloss` replaces chalk for all styled terminal output — colors, bold, dim, padding, borders. It respects `NO_COLOR` and handles non-TTY detection. The dashboard uses lipgloss for box drawing and layout. JSON output uses `encoding/json`.

### Templates: go:embed
Agent instructions (AGENTS.md), project template, and editor-specific files are embedded at compile time via `go:embed`. No runtime file resolution needed.

### No Config Subsystem (initially)
The TypeScript version has `~/.openspec/config.json` with feature flags but only uses it for internal development flags. Skip this entirely for v1. Add it when there's a real use case.

## Risks / Trade-offs

- **Risk**: Format drift between Go and TypeScript parsers during transition
  - Mitigation: Share testdata fixtures; run both parsers against same inputs in CI
- **Risk**: Missing edge cases in markdown parsing
  - Mitigation: Extract test fixtures from real OpenSpec projects; property-based testing for parser
- **Risk**: Three-editor limitation frustrates users of other tools
  - Mitigation: AGENTS.md works as universal fallback; editor interface is trivial to extend

## Migration Plan
1. Build Go CLI with format-compatible parsing
2. Validate against existing OpenSpec projects (same directory structure)
3. Publish as `openspec` binary (separate from npm package)
4. Users switch by replacing `npx openspec` with `openspec` binary
5. No data migration needed — reads same `openspec/` directory

## Open Questions
- None remaining. Init supports both interactive (huh wizard) and non-interactive (`--tools` flag). Validation is concurrent with `--concurrency` flag.
