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
	allWords   []string // complete word list for scoring non-candidates
	scorer     scoring.Scorer
	hardMode   bool
	guesses    []Guess // history for hard-mode enforcement
}

// Option configures a Solver.
type Option func(*Solver)

// WithScorer sets the scoring algorithm.
func WithScorer(s scoring.Scorer) Option {
	return func(sol *Solver) {
		sol.scorer = s
	}
}

// WithHardMode enables hard mode constraint enforcement on suggestions.
func WithHardMode(enabled bool) Option {
	return func(sol *Solver) {
		sol.hardMode = enabled
	}
}

// WithAllWords sets the complete word list for scoring non-candidate guesses.
func WithAllWords(words []string) Option {
	return func(sol *Solver) {
		sol.allWords = words
	}
}

// New creates a solver with the given word list and options.
func New(words []string, opts ...Option) *Solver {
	cp := make([]string, len(words))
	copy(cp, words)
	s := &Solver{
		candidates: cp,
		scorer:     scoring.FrequencyScorer{},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// HardMode returns whether hard mode is enabled.
func (s *Solver) HardMode() bool {
	return s.hardMode
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
// If the scorer implements AllWordsScorer and allWords is set, it will
// consider all valid words as potential guesses.
func (s *Solver) Suggest() string {
	// If scorer supports AllWords, set it before scoring.
	if len(s.allWords) > 0 {
		if aws, ok := s.scorer.(scoring.AllWordsScorer); ok {
			aws.SetAllWords(s.allWords)
		}
	}

	if s.hardMode && len(s.guesses) > 0 {
		valid := s.hardModeFilter(s.candidates)
		if len(valid) > 0 {
			return scoring.Best(valid, s.scorer)
		}
	}
	return scoring.Best(s.candidates, s.scorer)
}

// ValidateHardMode checks if a word satisfies hard-mode constraints.
// Returns true if the word is valid under current constraints.
func (s *Solver) ValidateHardMode(word string) bool {
	for _, g := range s.guesses {
		if !satisfiesHardMode(word, g) {
			return false
		}
	}
	return true
}

// hardModeFilter returns only candidates that satisfy hard-mode constraints.
func (s *Solver) hardModeFilter(words []string) []string {
	var result []string
	for _, w := range words {
		if s.ValidateHardMode(w) {
			result = append(result, w)
		}
	}
	return result
}

// satisfiesHardMode checks if a word uses all revealed hints from a guess.
func satisfiesHardMode(word string, g Guess) bool {
	for i := 0; i < 5; i++ {
		switch g.Feedback[i] {
		case Green:
			if word[i] != g.Word[i] {
				return false
			}
		case Yellow:
			if !contains(word, g.Word[i]) {
				return false
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
			if !hasGreenOrYellow(g, ch) && contains(word, ch) {
				return false
			}
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
