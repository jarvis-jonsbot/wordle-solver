package scoring

import (
	"math"
)

// WeightedEntropyScorer scores words by expected information gain using a
// non-uniform prior over candidate words. Words in the priority set are
// assigned weight 1.0; all other words are assigned PriorityWeight (β).
//
// This encodes the belief that historically-curated answer words are more
// likely to be chosen than obscure-but-valid words, without hard-excluding
// the latter (since the NYT editor has added new words not in the original list).
//
// β=1.0  → uniform prior (equivalent to EntropyScorer)
// β=0.0  → only priority words are considered possible answers
// β=0.1  → priority words are 10× more likely (recommended default)
type WeightedEntropyScorer struct {
	// PriorityWords is the set of words with weight 1.0.
	PriorityWords map[string]struct{}
	// PriorityWeight (β) is the weight assigned to non-priority words.
	// Must be in [0, 1]. Default 0.1.
	PriorityWeight float64
	// AllWords is the complete word list to consider as potential guesses.
	// If empty, only candidates are scored.
	AllWords []string
}

// SetAllWords implements AllWordsScorer.
func (s *WeightedEntropyScorer) SetAllWords(allWords []string) {
	s.AllWords = allWords
}

// weight returns the prior probability weight for a candidate word.
func (s WeightedEntropyScorer) weight(word string) float64 {
	if _, ok := s.PriorityWords[word]; ok {
		return 1.0
	}
	beta := s.PriorityWeight
	if beta <= 0 {
		return 0.0
	}
	return beta
}

// Score computes weighted-entropy scores for each candidate.
// Each candidate's contribution to a feedback bucket is its weight, not 1.
// The entropy is then computed over the resulting weighted distribution.
// If AllWords is set, all words are considered as potential guesses; otherwise
// only candidates are scored.
func (s WeightedEntropyScorer) Score(candidates []string) map[string]float64 {
	if len(candidates) == 0 {
		return nil
	}

	// Total weight across all candidates (normalisation constant).
	totalWeight := 0.0
	for _, w := range candidates {
		totalWeight += s.weight(w)
	}
	if totalWeight == 0 {
		return nil
	}

	// Determine which words to score as guesses.
	guessPool := candidates
	if len(s.AllWords) > 0 {
		guessPool = s.AllWords
	}

	scores := make(map[string]float64, len(guessPool))

	for _, guess := range guessPool {
		// Accumulate weighted probability mass per feedback pattern.
		buckets := make(map[[5]byte]float64)
		for _, answer := range candidates {
			pattern := computeFeedback(guess, answer)
			buckets[pattern] += s.weight(answer)
		}

		// Shannon entropy: -Σ p·log2(p) where p = bucket_weight / total_weight
		var entropy float64
		for _, bw := range buckets {
			p := bw / totalWeight
			if p > 0 {
				entropy -= p * math.Log2(p)
			}
		}
		scores[guess] = entropy
	}

	return scores
}
