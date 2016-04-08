package opendmm

import (
	"testing"
)

func TestMgstage(t *testing.T) {
	queries := []string{
		"SIRO-1715",
		"200GANA-894",
	}
	assertSearchable(t, queries, mgsSearch)
}
