package opendmm

import (
	"testing"
)

func TestOnepondo(t *testing.T) {
	queries := []string{
	//"1pondo 120115_199",
	}
	assertSearchableWithDB(t, queries, opdSearch)
}
