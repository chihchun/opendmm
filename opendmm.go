package opendmm

import (
	"sync"

	"github.com/deckarep/golang-set"
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

var workerPoolSize = 1000
var workerPool = make(chan int, workerPoolSize)

func init() {
	for i := 1; i <= workerPoolSize; i++ {
		workerPool <- 1
	}
}

// SearchFunc is the interface of each engine's search function
type SearchFunc func(string, chan MovieMeta) *sync.WaitGroup

// Search for movies based on query and return a channel of MovieMeta
func Search(query string) chan MovieMeta {
	metach := make(chan MovieMeta)
	var wgs [](*sync.WaitGroup)
	for _, handler := range []SearchFunc{
		aveSearch,
		caribSearch,
		caribprSearch,
		dmmSearch,
		fc2Search,
		heyzoSearch,
		javSearch,
		mgsSearch,
		niceageSearch,
		scuteSearch,
		tkhSearch,
	} {
		wgs = append(wgs, handler(query, metach))
	}

	go func() {
		for _, wg := range wgs {
			wg.Wait()
		}
		close(metach)
	}()
	return postprocess(metach)
}

// GuessFunc is the interface of each engine's guess function
type GuessFunc func(string) mapset.Set

// Guess possible movie codes from query string
func Guess(query string) mapset.Set {
	keywords := mapset.NewSet()
	keywords = keywords.Union(aveGuess(query))
	keywords = keywords.Union(caribGuessFull(query))
	keywords = keywords.Union(caribprGuessFull(query))
	keywords = keywords.Union(dmmGuess(query))
	keywords = keywords.Union(heyzoGuessFull(query))
	keywords = keywords.Union(tkhGuessFull(query))
	keywords = keywords.Union(scuteGuessFull(query))
	keywords = keywords.Union(fc2GuessFull(query))
	return keywords
}
