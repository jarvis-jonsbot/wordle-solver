// Package types contains shared types used across multiple commands.
package types

// RoundOutput is the JSON structure for each round in --json mode.
type RoundOutput struct {
	Round      int    `json:"round"`
	Candidates int    `json:"candidates"`
	Suggestion string `json:"suggestion"`
	Guess      string `json:"guess"`
	Feedback   string `json:"feedback"`
	Remaining  int    `json:"remaining"`
	Solved     bool   `json:"solved,omitempty"`
}
