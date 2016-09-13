package opendmm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/golang/glog"
)

func libreSearch(query string, metach chan MovieMeta) *sync.WaitGroup {
	glog.Info("[LIBRE] Query: ", query)
	wg := new(sync.WaitGroup)
	keywords := dmmGuess(query)
	for keyword := range keywords.Iter() {
		wg.Add(1)
		go func(keyword string) {
			defer wg.Done()
			libreSearchKeyword(keyword, metach)
		}(keyword.(string))
	}
	return wg
}

func libreSearchKeyword(keyword string, metach chan MovieMeta) {
	glog.Info("[CARIB] Keyword: ", keyword)
	urlstr := fmt.Sprintf(
		"http://www.libredmm.com/products/%s.json",
		url.QueryEscape(keyword),
	)
	libreParse(urlstr, keyword, metach)
}

func libreParse(urlstr string, keyword string, metach chan MovieMeta) {
	glog.Info("[LIBRE] Product page: ", urlstr)
	resp, err := http.Get(urlstr)
	if err != nil {
		glog.Warningf("[LIBRE] Error getting %s: %v", urlstr, err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Warningf("[LIBRE] Error reading %s: %v", urlstr, err)
		return
	}
	var obj map[string](interface{})
	json.Unmarshal(body, &obj)

	var meta MovieMeta
	meta.Actresses = toStringArray(obj["actresses"])
	meta.ActressTypes = toStringArray(obj["actress_types"])
	meta.Categories = toStringArray(obj["categories"])
	meta.Code, _ = obj["code"].(string)
	meta.CoverImage, _ = obj["cover_image"].(string)
	meta.Description, _ = obj["description"].(string)
	meta.Directors = toStringArray(obj["directors"])
	meta.Genres = toStringArray(obj["genres"])
	meta.Label, _ = obj["label"].(string)
	meta.Maker, _ = obj["maker"].(string)
	meta.MovieLength, _ = obj["movie_length"].(string)
	meta.Page, _ = obj["page"].(string)
	meta.ReleaseDate, _ = obj["release_date"].(string)
	meta.SampleImages = toStringArray(obj["sample_images"])
	meta.Series, _ = obj["series"].(string)
	meta.Tags = toStringArray(obj["tags"])
	meta.ThumbnailImage, _ = obj["thumbnail_image"].(string)
	meta.Title, _ = obj["title"].(string)

	metach <- meta
}

func toStringArray(field interface{}) []string {
	var strs []string
	objs, ok := field.([]interface{})
	if ok {
		for _, obj := range objs {
			strs = append(strs, obj.(string))
		}
	}
	return strs
}
