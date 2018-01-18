package opendmm

import (
	"testing"
)

func TestFc2(t *testing.T) {
	queries := []string{
		"FC2 749615",
		"FC2-PPV 749615",
	}
	assertSearchable(t, queries, fc2Search)
}
