# wordle-solver

A Wordle solver engine in Go. Suggests optimal guesses using letter-frequency scoring and constraint elimination.

## Quick Start

```bash
make build
./bin/wordle-solver
```

## Usage

The solver runs interactively:

1. It suggests a word
2. You enter your guess and the feedback (G=green, Y=yellow, .=gray)
3. It narrows candidates and suggests the next guess
4. Repeat until solved

## Commands

```bash
make build   # Build binary to bin/
make test    # Run tests
make lint    # Run linter
make dev     # Run interactively
```

## How It Works

- Loads a 14,800+ word list of valid 5-letter words
- Scores candidates by per-position letter frequency (preferring diverse letters)
- Applies green/yellow/gray constraints to eliminate impossible words
- Suggests the highest-scoring remaining candidate each round

## License

MIT
