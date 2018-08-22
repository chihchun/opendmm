package opendmm

import (
	"testing"

	"github.com/benbjohnson/phantomjs"
)

func TestJavlibrary(t *testing.T) {
	phantomjs.DefaultProcess.Open()
	defer phantomjs.DefaultProcess.Close()

	queries := []string{
		"MIAD-617",
	}
	assertSearchable(t, queries, javSearch)
}
