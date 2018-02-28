package opendmm

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/deckarep/golang-set"
	"github.com/golang/glog"
)

func heyzoSearch(query string, metach chan MovieMeta) {
	glog.Info("[HEYZO] Query: ", query)
	keywords := heyzoGuess(query)
	wg := new(sync.WaitGroup)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			heyzoSearchKeyword(keyword, metach)
		}(keyword.(string))
	}
	wg.Wait()
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
	glog.Info("[HEYZO] Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"http://www.heyzo.com/moviepages/%s/index.html",
		url.QueryEscape(keyword),
	)
	heyzoParse(urlstr, keyword, metach)
}

func heyzoParse(urlstr string, keyword string, metach chan MovieMeta) {
	glog.Info("[HEYZO] Product page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, http.Get)
	if err != nil {
		glog.Warningf("[HEYZO] Error parsing %s: %v", urlstr, err)
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
	meta.ReleaseDate = doc.Find("#movie > div.info-bg.info-bgWide > div > span.release-day + *").Text()
	meta.Actresses = doc.Find("#movie > div.info-bg.info-bgWide > div > span.actor + *").Find("a").Map(
		func(i int, a *goquery.Selection) string {
			return a.Text()
		})
	meta.Label = strings.Replace(
		doc.Find("#movie > div.info-bg.info-bgWide > div > span.label + *").Text(), "-", "", -1)
	meta.ActressTypes = doc.Find("#movie > div.info-bg.info-bgWide > div > div.actor-type > span").Map(
		func(i int, span *goquery.Selection) string {
			return span.Text()
		})
	meta.Tags = doc.Find("#movie > div.info-bg.info-bgWide > div > div.tag_cloud > ul > li").Map(
		func(i int, li *goquery.Selection) string {
			return li.Text()
		})
	descNode := doc.Find("#movie > div.info-bg.info-bgWide > div > p > *").Nodes
	if len(descNode) > 0 {
		meta.Description = descNode[0].Data
	}

	metach <- meta
}
