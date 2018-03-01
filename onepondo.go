package opendmm

import (
	"fmt"
	"net/url"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/deckarep/golang-set"
	"github.com/golang/glog"
	"github.com/junzh0u/httpx"
)

func opdSearch(query string, wg *sync.WaitGroup, metach chan MovieMeta) {
	keywords := opdGuess(query)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			opdSearchKeyword(keyword, metach)
		}(keyword.(string))
	}
}

func opdGuess(query string) mapset.Set {
	keywords := mapset.NewSet()
	matched, _ := regexp.MatchString("(?i)(1|one)pon(do)?", query)
	if !matched {
		return keywords
	}

	re := regexp.MustCompile("(\\d{6})[-_](\\d{3})")
	matches := re.FindAllStringSubmatch(query, -1)
	for _, match := range matches {
		keywords.Add(fmt.Sprintf("%s_%s", match[1], match[2]))
	}
	return keywords
}

func opdGuessFull(query string) mapset.Set {
	keywords := mapset.NewSet()
	for keyword := range opdGuess(query).Iter() {
		keywords.Add(fmt.Sprintf("1pondo %s", keyword))
	}
	return keywords
}

func opdSearchKeyword(keyword string, metach chan MovieMeta) {
	glog.Info("Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"https://www.1pondo.tv/movies/%s/",
		url.QueryEscape(keyword),
	)
	opdParse(urlstr, keyword, metach)
}

func opdParse(urlstr string, keyword string, metach chan MovieMeta) {
	glog.V(2).Info("Product page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, httpx.GetFullPage)
	if err != nil {
		glog.V(2).Infof("Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	meta.Code = fmt.Sprintf("1pondo %s", keyword)
	meta.Page = urlstr
	meta.CoverImage = fmt.Sprintf("https://www.1pondo.tv/assets/sample/%s/str.jpg", keyword)
	meta.Description = doc.Find(".box-comment div:nth-child(1)").Text()
	doc.Find(".tag-area a").Each(
		func(i int, a *goquery.Selection) {
			meta.Tags = append(meta.Tags, a.Text())
		})
	meta.Title = doc.Find("dl.movie-title dd").Text()
	doc.Find("dl.actress-name dd a").Each(
		func(i int, a *goquery.Selection) {
			meta.Actresses = append(meta.Actresses, a.Text())
		})
	meta.ReleaseDate = doc.Find("dl.release-date dd").Text()
	meta.MovieLength = doc.Find("dl.duration dd").Text()
	meta.Series = doc.Find("dl.series dd").Text()

	metach <- meta
}
