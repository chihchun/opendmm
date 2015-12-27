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

func heyzoParse(keyword string, urlstr string, metach chan MovieMeta) {
  glog.Info("[HEYZO] Parse: ", urlstr)
  doc, err := newUtf8Document(urlstr)
  if err != nil {
    glog.Error("[HEYZO] Error: ", err)
    return
  }

  var meta MovieMeta
  meta.Maker = "Heyzo"
  meta.Code = fmt.Sprintf("Heyzo %s", keyword)
  meta.Page = urlstr

  var urlbase *url.URL
  urlbase, err = url.Parse(urlstr)
  if err != nil {
    return
  }
  var urlcover *url.URL
  urlcover, err = urlbase.Parse(
    fmt.Sprintf("/contents/3000/%s/images/player_thumbnail_450.jpg", keyword))
  if err == nil {
    meta.CoverImage = urlcover.String()
  }

  meta.Title = doc.Find("#movie > h1").Text()
  meta.ReleaseDate = doc.Find("#movie > div.info-bg.info-bgWide > div > span.release-day + *").Text()
  meta.Actresses = doc.Find("#movie > div.info-bg.info-bgWide > div > span.actor + *").Find("a").Map(
    func(i int, a *goquery.Selection) string {
      return a.Text()
    })
  meta.Label = strings.Replace(
    doc.Find("#movie > div.info-bg.info-bgWide > div > span.label + *").Text(), "-", "", -1)
  meta.ActressTypes = doc.Find("#movie > div.info-bg.info-bgWide > div > div.actor-type > span").Map(
    func(i int, span *goquery.Selection) string {
      return span.Text()
    })
  meta.Tags = doc.Find("#movie > div.info-bg.info-bgWide > div > div.tag_cloud > ul > li").Map(
    func(i int, li *goquery.Selection) string {
      return li.Text()
    })
  meta.Description = doc.Find("#movie > div.info-bg.info-bgWide > div > p > *").Nodes[0].Data

  metach <- meta
}

func heyzoSearchKeyword(keyword string, metach chan MovieMeta) {
  glog.Info("[HEYZO] Keyword: ", keyword)
  urlstr := fmt.Sprintf(
    "http://www.heyzo.com/moviepages/%s/index.html",
    url.QueryEscape(keyword),
  )
  heyzoParse(keyword, urlstr, metach)
}

func heyzoSearch(query string, metach chan MovieMeta) *sync.WaitGroup {
  glog.Info("[HEYZO] Query: ", query)
  wg := new(sync.WaitGroup)
  matched, err := regexp.Match("(?i)heyzo", []byte(query))
  if err != nil {
    glog.Error("[HEYZO] Error: ", err)
    return wg
  }
  if !matched {
    return wg
  }

  re := regexp.MustCompile("\\d{3,4}")
  matches := re.FindAllString(query, -1)
  for _, match := range matches {
    keyword := fmt.Sprintf("%04s", match)
    wg.Add(1)
    go func() {
      defer wg.Done()
      heyzoSearchKeyword(keyword, metach)
    }()
  }
  return wg
}
