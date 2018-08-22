package opendmm

import (
	"testing"

	"github.com/benbjohnson/phantomjs"
)

func TestOpd(t *testing.T) {
	phantomjs.DefaultProcess.Open()
	defer phantomjs.DefaultProcess.Close()

	queries := []string{"1pondo 022018_648"}
	assertSearchable(t, queries, opdSearch)
}
