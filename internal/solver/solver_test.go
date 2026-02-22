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
