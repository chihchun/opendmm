package opendmm

import (
	"testing"
)

func TestCarib(t *testing.T) {
	queries := []string{"Carib 081215-945"}
	assertSearchable(t, queries, caribSearch)
}
