// Package scoring ranks candidate words using pluggable scoring algorithms.
package scoring

import (
	"math"
)

// Scorer scores candidate words and returns a map of word → score.
// Higher scores indicate better guesses.
type Scorer interface {
	Score(candidates []string) map[string]float64
}

// Best returns the highest-scoring word from candidates using the given scorer.
func Best(candidates []string, scorer Scorer) string {
	scores := scorer.Score(candidates)
	if len(scores) == 0 {
		return ""
	}
	best := candidates[0]
	for _, w := range candidates {
		if scores[w] > scores[best] {
			best = w
		}
	}
	return best
}

// FrequencyScorer scores words by per-position letter frequency.
type FrequencyScorer struct{}

// Score computes frequency-based scores for each candidate.
func (FrequencyScorer) Score(candidates []string) map[string]float64 {
	var freq [5][26]float64
	for _, w := range candidates {
		for i, ch := range w {
			freq[i][ch-'A']++
		}
	}
	n := float64(len(candidates))
	if n == 0 {
		return nil
	}

	scores := make(map[string]float64, len(candidates))
	for _, w := range candidates {
		var s float64
		seen := [26]bool{}
		for i, ch := range w {
			idx := ch - 'A'
			if !seen[idx] {
				s += freq[i][idx] / n
				seen[idx] = true
			}
		}
		scores[w] = s
	}
	return scores
}

// EntropyScorer scores words by expected information gain (Shannon entropy
// over the distribution of possible feedback patterns).
type EntropyScorer struct{}

// Score computes entropy-based scores for each candidate.
// For each potential guess, it simulates feedback against every possible answer
// and measures how evenly the guess partitions the remaining candidates.
func (EntropyScorer) Score(candidates []string) map[string]float64 {
	if len(candidates) == 0 {
		return nil
	}
	scores := make(map[string]float64, len(candidates))
	n := float64(len(candidates))

	for _, guess := range candidates {
		// Count how many candidates produce each feedback pattern.
		buckets := make(map[[5]byte]int)
		for _, answer := range candidates {
			pattern := computeFeedback(guess, answer)
			buckets[pattern]++
		}
		// Shannon entropy: -Σ p·log2(p)
		var entropy float64
		for _, count := range buckets {
			p := float64(count) / n
			if p > 0 {
				entropy -= p * math.Log2(p)
			}
		}
		scores[guess] = entropy
	}
	return scores
}

// computeFeedback returns the Wordle feedback pattern for a guess against an answer.
// 'G' = green, 'Y' = yellow, '.' = gray.
func computeFeedback(guess, answer string) [5]byte {
	var pattern [5]byte
	used := [5]bool{}

	// First pass: greens.
	for i := 0; i < 5; i++ {
		if guess[i] == answer[i] {
			pattern[i] = 'G'
			used[i] = true
		}
	}

	// Second pass: yellows.
	for i := 0; i < 5; i++ {
		if pattern[i] == 'G' {
			continue
		}
		for j := 0; j < 5; j++ {
			if !used[j] && guess[i] == answer[j] {
				pattern[i] = 'Y'
				used[j] = true
				break
			}
		}
		if pattern[i] == 0 {
			pattern[i] = '.'
		}
	}
	return pattern
}
