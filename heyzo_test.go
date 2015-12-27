package opendmm

import (
  "testing"
)

func TestHeyzo(t *testing.T) {
  queries := []string {"Heyzo 1021", "Heyzo 001"}
  testEngine(t, queries, heyzoSearch)
}
