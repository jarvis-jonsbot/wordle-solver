package wordlist

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	input := "crane\nskull\nhi\nplant\n"
	words, err := Load(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(words) != 3 {
		t.Fatalf("expected 3 words, got %d", len(words))
	}
	if words[0] != "CRANE" {
		t.Errorf("expected CRANE, got %s", words[0])
	}
}
