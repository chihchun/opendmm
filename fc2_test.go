package opendmm

import (
	"testing"
)

func TestFc2(t *testing.T) {
	queries := []string{
		"FC2 917114",
		"FC2-PPV 917114",
		"FC2-PPV 1055613",
	}
	assertSearchable(t, queries, fc2Search)

	blackhole := []string {
		"FC2 749615",
	}
	assertUnsearchable(t, blackhole, fc2Search)

}
