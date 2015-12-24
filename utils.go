package opendmm

import (
  "bytes"
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
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }
  encoding, name, certain := charset.DetermineEncoding(body, "")
  if certain {
    glog.Infof("[Utils] Encoding of %v is %v", url, name)
  } else {
    glog.Infof("[Utils] Guess encoding of %v is %v", url, name)
  }
  utfBody, err := encoding.NewDecoder().Bytes(body)
  if err != nil {
    return nil, err
  }
  doc, err := goquery.NewDocumentFromReader(bytes.NewReader(utfBody))
  return doc, err
}
