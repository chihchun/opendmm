package opendmm

import (
	"regexp"
	"strconv"

	"github.com/golang/glog"
)

func isCodeEqual(lcode, rcode string) bool {
	re := regexp.MustCompile("(?i)([a-z]+)-(\\d+)")
	lmeta := re.FindStringSubmatch(lcode)
	rmeta := re.FindStringSubmatch(rcode)
	glog.Info(lmeta)
	glog.Info(rmeta)
	if lmeta == nil || rmeta == nil {
		return false
	}
	if lmeta[1] != rmeta[1] {
		return false
	}
	lnum, err := strconv.Atoi(lmeta[2])
	if err != nil {
		glog.Errorf("[DMM] %s", err)
		return false
	}
	rnum, err := strconv.Atoi(rmeta[2])
	if err != nil {
		glog.Errorf("[DMM] %s", err)
		return false
	}
	return lnum == rnum
}
