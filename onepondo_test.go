package opendmm

import (
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

func TestOnepondo(t *testing.T) {
	httpCache, err := leveldb.OpenFile("/tmp/opendmm.http.cache", nil)
	if err != nil {
		t.Fatal(err)
	}
	movieCache, err := leveldb.OpenFile("/tmp/opendmm.movie.cache", nil)
	if err != nil {
		t.Fatal(err)
	}

	metach := make(chan MovieMeta)
	wg := opdCrawl(httpCache, metach)
	go func() {
		wg.Wait()
		close(metach)
	}()
	for _ = range cacheIntoDB(movieCache)(postprocess(metach)) {
	}

	queries := []string{
		"1pondo 050704_429",
	}
	assertSearchable(t, queries, opdSearch(movieCache))
}
