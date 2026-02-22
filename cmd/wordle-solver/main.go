// Command wordle-solver is an interactive Wordle solving CLI.
package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jarvis-jonsbot/wordle-solver/internal/solver"
	"github.com/jarvis-jonsbot/wordle-solver/internal/wordlist"
)

//go:embed words.txt
var wordsFS embed.FS

func run() error {
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

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
