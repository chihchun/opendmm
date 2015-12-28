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

func deduplicate(in chan MovieMeta) chan MovieMeta {
  out := make(chan MovieMeta)
  go func() {
    defer close(out)
    for meta := range in {
      glog.Info("[STAGE] Deduplicate")
      segments := regexp.MustCompile("\\s").Split(meta.Title, -1)
      for i, segment := range segments {
        if segment == meta.Code {
          segments[i] = ""
        } else {
          for _, actress := range meta.Actresses {
            if segment == actress {
              segments[i] = ""
              break
            }
          }
        }
      }
      meta.Title = strings.Join(segments, " ")
      out <- meta
    }
  }()
  return out
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
  metach := make(chan MovieMeta)
  var wgs [](*sync.WaitGroup)
  wgs = append(wgs, aveSearch(query, metach))
  wgs = append(wgs, caribSearch(query, metach))
  wgs = append(wgs, caribprSearch(query, metach))
  wgs = append(wgs, dmmSearch(query, metach))
  wgs = append(wgs, heyzoSearch(query, metach))
  wgs = append(wgs, javSearch(query, metach))
  wgs = append(wgs, tkhSearch(query, metach))
  go func() {
    for _, wg := range wgs {
      wg.Wait()
    }
    close(metach)
  }()
  return validateFields(trimSpaces(deduplicate(metach)))
}
