package opendmm

import (
  "fmt"
  "net/url"
  "regexp"
  "sync"

  "github.com/golang/glog"
  "github.com/PuerkitoBio/goquery"
)

func javParse(urlstr string, metach chan MovieMeta) {
  glog.Info("[JAV] Parse: ", urlstr)
  doc, err := newUtf8Document(urlstr)
  if err != nil {
    glog.Error("[JAV] Error: ", err)
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
    metach <- meta
  } else {
    urlbase, err := url.Parse(urlstr)
    if err != nil {
      return
    }
    href, ok := doc.Find("div.videothumblist > div.videos > div.video > a").First().Attr("href")
    if !ok {
      return
    }
    urlhref, err := urlbase.Parse(href)
    if err != nil {
      return
    }
    javParse(urlhref.String(), metach)
  }
}

func javSearchKeyword(keyword string, metach chan MovieMeta) {
  glog.Info("[JAV] Keyword: ", keyword)
  urlstr := fmt.Sprintf(
    "http://www.javlibrary.com/ja/vl_searchbyid.php?keyword=%s",
    url.QueryEscape(keyword),
  )
  javParse(urlstr, metach)
}

func javSearch(query string, metach chan MovieMeta) *sync.WaitGroup {
  glog.Info("[JAV] Query: ", query)
  wg := new(sync.WaitGroup)
  re := regexp.MustCompile("(?i)([a-z]{2,6})-?(\\d{2,5})")
  matches := re.FindAllStringSubmatch(query, -1)
  for _, match := range matches {
    keyword := fmt.Sprintf("%s-%s", match[1], match[2])
    wg.Add(1)
    go func() {
      defer wg.Done()
      javSearchKeyword(keyword, metach)
    }()
  }
  return wg
}
