package opendmm

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "net/http"

  "github.com/golang/glog"
  "github.com/PuerkitoBio/goquery"
  "golang.org/x/net/html/charset"
)

func newUtf8Document(url string) (*goquery.Document, error) {
  res, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  if res.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("Unexpected status code %d", res.StatusCode)
  }

  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }
  rawdoc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
  if err != nil {
    return nil, err
  }
  contenttype, _ := rawdoc.Find("meta[http-equiv=content-type]").Attr("content")
  encoding, name, certain := charset.DetermineEncoding(body, contenttype)
  if certain {
    glog.Infof("[UTILS] Encoding of %v is %v", url, name)
  } else {
    glog.Infof("[UTILS] Guess encoding of %v is %v", url, name)
  }
  utfBody, err := encoding.NewDecoder().Bytes(body)
  if err != nil {
    return nil, err
  }
  doc, err := goquery.NewDocumentFromReader(bytes.NewReader(utfBody))
  return doc, err
}
