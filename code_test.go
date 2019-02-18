package opendmm

import (
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
)

func TestCodeToString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(code{series: "DUMMY", number: 10}.toString(), "DUMMY-010")
	assert.Equal(code{series: "DUMMY", number: 100}.toString(), "DUMMY-100")
	assert.Equal(code{series: "DUMMY", number: 1000}.toString(), "DUMMY-1000")
	assert.Equal(code{series: "DUMMY", prefix: "S", number: 1}.toString(), "DUMMY-S01")
	assert.Equal(code{series: "DUMMY", prefix: "S", number: 10}.toString(), "DUMMY-S10")
	assert.Equal(code{series: "DUMMY", prefix: "S", number: 100}.toString(), "DUMMY-S100")
}

func TestCodeVariations(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(
		code{series: "MIDE", number: 29}.variations(),
		mapset.NewSet(
			"MIDE-029",
			"MIDE-0029",
			"MIDE-00029",
		),
	)
	assert.Equal(
		code{series: "3DSVR", number: 106}.variations(),
		mapset.NewSet(
			"3DSVR-106",
			"3DSVR-0106",
			"3DSVR-00106",
		),
	)
	assert.Equal(
		code{series: "XV", number: 1001}.variations(),
		mapset.NewSet(
			"XV-1001",
			"XV-01001",
		),
	)
	assert.Equal(
		code{series: "MKBD", prefix: "S", number: 96}.variations(),
		mapset.NewSet(
			"MKBD-S96",
		),
	)
	assert.Equal(
		code{series: "MKBD", prefix: "S", number: 100}.variations(),
		mapset.NewSet(
			"MKBD-S100",
		),
	)
}

func TestGuessCodes(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(
		guessCodes("MIDE-029"),
		mapset.NewSet(
			code{
				series: "MIDE",
				number: 29,
			},
		),
	)
	assert.Equal(
		guessCodes("XV-1001"),
		mapset.NewSet(
			code{
				series: "XV",
				number: 1001,
			},
		),
	)
	assert.Equal(
		guessCodes("IPZ687"),
		mapset.NewSet(
			code{
				series: "IPZ",
				number: 687,
			},
		),
	)
	assert.Equal(
		guessCodes("MMGH00010"),
		mapset.NewSet(
			code{
				series: "MMGH",
				number: 10,
			},
		),
	)
	assert.Equal(
		guessCodes("140c02202"),
		mapset.NewSet(
			code{
				series: "C",
				number: 2202,
			},
			code{
				series: "140C",
				number: 2202,
			},
		),
	)
	assert.Equal(
		guessCodes("3DSVR-100"),
		mapset.NewSet(
			code{
				series: "3DSVR",
				number: 100,
			},
		),
	)
	assert.Equal(
		guessCodes("200GANA-894"),
		mapset.NewSet(
			code{
				series: "200GANA",
				number: 894,
			},
			code{
				series: "GANA",
				number: 894,
			},
		),
	)
	assert.Equal(
		guessCodes("CW3D2DBD-30"),
		mapset.NewSet(
			code{
				series: "CW3D2DBD",
				number: 30,
			},
		),
	)
	assert.Equal(
		guessCodes("MKBD-S97"),
		mapset.NewSet(
			code{
				series: "MKBD",
				prefix: "S",
				number: 97,
			},
		),
	)
	assert.Equal(
		guessCodes("CWPBD_77"),
		mapset.NewSet(
			code{
				series: "CWPBD",
				number: 77,
			},
		),
	)
	assert.Equal(
		guessCodes("140c02202 3DSVR-100"),
		mapset.NewSet(
			code{
				series: "C",
				number: 2202,
			},
			code{
				series: "140C",
				number: 2202,
			},
			code{
				series: "3DSVR",
				number: 100,
			},
		),
	)
}
