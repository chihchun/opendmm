package opendmm

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/deckarep/golang-set"
	"github.com/golang/glog"
)

func dmmSearch(query string, metach chan MovieMeta) *sync.WaitGroup {
	glog.Info("[DMM] Query: ", query)
	keywords := dmmGuess(query)
	wg := new(sync.WaitGroup)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			dmmSearchKeyword(keyword, wg, metach)
		}(keyword.(string))
	}
	return wg
}

func dmmGuess(query string) mapset.Set {
	re := regexp.MustCompile("(?i)([a-z][a-z0-9]{1,5}?)[-_]?(\\d{2,5})")
	matches := re.FindAllStringSubmatch(query, -1)
	keywords := mapset.NewSet()
	for _, match := range matches {
		keywords.Add(fmt.Sprintf("%s-%s", strings.ToUpper(match[1]), match[2]))
	}
	return keywords
}

func dmmIsCodeEqual(lcode, rcode string) bool {
	re := regexp.MustCompile("(?i)([a-z]+)-(\\d+)")
	lmeta := re.FindStringSubmatch(lcode)
	rmeta := re.FindStringSubmatch(rcode)
	glog.Info(lmeta)
	glog.Info(rmeta)
	if lmeta == nil || rmeta == nil {
		return false
	}
	if lmeta[1] != rmeta[1] {
		return false
	}
	lnum, err := strconv.Atoi(lmeta[2])
	if err != nil {
		glog.Errorf("[DMM] %s", err)
		return false
	}
	rnum, err := strconv.Atoi(rmeta[2])
	if err != nil {
		glog.Errorf("[DMM] %s", err)
		return false
	}
	return lnum == rnum
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

	if !dmmIsCodeEqual(keyword, meta.Code) {
		glog.Warningf("[DMM] Code mismatch: Expected %s, got %s", keyword, meta.Code)
	} else {
		metach <- meta
	}
}

func dmmParseCode(code string) string {
	re := regexp.MustCompile("(?i)([a-z]+)(\\d+)")
	meta := re.FindStringSubmatch(code)
	if meta != nil {
		return fmt.Sprintf("%s-%s", strings.ToUpper(meta[1]), meta[2])
	}
	return code
}
