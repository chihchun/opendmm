package opendmm

import (
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
}

var reqCh chan searchRequest

func init() {
	reqCh = make(chan searchRequest, 100)
	var minionChs [](chan searchRequest)
	for _, minion := range []searchFunc{
		aveSearch,
		caribSearch,
		caribprSearch,
		dmmSearch,
		fc2Search,
		heyzoSearch,
		javSearch,
		mgsSearch,
		niceageSearch,
		opdSearch,
		scuteSearch,
		tkhSearch,
	} {
		minionCh := make(chan searchRequest, 100)
		minionChs = append(minionChs, minionCh)
		go func(minion searchFunc) {
			for req := range minionCh {
				minion(req.query, req.out)
			}
		}(minion)
	}
	go func() {
		for req := range reqCh {
			for _, minionCh := range minionChs {
				minionCh <- req
			}
		}
	}()
}

// Search for movies based on query and return a channel of MovieMeta
func Search(query string) chan MovieMeta {
	out := make(chan MovieMeta, 100)
	req := searchRequest{query, out}
	reqCh <- req
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
