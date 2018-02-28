package opendmm

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
)

func niceageSearch(query string, metach chan MovieMeta) {
	glog.Info("Query: ", query)
	keywords := dmmGuess(query)
	wg := new(sync.WaitGroup)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			niceageSearchKeyword(keyword, metach)
		}(keyword.(string))
	}
	wg.Wait()
}

func niceageSearchKeyword(keyword string, metach chan MovieMeta) {
	glog.Info("Keyword: ", keyword)
	urlstr := fmt.Sprintf("http://nice-age.net/%s.html", keyword)
	niceageParse(urlstr, metach)
}

func niceageParse(urlstr string, metach chan MovieMeta) {
	glog.Info("Product page: ", urlstr)
	doc, err := newDocumentInUTF8(urlstr, http.Get)
	if err != nil {
		glog.Warningf("Error parsing %s: %v", urlstr, err)
		return
	}

	var meta MovieMeta
	urlbase, err := url.Parse(urlstr)
	if err != nil {
		glog.Error(err)
		return
	}
	imageHref, ok := doc.Find("#detail > div > a > img").Attr("src")
	if !ok {
		glog.Errorf("no cover image")
		return
	}
	urlimage, err := urlbase.Parse(imageHref)
	if err != nil {
		glog.Error(err)
		return
	}
	meta.CoverImage = urlimage.String()

	infoTable := doc.Find("table.product")
	infoTable.Find("th").Each(
		func(i int, th *goquery.Selection) {
			td := th.Next()
			glog.Info(th)
			glog.Info(td)
			if strings.Contains(th.Text(), "タイトル") {
				meta.Title = td.Text()
			} else if strings.Contains(th.Text(), "出演") {
				meta.Actresses = strings.Split(strings.TrimSpace(td.Text()), " ")
			} else if strings.Contains(th.Text(), "型番") {
				meta.Code = niceageParseCode(td.Text())
			} else if strings.Contains(th.Text(), "発売日") {
				meta.ReleaseDate = td.Text()
			} else if strings.Contains(th.Text(), "レーベル") {
				meta.Label = td.Text()
			} else if strings.Contains(th.Text(), "収録時間") {
				meta.MovieLength = td.Text()
			}
		})

	metach <- meta
}

func niceageParseCode(code string) string {
	re := regexp.MustCompile("(?i)([a-z]+)-(\\d+)")
	meta := re.FindStringSubmatch(code)
	if meta != nil {
		return fmt.Sprintf("%s-%s", strings.ToUpper(meta[1]), meta[2])
	}
	return code
}
