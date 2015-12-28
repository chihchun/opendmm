package opendmm

import (
  "testing"
)

func TestCaribpr(t *testing.T) {
  queries := []string {"Caribpr 031513_530"}
  assertSearchable(t, queries, caribprSearch)
}
