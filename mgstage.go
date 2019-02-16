package opendmm

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	mapset "github.com/deckarep/golang-set"
	"github.com/golang/glog"
	"github.com/junzh0u/httpx"
)

var mgsCookie = http.Cookie{
	Name:     "adc",
	Value:    "1",
	Domain:   ".mgstage.com",
	Path:     "/",
	HttpOnly: false,
	Secure:   false,
	Expires:  time.Now().Add(1000 * time.Hour),
}

func mgsSearch(query string, wg *sync.WaitGroup, metach chan MovieMeta) {
	keywords := mgsGuess(query)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			mgsSearchKeyword(keyword, wg, metach)
		}(keyword.(string))
	}
}

func mgsGuess(query string) mapset.Set {
	re := regexp.MustCompile("(?i)([a-z0-9]{2,7}?)-?(\\d{2,5})")
	matches := re.FindAllStringSubmatch(query, -1)
	keywords := mapset.NewSet()
	for _, match := range matches {
		series := strings.ToUpper(match[1])
		num := match[2]
		keywords.Add(fmt.Sprintf("%s-%s", series, num))
		keywords.Add(fmt.Sprintf("%s-%04s", series, num))
	}
	return keywords
}

func mgsSearchKeyword(keyword string, wg *sync.WaitGroup, metach chan MovieMeta) {
	glog.Info("Keyword: ", keyword)
	urlstrs := []string{
		fmt.Sprintf(
			"http://www.mgstage.com/search/search.php?search_word=%s&search_shop_id=shiroutotv",
			url.QueryEscape(keyword),
		),
		fmt.Sprintf(
			"http://www.mgstage.com/search/search.php?search_word=%s&search_shop_id=nanpatv",
			url.QueryEscape(keyword),
		),
		fmt.Sprintf(
			"http://www.mgstage.com/search/search.php?search_word=%s",
			url.QueryEscape(keyword),
		),
	}
	for _, urlstr := range urlstrs {
		wg.Add(1)
		go func(urlstr string) {
			defer wg.Done()
			mgsParseSearchPage(keyword, urlstr, wg, metach)
		}(urlstr)
	}
}

func mgsParseSearchPage(keyword string, urlstr string, wg *sync.WaitGroup, metach chan MovieMeta) {
	glog.V(2).Info("Search page: ", urlstr)
	doc, err := newDocument(urlstr, httpx.GetContentViaPhantomJS([]*http.Cookie{&mgsCookie}, 0, "center_column", ""))
	if err != nil {
		glog.V(2).Infof("Error parsing %s: %v", urlstr, err)
		return
	}

	urlbase, err := url.Parse(urlstr)
	doc.Find("#center_column > div.all_search_list > ul > li > a, #center_column > div.search_list > div > ul > li > a").Each(
		func(i int, a *goquery.Selection) {
			wg.Add(1)
			go func(a *goquery.Selection) {
				defer wg.Done()
				href, ok := a.Attr("href")
				if !ok {
					return
				}
				urlhref, err := urlbase.Parse(href)
				if err != nil {
					glog.V(2).Info(err)
					return
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					mgsParseProductPage(urlhref.String(), keyword, metach)
				}()
			}(a)
		})
}

func mgsParseProductPage(urlstr string, keyword string, metach chan MovieMeta) {
	glog.V(2).Info("Product page: ", urlstr)
	doc, err := newDocument(urlstr, httpx.GetContentViaPhantomJS([]*http.Cookie{&mgsCookie}, 0, "center_column", ""))
	if err != nil {
		glog.V(2).Infof("Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	meta.Page = urlstr
	meta.Title = doc.Find("#center_column h1.tag").Text()
	meta.CoverImage, _ = doc.Find("#center_column a#EnlargeImage").Attr("href")

	dateRe := regexp.MustCompile("\\d+/\\d+/\\d+")
	meta.ReleaseDate = dateRe.FindString(doc.Find("span.date").Text())

	doc.Find("ul.detail_txt > li").Each(
		func(i int, li *goquery.Selection) {
			text := li.Text()
			if strings.Contains(text, "シリーズ") {
				meta.Series = li.Find("a").Text()
			} else if strings.Contains(text, "出演者") {
				meta.Actresses = li.Find("a").Map(
					func(i int, a *goquery.Selection) string {
						return a.Text()
					})
			} else if strings.Contains(text, "ジャンル") {
				meta.Genres = li.Find("a").Map(
					func(i int, a *goquery.Selection) string {
						return a.Text()
					})
			} else {
				re := regexp.MustCompile("・([^・]*)：([^・]*)")
				matches := re.FindAllStringSubmatch(text, -1)
				for _, match := range matches {
					key := strings.TrimSpace(match[1])
					value := strings.TrimSpace(match[2])
					switch key {
					case "収録時間":
						meta.MovieLength = value
					case "品番":
						meta.Code = value
					}
				}
			}
		})

	doc.Find("div.detail_data > table > tbody > tr").Each(
		func(i int, tr *goquery.Selection) {
			th := tr.Find("th").First()
			k := strings.TrimSpace(th.Text())
			td := tr.Find("td").First()
			if strings.Contains(k, "出演") {
				meta.Actresses = td.Find("a").Map(
					func(i int, a *goquery.Selection) string {
						return a.Text()
					})
			} else if strings.Contains(k, "メーカー") {
				meta.Maker = td.Text()
			} else if strings.Contains(k, "収録時間") {
				meta.MovieLength = td.Text()
			} else if strings.Contains(k, "品番") {
				meta.Code = td.Text()
			} else if strings.Contains(k, "配信開始日") {
				meta.ReleaseDate = td.Text()
			} else if strings.Contains(k, "シリーズ") {
				meta.Series = td.Text()
			} else if strings.Contains(k, "レーベル") {
				meta.Label = td.Text()
			} else if strings.Contains(k, "ジャンル") {
				meta.Genres = td.Find("a").Map(
					func(i int, a *goquery.Selection) string {
						return a.Text()
					})
			}
		})

	if !dmmIsCodeEqual(keyword, meta.Code) {
		glog.V(2).Infof("Code mismatch: Expected %s, got %s", keyword, meta.Code)
	} else {
		metach <- meta
	}
}
