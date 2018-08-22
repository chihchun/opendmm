package opendmm

import (
	"testing"

	"github.com/benbjohnson/phantomjs"
)

func TestScute(t *testing.T) {
	phantomjs.DefaultProcess.Open()
	defer phantomjs.DefaultProcess.Close()

	queries := []string{
		"S-Cute 426 Aya #3",
		"SCute 426 Aya #3",
		"485_wakaba_02_hd",
		"258_erina_02",
	}
	assertSearchable(t, queries, scuteSearch)

	blackhole := []string{
		"258_erina_20",
	}
	assertUnsearchable(t, blackhole, scuteSearch)
}
