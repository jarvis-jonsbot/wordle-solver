package solver

import "testing"

func TestApplyGreen(t *testing.T) {
	s := New([]string{"CRANE", "CRATE", "CRAFT", "PLANT"})
	s.Apply(Guess{
		Word:     "CRANE",
		Feedback: [5]Feedback{Green, Green, Green, Gray, Gray},
	})
	for _, w := range s.Candidates() {
		if w[0] != 'C' || w[1] != 'R' || w[2] != 'A' {
			t.Errorf("candidate %s doesn't match green constraints", w)
		}
	}
}

func TestApplyYellow(t *testing.T) {
	s := New([]string{"CRANE", "REACT", "TRACE", "PLANT"})
	s.Apply(Guess{
		Word:     "CRANE",
		Feedback: [5]Feedback{Yellow, Yellow, Yellow, Gray, Gray},
	})
	// CRANE itself should be excluded (C at pos 0 is yellow = C not at pos 0)
	for _, w := range s.Candidates() {
		if w == "CRANE" {
			t.Error("CRANE should be excluded by yellow feedback")
		}
	}
}

func TestSuggest(t *testing.T) {
	s := New([]string{"CRANE", "SKULL"})
	suggestion := s.Suggest()
	if suggestion != "CRANE" && suggestion != "SKULL" {
		t.Errorf("unexpected suggestion: %s", suggestion)
	}
}

func TestApplyReducesCandidates(t *testing.T) {
	words := []string{"CRANE", "CRATE", "SKULL", "PLANT", "TRACE"}
	s := New(words)
	before := len(s.Candidates())
	s.Apply(Guess{
		Word:     "CRANE",
		Feedback: [5]Feedback{Green, Green, Green, Gray, Gray},
	})
	after := len(s.Candidates())
	if after >= before {
		t.Errorf("expected fewer candidates after apply, got %d -> %d", before, after)
	}
}

func TestHardModeGreenConstraint(t *testing.T) {
	words := []string{"CRANE", "CRATE", "CRAFT", "PLANT", "SKULL"}
	s := New(words)
	s.HardMode = true

	// Apply guess with green C at position 0
	s.Apply(Guess{
		Word:     "CRANE",
		Feedback: [5]Feedback{Green, Gray, Gray, Gray, Gray},
	})

	suggestion := s.Suggest()
	// Suggestion must have C at position 0
	if len(suggestion) > 0 && suggestion[0] != 'C' {
		t.Errorf("hard mode suggestion %s violates green constraint (C at pos 0)", suggestion)
	}
}

func TestHardModeYellowConstraint(t *testing.T) {
	words := []string{"CRANE", "REACT", "TRACE", "PLANT", "ACRES"}
	s := New(words)
	s.HardMode = true

	// Apply guess with yellow C, R, A (must appear in word)
	s.Apply(Guess{
		Word:     "CRANE",
		Feedback: [5]Feedback{Yellow, Yellow, Yellow, Gray, Gray},
	})

	suggestion := s.Suggest()
	// Suggestion must contain C, R, and A
	if len(suggestion) > 0 {
		hasC := contains(suggestion, 'C')
		hasR := contains(suggestion, 'R')
		hasA := contains(suggestion, 'A')
		if !hasC || !hasR || !hasA {
			t.Errorf("hard mode suggestion %s missing yellow letters (C,R,A)", suggestion)
		}
	}
}

func TestHardModeOffAllowsAnyCandidate(t *testing.T) {
	// Simple test: verify hard mode off allows suggestions from remaining candidates
	words := []string{"BATCH", "CATCH", "WATCH"}
	s := New(words)
	s.HardMode = false

	s.Apply(Guess{
		Word:     "LUMPY",
		Feedback: [5]Feedback{Gray, Gray, Gray, Gray, Gray},
	})

	// All gray - L,U,M,P,Y not in word
	// All three candidates (BATCH, CATCH, WATCH) should remain
	suggestion := s.Suggest()
	if suggestion == "" {
		t.Error("expected non-empty suggestion")
	}
	if len(s.Candidates()) != 3 {
		t.Errorf("expected 3 candidates, got %d", len(s.Candidates()))
	}
}
