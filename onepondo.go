package opendmm

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/junzh0u/httpx"
	"github.com/syndtr/goleveldb/leveldb"
)

func opdParse(urlstr string, httpCache *leveldb.DB) (MovieMeta, error) {
	var meta MovieMeta

	glog.Info("[OPD] Product page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, httpx.Cached(httpCache)(httpx.GetFullPage))
	if err != nil {
		glog.Warningf("[OPD] Error parsing %s: %v", urlstr, err)
		return meta, err
	}

	urlbase, err := url.Parse(urlstr)
	if err != nil {
		glog.Errorf("[OPD] %s", err)
		return meta, err
	}

	meta.Page = urlstr

	poster, ok := doc.Find("video").Attr("poster")
	if !ok {
		err = fmt.Errorf("No poster found")
		glog.Errorf("[OPD] %s", err)
		return meta, err
	}
	urlposter, err := urlbase.Parse(poster)
	if err != nil {
		glog.Errorf("[OPD] %s", err)
		return meta, err
	}
	meta.CoverImage = urlposter.String()
	meta.Title = doc.Find("h2.m-title").Text()
	meta.Actresses = doc.Find("div.video-actor > a").Map(
		func(i int, a *goquery.Selection) string {
			return strings.TrimSpace(strings.Replace(a.Text(), "の作品一覧", "", -1))
		})
	meta.ReleaseDate = doc.Find("dd.m-release").Text()
	meta.MovieLength = doc.Find("dd.m-duration").Text()
	meta.Tags = doc.Find("div.m-tag > a").Map(
		func(i int, a *goquery.Selection) string {
			return a.Text()
		})
	meta.Description = doc.Find("div.m-comment").Text()

	re := regexp.MustCompile("\\d{6}_\\d{3}")
	meta.Code = "1pondo " + re.FindString(meta.CoverImage)

	return meta, nil
}

func opdCrawl(httpCache *leveldb.DB, metach chan MovieMeta) *sync.WaitGroup {
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errcnt := 0
		for id := 1; errcnt < 100; id++ {
			urlstr := fmt.Sprintf("http://m.1pondo.tv/movies/%d/", id)
			meta, err := opdParse(urlstr, httpCache)
			if err != nil {
				glog.Infof("[OPD] Error Count: %d", errcnt)
				httpCache.Delete([]byte(urlstr), nil)
				errcnt++
			} else {
				metach <- meta
				errcnt = 0
			}
		}
	}()
	return wg
}

func opdSearchKeyword(
	keyword string,
	movieCache *leveldb.DB,
	metach chan MovieMeta,
) {
	bdata, err := movieCache.Get([]byte("1pondo "+keyword), nil)
	if err != nil {
		glog.Warningf("[OPD] %s not found in cache", keyword)
		return
	}
	var meta MovieMeta
	err = json.Unmarshal(bdata, &meta)
	if err != nil {
		glog.Warningf(
			"[OPD] Error unmarshaling data for %s:\n%s", keyword, string(bdata))
		return
	}
	metach <- meta
}

func opdSearch(movieCache *leveldb.DB) SearchFunc {
	return func(query string, metach chan MovieMeta) *sync.WaitGroup {
		glog.Info("[OPD] Query: ", query)
		wg := new(sync.WaitGroup)
		re := regexp.MustCompile("(\\d{6})[-_](\\d{3})")
		matches := re.FindAllStringSubmatch(query, -1)
		for _, match := range matches {
			keyword := fmt.Sprintf("%s_%s", match[1], match[2])
			wg.Add(1)
			go func() {
				defer wg.Done()
				opdSearchKeyword(keyword, movieCache, metach)
			}()
		}
		return wg
	}
}
