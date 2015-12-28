package opendmm

import (
  "testing"
)

func TestDmm(t *testing.T) {
  queries := []string {
    "MIDE-029",
    "mide-029",
    "XV-100",
    "XV-1001",
  }
  assertSearchable(t, queries, dmmSearch)
  blackhole := []string {"MCB3DBD-25"}
  assertUnsearchable(t, blackhole, dmmSearch)
}
