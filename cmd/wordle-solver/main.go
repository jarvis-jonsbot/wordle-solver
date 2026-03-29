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
	"github.com/jarvis-jonsbot/wordle-solver/internal/types"
	"github.com/jarvis-jonsbot/wordle-solver/internal/wordlist"
)

//go:embed words.txt priority_words.txt
var wordsFS embed.FS

func run() error {
	hardMode := flag.Bool("hard-mode", false, "Enable hard mode (must use revealed hints)")
	jsonOutput := flag.Bool("json", false, "Output each round as JSONL")
	scorerName := flag.String("scorer", "weighted-entropy", "Scoring algorithm: frequency, entropy, or weighted-entropy")
	priorityWeight := flag.Float64("priority-weight", 0.1, "Weight for non-priority words in weighted-entropy scorer (0=pure priority, 1=uniform)")
	opener := flag.String("opener", "SALET", "Opening word (leave empty to compute)")
	flag.Parse()

	// Load priority word list (used by weighted-entropy scorer).
	pf, err := wordsFS.Open("priority_words.txt")
	if err != nil {
		return fmt.Errorf("opening priority word list: %w", err)
	}
	priorityWords, err := wordlist.Load(pf)
	_ = pf.Close()
	if err != nil {
		return fmt.Errorf("loading priority words: %w", err)
	}
	prioritySet := make(map[string]struct{}, len(priorityWords))
	for _, w := range priorityWords {
		prioritySet[w] = struct{}{}
	}

	var sc scoring.Scorer
	switch strings.ToLower(*scorerName) {
	case "frequency":
		sc = scoring.FrequencyScorer{}
	case "entropy":
		sc = scoring.EntropyScorer{}
	case "weighted-entropy":
		sc = scoring.WeightedEntropyScorer{
			PriorityWords:  prioritySet,
			PriorityWeight: *priorityWeight,
		}
	default:
		return fmt.Errorf("unknown scorer %q (use 'frequency', 'entropy', or 'weighted-entropy')", *scorerName)
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

	s := solver.New(words, solver.WithScorer(sc), solver.WithHardMode(*hardMode), solver.WithAllWords(words))

	// Determine the opening suggestion.
	var openingSuggestion string
	if *opener != "" {
		openingSuggestion = strings.ToUpper(*opener)
	} else {
		openingSuggestion = s.Suggest()
	}

	if !*jsonOutput {
		fmt.Printf("Loaded %d words (%d priority).\n", len(s.Candidates()), len(priorityWords))
		if *hardMode {
			fmt.Println("Hard mode enabled.")
		}
		fmt.Printf("Scorer: %s\n", *scorerName)
		if strings.ToLower(*scorerName) == "weighted-entropy" {
			fmt.Printf("Priority weight (β): %.2f\n", *priorityWeight)
		}
		fmt.Printf("Best opener: %s\n\n", openingSuggestion)
	}

	enc := json.NewEncoder(os.Stdout)

	for round := 1; round <= 6; round++ {
		var suggestion string
		if round == 1 {
			suggestion = openingSuggestion
		} else {
			suggestion = s.Suggest()
		}
		if !*jsonOutput {
			fmt.Printf("Round %d (%d candidates remaining)\n", round, len(s.Candidates()))
			fmt.Printf("Suggestion: %s\n", suggestion)
		}

		// Read and validate guess.
		var guess string
		for {
			if !*jsonOutput {
				fmt.Print("Enter guess: ")
			}
			if _, err := fmt.Scanln(&guess); err != nil {
				return fmt.Errorf("reading guess: %w", err)
			}
			guess = strings.ToUpper(strings.TrimSpace(guess))
			if len(guess) != 5 {
				if !*jsonOutput {
					fmt.Println("Guess must be 5 letters.")
				}
				continue
			}
			if *hardMode && !s.ValidateHardMode(guess) {
				if !*jsonOutput {
					fmt.Println("Invalid guess — hard mode requires using all revealed hints.")
				}
				continue
			}
			break
		}

		// Read and validate feedback.
		var fb string
		for {
			if !*jsonOutput {
				fmt.Print("Enter feedback (G=green, Y=yellow, .=gray): ")
			}
			if _, err := fmt.Scanln(&fb); err != nil {
				return fmt.Errorf("reading feedback: %w", err)
			}
			fb = strings.ToUpper(strings.TrimSpace(fb))
			if len(fb) != 5 {
				if !*jsonOutput {
					fmt.Println("Feedback must be 5 characters.")
				}
				continue
			}
			break
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
			out := types.RoundOutput{
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
