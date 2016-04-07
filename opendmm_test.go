package opendmm

import (
	"testing"
)

func TestOpendmmSearch(t *testing.T) {
	queries := []string{
		"SDDE-201",
	}
	for _, query := range queries {
		metach := Search(query, nil)
		meta, ok := <-metach
		if !ok {
			t.Errorf("%s not found", query)
		} else {
			t.Logf("%s -> %+v", query, meta)
		}
	}
}
