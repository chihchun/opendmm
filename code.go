package opendmm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

type code struct {
	series string
	prefix string
	number int
}

func (c code) toString() string {
	if c.prefix == "" {
		return fmt.Sprintf("%s-%03d", c.series, c.number)
	}
	return fmt.Sprintf("%s-%s%02d", c.series, c.prefix, c.number)
}

func (c code) variations() mapset.Set {
	if c.prefix == "" {
		return mapset.NewSet(
			fmt.Sprintf("%s-%03d", c.series, c.number),
			fmt.Sprintf("%s-%04d", c.series, c.number),
			fmt.Sprintf("%s-%05d", c.series, c.number),
		)
	}
	return mapset.NewSet(
		fmt.Sprintf("%s-%s%02d", c.series, c.prefix, c.number),
	)
}

func guessCodes(query string) mapset.Set {
	re := regexp.MustCompile(
		"(?i)(\\d{3})?((?:3d|2d|s2|[a-z]){1,7}?)[-_]?([sm]?)(0*(\\d{2,5}))",
	)
	matches := re.FindAllStringSubmatch(query, -1)
	codes := mapset.NewSet()
	for _, match := range matches {
		n, err := strconv.Atoi(match[4])
		if err == nil {
			codes.Add(code{
				series: strings.ToUpper(match[2]),
				prefix: strings.ToUpper(match[3]),
				number: n,
			})
			codes.Add(code{
				series: strings.ToUpper(match[1] + match[2]),
				prefix: strings.ToUpper(match[3]),
				number: n,
			})
		}
	}
	return codes
}
