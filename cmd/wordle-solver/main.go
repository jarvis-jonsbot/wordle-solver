// Command wordle-solver is an interactive Wordle solving CLI.
package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jarvis-jonsbot/wordle-solver/internal/scoring"
	"github.com/jarvis-jonsbot/wordle-solver/internal/solver"
	"github.com/jarvis-jonsbot/wordle-solver/internal/wordlist"
)

//go:embed words.txt
var wordsFS embed.FS

// RoundOutput is the JSON structure for each round in --json mode.
type RoundOutput struct {
	Round      int    `json:"round"`
	Candidates int    `json:"candidates"`
	Suggestion string `json:"suggestion"`
	Guess      string `json:"guess"`
	Feedback   string `json:"feedback"`
	Remaining  int    `json:"remaining"`
	Solved     bool   `json:"solved,omitempty"`
}

func run() error {
	hardMode := flag.Bool("hard-mode", false, "Enable hard mode (must use revealed hints)")
	jsonOutput := flag.Bool("json", false, "Output each round as JSONL")
	scorerName := flag.String("scorer", "frequency", "Scoring algorithm: frequency or entropy")
	flag.Parse()

	var sc scoring.Scorer
	switch strings.ToLower(*scorerName) {
	case "frequency":
		sc = scoring.FrequencyScorer{}
	case "entropy":
		sc = scoring.EntropyScorer{}
	default:
		return fmt.Errorf("unknown scorer %q (use 'frequency' or 'entropy')", *scorerName)
	}

	f, err := wordsFS.Open("words.txt")
	if err != nil {
		return fmt.Errorf("opening word list: %w", err)
	}
	defer func() { _ = f.Close() }()

	words, err := wordlist.Load(f)
	if err != nil {
		return fmt.Errorf("loading words: %w", err)
	}

	s := solver.New(words, solver.WithScorer(sc), solver.WithHardMode(*hardMode))

	if !*jsonOutput {
		fmt.Printf("Loaded %d words.\n", len(s.Candidates()))
		if *hardMode {
			fmt.Println("Hard mode enabled.")
		}
		fmt.Printf("Scorer: %s\n", *scorerName)
		fmt.Printf("Best opener: %s\n\n", s.Suggest())
	}

	enc := json.NewEncoder(os.Stdout)

	for round := 1; round <= 6; round++ {
		suggestion := s.Suggest()
		if !*jsonOutput {
			fmt.Printf("Round %d (%d candidates remaining)\n", round, len(s.Candidates()))
			fmt.Printf("Suggestion: %s\n", suggestion)
			fmt.Print("Enter guess: ")
		}

		var guess string
		if _, err := fmt.Scanln(&guess); err != nil {
			return fmt.Errorf("reading guess: %w", err)
		}
		guess = strings.ToUpper(strings.TrimSpace(guess))
		if len(guess) != 5 {
			if !*jsonOutput {
				fmt.Println("Guess must be 5 letters.")
			}
			round--
			continue
		}

		// Validate hard mode.
		if *hardMode && !s.ValidateHardMode(guess) {
			if !*jsonOutput {
				fmt.Println("Invalid guess — hard mode requires using all revealed hints.")
			}
			round--
			continue
		}

		if !*jsonOutput {
			fmt.Print("Enter feedback (G=green, Y=yellow, .=gray): ")
		}
		var fb string
		if _, err := fmt.Scanln(&fb); err != nil {
			return fmt.Errorf("reading feedback: %w", err)
		}
		fb = strings.ToUpper(strings.TrimSpace(fb))
		if len(fb) != 5 {
			if !*jsonOutput {
				fmt.Println("Feedback must be 5 characters.")
			}
			round--
			continue
		}

		var feedback [5]solver.Feedback
		for i, ch := range fb {
			switch ch {
			case 'G':
				feedback[i] = solver.Green
			case 'Y':
				feedback[i] = solver.Yellow
			default:
				feedback[i] = solver.Gray
			}
		}

		solved := fb == "GGGGG"
		s.Apply(solver.Guess{Word: guess, Feedback: feedback})

		if *jsonOutput {
			out := RoundOutput{
				Round:      round,
				Candidates: len(s.Candidates()),
				Suggestion: suggestion,
				Guess:      guess,
				Feedback:   fb,
				Remaining:  len(s.Candidates()),
				Solved:     solved,
			}
			if err := enc.Encode(out); err != nil {
				return fmt.Errorf("encoding JSON: %w", err)
			}
		}

		if solved {
			if !*jsonOutput {
				fmt.Println("🎉 Solved!")
			}
			return nil
		}

		if !*jsonOutput {
			fmt.Printf("Remaining: %d candidates\n\n", len(s.Candidates()))
		}

		if len(s.Candidates()) == 0 {
			if !*jsonOutput {
				fmt.Println("No candidates remaining — word not in list?")
			}
			return nil
		}
	}

	if !*jsonOutput {
		fmt.Println("Out of guesses!")
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
