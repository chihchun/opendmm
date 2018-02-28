package opendmm

import (
	"testing"
)

func assertSearchable(t *testing.T, queries []string, search searchFunc) {
	for _, query := range queries {
		metach := make(chan MovieMeta)
		go func() {
			search(query, metach)
			close(metach)
		}()
		meta, ok := <-postprocess(metach)
		if !ok {
			t.Errorf("%s not found", query)
		} else {
			t.Logf("%s -> %s %s", query, meta.Code, meta.Title)
		}
	}
}

func assertUnsearchable(t *testing.T, queries []string, search searchFunc) {
	for _, query := range queries {
		metach := make(chan MovieMeta)
		go func() {
			search(query, metach)
			close(metach)
		}()
		meta, ok := <-postprocess(metach)
		if ok {
			t.Errorf("Unexpected: %s -> %+v", query, meta)
		}
	}
}
