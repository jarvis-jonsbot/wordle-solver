// Command wordle-solver is an interactive Wordle solving CLI.
package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jarvis-jonsbot/wordle-solver/internal/solver"
	"github.com/jarvis-jonsbot/wordle-solver/internal/wordlist"
)

//go:embed words.txt
var wordsFS embed.FS

// RoundOutput represents the JSON output for a single round.
type RoundOutput struct {
	Round      int    `json:"round"`
	Candidates int    `json:"candidates"`
	Suggestion string `json:"suggestion"`
	Guess      string `json:"guess,omitempty"`
	Feedback   string `json:"feedback,omitempty"`
	Remaining  int    `json:"remaining,omitempty"`
	Solved     *bool  `json:"solved,omitempty"`
}

func run(hardMode, jsonMode bool) error {
	f, err := wordsFS.Open("words.txt")
	if err != nil {
		return fmt.Errorf("opening word list: %w", err)
	}
	defer func() { _ = f.Close() }()

	words, err := wordlist.Load(f)
	if err != nil {
		return fmt.Errorf("loading words: %w", err)
	}

	s := solver.New(words)
	s.HardMode = hardMode

	if jsonMode {
		return runJSON(s)
	}
	return runInteractive(s)
}

func runInteractive(s *solver.Solver) error {
	fmt.Printf("Loaded %d words.\n", len(s.Candidates()))
	fmt.Printf("Best opener: %s\n\n", s.Suggest())

	for round := 1; round <= 6; round++ {
		fmt.Printf("Round %d (%d candidates remaining)\n", round, len(s.Candidates()))
		fmt.Printf("Suggestion: %s\n", s.Suggest())
		fmt.Print("Enter guess: ")

		var guess string
		if _, err := fmt.Scanln(&guess); err != nil {
			return fmt.Errorf("reading guess: %w", err)
		}
		guess = strings.ToUpper(strings.TrimSpace(guess))
		if len(guess) != 5 {
			fmt.Println("Guess must be 5 letters.")
			round--
			continue
		}

		fmt.Print("Enter feedback (G=green, Y=yellow, .=gray): ")
		var fb string
		if _, err := fmt.Scanln(&fb); err != nil {
			return fmt.Errorf("reading feedback: %w", err)
		}
		fb = strings.ToUpper(strings.TrimSpace(fb))
		if len(fb) != 5 {
			fmt.Println("Feedback must be 5 characters.")
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

		if fb == "GGGGG" {
			fmt.Println("🎉 Solved!")
			return nil
		}

		s.Apply(solver.Guess{Word: guess, Feedback: feedback})
		fmt.Printf("Remaining: %d candidates\n\n", len(s.Candidates()))

		if len(s.Candidates()) == 0 {
			fmt.Println("No candidates remaining — word not in list?")
			return nil
		}
	}
	fmt.Println("Out of guesses!")
	return nil
}

func runJSON(s *solver.Solver) error {
	scanner := bufio.NewScanner(os.Stdin)

	for round := 1; round <= 6; round++ {
		candidateCount := len(s.Candidates())
		suggestion := s.Suggest()

		// Read guess from stdin
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("reading input: %w", err)
			}
			// EOF - no more input
			break
		}
		guess := strings.ToUpper(strings.TrimSpace(scanner.Text()))
		if len(guess) != 5 {
			return fmt.Errorf("invalid guess length: %s", guess)
		}

		// Read feedback from stdin
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("reading input: %w", err)
			}
			return fmt.Errorf("missing feedback for guess %s", guess)
		}
		fb := strings.ToUpper(strings.TrimSpace(scanner.Text()))
		if len(fb) != 5 {
			return fmt.Errorf("invalid feedback length: %s", fb)
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

		// Apply feedback before outputting remaining count
		s.Apply(solver.Guess{Word: guess, Feedback: feedback})
		remaining := len(s.Candidates())

		output := RoundOutput{
			Round:      round,
			Candidates: candidateCount,
			Suggestion: suggestion,
			Guess:      guess,
			Feedback:   fb,
			Remaining:  remaining,
		}

		if solved {
			trueVal := true
			output.Solved = &trueVal
		}

		if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
			return fmt.Errorf("encoding JSON: %w", err)
		}

		if solved {
			return nil
		}

		if remaining == 0 {
			falseVal := false
			finalOutput := RoundOutput{
				Round:  round,
				Solved: &falseVal,
			}
			if err := json.NewEncoder(os.Stdout).Encode(finalOutput); err != nil {
				return fmt.Errorf("encoding JSON: %w", err)
			}
			return nil
		}
	}

	// Out of rounds
	falseVal := false
	finalOutput := RoundOutput{
		Solved: &falseVal,
	}
	return json.NewEncoder(os.Stdout).Encode(finalOutput)
}

func main() {
	hardMode := flag.Bool("hard-mode", false, "Enable hard mode (revealed hints must be used in subsequent guesses)")
	jsonMode := flag.Bool("json", false, "Output JSON lines instead of interactive mode")
	flag.Parse()

	if err := run(*hardMode, *jsonMode); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
