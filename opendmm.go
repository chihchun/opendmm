package opendmm

import (
  "reflect"
  "strings"
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

func trimSpaces(in chan MovieMeta) chan MovieMeta {
  out := make(chan MovieMeta)
  go func() {
    defer close(out)
    meta, ok := <-in
    if !ok {
      return
    }
    glog.Info("[STAGE] Trim spaces")

    value := reflect.ValueOf(&meta).Elem()
    for fi := 0; fi < value.NumField(); fi++ {
      field := value.Field(fi)
      switch field.Interface().(type) {
      case string:
        field.SetString(strings.TrimSpace(field.String()))
      case []string:
        for ei := 0; ei < field.Len(); ei++ {
          elem := field.Index(ei)
          elem.SetString(strings.TrimSpace(elem.String()))
        }
      }
    }

    out <- meta
  }()
  return out
}

func Search(query string) chan MovieMeta {
  var wg sync.WaitGroup
  metach := make(chan MovieMeta)
  wg.Add(1)
  go func() {
    defer wg.Done()
    dmmSearch(query, metach, &wg)
  }()
  wg.Add(1)
  go func() {
    defer wg.Done()
    javSearch(query, metach, &wg)
  }()
  go func() {
    wg.Wait()
    close(metach)
  }()
  return trimSpaces(metach)
}
