package opendmm

import (
	"testing"
)

func TestScute(t *testing.T) {
	queries := []string{
		"S-Cute 426 Aya #3",
		"SCute 426 Aya #3",
		"scute 485_wakaba_02_hd",
		"s_cute 258_erina_02",
	}
	assertSearchable(t, queries, scuteSearch)
}
