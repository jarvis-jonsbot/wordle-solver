// Package wordlist handles loading and filtering 5-letter word lists.
package wordlist

import (
	"bufio"
	"io"
	"strings"
)

// Load reads a newline-delimited list of 5-letter words from r.
func Load(r io.Reader) ([]string, error) {
	var words []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		w := strings.TrimSpace(scanner.Text())
		if len(w) == 5 {
			words = append(words, strings.ToUpper(w))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil
}
