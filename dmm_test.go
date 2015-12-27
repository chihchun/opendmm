package opendmm

import (
  "testing"
)

func TestDmm(t *testing.T) {
  queries := []string {"MIDE-029"}
  testEngine(t, queries, caribprSearch)
}
