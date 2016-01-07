package opendmm

import (
  "fmt"
  "net/http"
  "strings"

  "github.com/boltdb/bolt"
  "github.com/junzh0u/httpx"
  "github.com/PuerkitoBio/goquery"
)

func newDocumentInUTF8(url string, getfunc func(string) (*http.Response, error)) (*goquery.Document, error) {
  resp, err := getfunc(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("Unexpected status code %d from %s", resp.StatusCode, url)
  }
  body, err := httpx.RespBodyInUTF8(resp)
  if err != nil {
    return nil, err
  }
  return goquery.NewDocumentFromReader(strings.NewReader(body))
}

func openDB(path string) (*bolt.DB, error) {
  db, err := bolt.Open(path, 0600, nil)
  if err != nil {
    return nil, err
  }
  err = db.Update(func(tx *bolt.Tx) error {
    _, err := tx.CreateBucketIfNotExists([]byte("MovieMeta"))
    if err != nil {
      return fmt.Errorf("Error creating bucket MovieMeta: %s", err)
    }
    return nil
  })
  if err != nil {
    return nil, err
  }
  return db, nil
}
