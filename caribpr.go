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

func caribprParse(keyword string, urlstr string, metach chan MovieMeta, wg *sync.WaitGroup) {
  glog.Info("[CARIBPR] Parse: ", urlstr)
  doc, err := newUtf8Document(urlstr)
  if err != nil {
    glog.Error("[CARIBPR] Error: ", err)
    return
  }

  var meta MovieMeta
  meta.Code = fmt.Sprintf("Caribpr %s", keyword)
  meta.Page = urlstr

  var urlbase *url.URL
  urlbase, err = url.Parse(urlstr)
  if err != nil {
    return
  }
  var urlcover *url.URL
  urlcover, err = urlbase.Parse("./images/l_l.jpg")
  if err == nil {
    meta.CoverImage = urlcover.String()
  }
  var urlthumbnail *url.URL
  urlthumbnail, err = urlbase.Parse("./images/main_b.jpg")
  if err == nil {
    meta.ThumbnailImage = urlthumbnail.String()
  }

  meta.Title = doc.Find("#main-content > div.main-content-movieinfo > div.video-detail").Text()
  meta.Description = doc.Find("#main-content > div.main-content-movieinfo > div.movie-comment").Text()
  doc.Find("#main-content > div.detail-content.detail-content-gallery > ul > li > div > a").Each(
    func(i int, a *goquery.Selection) {
      href, ok := a.Attr("href")
      if ok {
        if !strings.Contains(href, "/member/") {
          meta.SampleImages = append(meta.SampleImages, href)
        }
      }
    })

  doc.Find("#main-content > div.main-content-movieinfo > div.movie-info > dl").Each(
    func(i int, dl *goquery.Selection) {
      dt := dl.Find("dt")
      if strings.Contains(dt.Text(), "出演") {
        meta.Actresses = dl.Find("dd a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(dt.Text(), "カテゴリー") {
        meta.Categories = dl.Find("dd a").Map(
          func(i int, a *goquery.Selection) string {
            return a.Text()
          })
      } else if strings.Contains(dt.Text(), "販売日") {
        meta.ReleaseDate = dl.Find("dd").Text()
      } else if strings.Contains(dt.Text(), "再生時間") {
        meta.MovieLength = dl.Find("dd").Text()
      } else if strings.Contains(dt.Text(), "スタジオ") {
        meta.Maker = dl.Find("dd").Text()
      } else if strings.Contains(dt.Text(), "シリーズ") {
        meta.Series = dl.Find("dd").Text()
      }
    })

  metach <- meta
}

func caribprSearchKeyword(keyword string, metach chan MovieMeta, wg *sync.WaitGroup) {
  glog.Info("[CARIBPR] Keyword: ", keyword)
  urlstr := fmt.Sprintf(
    "http://www.caribbeancompr.com/moviepages/%s/index.html",
    url.QueryEscape(keyword),
  )
  caribprParse(keyword, urlstr, metach, wg)
}

func caribprSearch(query string, metach chan MovieMeta, wg *sync.WaitGroup) {
  glog.Info("[CARIBPR] Query: ", query)
  re := regexp.MustCompile("(\\d{6})[-_](\\d{3})")
  matches := re.FindAllStringSubmatch(query, -1)
  for _, match := range matches {
    keyword := fmt.Sprintf("%s_%s", match[1], match[2])
    wg.Add(1)
    go func() {
      defer wg.Done()
      caribprSearchKeyword(keyword, metach, wg)
    }()
  }
}
