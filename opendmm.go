package opendmm

import (
  "reflect"
  "regexp"
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
    for meta := range in {
      glog.Info("[STAGE] Trim spaces")

      value := reflect.ValueOf(&meta).Elem()
      for fi := 0; fi < value.NumField(); fi++ {
        field := value.Field(fi)
        switch field.Interface().(type) {
        case string:
          str := field.String()
          str = strings.TrimSpace(str)
          str = regexp.MustCompile("\\s+").ReplaceAllString(str, " ")
          field.SetString(str)
        case []string:
          for ei := 0; ei < field.Len(); ei++ {
            elem := field.Index(ei)
            str := elem.String()
            str = strings.TrimSpace(str)
            str = regexp.MustCompile("\\s+").ReplaceAllString(str, " ")
            elem.SetString(str)
          }
        }
      }
      out <- meta
    }
  }()
  return out
}

func validateFields(in chan MovieMeta) chan MovieMeta {
  out := make(chan MovieMeta)
  go func() {
    defer close(out)
    for meta := range in {
      glog.Info("[STAGE] Validate fields")
      if meta.Code == "" || meta.Title == "" || meta.CoverImage == "" {
        glog.Warning("[STAGE] Validate failed")
      } else {
        out <- meta
      }
    }
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
  wg.Add(1)
  go func() {
    defer wg.Done()
    caribSearch(query, metach, &wg)
  }()
  wg.Add(1)
  go func() {
    defer wg.Done()
    caribprSearch(query, metach, &wg)
  }()
  wg.Add(1)
  go func() {
    defer wg.Done()
    aveSearch(query, metach, &wg)
  }()
  wg.Add(1)
  go func() {
    defer wg.Done()
    heyzoSearch(query, metach, &wg)
  }()

  go func() {
    wg.Wait()
    close(metach)
  }()
  return validateFields(trimSpaces(metach))
}
