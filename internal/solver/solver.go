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
	HardMode   bool
	guesses    []Guess
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
	s.guesses = append(s.guesses, g)
	var filtered []string
	for _, w := range s.candidates {
		if matches(w, g) {
			filtered = append(filtered, w)
		}
	}
	s.candidates = filtered
}

// Suggest returns the best next guess from remaining candidates.
// In hard mode, filters candidates to only those satisfying all known constraints.
func (s *Solver) Suggest() string {
	candidates := s.candidates
	if s.HardMode && len(s.guesses) > 0 {
		candidates = s.filterByConstraints(candidates)
	}
	return scoring.Best(candidates)
}

// filterByConstraints returns only words that satisfy all known constraints from previous guesses.
func (s *Solver) filterByConstraints(words []string) []string {
	var filtered []string
	for _, w := range words {
		if s.satisfiesConstraints(w) {
			filtered = append(filtered, w)
		}
	}
	return filtered
}

// satisfiesConstraints checks if a word satisfies all known green and yellow constraints.
func (s *Solver) satisfiesConstraints(word string) bool {
	for _, g := range s.guesses {
		for i := 0; i < 5; i++ {
			ch := g.Word[i]
			switch g.Feedback[i] {
			case Green:
				// Green: letter must be at this position
				if word[i] != ch {
					return false
				}
			case Yellow:
				// Yellow: letter must appear somewhere in the word
				if !contains(word, ch) {
					return false
				}
			}
		}
	}
	return true
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
