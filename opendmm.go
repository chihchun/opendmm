package opendmm

import (
  "reflect"
  "strings"

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
    m := <-in
    glog.Info("[STAGE] Trim spaces")

    v := reflect.ValueOf(&m).Elem()
    for fi := 0; fi < v.NumField(); fi++ {
      f := v.Field(fi)
      switch f.Interface().(type) {
      case string:
        f.SetString(strings.TrimSpace(f.String()))
      case []string:
        for si := 0; si < f.Len(); si++ {
          sf := f.Index(si)
          sf.SetString(strings.TrimSpace(sf.String()))
        }
      }
    }

    out <- m
  }()
  return out
}

func Search(query string) chan MovieMeta {
  meta := make(chan MovieMeta)
  go dmmSearch(query, meta)
  go javSearch(query, meta)
  return trimSpaces(meta)
}
