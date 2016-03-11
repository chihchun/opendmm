package opendmm

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

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
				glog.Warning("[STAGE] Validate failed: ", meta)
			} else {
				out <- meta
			}
		}
	}()
	return out
}

func postprocess(in chan MovieMeta) chan MovieMeta {
	return validateFields(trimSpaces(deduplicate(in)))
}

func saveToDB(metach chan MovieMeta, db *bolt.DB) {
	for meta := range metach {
		glog.Info("[STAGE] Save to DB")
		err := writeMetaToDB(meta, db)
		if err != nil {
			glog.Errorf("[STAGE] Failed to save %+v", meta)
		}
	}
}
