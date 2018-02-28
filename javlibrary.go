package opendmm

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
)

func javSearch(query string, metach chan MovieMeta) {
	glog.Info("[JAV] Query: ", query)
	keywords := dmmGuess(query)
	wg := new(sync.WaitGroup)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			javSearchKeyword(keyword, wg, metach)
		}(keyword.(string))
	}
	wg.Wait()
}

func javSearchKeyword(keyword string, wg *sync.WaitGroup, metach chan MovieMeta) {
	glog.Info("[JAV] Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"http://www.javlibrary.com/ja/vl_searchbyid.php?keyword=%s",
		url.QueryEscape(keyword),
	)
	javParse(urlstr, keyword, wg, metach)
}

func javParse(urlstr string, keyword string, wg *sync.WaitGroup, metach chan MovieMeta) {
	glog.Info("[JAV] Product/Search page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, http.Get)
	if err != nil {
		glog.Warningf("[JAV] Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	var ok bool
	meta.Page, ok = doc.Find("link[rel=shortlink]").Attr("href")
	if ok {
		meta.Code = doc.Find("#video_id .text").Text()
		meta.Title = doc.Find("#video_title > h3").Text()
		meta.CoverImage, _ = doc.Find("#video_jacket > img").Attr("src")
		meta.ReleaseDate = doc.Find("#video_date .text").Text()
		meta.MovieLength = doc.Find("#video_length .text").Text()
		meta.Directors = doc.Find("#video_director .text span.director").Map(
			func(i int, span *goquery.Selection) string {
				return span.Text()
			})
		meta.Maker = doc.Find("#video_maker .text").Text()
		meta.Label = doc.Find("#video_label .text").Text()
		meta.Genres = doc.Find("#video_genres .text span.genre").Map(
			func(i int, span *goquery.Selection) string {
				return span.Text()
			})
		meta.Actresses = doc.Find("#video_cast .text span.cast span.star").Map(
			func(i int, span *goquery.Selection) string {
				return span.Text()
			})

		if !dmmIsCodeEqual(keyword, meta.Code) {
			glog.Warningf("[JAV] Code mismatch: Expected %s, got %s", keyword, meta.Code)
		} else {
			metach <- meta
		}
	} else {
		urlbase, err := url.Parse(urlstr)
		if err != nil {
			glog.Errorf("[JAV] %s", err)
			return
		}
		doc.Find("div.videothumblist > div.videos > div.video > a").Each(
			func(i int, a *goquery.Selection) {
				href, ok := a.Attr("href")
				if !ok {
					return
				}
				urlhref, err := urlbase.Parse(href)
				if err != nil {
					glog.Errorf("[JAV] %s", err)
					return
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					javParse(urlhref.String(), keyword, wg, metach)
				}()
			})
	}
}
