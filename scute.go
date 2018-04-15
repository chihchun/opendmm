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

func scuteSearch(query string, wg *sync.WaitGroup, metach chan MovieMeta) {
	keywords := scuteGuess(query)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			scuteSearchKeyword(keyword, metach)
		}(keyword.(string))
	}
}

func scuteGuess(query string) mapset.Set {
	keywords := mapset.NewSet()
	re := regexp.MustCompile("(?i)(\\d{3})[_ ]([a-z]+)[_ ]#?(\\d{1,2})")
	matches := re.FindAllStringSubmatch(query, -1)
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
	glog.Info("Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"http://www.s-cute.com/contents/%s/",
		url.QueryEscape(keyword),
	)
	scuteParse(urlstr, keyword, metach)
}

func scuteParse(urlstr string, keyword string, metach chan MovieMeta) {
	glog.V(2).Info("Product page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, httpx.GetWithPhantomJS(savePageJS))
	if err != nil {
		glog.V(2).Infof("Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	meta.Code = fmt.Sprintf("S-Cute %s", keyword)
	meta.Page = urlstr

	meta.Title = doc.Find("div.blog-single > h3").Text()
	meta.Description = doc.Find("div.blog-single > p").Text()

	meta.CoverImage, _ = doc.Find("div.content-cover > img:nth-child(1)").Attr("src")

	doc.Find("div.about-author > a > h5").Each(
		func(i int, h *goquery.Selection) {
			actress := strings.SplitN(h.Text(), " ", 2)[1]
			meta.Actresses = append(meta.Actresses, actress)
		})

	dateRe := regexp.MustCompile("\\d+/\\d+/\\d+")
	meta.ReleaseDate = dateRe.FindString(doc.Find("div.blog-single span.date").Text())

	metach <- meta
}
