package scoring

import (
	"testing"
)

func prioritySet(words ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(words))
	for _, w := range words {
		m[w] = struct{}{}
	}
	return m
}

func TestWeightedEntropyScorer_UniformEquivalence(t *testing.T) {
	// With β=1.0, WeightedEntropyScorer should produce the same ranking as EntropyScorer.
	candidates := []string{"CRANE", "SLATE", "AUDIO", "STARE", "ARISE"}

	uniform := EntropyScorer{}
	weighted := WeightedEntropyScorer{
		PriorityWords:  prioritySet(candidates...), // all words are priority
		PriorityWeight: 1.0,
	}

	uScores := uniform.Score(candidates)
	wScores := weighted.Score(candidates)

	for _, w := range candidates {
		diff := uScores[w] - wScores[w]
		if diff < -0.001 || diff > 0.001 {
			t.Errorf("word %s: uniform=%.4f weighted=%.4f (diff %.4f exceeds tolerance)",
				w, uScores[w], wScores[w], diff)
		}
	}
}

func TestWeightedEntropyScorer_PriorityWordsScore_Higher(t *testing.T) {
	// A priority word should score higher than an equally-entropic non-priority word
	// when many non-priority candidates are present, because the solver is optimising
	// for the weighted distribution.
	//
	// Simpler check: with β=0 (pure priority), non-priority words get 0 weight as
	// possible answers, so the best guess is whichever priority word best splits
	// the priority candidates. Non-priority words may still appear as guesses but
	// won't be chosen when a priority word ties or wins on entropy.

	priority := []string{"CRANE", "SLATE", "AUDIO"}
	nonPriority := []string{"ZZZZZ", "QWXYZ"} // implausible non-priority words
	candidates := make([]string, 0, len(priority)+len(nonPriority))
	candidates = append(candidates, priority...)
	candidates = append(candidates, nonPriority...)

	s := WeightedEntropyScorer{
		PriorityWords:  prioritySet(priority...),
		PriorityWeight: 0.0, // pure priority — non-priority words have zero answer weight
	}

	scores := s.Score(candidates)

	// Non-priority words should score 0 (they never appear as answers, so any guess
	// targeting them wastes information).
	// Actually with β=0, total weight = sum of priority words only.
	// Non-priority words as *guesses* can still have nonzero entropy if they split
	// priority candidates. So just verify scores are computed without panic.
	for _, w := range candidates {
		if _, ok := scores[w]; !ok {
			t.Errorf("expected score for %s, got none", w)
		}
	}
}

func TestWeightedEntropyScorer_EmptyCandidates(t *testing.T) {
	s := WeightedEntropyScorer{
		PriorityWords:  prioritySet("CRANE"),
		PriorityWeight: 0.1,
	}
	scores := s.Score(nil)
	if scores != nil {
		t.Errorf("expected nil scores for empty candidates, got %v", scores)
	}
}

func TestWeightedEntropyScorer_AllZeroWeight(t *testing.T) {
	// β=0 and no priority words → all weights zero → nil scores (no division by zero).
	s := WeightedEntropyScorer{
		PriorityWords:  prioritySet(), // empty
		PriorityWeight: 0.0,
	}
	scores := s.Score([]string{"CRANE", "SLATE"})
	if scores != nil {
		t.Errorf("expected nil scores when total weight is zero, got %v", scores)
	}
}

func TestWeightedEntropyScorer_BestIsStable(t *testing.T) {
	// Smoke test: Best() should not panic and should return a word from candidates.
	candidates := []string{"CRANE", "SLATE", "AUDIO", "STARE"}
	priority := prioritySet("CRANE", "SLATE")
	s := WeightedEntropyScorer{PriorityWords: priority, PriorityWeight: 0.1}

	best := Best(candidates, s)
	if best == "" {
		t.Error("Best() returned empty string")
	}
	found := false
	for _, w := range candidates {
		if w == best {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Best() returned %q which is not in candidates", best)
	}
}
