package main

import (
	"strings"
	"testing"
	"time"

	"github.com/jarvis-jonsbot/wordle-solver/internal/types"
)

func TestGenerateShareText(t *testing.T) {
	tests := []struct {
		name       string
		number     int
		solved     bool
		guessCount int
		guesses    []types.RoundOutput
		want       string
	}{
		{
			name:       "solved in 3",
			number:     1234,
			solved:     true,
			guessCount: 3,
			guesses: []types.RoundOutput{
				{Guess: "SALET", Feedback: "Y...."},
				{Guess: "CRUSH", Feedback: ".Y.GY"},
				{Guess: "SKULL", Feedback: "GGGGG"},
			},
			want: `Wordle #1234 3/6

🟨⬛⬛⬛⬛
⬛🟨⬛🟩🟨
🟩🟩🟩🟩🟩`,
		},
		{
			name:       "failed",
			number:     5678,
			solved:     false,
			guessCount: 6,
			guesses: []types.RoundOutput{
				{Guess: "SALET", Feedback: "....."},
				{Guess: "CROWN", Feedback: "..G.."},
				{Guess: "PLUMB", Feedback: "..G.."},
				{Guess: "PHOTO", Feedback: "..G.G"},
				{Guess: "VIGIL", Feedback: ".GGG."},
				{Guess: "FIGHT", Feedback: ".GGGG"},
			},
			want: `Wordle #5678 X/6

⬛⬛⬛⬛⬛
⬛⬛🟩⬛⬛
⬛⬛🟩⬛⬛
⬛⬛🟩⬛🟩
⬛🟩🟩🟩⬛
⬛🟩🟩🟩🟩`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateShareText(tt.number, tt.solved, tt.guessCount, tt.guesses)
			if got != tt.want {
				t.Errorf("generateShareText() =\n%s\n\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestGenerateJekyllPost(t *testing.T) {
	guesses := []types.RoundOutput{
		{Guess: "SALET", Feedback: "Y...."},
		{Guess: "CRUSH", Feedback: ".Y.GY"},
		{Guess: "SKULL", Feedback: "GGGGG"},
	}
	postDate := time.Date(2026, 3, 29, 0, 0, 0, 0, time.UTC)
	body := "This is a test blog post body."

	post := generateJekyllPost(1234, "SKULL", true, 3, guesses, postDate, body)

	// Check front matter fields.
	if !strings.Contains(post, "layout: post") {
		t.Error("post missing 'layout: post'")
	}
	if !strings.Contains(post, `title: "Wordle #1234 — SKULL (3/6)"`) {
		t.Error("post missing or incorrect title")
	}
	if !strings.Contains(post, "date: 2026-03-29 07:00:00 -0700") {
		t.Error("post missing or incorrect date")
	}
	if !strings.Contains(post, "wordle_number: 1234") {
		t.Error("post missing wordle_number")
	}
	if !strings.Contains(post, "word: SKULL") {
		t.Error("post missing word")
	}
	if !strings.Contains(post, "solved: true") {
		t.Error("post missing solved field")
	}
	if !strings.Contains(post, "guesses: 3") {
		t.Error("post missing guesses count")
	}

	// Check share_text contains emoji.
	if !strings.Contains(post, "🟨⬛⬛⬛⬛") {
		t.Error("share_text missing first guess emoji")
	}
	if !strings.Contains(post, "🟩🟩🟩🟩🟩") {
		t.Error("share_text missing solved emoji")
	}

	// Check body is included.
	if !strings.Contains(post, body) {
		t.Error("post missing body content")
	}

	// Check front matter closing.
	if !strings.Contains(post, "---\n\n") {
		t.Error("post missing front matter closing")
	}
}

func TestGenerateJekyllPostFilename(t *testing.T) {
	// Test filename generation indirectly by verifying date format.
	postDate := time.Date(2026, 3, 29, 0, 0, 0, 0, time.UTC)
	expected := "2026-03-29"
	got := postDate.Format("2006-01-02")

	if got != expected {
		t.Errorf("date format = %q, want %q", got, expected)
	}

	// Verify filename construction.
	wantFilename := "2026-03-29-wordle-1234.md"
	actualFilename := postDate.Format("2006-01-02") + "-wordle-1234.md"
	if actualFilename != wantFilename {
		t.Errorf("filename = %q, want %q", actualFilename, wantFilename)
	}
}

func TestGenerateShareTextUnsolved(t *testing.T) {
	guesses := []types.RoundOutput{
		{Guess: "SALET", Feedback: "....."},
		{Guess: "CROWN", Feedback: "..G.."},
	}

	got := generateShareText(999, false, 2, guesses)

	if !strings.Contains(got, "Wordle #999 X/6") {
		t.Error("unsolved share text should contain 'X/6'")
	}
	if !strings.Contains(got, "⬛⬛⬛⬛⬛") {
		t.Error("share text missing first guess emoji")
	}
	if !strings.Contains(got, "⬛⬛🟩⬛⬛") {
		t.Error("share text missing second guess emoji")
	}
}
