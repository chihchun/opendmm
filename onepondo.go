package opendmm

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/golang/glog"
	"github.com/syndtr/goleveldb/leveldb"
)

func opdSearchKeyword(
	keyword string,
	movieCache *leveldb.DB,
	metach chan MovieMeta,
) {
	bdata, err := movieCache.Get([]byte(keyword), nil)
	if err != nil {
		glog.Warningf("[OPD] %s not found in cache", keyword)
		return
	}
	var meta MovieMeta
	err = json.Unmarshal(bdata, &meta)
	if err != nil {
		glog.Warningf(
			"[OPD] Error unmarshaling data for %s:\n%s", keyword, string(bdata))
		return
	}
	metach <- meta
}

func opdSearch(
	query string,
	httpCache *leveldb.DB,
	movieCache *leveldb.DB,
	metach chan MovieMeta,
) *sync.WaitGroup {
	glog.Info("[OPD] Query: ", query)
	wg := new(sync.WaitGroup)
	re := regexp.MustCompile("(\\d{6})[-_](\\d{3})")
	matches := re.FindAllStringSubmatch(query, -1)
	for _, match := range matches {
		keyword := fmt.Sprintf("%s-%s", match[1], match[2])
		wg.Add(1)
		go func() {
			defer wg.Done()
			opdSearchKeyword(keyword, movieCache, metach)
		}()
	}
	return wg
}
