package opendmm

import (
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

func TestOnepondo(t *testing.T) {
	url2code := map[string]string{
		"http://m.1pondo.tv/movies/1/": "1pondo 050704_429",
	}

	httpCache, err := leveldb.OpenFile("/tmp/opendmm.http.cache", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer httpCache.Close()
	for url, code := range url2code {
		meta, err := opdParse(url, httpCache)
		if err != nil {
			t.Fatal(err)
		}
		if meta.Code != code {
			t.Fatalf("Expecting code %s, get %s instead", code, meta.Code)
		}
	}
}
