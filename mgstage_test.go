package opendmm

import (
	"testing"
)

func TestMgs(t *testing.T) {
	queries := []string{
		"SIRO-1715",
	}
	assertSearchable(t, queries, mgsSearch)
}
