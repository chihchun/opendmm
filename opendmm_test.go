package opendmm

import (
	"sync"
	"testing"
)

func TestSearch(t *testing.T) {
	queries := []string{
		"SDDE-201",
		"200GANA-894",
	}
	wg := new(sync.WaitGroup)
	for _, query := range queries {
		wg.Add(1)
		go func(query string) {
			defer wg.Done()
			metach := Search(query)
			meta, ok := <-metach
			if !ok {
				t.Errorf("%s not found", query)
			} else {
				t.Logf("%s -> %+v", query, meta)
			}
		}(query)
	}
	wg.Wait()
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
