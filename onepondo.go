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
	doc, err := newDocument(urlstr, httpx.GetContentViaPhantomJS([]*http.Cookie{&mgsCookie}, 0, "<body", ""))
	if err != nil {
		glog.V(2).Infof("Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	meta.Code = fmt.Sprintf("1pondo %s", keyword)
	meta.Page = urlstr
	meta.CoverImage = fmt.Sprintf("https://www.1pondo.tv/assets/sample/%s/str.jpg", keyword)
	meta.Title = doc.Find("div.movie-overview h1").Text()
	meta.Description = doc.Find("div.movie-detail p").Text()
	doc.Find("li.movie-detail__spec").Each(
		func(i int, li *goquery.Selection) {
			title := li.Find("span.spec-title").Text()
			content := li.Find("span.spec-content")
			if strings.Contains(title, "配信日") {
				meta.ReleaseDate = content.Text()
			} else if strings.Contains(title, "出演") {
				content.Find("a").Each(
					func(i int, a *goquery.Selection) {
						meta.Actresses = append(meta.Actresses, a.Text())
					})
			} else if strings.Contains(title, "シリーズ") {
				meta.Series = content.Text()
			} else if strings.Contains(title, "再生時間") {
				meta.MovieLength = content.Text()
			} else if strings.Contains(title, "タグ") {
				content.Find("a").Each(
					func(i int, a *goquery.Selection) {
						meta.Tags = append(meta.Tags, a.Text())
					})
			}
		})
	doc.Find("img.gallery-image").Each(
		func(i int, img *goquery.Selection) {
			src, ok := img.Attr("data-vue-img-src")
			if ok {
				meta.SampleImages = append(meta.SampleImages, src)
			}
		})

	metach <- meta
}
