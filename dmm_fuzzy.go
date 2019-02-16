package opendmm

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	mapset "github.com/deckarep/golang-set"
)

func dmmFuzzySearch(query string, wg *sync.WaitGroup, metach chan MovieMeta) {
	keywords := dmmFuzzyGuess(query)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			dmmSearchKeyword(keyword, wg, metach)
		}(keyword.(string))
	}
}

func dmmFuzzyGuess(query string) mapset.Set {
	re := regexp.MustCompile("(?i)([a-z][a-z0-9]{0,6}?)[-_]?(0*(\\d{2,5}))")
	matches := re.FindAllStringSubmatch(query, -1)
	keywords := mapset.NewSet()
	for _, match := range matches {
		series := strings.ToUpper(match[1])
		num := match[2]
		keywords.Add(fmt.Sprintf("%s-%03s", series, num))
		keywords.Add(fmt.Sprintf("%s-%04s", series, num))
		keywords.Add(fmt.Sprintf("%s-%05s", series, num))
	}
	return keywords
}
