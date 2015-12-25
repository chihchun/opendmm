package opendmm

import (
  "fmt"
  "net/url"
  "regexp"
  "strings"

  "github.com/golang/glog"
  "github.com/PuerkitoBio/goquery"
)

func javParse(urlstr string, meta chan MovieMeta) {
  glog.Info("[JAV] Parse: ", urlstr)
  doc, err := newUtf8Document(urlstr)
  if err != nil {
    glog.Error("[JAV] Error: ", err)
    return
  }

  var m MovieMeta
  var ok bool
  m.Page, ok = doc.Find("link[rel=shortlink]").Attr("href")
  if !ok {
    base, err := url.Parse(urlstr)
    if err == nil {
      doc.Find("div.videothumblist > div.videos > div.video > a").Each(
        func(i int, s *goquery.Selection) {
          href, ok := s.Attr("href")
          if ok {
            hrefurl, err := base.Parse(href)
            if err == nil {
              go javParse(hrefurl.String(), meta)
            }
          }
        })
    }
    return
  }

  m.Code = doc.Find("#video_id .text").Text()
  m.Title = strings.Replace(doc.Find("#video_title > h3").Text(), m.Code, "", -1)
  m.CoverImage, ok = doc.Find("#video_jacket > img").Attr("src")
  m.ReleaseDate = doc.Find("#video_date .text").Text()
  m.MovieLength = doc.Find("#video_length .text").Text()
  m.Directors = doc.Find("#video_director .text span.director").Map(
    func(i int, span *goquery.Selection) string {
      return span.Text()
    })
  m.Maker = doc.Find("#video_maker .text").Text()
  m.Label = doc.Find("#video_label .text").Text()
  m.Genres = doc.Find("#video_genres .text span.genre").Map(
    func(i int, span *goquery.Selection) string {
      return span.Text()
    })
  m.Actresses = doc.Find("#video_cast .text span.cast span.star").Map(
    func(i int, span *goquery.Selection) string {
      return span.Text()
    })
  meta <- m
}

func javSearchKeyword(kw string, meta chan MovieMeta) {
  glog.Info("[JAV] Keyword: ", kw)
  url := fmt.Sprintf(
    "http://www.javlibrary.com/ja/vl_searchbyid.php?keyword=%s",
    url.QueryEscape(kw),
  )
  go javParse(url, meta)
}

func javSearch(q string, meta chan MovieMeta) {
  glog.Info("[JAV] Query: ", q)
  re := regexp.MustCompile("(?i)([a-z]{2,6})-?(\\d{2,5})")
  ms := re.FindAllStringSubmatch(q, -1)
  for _, m := range ms {
    kw := fmt.Sprintf("%s-%s", m[1], m[2])
    go javSearchKeyword(kw, meta)
  }
}
