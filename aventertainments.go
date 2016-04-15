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

func aveSearch(query string, metach chan MovieMeta) *sync.WaitGroup {
	glog.Info("[AVE] Query: ", query)
	keywords := aveGuess(query)
	wg := new(sync.WaitGroup)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			aveSearchKeyword(keyword, wg, metach)
		}(keyword.(string))
	}
	return wg
}

func aveGuess(query string) mapset.Set {
	re := regexp.MustCompile("(?i)([a-z2-3]{2,8})[-_]?([sm]?)(\\d{2,5})")
	matches := re.FindAllStringSubmatch(query, -1)
	keywords := mapset.NewSet()
	for _, match := range matches {
		series := strings.ToUpper(match[1])
		prefix := strings.ToUpper(match[2])
		keywords.Add(fmt.Sprintf("%s-%s%s", series, prefix, match[3]))
		number, _ := strconv.Atoi(match[3])
		keywords.Add(fmt.Sprintf("%s-%s%d", series, prefix, number))
	}
	return keywords
}

func aveSearchKeyword(keyword string, wg *sync.WaitGroup, metach chan MovieMeta) {
	glog.Info("[AVE] Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"http://www.aventertainments.com/search_Products.aspx?keyword=%s",
		url.QueryEscape(keyword),
	)
	glog.Info("[AVE] Search page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, http.Get)
	if err != nil {
		glog.Warningf("[AVE] Error parsing %s: %v", urlstr, err)
		return
	}

	doc.Find("div.main-unit2 > table a").Each(
		func(i int, a *goquery.Selection) {
			href, ok := a.Attr("href")
			if ok {
				wg.Add(1)
				go func() {
					defer wg.Done()
					aveParse(href, keyword, metach)
				}()
			}
		})
}

func aveParse(urlstr string, keyword string, metach chan MovieMeta) {
	glog.Info("[AVE] Product page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, http.Get)
	if err != nil {
		glog.Warningf("[AVE] Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	var ok bool
	meta.Page = urlstr
	meta.Title = doc.Find("#mini-tabet > h2").Text()
	meta.CoverImage, ok = doc.Find("#titlebox > div.list-cover > img").Attr("src")
	if ok {
		meta.CoverImage = strings.Replace(meta.CoverImage, "jacket_images", "bigcover", -1)
	}
	meta.Code = strings.Replace(doc.Find("#mini-tabet > div").Text(), "商品番号:", "", -1)
	doc.Find("#titlebox > ul > li").Each(
		func(i int, li *goquery.Selection) {
			k := li.Find("span").Text()
			if strings.Contains(k, "主演女優") {
				meta.Actresses = li.Find("a").Map(
					func(i int, a *goquery.Selection) string {
						return a.Text()
					})
			} else if strings.Contains(k, "スタジオ") {
				meta.Maker = li.Find("a").Text()
			} else if strings.Contains(k, "シリーズ") {
				meta.Series = li.Find("a").Text()
			} else if strings.Contains(k, "発売日") {
				meta.ReleaseDate = li.Text()
			} else if strings.Contains(k, "収録時間") {
				meta.MovieLength = li.Text()
			}
		})

	if strings.TrimSpace(meta.Code) != keyword {
		glog.Warningf("[AVE] Code mismatch: Expected %s, got %s", keyword, meta.Code)
	} else {
		metach <- meta
	}
}
