package scoring

import "testing"

func TestBest(t *testing.T) {
	candidates := []string{"CRANE", "SKULL", "PLANT"}
	best := Best(candidates)
	if best == "" {
		t.Fatal("expected a best word, got empty")
	}
	// Just verify it returns one of the candidates.
	found := false
	for _, w := range candidates {
		if w == best {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("best %q not in candidates", best)
	}
}

func TestFrequencyScoreEmpty(t *testing.T) {
	scores := FrequencyScore(nil)
	if scores != nil {
		t.Errorf("expected nil for empty input, got %v", scores)
	}
}
