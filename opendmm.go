package opendmm

import (
  "sync"

  "github.com/boltdb/bolt"
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

func Search(query string, db *bolt.DB) chan MovieMeta {
  metach := make(chan MovieMeta)

  var wgs [](*sync.WaitGroup)
  wgs = append(wgs, aveSearch(query, metach))
  wgs = append(wgs, caribSearch(query, metach))
  wgs = append(wgs, caribprSearch(query, metach))
  wgs = append(wgs, dmmSearch(query, metach))
  wgs = append(wgs, heyzoSearch(query, metach))
  wgs = append(wgs, javSearch(query, metach))
  wgs = append(wgs, tkhSearch(query, metach))
  if db != nil {
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

func Crawl(db *bolt.DB) {
  if db == nil {
    glog.Fatal("Must provide a db to start crawling")
  }
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
