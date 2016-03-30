package opendmm

import (
	"testing"
)

func TestMgstage(t *testing.T) {
	queries := []string{
		"SIRO-1715",
	}
	assertSearchable(t, queries, mgsSearch)
}
