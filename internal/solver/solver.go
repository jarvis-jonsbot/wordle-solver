// Package solver implements the core Wordle constraint logic.
package solver

import (
	"github.com/jarvis-jonsbot/wordle-solver/internal/scoring"
)

// Feedback represents the result of a single letter guess.
type Feedback int

const (
	Gray   Feedback = iota // Letter not in word
	Yellow                 // Letter in word but wrong position
	Green                  // Letter in correct position
)

// Guess holds a guessed word and its per-letter feedback.
type Guess struct {
	Word     string
	Feedback [5]Feedback
}

// Solver tracks constraints and filters candidates.
type Solver struct {
	candidates []string
}

// New creates a solver with the given word list.
func New(words []string) *Solver {
	cp := make([]string, len(words))
	copy(cp, words)
	return &Solver{candidates: cp}
}

// Candidates returns the current candidate list.
func (s *Solver) Candidates() []string {
	return s.candidates
}

// Apply narrows candidates based on a guess and its feedback.
func (s *Solver) Apply(g Guess) {
	var filtered []string
	for _, w := range s.candidates {
		if matches(w, g) {
			filtered = append(filtered, w)
		}
	}
	s.candidates = filtered
}

// Suggest returns the best next guess from remaining candidates.
func (s *Solver) Suggest() string {
	return scoring.Best(s.candidates)
}

// matches checks if a candidate word is consistent with the given guess feedback.
func matches(word string, g Guess) bool {
	for i := 0; i < 5; i++ {
		ch := g.Word[i]
		switch g.Feedback[i] {
		case Green:
			if word[i] != ch {
				return false
			}
		case Yellow:
			if word[i] == ch {
				return false
			}
			if !contains(word, ch) {
				return false
			}
		case Gray:
			// Gray means the letter doesn't appear in any non-green/yellow position.
			// Simple check: if it's not green/yellow elsewhere, it shouldn't be in the word.
			if !hasGreenOrYellow(g, ch) && contains(word, ch) {
				return false
			}
			// If the letter IS green/yellow elsewhere, just ensure it's not at this position.
			if hasGreenOrYellow(g, ch) && word[i] == ch {
				return false
			}
		}
	}
	return true
}

func contains(word string, ch byte) bool {
	for i := 0; i < len(word); i++ {
		if word[i] == ch {
			return true
		}
	}
	return false
}

func hasGreenOrYellow(g Guess, ch byte) bool {
	for i := 0; i < 5; i++ {
		if g.Word[i] == ch && (g.Feedback[i] == Green || g.Feedback[i] == Yellow) {
			return true
		}
	}
	return false
}
