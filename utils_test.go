package opendmm

import (
	"testing"
)

func assertSearchable(t *testing.T, queries []string, search SearchFunc) {
	for _, query := range queries {
		metach := make(chan MovieMeta)
		wg := search(query, metach)
		go func() {
			wg.Wait()
			close(metach)
		}()
		meta, ok := <-postprocess(metach)
		if !ok {
			t.Errorf("%s not found", query)
		} else {
			t.Logf("%s -> %+v", query, meta)
		}
	}
}

func assertUnsearchable(t *testing.T, queries []string, search SearchFunc) {
	for _, query := range queries {
		metach := make(chan MovieMeta)
		wg := search(query, metach)
		go func() {
			wg.Wait()
			close(metach)
		}()
		meta, ok := <-postprocess(metach)
		if ok {
			t.Errorf("Unexpected: %s -> %+v", query, meta)
		}
	}
}
