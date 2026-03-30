# Change: Rewrite OpenSpec CLI in Go

## Why
The TypeScript OpenSpec CLI (v0.17.2) requires Node.js >= 20, ships ~30MB of node_modules for what is fundamentally file I/O and string parsing, and couples tightly to npm distribution. A Go rewrite produces a single static binary with zero runtime dependencies, sub-millisecond startup, and trivial cross-platform distribution. The spec format is structured markdown — mechanical to parse without an AST — making Go's stdlib (`os`, `strings`, `regexp`, `filepath`) a natural fit.

## What Changes
- **BREAKING**: Drop all 20 editor configurators except OpenCode, Claude Code, Codex, and Goose
- Replace `@inquirer/prompts` interactive wizard with `charmbracelet/huh` (equivalent Go TUI forms)
- Replace `chalk` terminal styling with `charmbracelet/lipgloss`
- **BREAKING**: Drop shell completion subcommand (cobra has built-in completion)
- **BREAKING**: Drop deprecated noun-first commands (`change show`, `spec show`, etc.)
- Replace Node.js CLI with single Go binary using cobra
- Replace chalk/ora terminal output with Go stdlib + minimal ANSI helpers
- Replace zod validation with Go struct validation
- Replace commander argument parsing with cobra
- Maintain identical spec format, directory structure, and workflow semantics
- Add comprehensive test suite with test-to-implementation cross-reference matrix

## Impact
- Affected specs: cli, parsing, validation, archive, dashboard, editor-integration, upstream-tracking, testing
- Affected code: entire codebase (greenfield rewrite)
- Distribution: single binary via `go install`, Homebrew tap, or direct download
- Backward compat: reads/writes identical `openspec/` directory structure — drop-in replacement
