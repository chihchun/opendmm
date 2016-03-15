package opendmm

import (
	"sync"
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
func Search(query string, dbpath string) chan MovieMeta {
	metach := make(chan MovieMeta)

	var wgs [](*sync.WaitGroup)
	wgs = append(wgs, aveSearch(query, metach))
	wgs = append(wgs, caribSearch(query, metach))
	wgs = append(wgs, caribprSearch(query, metach))
	wgs = append(wgs, dmmSearch(query, metach))
	wgs = append(wgs, heyzoSearch(query, metach))
	wgs = append(wgs, javSearch(query, metach))
	wgs = append(wgs, tkhSearch(query, metach))
	wgs = append(wgs, opdSearch(query, dbpath, metach))
	go func() {
		for _, wg := range wgs {
			wg.Wait()
		}
		close(metach)
	}()
	return postprocess(metach)
}
