package opendmm

import (
	"testing"
)

func TestLibre(t *testing.T) {
	queries := []string{
		"KV-176",
	}
	assertSearchable(t, queries, libreSearch)
}
