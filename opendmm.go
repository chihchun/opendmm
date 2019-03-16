package opendmm

import (
	"sync"

	mapset "github.com/deckarep/golang-set"
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

type searchFunc func(string, *sync.WaitGroup, chan MovieMeta)

// Search for movies based on query and return a channel of MovieMeta
func Search(query string) chan MovieMeta {
	out := make(chan MovieMeta)
	go func(out chan MovieMeta) {
		defer close(out)
		batches := [][]searchFunc{[]searchFunc{
			aveSearch,
			dmmSearch,
		}, []searchFunc{
			javSearch,
			mgsSearch,
		}}
		for _, engines := range batches {
			batchOut := searchWithEngines(query, engines)
			meta, ok := <-batchOut
			if ok {
				out <- meta
				break
			}
		}
	}(out)
	return out
}

func searchWithEngines(query string, engines []searchFunc) chan MovieMeta {
	wg := new(sync.WaitGroup)
	out := make(chan MovieMeta, 100)
	for _, engine := range engines {
		engine(query, wg, out)
	}
	go func(wg *sync.WaitGroup, out chan MovieMeta) {
		wg.Wait()
		close(out)
	}(wg, out)
	return postprocess(out)
}

// Guess possible movie codes from query string
func Guess(query string) mapset.Set {
	keywords := mapset.NewSet()
	keywords = keywords.Union(aveGuess(query))
	keywords = keywords.Union(dmmGuess(query))
	keywords = keywords.Union(mgsGuess(query))
	return keywords
}
