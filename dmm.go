package opendmm

import (
  "fmt"
  "net/url"
  "regexp"
  "strings"

  "github.com/golang/glog"
  "github.com/PuerkitoBio/goquery"
)

func dmmParseCode(code string) string {
  re := regexp.MustCompile("(?i)([a-z]+)(\\d+)")
  meta := re.FindStringSubmatch(code)
  if meta != nil {
    return fmt.Sprintf("%s-%s", strings.ToUpper(meta[1]), meta[2])
  }
  return code
}

func dmmParse(murl string, metach chan MovieMeta) {
  glog.Info("[DMM] Parse: ", murl)
  doc, err := newUtf8Document(murl)
  if err != nil {
    glog.Error("[DMM] Error: ", err)
    return
  }

  var meta MovieMeta
  var ok bool
  meta.Page = murl
  meta.Title = doc.Find("#title").Text()
  meta.ThumbnailImage, ok = doc.Find("#sample-video img").Attr("src")
  meta.CoverImage, ok = doc.Find("#sample-video a").Attr("href")
  if !ok {
    meta.CoverImage = meta.ThumbnailImage
  }
  doc.Find("div.page-detail > table > tbody > tr > td > table > tbody > tr").Each(
    func(i int, tr *goquery.Selection) {
      td := tr.Find("td").First()
      k := td.Text()
      v := td.Next()
      if strings.Contains(k, "配信開始日") {
        meta.ReleaseDate = v.Text()
      } else if strings.Contains(k, "収録時間") {
        meta.MovieLength = v.Text()
      } else if strings.Contains(k, "出演者") {
        meta.Actresses = v.Find("a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(k, "監督") {
        meta.Directors = v.Find("a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(k, "シリーズ") {
        meta.Series = v.Text()
      } else if strings.Contains(k, "メーカー") {
        meta.Maker = v.Text()
      } else if strings.Contains(k, "レーベル") {
        meta.Label = v.Text()
      } else if strings.Contains(k, "ジャンル") {
        meta.Genres = v.Find("a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(k, "品番") {
        meta.Code = dmmParseCode(v.Text())
      }
    })
  metach <- meta
}

func dmmSearchKeyword(kw string, metach chan MovieMeta) {
  glog.Info("[DMM] Keyword: ", kw)
  urlstr := fmt.Sprintf(
    "http://www.dmm.co.jp/search/=/searchstr=%s",
    url.QueryEscape(kw),
  )
  glog.Info("[DMM] Search: ", urlstr)
  doc, err := newUtf8Document(urlstr)
  if (err != nil) {
    glog.Error("[DMM] Error: ", err)
    return
  }

  doc.Find("#list > li > div > p.tmb > a").Each(
    func(i int, s *goquery.Selection) {
      href, ok := s.Attr("href")
      if ok {
        dmmParse(href, metach)
      }
    })
}

func dmmSearch(q string, metach chan MovieMeta) {
  glog.Info("[DMM] Query: ", q)
  re := regexp.MustCompile("(?i)([a-z]{2,6})-?(\\d{2,5})")
  matches := re.FindAllStringSubmatch(q, -1)
  for _, match := range matches {
    kw := fmt.Sprintf("%s-%s", match[1], match[2])
    go dmmSearchKeyword(kw, metach)
  }
}
