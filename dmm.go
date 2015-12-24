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
  m := re.FindStringSubmatch(code)
  if m != nil {
    return fmt.Sprintf("%s-%s", strings.ToUpper(m[1]), m[2])
  }
  return code
}

func dmmParse(urlstr string, meta chan MovieMeta) {
  glog.Info("[DMM] Parse: ", urlstr)
  doc, err := newUtf8Document(urlstr)
  if err != nil {
    glog.Error("[DMM] Error: ", err)
    return
  }

  var m MovieMeta
  var ok bool
  m.Page = urlstr
  m.Title = doc.Find("#title").Text()
  m.ThumbnailImage, ok = doc.Find("#sample-video img").Attr("src")
  m.CoverImage, ok = doc.Find("#sample-video a").Attr("href")
  if !ok {
    m.CoverImage = m.ThumbnailImage
  }
  doc.Find("div.page-detail > table > tbody > tr > td > table > tbody > tr").Each(
    func(i int, tr *goquery.Selection) {
      td := tr.Find("td").First()
      k := td.Text()
      v := td.Next()
      if strings.Contains(k, "配信開始日") {
        m.ReleaseDate = v.Text()
      } else if strings.Contains(k, "収録時間") {
        m.MovieLength = v.Text()
      } else if strings.Contains(k, "出演者") {
        m.Actresses = v.Find("a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(k, "監督") {
        m.Directors = v.Find("a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(k, "シリーズ") {
        m.Series = v.Text()
      } else if strings.Contains(k, "メーカー") {
        m.Maker = v.Text()
      } else if strings.Contains(k, "レーベル") {
        m.Label = v.Text()
      } else if strings.Contains(k, "ジャンル") {
        m.Genres = v.Find("a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(k, "品番") {
        m.Code = dmmParseCode(v.Text())
      }
    })
  meta <- m
}

func dmmSearchKeyword(kw string, meta chan MovieMeta) {
  glog.Info("[DMM] Keyword: ", kw)
  urlstr := fmt.Sprintf(
    "http://www.dmm.co.jp/search/=/searchstr=%s",
    url.QueryEscape(kw),
  )
  glog.Info("[DMM] Search: ", urlstr)
  sdoc, err := newUtf8Document(urlstr)
  if (err != nil) {
    glog.Error("[DMM] Error: ", err)
  }

  sdoc.Find("#list > li > div > p.tmb > a").Each(func(i int, s *goquery.Selection) {
    href, ok := s.Attr("href")
    if ok {
      go dmmParse(href, meta)
    }
  })
}

func dmmSearch(q string, meta chan MovieMeta) {
  glog.Info("[DMM] Query: ", q)
  re := regexp.MustCompile("(?i)([a-z]{2,6})-?(\\d{2,5})")
  ms := re.FindAllStringSubmatch(q, -1)
  for _, m := range ms {
    kw := fmt.Sprintf("%s-%s", m[1], m[2])
    go dmmSearchKeyword(kw, meta)
  }
}
