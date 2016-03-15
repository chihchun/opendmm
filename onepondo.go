package opendmm

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

func opdSearchKeyword(keyword string, db *bolt.DB, metach chan MovieMeta) {
	glog.Info("[OPD] Keyword: ", keyword)
	var meta MovieMeta
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("MovieMeta"))
		bdata := bucket.Get([]byte(keyword))
		if bdata == nil {
			return fmt.Errorf("[OPD] Not found: %s", keyword)
		}
		return json.Unmarshal(bdata, &meta)
	})
	if err == nil {
		metach <- meta
	}
}

func opdSearch(query string, dbpath string, metach chan MovieMeta) *sync.WaitGroup {
	glog.Info("[OPD] Query: ", query)
	wg := new(sync.WaitGroup)
	db, err := newDB(dbpath)
	if err != nil {
		glog.Error("[OPD] Can't open DB")
		return wg
	}

	re := regexp.MustCompile("(\\d{6})[-_](\\d{3})")
	matches := re.FindAllStringSubmatch(query, -1)
	for _, match := range matches {
		keyword := fmt.Sprintf("1pondo %s_%s", match[1], match[2])
		wg.Add(1)
		go func() {
			defer wg.Done()
			opdSearchKeyword(keyword, db, metach)
		}()
	}
	return wg
}
