package opendmm

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"

	"github.com/golang/glog"
	"github.com/syndtr/goleveldb/leveldb"
)

// ProcessStage is a pipe of MovieMeta
type ProcessStage func(chan MovieMeta) chan MovieMeta

func deduplicate(in chan MovieMeta) chan MovieMeta {
	out := make(chan MovieMeta)
	go func() {
		defer close(out)
		for meta := range in {
			glog.Infof("[STAGE] Deduplicate: %s", meta.Code)
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
			glog.Infof("[STAGE] Trim spaces: %s", meta.Code)

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
			glog.Infof("[STAGE] Validate fields: %s", meta.Code)
			if meta.Code == "" || meta.Title == "" || meta.CoverImage == "" {
				glog.Warningf("[STAGE] Validate failed: %+v", meta)
			} else {
				out <- meta
			}
		}
	}()
	return out
}

func normalizeURLFields(in chan MovieMeta) chan MovieMeta {
	out := make(chan MovieMeta)
	go func() {
		defer close(out)
		for meta := range in {
			meta.CoverImage = normalizeURL(meta.CoverImage)
			meta.Page = normalizeURL(meta.Page)
			meta.SampleImages = normalizeURLs(meta.SampleImages)
			meta.ThumbnailImage = normalizeURL(meta.ThumbnailImage)
			out <- meta
		}
	}()
	return out
}

func postprocess(in chan MovieMeta) chan MovieMeta {
	return normalizeURLFields(validateFields(trimSpaces(deduplicate(in))))
}

func cacheIntoDB(db *leveldb.DB) ProcessStage {
	return func(in chan MovieMeta) chan MovieMeta {
		out := make(chan MovieMeta)
		go func() {
			defer close(out)
			for meta := range in {
				glog.Infof("[STAGE] Caching into DB: %s", meta.Code)
				bdata, err := json.Marshal(meta)
				if err != nil {
					glog.Errorf("[STAGE] Cache failed (%s): %+v", err, meta)
					continue
				}
				db.Put([]byte(meta.Code), bdata, nil)
				out <- meta
			}
		}()
		return out
	}
}
