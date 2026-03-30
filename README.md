# openspec-go

A Go rewrite of [OpenSpec](https://github.com/fission-ai/openspec) — the spec-driven development tool for managing structured markdown specifications and change proposals.

This is a scratch-an-itch project. The upstream TypeScript CLI works fine. This rewrite exists because I wanted a single static binary without a Node.js runtime, and because rewriting a well-specified tool in a new language is a good way to learn both the tool and the language deeply.

**This is not a fork or replacement for upstream.** If you're looking for the real thing, go use [@fission-ai/openspec](https://www.npmjs.com/package/@fission-ai/openspec).

## What's different

- Single static binary, zero runtime dependencies
- Four editor integrations (Claude Code, OpenCode, Codex, Goose) instead of twenty
- No deprecated noun-first commands
- No plugin system or config file — just the CLI

## What's the same

- Identical spec format and directory structure
- Same workflow: init → write specs → propose changes → validate → archive
- Drop-in replacement — reads/writes the same `openspec/` directory

## Install

```sh
go install github.com/chuck/openspec-go/cmd/openspec@latest
```

Or build from source:

```sh
make build
./bin/openspec --version
```

Requires [devbox](https://www.jetify.com/devbox) for toolchain management.

## Usage

```sh
openspec init --tools claude-code,goose
openspec list
openspec show my-change --json
openspec validate --all --strict
openspec archive my-change --yes
openspec view
```

## Upstream tracking

This project tracks `@fission-ai/openspec@0.17.2` as its baseline. A containerized weekly CI job checks for new upstream releases and opens GitHub issues for review. See [UPSTREAM.md](UPSTREAM.md) for details.

## Tests

```sh
make test        # run tests
make cover       # run tests with coverage report
make verify      # build Wolfi container, run cross-reference audit + format compatibility check
```

## License

Same as upstream.
