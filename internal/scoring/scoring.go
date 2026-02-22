// Package scoring ranks candidate words by letter frequency.
package scoring

// FrequencyScore scores each word by the sum of per-position letter frequencies
// across the candidate list. Higher scores mean more common letters.
func FrequencyScore(candidates []string) map[string]float64 {
	// Count letter frequency at each position.
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
			// Only count each letter once to prefer diverse letters.
			if !seen[idx] {
				s += freq[i][idx] / n
				seen[idx] = true
			}
		}
		scores[w] = s
	}
	return scores
}

// Best returns the highest-scoring word from the candidates.
func Best(candidates []string) string {
	scores := FrequencyScore(candidates)
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
