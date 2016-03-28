package opendmm

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
)

func dmmParseCode(code string) string {
	re := regexp.MustCompile("(?i)([a-z]+)(\\d+)")
	meta := re.FindStringSubmatch(code)
	if meta != nil {
		return fmt.Sprintf("%s-%s", strings.ToUpper(meta[1]), meta[2])
	}
	return code
}

func dmmParse(urlstr string, keyword string, metach chan MovieMeta) {
	glog.Info("[DMM] Product page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, http.Get)
	if err != nil {
		glog.Warningf("[DMM] Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	var ok bool
	meta.Page = urlstr
	meta.Title = doc.Find("#title").Text()
	meta.ThumbnailImage, _ = doc.Find("#sample-video img").Attr("src")
	meta.CoverImage, ok = doc.Find("#sample-video a").Attr("href")
	if !ok {
		meta.CoverImage = meta.ThumbnailImage
	}
	doc.Find("div.page-detail > table > tbody > tr > td > table > tbody > tr").Each(
		func(i int, tr *goquery.Selection) {
			td := tr.Find("td").First()
			k := td.Text()
			v := td.Next()
			if strings.Contains(k, "配信開始日") {
				meta.ReleaseDate = v.Text()
			} else if strings.Contains(k, "収録時間") {
				meta.MovieLength = v.Text()
			} else if strings.Contains(k, "出演者") {
				meta.Actresses = v.Find("a").Map(
					func(i int, a *goquery.Selection) string {
						return a.Text()
					})
			} else if strings.Contains(k, "監督") {
				meta.Directors = v.Find("a").Map(
					func(i int, a *goquery.Selection) string {
						return a.Text()
					})
			} else if strings.Contains(k, "シリーズ") {
				meta.Series = v.Text()
			} else if strings.Contains(k, "メーカー") {
				meta.Maker = v.Text()
			} else if strings.Contains(k, "レーベル") {
				meta.Label = v.Text()
			} else if strings.Contains(k, "ジャンル") {
				meta.Genres = v.Find("a").Map(
					func(i int, a *goquery.Selection) string {
						return a.Text()
					})
			} else if strings.Contains(k, "品番") {
				meta.Code = dmmParseCode(v.Text())
			}
		})

	if strings.TrimSpace(meta.Code) != keyword {
		glog.Warningf("[DMM] Code mismatch: Expected %s, got %s", keyword, meta.Code)
	} else {
		metach <- meta
	}
}

func dmmSearchKeyword(keyword string, wg *sync.WaitGroup, metach chan MovieMeta) {
	glog.Info("[DMM] Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"http://www.dmm.co.jp/search/=/searchstr=%s",
		url.QueryEscape(keyword),
	)
	glog.Info("[DMM] Search page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, http.Get)
	if err != nil {
		glog.Warningf("[DMM] Error parsing %s: %v", urlstr, err)
		return
	}

	doc.Find("#list > li > div > p.tmb > a").Each(
		func(i int, a *goquery.Selection) {
			href, ok := a.Attr("href")
			if ok {
				wg.Add(1)
				go func() {
					defer wg.Done()
					dmmParse(href, keyword, metach)
				}()
			}
		})
}

func dmmSearch(query string, metach chan MovieMeta) *sync.WaitGroup {
	glog.Info("[DMM] Query: ", query)
	wg := new(sync.WaitGroup)
	re := regexp.MustCompile("(?i)(\\w{2,6}?)-?(\\d{2,5})")
	matches := re.FindAllStringSubmatch(query, -1)
	for _, match := range matches {
		keyword := fmt.Sprintf("%s-%s", strings.ToUpper(match[1]), match[2])
		wg.Add(1)
		go func() {
			defer wg.Done()
			dmmSearchKeyword(keyword, wg, metach)
		}()
	}
	return wg
}
