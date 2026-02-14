# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
go build -o grin .                    # Build binary
go test ./...                         # Run all tests
go test -v ./...                      # Run all tests (verbose)
go test -run TestGrinFixtures ./...   # Run a single test by name
go test -race ./...                   # Run tests with race detector
go test -coverprofile=c.out ./...     # Run tests with coverage
golangci-lint run ./...               # Run linter
```

## Pre-commit Hooks

Git hooks live in `.githooks/` (tracked in repo). New clones must configure:

```bash
git config core.hooksPath .githooks
```

The pre-commit hook runs `go vet`, `go test -race`, and `golangci-lint` before each commit.

## Architecture

grin is a flat `package main` Go CLI tool modeled after [gron](https://github.com/tomnomnom/gron). It converts INI files into discrete greppable assignments and back.

**Data flow (grin direction):** INI text → `statementsFromINI()` → `[]statement` → sort → render via `statementconv` → output

**Data flow (ungrin direction):** assignment text → `lexStatement()` per line → `[]statement` → `extractPathAndValue()` → `ungrinFromStatements()` → INI text

### Core types (token.go → statements.go)

- `token{text, typ}` — atomic unit (bare identifier, dot, equals, semicolon, quoted string, `{}`)
- `statement` (`[]token`) — one assignment line, e.g. `ini.section.key = "value";`
- `statements` (`[]statement`) — sortable slice; `Less()` sorts lexicographically with `=` tokens first so `ini = {};` always precedes `ini.x = ...`

Statements are built via chainable helpers: `withBare()`, `withPath()`, `withStringValue()`, `withEmptyObject()`.

### Key files

- **ini.go** — INI parser. Line-oriented: skips comments/blanks, tracks sections (including dotted `[a.b.c]` nesting), emits `{}` for intermediate section objects, validates keys via `validIdentifier()`.
- **ungrin.go** — Lexer (`lexStatement`) parses assignment strings back to tokens. `ungrinFromStatements()` groups by section path and emits INI.
- **main.go** — CLI entry point. Three action functions: `grinAction`, `ungrinAction`, `grinValuesAction`. Flags: `-u`, `-v`, `-c`, `-m`, `--no-sort`, `--version`.
- **identifier.go** — `validIdentifier()`: letters/digits/underscores/hyphens, must start with letter or underscore.

### Test structure

Tests use table-driven patterns and file-based fixtures in `testdata/` (`.ini`/`.grin` pairs). `TestRoundTrip` verifies INI→grin→ungrin→grin stability. `TestGrinFixtures` compares output against golden `.grin` files.

## Commit Conventions

- Use text-form gitmoji in commit subjects (e.g. `:sparkles:`, `:bug:`, `:robot_face:`)
- Use `:robot_face:` for changes to CLAUDE.md and `.claude/`
- Commit body describes the "why", not the "what"
- Include trailer: `Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>`

## Release

Tag with `v*` and push to trigger goreleaser via GitHub Actions. Version is injected via `-ldflags "-X main.grinVersion=..."`.
