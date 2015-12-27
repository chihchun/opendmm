package opendmm

import (
  "testing"
)

func TestJavlibrary(t *testing.T) {
  queries := []string {"MIDE-029"}
  testEngine(t, queries, javSearch)
}
