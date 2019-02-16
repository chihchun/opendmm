package opendmm

import (
	"testing"
)

func TestDmmFuzzy(t *testing.T) {
	queries := []string{
		"3DSVR-106",
	}
	assertSearchable(t, queries, dmmFuzzySearch)
}
