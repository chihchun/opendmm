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

func tkhSearch(query string, metach chan MovieMeta) *sync.WaitGroup {
	glog.Info("[TKH] Query: ", query)

	wg := new(sync.WaitGroup)
	keywords := tkhGuess(query)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			tkhSearchKeyword(keyword, metach)
		}(keyword.(string))
	}
	return wg
}

func tkhGuess(query string) mapset.Set {
	re := regexp.MustCompile("(?i)(tokyo.*hot.*|^)(k|n)(\\d{3,4})")
	matches := re.FindAllStringSubmatch(query, -1)
	keywords := mapset.NewSet()
	for _, match := range matches {
		keywords.Add(fmt.Sprintf("%s%04s", strings.ToLower(match[2]), match[3]))
	}
	return keywords
}

func tkhGuessFull(query string) mapset.Set {
	keywords := mapset.NewSet()
	for keyword := range tkhGuess(query).Iter() {
		keywords.Add(fmt.Sprintf("Tokyo Hot %s", keyword))
	}
	return keywords
}

func tkhSearchKeyword(keyword string, metach chan MovieMeta) {
	glog.Info("[TKH] Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"http://www.tokyo-hot.com/product/?q=%s",
		url.QueryEscape(keyword),
	)
	glog.Info("[TKH] Search page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, http.Get)
	if err != nil {
		glog.Warningf("[TKH] Error parsing %s: %v", urlstr, err)
		return
	}

	href, ok := doc.Find("#main > ul > li > a").First().Attr("href")
	if ok {
		urlbase, err := url.Parse(urlstr)
		if err != nil {
			return
		}
		urlhref, err := urlbase.Parse(href)
		if err != nil {
			return
		}
		tkhParse(urlhref.String(), metach)
	}
}

func tkhParse(urlstr string, metach chan MovieMeta) {
	glog.Info("[TKH] Product page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, http.Get)
	if err != nil {
		glog.Warningf("[TKH] Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	meta.Page = urlstr
	meta.Title = doc.Find("#container > div.pagetitle > h2").Text()
	meta.CoverImage, _ = doc.Find("#container > div.movie.cf > div.in > div.flowplayer > video").Attr("poster")
	meta.SampleImages = doc.Find("#main > div.contents > div.scap > a, #main > div.contents > div.vcap > a").Map(
		func(i int, a *goquery.Selection) string {
			href, _ := a.Attr("href")
			return href
		})
	doc.Find("#main > div.contents > div.infowrapper > dl > dt").Each(
		func(i int, dt *goquery.Selection) {
			k := dt.Text()
			dd := dt.Next()
			if strings.Contains(k, "出演者") {
				meta.Actresses = dd.Find("a").Map(
					func(i int, a *goquery.Selection) string {
						return a.Text()
					})
			} else if strings.Contains(k, "シリーズ") {
				meta.MovieLength = dd.Text()
			} else if strings.Contains(k, "カテゴリ") {
				meta.Categories = dd.Find("a").Map(
					func(i int, a *goquery.Selection) string {
						return a.Text()
					})
			} else if strings.Contains(k, "配信開始日") {
				meta.ReleaseDate = dd.Text()
			} else if strings.Contains(k, "収録時間") {
				meta.MovieLength = dd.Text()
			} else if strings.Contains(k, "作品番号") {
				meta.Code = fmt.Sprintf("Tokyo Hot %s", dd.Text())
			}
		})
	metach <- meta
}
