package opendmm

import (
  "fmt"
  "net/url"
  "regexp"
  "strings"
  "sync"

  "github.com/golang/glog"
  "github.com/PuerkitoBio/goquery"
)

func tkhParseCode(code string) string {
  re := regexp.MustCompile("(?i)([a-z]+)(\\d+)")
  meta := re.FindStringSubmatch(code)
  if meta != nil {
    return fmt.Sprintf("%s-%s", strings.ToUpper(meta[1]), meta[2])
  }
  return code
}

func tkhParse(murl string, metach chan MovieMeta) {
  glog.Info("[TKH] Parse: ", murl)
  doc, err := newUtf8Document(murl)
  if err != nil {
    glog.Error("[TKH] Error: ", err)
    return
  }

  var meta MovieMeta
  meta.Page = murl
  meta.Title = doc.Find("#container > div.pagetitle > h2").Text()
  meta.CoverImage, _ = doc.Find("#container > div.movie.cf > div.in > div.flowplayer > video").Attr("poster")
  meta.SampleImages = doc.Find("#main > div.contents > div.scap > a, #main > div.contents > div.vcap > a").Map(
    func(i int, a *goquery.Selection) string {
      href, _ := a.Attr("href")
      return href
    })
  doc.Find("#main > div.contents > div.infowrapper > dl > dt").Each(
    func(i int, dt *goquery.Selection) {
      k := dt.Text()
      dd := dt.Next()
      if strings.Contains(k, "出演者") {
        meta.Actresses = dd.Find("a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(k, "シリーズ") {
        meta.MovieLength = dd.Text()
      } else if strings.Contains(k, "カテゴリ") {
        meta.Categories = dd.Find("a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(k, "配信開始日") {
        meta.ReleaseDate = dd.Text()
      } else if strings.Contains(k, "収録時間") {
        meta.MovieLength = dd.Text()
      } else if strings.Contains(k, "作品番号") {
        meta.Code = fmt.Sprintf("Tokyo Hot %s", dd.Text())
      }
    })
  metach <- meta
}

func tkhSearchKeyword(keyword string, metach chan MovieMeta) {
  glog.Info("[TKH] Keyword: ", keyword)
  urlstr := fmt.Sprintf(
    "http://www.tokyo-hot.com/product/?q=%s",
    url.QueryEscape(keyword),
  )
  glog.Info("[TKH] Search: ", urlstr)
  doc, err := newUtf8Document(urlstr)
  if (err != nil) {
    glog.Error("[TKH] Error: ", err)
    return
  }

  href, ok := doc.Find("#main > ul > li > a").First().Attr("href")
  if ok {
    urlbase, err := url.Parse(urlstr)
    if err != nil {
      return
    }
    urlhref, err := urlbase.Parse(href)
    if err != nil {
      return
    }
    tkhParse(urlhref.String(), metach)
  }
}

func tkhSearch(query string, metach chan MovieMeta) *sync.WaitGroup {
  glog.Info("[TKH] Query: ", query)
  wg := new(sync.WaitGroup)
  re := regexp.MustCompile("(?i)(k|n)(\\d{3,4})")
  matches := re.FindAllStringSubmatch(query, -1)
  for _, match := range matches {
    keyword := fmt.Sprintf("%s%04s", strings.ToLower(match[1]), match[2])
    wg.Add(1)
    go func() {
      defer wg.Done()
      tkhSearchKeyword(keyword, metach)
    }()
  }
  return wg
}
