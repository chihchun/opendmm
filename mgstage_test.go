package opendmm

import (
	"testing"

	"github.com/benbjohnson/phantomjs"
)

func TestMgstage(t *testing.T) {
	phantomjs.DefaultProcess.Open()
	defer phantomjs.DefaultProcess.Close()

	queries := []string{
		"SIRO-1715",
		"200GANA-894",
		"259LUXU-011",
		"3DSVR-200",
	}
	assertSearchable(t, queries, mgsSearch)
}
