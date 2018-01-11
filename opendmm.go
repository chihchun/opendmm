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

// SearchFunc is the interface of each engine's search function
type SearchFunc func(string, chan MovieMeta) *sync.WaitGroup

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
	wgs = append(wgs, niceageSearch(query, metach))
	wgs = append(wgs, tkhSearch(query, metach))
	wgs = append(wgs, scuteSearch(query, metach))

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
	keywords = keywords.Union(tkhGuessFull(query))
	keywords = keywords.Union(scuteGuessFull(query))
	return keywords
}
