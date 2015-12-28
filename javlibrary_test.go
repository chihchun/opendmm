package opendmm

import (
  "testing"
)

func TestJavlibrary(t *testing.T) {
  queries := []string {
    "MIDE-029",
    "mide-029",
    "SDDE-001",
    "XV-100",
    "XV-1001",
  }
  assertSearchable(t, queries, javSearch)
}
