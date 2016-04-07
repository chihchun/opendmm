package opendmm

import (
	"sync"

	"github.com/deckarep/golang-set"
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

// SearchFunc is the interface of each engine's search function
type SearchFunc func(string, chan MovieMeta) *sync.WaitGroup

// Search for movies based on query and return a channel of MovieMeta
// Takes an additional leveldb DB pointer to perform search on some special
// search engines, such as onepondo.
// The DB can be generated using the Crawl API.
// Passing a empty DB pointer will just skip those engines.
func Search(query string, cache *leveldb.DB) chan MovieMeta {
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
	if cache != nil {
		wgs = append(wgs, opdSearch(cache)(query, metach))
	}

	go func() {
		for _, wg := range wgs {
			wg.Wait()
		}
		close(metach)
	}()
	return postprocess(metach)
}

// Guess possible movie codes from query string
func Guess(query string) mapset.Set {
	keywords := mapset.NewSet()
	keywords = keywords.Union(aveGuess(query))
	keywords = keywords.Union(caribGuessFull(query))
	keywords = keywords.Union(caribprGuessFull(query))
	keywords = keywords.Union(dmmGuess(query))
	keywords = keywords.Union(heyzoGuessFull(query))
	keywords = keywords.Union(opdGuessFull(query))
	keywords = keywords.Union(tkhGuessFull(query))
	return keywords
}

// Crawl movies that aren't searchable directly and cache into DB
func Crawl(movieCachePath, httpCachePath string) {
	httpCache, err := leveldb.OpenFile(httpCachePath, nil)
	if err != nil {
		glog.Fatal(err)
	}
	defer httpCache.Close()

	metach := make(chan MovieMeta)
	var wgs [](*sync.WaitGroup)
	wgs = append(wgs, opdCrawl(httpCache, metach))
	go func() {
		for _, wg := range wgs {
			wg.Wait()
		}
		close(metach)
		httpCache.Close()
	}()

	movieCache, err := leveldb.OpenFile(movieCachePath, nil)
	if err != nil {
		glog.Fatal(err)
	}
	defer movieCache.Close()
	for _ = range cacheIntoDB(movieCache)(postprocess(metach)) {
	}
}
