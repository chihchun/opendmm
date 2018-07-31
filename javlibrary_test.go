package opendmm

import (
	"testing"
)

func TestJavlibrary(t *testing.T) {
	queries := []string{
		"MIAD-617",
	}
	assertSearchable(t, queries, javSearch)
}
