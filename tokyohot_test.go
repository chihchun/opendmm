package opendmm

import (
	"testing"
)

func TestTokyohot(t *testing.T) {
	queries := []string{
		"Tokyo Hot n110",
		"tokyo-hot n110",
		"Tokyo Hot n1108",
	}
	assertSearchable(t, queries, tkhSearch)
}
