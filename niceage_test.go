package opendmm

import (
	"testing"
)

func TestNiceage(t *testing.T) {
	queries := []string{
		"NMNS-002",
		"NMNS-002B",
	}
	assertSearchable(t, queries, niceageSearch)
}
