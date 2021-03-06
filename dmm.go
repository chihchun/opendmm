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
	mapset "github.com/deckarep/golang-set"
	"github.com/golang/glog"
	"github.com/junzh0u/httpx"
)

func dmmSearch(query string, wg *sync.WaitGroup, metach chan MovieMeta) {
	keywords := dmmGuess(query)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			dmmSearchKeyword(keyword, wg, metach)
		}(keyword.(string))
	}
}

func dmmRe() *regexp.Regexp {
	return regexp.MustCompile("(?i)((?:t28|(?:3d|2d|s2|[a-z]){1,7}?))[-_]?(0*(\\d{2,5}))")
}

func dmmGuess(query string) mapset.Set {
	matches := dmmRe().FindAllStringSubmatch(query, -1)
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

func dmmIsCodeEqual(lcode, rcode string) bool {
	lmeta := dmmRe().FindStringSubmatch(lcode)
	rmeta := dmmRe().FindStringSubmatch(rcode)
	if lmeta == nil || rmeta == nil {
		return false
	}
	if lmeta[1] != rmeta[1] {
		return false
	}
	lnum, err := strconv.Atoi(lmeta[2])
	if err != nil {
		glog.Error(err)
		return false
	}
	rnum, err := strconv.Atoi(rmeta[2])
	if err != nil {
		glog.Error(err)
		return false
	}
	return lnum == rnum
}

func dmmSearchKeyword(keyword string, wg *sync.WaitGroup, metach chan MovieMeta) {
	glog.Info("Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"http://www.dmm.co.jp/search/=/searchstr=%s",
		url.QueryEscape(
			regexp.MustCompile("(?i)[a-z].*").FindString(keyword),
		),
	)
	glog.V(2).Info("Search page: ", urlstr)
	doc, err := newDocument(urlstr, httpx.GetContentInUTF8(http.Get))
	if err != nil {
		glog.V(2).Infof("Error parsing %s: %v", urlstr, err)
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
	glog.V(2).Info("Product page: ", urlstr)
	doc, err := newDocument(urlstr, httpx.GetContentInUTF8(http.Get))
	if err != nil {
		glog.V(2).Infof("Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	var ok bool
	meta.Page = urlstr
	meta.Title = doc.Find(".area-headline h1").Text()
	meta.ThumbnailImage, _ = doc.Find("#sample-video img").Attr("src")
	meta.CoverImage, ok = doc.Find("#sample-video a").Attr("href")
	if !ok || strings.HasPrefix(meta.CoverImage, "javascript") {
		meta.CoverImage = meta.ThumbnailImage
	}
	doc.Find("div.page-detail > table > tbody > tr > td > table > tbody > tr").Each(
		func(i int, tr *goquery.Selection) {
			td := tr.Find("td").First()
			k := td.Text()
			v := td.Next()
			if strings.Contains(k, "開始日") || strings.Contains(k, "発売日") {
				date := strings.TrimSpace(v.Text())
				matched, _ := regexp.MatchString("^-+$", date)
				if !matched {
					meta.ReleaseDate = date
				}
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
		glog.V(2).Infof("Code mismatch: Expected %s, got %s", keyword, meta.Code)
	} else {
		metach <- meta
	}
}

func dmmParseCode(code string) string {
	re := regexp.MustCompile("(?i)((?:3d|2d|[a-z])+)(\\d+)")
	m := re.FindStringSubmatch(code)
	if m != nil {
		return fmt.Sprintf("%s-%s", strings.ToUpper(m[1]), m[2])
	}
	return code
}
