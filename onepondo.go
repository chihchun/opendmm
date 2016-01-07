package opendmm

import (
  "fmt"
  "regexp"
  "strings"
  "sync"

  "github.com/boltdb/bolt"
  "github.com/golang/glog"
  "github.com/junzh0u/httpx"
  "github.com/PuerkitoBio/goquery"
)

func opdParse(db *bolt.DB, urlstr string) {
  glog.Info("[OPD] Product page: ", urlstr)
  doc, err := newDocumentInUTF8(urlstr, httpx.GetMobile)
  if err != nil {
    glog.Warningf("[OPD] Error parsing %s: %v", urlstr, err)
    return
  }

  var meta MovieMeta

  rawhtml, _ := doc.Html()
  re := regexp.MustCompile("original_movie_id\\s*=\\s*(\\d{6}_\\d{3})")
  match := re.FindStringSubmatch(rawhtml)
  if match == nil {
    return
  }
  meta.Page = urlstr
  meta.Code = fmt.Sprintf("1pondo %s", match[1])
  meta.CoverImage = fmt.Sprintf("http://www.1pondo.tv/assets/sample/%s/str.jpg", match[1])
  meta.ThumbnailImage = fmt.Sprintf("http://www.1pondo.tv/assets/sample/%s/thum_b.jpg", match[1])

  meta.Title = doc.Find("#tab-1 > h1").Text()
  doc.Find("#tab-1 > table > tbody > tr").Each(
    func(i int, tr *goquery.Selection) {
      tds := tr.Find("td")
      k := tds.First().Text()
      v := tds.Last()
      if strings.Contains(k, "女優名：") {
        meta.Actresses = v.Find("li").Map(
          func(i int, li *goquery.Selection) string {
            return li.Text()
          })
      } else if strings.Contains(k, "配信日") {
        meta.ReleaseDate = v.Text()
      } else if strings.Contains(k, "プレイ時間") {
        meta.MovieLength = v.Text()
      } else if strings.Contains(k, "ジャンル") {
        meta.Genres = v.Find("li").Map(
          func(i int, li *goquery.Selection) string {
            return li.Text()
          })
      }
    })

  glog.Info(meta)
}

func opdCrawlList(db *bolt.DB, wg *sync.WaitGroup, page int) {
  glog.Info("[OPD] Crawling page ", page)
  urlstr := fmt.Sprintf("http://m.1pondo.tv/listpages/all_%d.html", page)
  doc, err := newDocumentInUTF8(urlstr, httpx.GetMobile)
  if err != nil {
    glog.Warningf("[OPD] Error parsing %s: %v", urlstr, err)
    return
  }
  as := doc.Find("#showList > li > div > a")
  if as.Length() == 0 {
    glog.Info("[OPD] Reached empty page at: ", urlstr)
    return
  }
  as.Each(func(i int, a *goquery.Selection) {
    href, ok := a.Attr("href")
    if !ok {
      return
    }
    wg.Add(1)
    go func() {
      defer wg.Done()
      opdParse(db, href)
    }()
  })
  // opdCrawlList(db, wg, page + 1)
}

func opdCrawl(db *bolt.DB) *sync.WaitGroup {
  glog.Info("[OPD] Crawling start")
  wg := new(sync.WaitGroup)
  wg.Add(1)
  go func() {
    defer wg.Done()
    opdCrawlList(db, wg, 1)
  }()
  return wg
}
