package opendmm

import (
	"testing"
)

func TestAventertainments(t *testing.T) {
	queries := []string{
		"SKYHD-001",
		"MCB3DBD-25",
		"lafbd-10",
		"CW3D2DBD-30",
		"MKBD-S97",
		"LAFBD-030",
		"BD-M08",
		"CWPBD_77",
	}
	assertSearchable(t, queries, aveSearch)
}
