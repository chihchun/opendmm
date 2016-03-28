package opendmm

import (
	"sync"

	"github.com/golang/glog"
	"github.com/syndtr/goleveldb/leveldb"
)

// MovieMeta contains meta data of movie
type MovieMeta struct {
	Actresses      []string
	ActressTypes   []string
	Categories     []string
	Code           string
	CoverImage     string
	Description    string
	Directors      []string
	Genres         []string
	Label          string
	Maker          string
	MovieLength    string
	Page           string
	ReleaseDate    string
	SampleImages   []string
	Series         string
	Tags           []string
	ThumbnailImage string
	Title          string
}

// Search for movies based on query and return a channel of MovieMeta
func Search(query string) chan MovieMeta {
	metach := make(chan MovieMeta)
	var wgs [](*sync.WaitGroup)
	wgs = append(wgs, aveSearch(query, metach))
	wgs = append(wgs, caribSearch(query, metach))
	wgs = append(wgs, caribprSearch(query, metach))
	wgs = append(wgs, dmmSearch(query, metach))
	wgs = append(wgs, heyzoSearch(query, metach))
	wgs = append(wgs, javSearch(query, metach))
	wgs = append(wgs, mgsSearch(query, metach))
	wgs = append(wgs, tkhSearch(query, metach))

	httpCache, err := leveldb.OpenFile("/tmp/opendmm.http.cache", nil)
	if err != nil {
		glog.Fatal(err)
	}
	movieCache, err := leveldb.OpenFile("/tmp/opendmm.movie.cache", nil)
	if err != nil {
		glog.Fatal(err)
	}
	wgs = append(wgs, opdSearch(query, httpCache, movieCache, metach))

	go func() {
		for _, wg := range wgs {
			wg.Wait()
		}
		close(metach)
		httpCache.Close()
		movieCache.Close()
	}()
	return postprocess(metach)
}
