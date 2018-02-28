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

func scuteSearch(query string, metach chan MovieMeta) {
	glog.Info("[S-Cute] Query: ", query)

	wg := new(sync.WaitGroup)
	keywords := scuteGuess(query)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			scuteSearchKeyword(keyword, metach)
		}(keyword.(string))
	}
	wg.Wait()
}

func scuteGuess(query string) mapset.Set {
	re := regexp.MustCompile("(?i)(\\d{3})[_ ]([a-z]+)[_ ]#?(\\d{1,2})")
	matches := re.FindAllStringSubmatch(query, -1)
	keywords := mapset.NewSet()
	for _, match := range matches {
		keywords.Add(fmt.Sprintf("%s_%s_%02s", strings.ToUpper(match[1]), strings.ToLower(match[2]), match[3]))
	}
	return keywords
}

func scuteGuessFull(query string) mapset.Set {
	keywords := mapset.NewSet()
	for keyword := range scuteGuess(query).Iter() {
		keywords.Add(fmt.Sprintf("S-Cute %s", keyword))
	}
	return keywords
}

func scuteSearchKeyword(keyword string, metach chan MovieMeta) {
	glog.Info("[S-Cute] Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"http://www.s-cute.com/contents/%s/",
		url.QueryEscape(keyword),
	)
	scuteParse(urlstr, keyword, metach)
}

func scuteParse(urlstr string, keyword string, metach chan MovieMeta) {
	glog.Info("[S-Cute] Product page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, httpx.GetWithPhantomJS(savePageJS))
	if err != nil {
		glog.Warningf("[S-Cute] Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	meta.Code = fmt.Sprintf("S-Cute %s", keyword)
	meta.Page = urlstr

	meta.Title = doc.Find("span.title").Text()
	meta.Description = doc.Find("div.detail > article > p:nth-child(4)").Text()

	bgStyle, ok := doc.Find("#js-sample > div.vjs-poster").Attr("style")
	if ok {
		re := regexp.MustCompile("(?i)url\\((.+)\\)")
		match := re.FindStringSubmatch(bgStyle)
		meta.CoverImage = match[1]
	} else {
		meta.CoverImage, _ = doc.Find("div.nosample > a > img").Attr("src")
	}

	doc.Find("div.cast > h5").Each(
		func(i int, h *goquery.Selection) {
			actress := strings.SplitN(h.Text(), " ", 2)[1]
			meta.Actresses = append(meta.Actresses, actress)
		})

	metach <- meta
}
