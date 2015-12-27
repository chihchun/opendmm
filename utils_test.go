package opendmm

import (
  "testing"
  "sync"
)

func testEngine(t *testing.T, queries []string, search func(string, chan MovieMeta, *sync.WaitGroup)) {
  for _, query := range queries {
    var wg sync.WaitGroup
    metach := make(chan MovieMeta)
    search(query, metach, &wg)
    metach = validateFields(trimSpaces(deduplicate(metach)))
    meta, ok := <-metach
    if !ok {
      t.Error("Not found")
    } else {
      t.Logf("%+v", meta)
    }
  }
}
