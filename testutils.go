package opendmm

import (
  "testing"
  "sync"
)

func assertSearchable(t *testing.T, queries []string, search func(string, chan MovieMeta) *sync.WaitGroup) {
  for _, query := range queries {
    metach := make(chan MovieMeta)
    wg := search(query, metach)
    go func() {
      wg.Wait()
      close(metach)
    }()
    meta, ok := <-validateFields(trimSpaces(deduplicate(metach)))
    if !ok {
      t.Errorf("%s not found", query)
    } else {
      t.Logf("%s -> %+v", query, meta)
    }
  }
}

func assertUnsearchable(t *testing.T, queries []string, search func(string, chan MovieMeta) *sync.WaitGroup) {
  for _, query := range queries {
    metach := make(chan MovieMeta)
    wg := search(query, metach)
    go func() {
      wg.Wait()
      close(metach)
    }()
    meta, ok := <-validateFields(trimSpaces(deduplicate(metach)))
    if ok {
      t.Error("Unexpected: %s -> %+v", query, meta)
    }
  }
}
