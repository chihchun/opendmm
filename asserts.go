package opendmm

import (
	"sync"
	"testing"
)

func assertSearchable(t *testing.T, queries []string, search searchFunc) {
	twg := new(sync.WaitGroup)
	for _, query := range queries {
		twg.Add(1)
		go func(query string) {
			defer twg.Done()
			out := make(chan MovieMeta, 100)
			wg := new(sync.WaitGroup)
			search(query, wg, out)
			wg.Wait()
			close(out)
			meta, ok := <-postprocess(out)
			if !ok {
				t.Errorf("%s not found", query)
			} else {
				t.Logf("%s -> %s %s", query, meta.Code, meta.Title)
			}
		}(query)
	}
	twg.Wait()
}

func assertUnsearchable(t *testing.T, queries []string, search searchFunc) {
	twg := new(sync.WaitGroup)
	for _, query := range queries {
		twg.Add(1)
		go func(query string) {
			defer twg.Done()
			out := make(chan MovieMeta, 100)
			wg := new(sync.WaitGroup)
			search(query, wg, out)
			wg.Wait()
			close(out)
			meta, ok := <-postprocess(out)
			if ok {
				t.Errorf("Unexpected: %s -> %+v", query, meta)
			}
		}(query)
	}
	twg.Wait()
}
