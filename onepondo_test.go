package opendmm

import (
	"testing"
)

func TestOpd(t *testing.T) {
	queries := []string{"1pondo 022018_648"}
	assertSearchable(t, queries, opdSearch)
}
