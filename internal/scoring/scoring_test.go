package scoring

import (
	"math"
	"testing"
)

func TestFrequencyScorerEmpty(t *testing.T) {
	s := FrequencyScorer{}
	got := s.Score(nil)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestFrequencyScorerScores(t *testing.T) {
	candidates := []string{"CRANE", "CRATE", "TRACE"}
	s := FrequencyScorer{}
	scores := s.Score(candidates)
	if len(scores) != 3 {
		t.Fatalf("expected 3 scores, got %d", len(scores))
	}
	for _, w := range candidates {
		if _, ok := scores[w]; !ok {
			t.Errorf("missing score for %s", w)
		}
	}
}

func TestFrequencyScorerBest(t *testing.T) {
	candidates := []string{"CRANE", "CRATE", "TRACE"}
	best := Best(candidates, FrequencyScorer{})
	if best == "" {
		t.Fatal("expected a best word")
	}
}

func TestEntropyScorerEmpty(t *testing.T) {
	s := EntropyScorer{}
	got := s.Score(nil)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestEntropyScorerSingleWord(t *testing.T) {
	s := EntropyScorer{}
	scores := s.Score([]string{"CRANE"})
	// Single word: only one feedback pattern (GGGGG), entropy = 0.
	if scores["CRANE"] != 0 {
		t.Errorf("expected 0 entropy for single candidate, got %f", scores["CRANE"])
	}
}

func TestEntropyScorerMultipleWords(t *testing.T) {
	candidates := []string{"CRANE", "CRATE", "TRACE", "BLUES", "MOUND"}
	s := EntropyScorer{}
	scores := s.Score(candidates)
	if len(scores) != 5 {
		t.Fatalf("expected 5 scores, got %d", len(scores))
	}
	// All scores should be >= 0.
	for w, sc := range scores {
		if sc < 0 {
			t.Errorf("negative entropy for %s: %f", w, sc)
		}
	}
}

func TestEntropyScorerBest(t *testing.T) {
	candidates := []string{"CRANE", "CRATE", "TRACE", "BLUES", "MOUND"}
	best := Best(candidates, EntropyScorer{})
	if best == "" {
		t.Fatal("expected a best word")
	}
}

func TestComputeFeedback(t *testing.T) {
	tests := []struct {
		guess, answer string
		want          [5]byte
	}{
		{"CRANE", "CRANE", [5]byte{'G', 'G', 'G', 'G', 'G'}},
		{"CRANE", "BLUES", [5]byte{'.', '.', '.', '.', 'Y'}},
		{"CRANE", "KNEEL", [5]byte{'.', '.', '.', 'Y', 'Y'}},
		{"SPEED", "ABIDE", [5]byte{'.', '.', 'Y', '.', 'Y'}},
		{"HELLO", "LLAMA", [5]byte{'.', '.', 'Y', 'Y', '.'}},
	}
	for _, tt := range tests {
		got := computeFeedback(tt.guess, tt.answer)
		if got != tt.want {
			t.Errorf("computeFeedback(%s, %s) = %s, want %s", tt.guess, tt.answer, string(got[:]), string(tt.want[:]))
		}
	}
}

func TestBestEmpty(t *testing.T) {
	got := Best(nil, FrequencyScorer{})
	if got != "" {
		t.Errorf("expected empty, got %s", got)
	}
}

func TestEntropyScorerDistinct(t *testing.T) {
	// Two identical words should have 0 entropy.
	candidates := []string{"AAAAA", "AAAAA"}
	s := EntropyScorer{}
	scores := s.Score(candidates)
	if math.Abs(scores["AAAAA"]) > 1e-9 {
		t.Errorf("expected ~0 entropy for identical words, got %f", scores["AAAAA"])
	}
}
