# wordle-solver

A Wordle solver engine in Go. Suggests optimal guesses using pluggable scoring algorithms and constraint elimination.

## Quick Start

```bash
make build
./bin/wordle-solver
```

## Usage

### Interactive Mode (default)

```text
$ ./bin/wordle-solver
Loaded 14855 words.
Scorer: frequency
Best opener: CRANE

Round 1 (14855 candidates remaining)
Suggestion: CRANE
Enter guess: CRANE
Enter feedback (G=green, Y=yellow, .=gray): .Y..Y
Remaining: 521 candidates

Round 2 (521 candidates remaining)
Suggestion: SUPER
...
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--hard-mode` | `false` | Enforce Wordle hard mode (must reuse revealed hints) |
| `--json` | `false` | Output each round as a JSON line (JSONL) |
| `--scorer` | `frequency` | Scoring algorithm: `frequency` or `entropy` |

### Hard Mode (`--hard-mode`)

In hard mode, every subsequent guess must use all revealed hints:

- Green letters must stay in their position
- Yellow letters must appear somewhere in the guess

```bash
./bin/wordle-solver --hard-mode
```

### JSON Output (`--json`)

Outputs one JSON object per round (JSONL format), useful for piping to other tools or skill integration:

```bash
echo -e "CRANE\n..YG.\nTRAIL\nGGGGG" | ./bin/wordle-solver --json
```

```json
{"round":1,"candidates":14855,"suggestion":"CRANE","guess":"CRANE","feedback":"..YG.","remaining":168}
{"round":2,"candidates":168,"suggestion":"TRAIL","guess":"TRAIL","feedback":"GGGGG","remaining":168,"solved":true}
```

Each line contains:

| Field | Description |
|-------|-------------|
| `round` | Round number (1-6) |
| `candidates` | Candidates remaining at start of round |
| `suggestion` | Solver's recommended guess |
| `guess` | The guess that was entered |
| `feedback` | Feedback string (G/Y/.) |
| `remaining` | Candidates remaining after applying feedback |
| `solved` | `true` on the final round if solved (omitted otherwise) |

### Scoring Algorithms (`--scorer`)

**`frequency`** (default) — Scores words by per-position letter frequency. Prefers words with diverse, commonly-occurring letters. Fast, good for most games.

**`entropy`** — Information-theoretic scoring. For each candidate, simulates all possible feedback patterns and picks the guess that maximizes Shannon entropy (expected information gain). Slower but generally optimal.

```bash
# Compare openers
./bin/wordle-solver --scorer frequency  # suggests CRANE
./bin/wordle-solver --scorer entropy    # may suggest differently
```

### Combining Flags

```bash
# Hard mode with entropy scoring, JSON output
echo -e "CRANE\n..YG.\nTRAIL\nGGGGG" | ./bin/wordle-solver --hard-mode --scorer entropy --json
```

## Commands

```bash
make build   # Build binary to bin/
make test    # Run tests with race detection + coverage
make lint    # Run golangci-lint
make dev     # Run interactively
```

## Project Structure

```text
cmd/wordle-solver/   CLI entrypoint + embedded word list
internal/solver/     Core constraint engine (green/yellow/gray filtering)
internal/scoring/    Pluggable scoring algorithms (frequency, entropy)
internal/wordlist/   Word list loader
```

## License

MIT
