# CLAUDE.md

## Project Overview

A Wordle solver engine written in Go. Given feedback from Wordle guesses (green/yellow/gray), it narrows down possible words and suggests optimal next guesses using letter-frequency scoring.

Designed to be used as a CLI tool or imported as a library. The solver maintains a word list, applies constraints from each guess, and ranks remaining candidates to minimize average guesses needed.

This will be packaged as an OpenClaw skill for automated daily Wordle solving.

## Tech Stack

- **Language**: Go 1.23+
- **Linting**: golangci-lint with gofumpt
- **Testing**: go test (stdlib)
- **Build**: Makefile

## Commands

```bash
make build       # Build the binary to bin/
make test        # Run tests with race detection + coverage
make lint        # Run golangci-lint
make dev         # Run the solver interactively
```

## Directory Structure

```text
wordle-solver/
├── cmd/
│   └── wordle-solver/
│       └── main.go          # CLI entrypoint
├── internal/
│   ├── solver/
│   │   ├── solver.go        # Core solver logic
│   │   └── solver_test.go
│   ├── wordlist/
│   │   ├── wordlist.go      # Word list loading + filtering
│   │   └── wordlist_test.go
│   └── scoring/
│       ├── scoring.go       # Letter frequency scoring
│       └── scoring_test.go
├── words/
│   └── words.txt            # 5-letter word list
├── .github/workflows/ci.yml
├── .gitignore
├── .golangci.yml
├── .claude/skills/
│   └── add-feature.md
├── CLAUDE.md
├── Makefile
├── README.md
└── go.mod
```

## Coding Conventions

- **Naming**: idiomatic Go — `camelCase` for unexported, `PascalCase` for exported, short receiver names
- **Error handling**: return errors, don't panic. Wrap with `fmt.Errorf("context: %w", err)`
- **Imports**: stdlib first, then external, then internal — gofumpt handles grouping
- **Packages**: `internal/` for all non-public packages. Keep packages small and focused.
- **Comments**: exported functions and types must have doc comments

## Commit Messages

Use conventional commits:

- `feat: <description>` — new feature
- `fix: <description>` — bug fix
- `chore: <description>` — maintenance, deps, CI
- `test: <description>` — adding/fixing tests
- `refactor: <description>` — code change that neither fixes nor adds

Include scope when useful: `feat(solver): add hard-mode constraint checking`

## Don't

- Don't use `os.Exit()` outside of `main.go`
- Don't use `panic` for recoverable errors
- Don't use `interface{}` / `any` without a strong reason
- Don't add external dependencies without justification — stdlib is preferred
- Don't skip error checks (errcheck linter will catch this)
- Don't hardcode the word list — load from embedded file or flag
- Don't use global mutable state

## Skills

See `.claude/skills/` for workflow guides:

- `add-feature.md` — how to add new solver features
