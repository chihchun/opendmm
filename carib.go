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

func caribSearch(query string, wg *sync.WaitGroup, metach chan MovieMeta) {
	keywords := caribGuess(query)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			caribSearchKeyword(keyword, metach)
		}(keyword.(string))
	}
}

func caribGuess(query string) mapset.Set {
	re := regexp.MustCompile("(\\d{6})-(\\d{3})")
	matches := re.FindAllStringSubmatch(query, -1)
	keywords := mapset.NewSet()
	for _, match := range matches {
		keywords.Add(fmt.Sprintf("%s-%s", match[1], match[2]))
	}
	return keywords
}

func caribGuessFull(query string) mapset.Set {
	keywords := mapset.NewSet()
	for keyword := range caribGuess(query).Iter() {
		keywords.Add(fmt.Sprintf("Carib %s", keyword))
	}
	return keywords
}

func caribSearchKeyword(keyword string, metach chan MovieMeta) {
	glog.Info("Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"http://www.caribbeancom.com/moviepages/%s/index.html",
		url.QueryEscape(keyword),
	)
	caribParse(urlstr, keyword, metach)
}

func caribParse(urlstr string, keyword string, metach chan MovieMeta) {
	glog.V(2).Info("Product page: ", urlstr)
	doc, err := newDocument(urlstr, httpx.GetContentInUTF8(http.Get))
	if err != nil {
		glog.V(2).Infof("Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	meta.Code = fmt.Sprintf("Carib %s", keyword)
	meta.Page = urlstr

	var urlbase *url.URL
	urlbase, err = url.Parse(urlstr)
	if err != nil {
		return
	}
	var urlcover *url.URL
	urlcover, err = urlbase.Parse("./images/l_l.jpg")
	if err == nil {
		meta.CoverImage = urlcover.String()
	}
	var urlthumbnail *url.URL
	urlthumbnail, err = urlbase.Parse("./images/main_b.jpg")
	if err == nil {
		meta.ThumbnailImage = urlthumbnail.String()
	}

	meta.Title = doc.Find("#moviepages h1").Text()
	meta.Description = doc.Find("#moviepages div.movie-info > p").Text()

	doc.Find("#moviepages li.movie-detail__spec").Each(
		func(i int, li *goquery.Selection) {
			title := li.Find("span.spec-title").Text()
			content := li.Find("span.spec-content")
			if strings.Contains(title, "出演") {
				meta.Actresses = content.Find("a").Map(
					func(i int, a *goquery.Selection) string {
						return a.Text()
					})
			} else if strings.Contains(title, "販売日") || strings.Contains(title, "配信日") {
				meta.ReleaseDate = content.Text()
			} else if strings.Contains(title, "再生時間") {
				meta.MovieLength = content.Text()
			} else if strings.Contains(title, "シリーズ") {
				meta.Series = content.Text()
			} else if strings.Contains(title, "タグ") {
				meta.Tags = content.Find("a").Map(
					func(i int, a *goquery.Selection) string {
						return a.Text()
					})
			}
		})

	metach <- meta
}
