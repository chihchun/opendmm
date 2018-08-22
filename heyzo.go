package opendmm

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/deckarep/golang-set"
	"github.com/golang/glog"
	"github.com/junzh0u/httpx"
)

func heyzoSearch(query string, wg *sync.WaitGroup, metach chan MovieMeta) {
	keywords := heyzoGuess(query)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			heyzoSearchKeyword(keyword, metach)
		}(keyword.(string))
	}
}

func heyzoGuess(query string) mapset.Set {
	keywords := mapset.NewSet()
	matched, _ := regexp.Match("(?i)heyzo", []byte(query))
	if !matched {
		return keywords
	}

	re := regexp.MustCompile("\\d{3,4}")
	matches := re.FindAllString(query, -1)
	for _, match := range matches {
		keywords.Add(fmt.Sprintf("%04s", match))
	}
	return keywords
}

func heyzoGuessFull(query string) mapset.Set {
	keywords := mapset.NewSet()
	for keyword := range heyzoGuess(query).Iter() {
		keywords.Add(fmt.Sprintf("Heyzo %s", keyword))
	}
	return keywords
}

func heyzoSearchKeyword(keyword string, metach chan MovieMeta) {
	glog.Info("Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"http://www.heyzo.com/moviepages/%s/index.html",
		url.QueryEscape(keyword),
	)
	heyzoParse(urlstr, keyword, metach)
}

func heyzoParse(urlstr string, keyword string, metach chan MovieMeta) {
	glog.V(2).Info("Product page: ", urlstr)
	doc, err := newDocument(urlstr, httpx.GetContentInUTF8(http.Get))
	if err != nil {
		glog.V(2).Infof("Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	meta.Maker = "Heyzo"
	meta.Code = fmt.Sprintf("Heyzo %s", keyword)
	meta.Page = urlstr

	var urlbase *url.URL
	urlbase, err = url.Parse(urlstr)
	if err != nil {
		return
	}
	var urlcover *url.URL
	urlcover, err = urlbase.Parse(
		fmt.Sprintf("/contents/3000/%s/images/player_thumbnail_450.jpg", keyword))
	if err == nil {
		meta.CoverImage = urlcover.String()
	}

	meta.Title = doc.Find("#movie > h1").Text()
	meta.ReleaseDate = doc.Find("tr.table-release-day > td:nth-child(2)").Text()
	meta.Actresses = doc.Find("tr.table-actor > td:nth-child(2) a").Map(
		func(i int, a *goquery.Selection) string {
			return a.Text()
		})
	meta.Series = doc.Find("tr.table-series > td:nth-child(2)").Text()
	meta.ActressTypes = doc.Find("tr.table-actor-type > td:nth-child(2) a").Map(
		func(i int, a *goquery.Selection) string {
			return a.Text()
		})
	meta.Tags = doc.Find("ul.tag-keyword-list li").Map(
		func(i int, li *goquery.Selection) string {
			return li.Text()
		})
	descNode := doc.Find("tr.table-memo p.memo").Nodes
	if len(descNode) > 0 {
		meta.Description = descNode[0].Data
	}

	metach <- meta
}
