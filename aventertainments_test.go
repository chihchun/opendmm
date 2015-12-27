package opendmm

import (
  "testing"
)

func TestAventertainments(t *testing.T) {
  queries := []string {"SKYHD-001"}
  testEngine(t, queries, aveSearch)
}
