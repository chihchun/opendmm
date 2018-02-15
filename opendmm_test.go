package opendmm

import (
	"testing"
)

func TestOpendmmSearch(t *testing.T) {
	queries := []string{
		"SDDE-201",
	}
	for _, query := range queries {
		metach := Search(query)
		meta, ok := <-metach
		if !ok {
			t.Errorf("%s not found", query)
		} else {
			t.Logf("%s -> %+v", query, meta)
		}
	}
}

func TestOpendmmGuess(t *testing.T) {
	queries := []string{
		"SDDE-201",
	}
	for _, query := range queries {
		if !Guess(query + "_suffix").Contains(query) {
			t.Errorf("Guessed wrong code for %s with suffix", query)
		}
		if !Guess("prefix_" + query).Contains(query) {
			t.Errorf("Guessed wrong code for %s with prefix", query)
		}
	}
}
