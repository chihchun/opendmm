package opendmm

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/deckarep/golang-set"
	"github.com/golang/glog"
	"github.com/junzh0u/httpx"
)

const (
	savePageJS string = `var system = require('system');
var page = require('webpage').create();

page.onError = function(msg, trace) {
	// do nothing
};

phantom.addCookie({
  'name'     : 'adc',
  'value'    : '1',
  'domain'   : '.mgstage.com',
  'path'     : '/',
  'httponly' : false,
  'secure'   : false,
  'expires'  : (new Date()).getTime() + (1000 * 60 * 60)
});

page.open(system.args[1], function(status) {
  console.log(page.content);
  phantom.exit();
});`
)

func mgsSearch(query string, metach chan MovieMeta) *sync.WaitGroup {
	glog.Info("[MGS] Query: ", query)
	keywords := mgsGuess(query)
	wg := new(sync.WaitGroup)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			mgsSearchKeyword(keyword, wg, metach)
		}(keyword.(string))
	}
	return wg
}

func mgsGuess(query string) mapset.Set {
	re := regexp.MustCompile("(?i)([a-z0-9]{2,7}?)-?(\\d{2,5})")
	matches := re.FindAllStringSubmatch(query, -1)
	keywords := mapset.NewSet()
	for _, match := range matches {
		keywords.Add(fmt.Sprintf("%s-%s", strings.ToUpper(match[1]), match[2]))
	}
	return keywords
}

func mgsSearchKeyword(keyword string, wg *sync.WaitGroup, metach chan MovieMeta) {
	glog.Info("[MGS] Keyword: ", keyword)
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
	glog.Info("[MGS] Search page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, httpx.GetWithPhantomJS(savePageJS))
	if err != nil {
		glog.Warningf("[MGS] Error parsing %s: %v", urlstr, err)
		return
	}

	urlbase, err := url.Parse(urlstr)
	doc.Find("ul > li > p.title > a").Each(
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
					glog.Warningf("[MGS] %s", err)
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
	glog.Info("[MGS] Product page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, httpx.GetWithPhantomJS(savePageJS))
	if err != nil {
		glog.Warningf("[MGS] Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	meta.Page = urlstr
	meta.Title = doc.Find("#center_column h1.tag").Text()
	meta.CoverImage, _ = doc.Find("#center_column a#EnlargeImage").Attr("href")

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
					glog.Infof("Key: %s; Value: %s", key, value)
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

	metach <- meta
}
