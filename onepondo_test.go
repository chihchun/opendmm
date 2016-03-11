package opendmm

import (
	"testing"
)

func TestOnepondo(t *testing.T) {
	needle := "1pondo 120115_199"

	metach := make(chan MovieMeta)
	wg := opdCrawl(metach)
	go func() {
		wg.Wait()
		close(metach)
	}()
	for meta := range postprocess(metach) {
		if meta.Code == needle {
			t.Logf("%s -> %+v", needle, meta)
			return
		}
	}
	t.Errorf("%s not found", needle)
}
