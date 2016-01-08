package opendmm

import (
  "sync"

  "github.com/golang/glog"
)

type MovieMeta struct {
  Actresses      []string
  ActressTypes   []string
  Categories     []string
  Code           string
  CoverImage     string
  Description    string
  Directors      []string
  Genres         []string
  Label          string
  Maker          string
  MovieLength    string
  Page           string
  ReleaseDate    string
  SampleImages   []string
  Series         string
  Tags           []string
  ThumbnailImage string
  Title          string
}

func Search(query string, dbpath string) chan MovieMeta {
  metach := make(chan MovieMeta)

  var wgs [](*sync.WaitGroup)
  wgs = append(wgs, aveSearch(query, metach))
  wgs = append(wgs, caribSearch(query, metach))
  wgs = append(wgs, caribprSearch(query, metach))
  wgs = append(wgs, dmmSearch(query, metach))
  wgs = append(wgs, heyzoSearch(query, metach))
  wgs = append(wgs, javSearch(query, metach))
  wgs = append(wgs, tkhSearch(query, metach))
  db, err := openDB(dbpath)
  if err != nil {
    glog.Error(err)
  } else {
    wgs = append(wgs, opdSearch(query, db, metach))
  }
  go func() {
    for _, wg := range wgs {
      wg.Wait()
    }
    close(metach)
  }()
  return postprocess(metach)
}

func Crawl(dbpath string) {
  db, err := openDB(dbpath)
  if err != nil {
    glog.Fatal(err)
  }
  defer db.Close()
  metach := make(chan MovieMeta)

  var wgs [](*sync.WaitGroup)
  wgs = append(wgs, opdCrawl(metach))
  go func() {
    for _, wg := range wgs {
      wg.Wait()
    }
    close(metach)
  }()

  saveToDB(postprocess(metach), db)
}
