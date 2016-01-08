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

func Search(query string) chan MovieMeta {
  db, err := openDB("/tmp/opendmm.boltdb")
  if err != nil {
    glog.Fatal(err)
  }
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
    wgs = append(wgs, opdSearch(db, query, metach))
  }
  go func() {
    for _, wg := range wgs {
      wg.Wait()
    }
    close(metach)
  }()
  return validateFields(trimSpaces(deduplicate(metach)))
}

func Crawl() {
  db, err := openDB("/tmp/opendmm.boltdb")
  if err != nil {
    glog.Fatal(err)
  }
  defer db.Close()
  metach := make(chan MovieMeta)

  var wgs [](*sync.WaitGroup)
  wgs = append(wgs, opdCrawl(db, metach))
  go func() {
    for _, wg := range wgs {
      wg.Wait()
    }
    close(metach)
  }()

  saveToDB(validateFields(trimSpaces(deduplicate(metach))), db)
}
