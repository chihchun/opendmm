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

type searchFunc func(string, chan MovieMeta)

type searchRequest struct {
	query string
	out   chan MovieMeta
	wg    *sync.WaitGroup
}

var reqs chan searchRequest

func init() {
	reqs = make(chan searchRequest, 10)
	var pipes [](chan searchRequest)

	// Fast minions
	for _, minion := range []searchFunc{
		aveSearch,
		caribSearch,
		caribprSearch,
		dmmSearch,
		fc2Search,
		heyzoSearch,
		javSearch,
		niceageSearch,
		tkhSearch,
	} {
		pipe := make(chan searchRequest, 10)
		pipes = append(pipes, pipe)
		for i := 0; i < 2; i++ {
			go func(minion searchFunc) {
				for req := range pipe {
					minion(req.query, req.out)
					req.wg.Done()
				}
			}(minion)
		}
	}

	// slow minions
	for _, minion := range []searchFunc{
		mgsSearch,
		opdSearch,
		scuteSearch,
	} {
		pipe := make(chan searchRequest, 100)
		pipes = append(pipes, pipe)
		for i := 0; i < 5; i++ {
			go func(minion searchFunc) {
				for req := range pipe {
					minion(req.query, req.out)
					req.wg.Done()
				}
			}(minion)
		}
	}

	// dispatcher
	go func() {
		for req := range reqs {
			for _, pipe := range pipes {
				req.wg.Add(1)
				pipe <- req
			}
			go func(req searchRequest) {
				req.wg.Wait()
				close(req.out)
			}(req)
		}
	}()
}

// Search for movies based on query and return a channel of MovieMeta
func Search(query string) chan MovieMeta {
	out := make(chan MovieMeta, 100)
	req := searchRequest{query, out, new(sync.WaitGroup)}
	reqs <- req
	return postprocess(out)
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
	keywords = keywords.Union(scuteGuessFull(query))
	keywords = keywords.Union(fc2GuessFull(query))
	return keywords
}
