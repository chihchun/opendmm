package opendmm

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/junzh0u/httpx"
)

func newDocument(url string, getcontent httpx.GetContentFunc) (*goquery.Document, error) {
	content, err := getcontent(url)
	if err != nil {
		return nil, err
	}
	glog.V(3).Infof("=======\nPage:%s\nContent:%s\n=======\n", url, content)
	return goquery.NewDocumentFromReader(strings.NewReader(content))
}

func normalizeURL(in string) string {
	if in == "" {
		return ""
	}
	u, _ := url.Parse(in)
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u.String()
}

func normalizeURLs(in []string) []string {
	var out []string
	for _, s := range in {
		out = append(out, normalizeURL(s))
	}
	return out
}

func unblockingSend(meta MovieMeta, metach chan MovieMeta) {
	select {
	case metach <- meta:
	default:
	}
}
