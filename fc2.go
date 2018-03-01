package opendmm

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sync"

	"github.com/deckarep/golang-set"
	"github.com/golang/glog"
)

func fc2Search(query string, wg *sync.WaitGroup, metach chan MovieMeta) {
	keywords := fc2Guess(query)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			fc2SearchKeyword(keyword, metach)
		}(keyword.(string))
	}
}

func fc2Guess(query string) mapset.Set {
	keywords := mapset.NewSet()
	matched, _ := regexp.MatchString("(?i)fc2", query)
	if !matched {
		return keywords
	}

	re := regexp.MustCompile("\\d{6}")
	matches := re.FindAllString(query, -1)
	for _, match := range matches {
		keywords.Add(fmt.Sprintf("%06s", match))
	}
	return keywords
}

func fc2GuessFull(query string) mapset.Set {
	keywords := mapset.NewSet()
	for keyword := range fc2Guess(query).Iter() {
		keywords.Add(fmt.Sprintf("FC2-PPV %s", keyword))
	}
	return keywords
}

func fc2SearchKeyword(keyword string, metach chan MovieMeta) {
	glog.Info("Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"http://adult.contents.fc2.com/article_search.php?id=%s",
		url.QueryEscape(keyword),
	)
	fc2Parse(urlstr, keyword, metach)
}

func fc2Parse(urlstr string, keyword string, metach chan MovieMeta) {
	glog.V(2).Info("Product page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, http.Get)
	if err != nil {
		glog.V(2).Infof("Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	meta.Code = fmt.Sprintf("FC2-PPV %s", keyword)
	meta.Page = urlstr

	meta.Title = doc.Find("section.detail > h2").Text()
	meta.CoverImage, _ = doc.Find("div.main_thum_img > a").Attr("href")
	meta.Description = doc.Find("section.explain > p").Text()

	metach <- meta
}
