module github.com/libredmm/opendmm

require (
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/andybalholm/cascadia v1.0.0
	github.com/benbjohnson/phantomjs v0.0.0-20170410144507-1b81734f7877
	github.com/deckarep/golang-set v1.7.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db
	github.com/heroku/x v0.0.0-20181102215100-85e5aa5e6aa1
	github.com/junzh0u/httpx v0.0.0-20181022012810-68ca3a6137e9
	github.com/junzh0u/ioutilx v0.0.0-20180827205417-117bc779d863
	github.com/syndtr/goleveldb v0.0.0-20181105012736-f9080354173f
	golang.org/x/net v0.0.0-20181114220301-adae6a3d119a
	golang.org/x/text v0.3.0
)

// +heroku goVersion go1.11
// +heroku install ./cmd/...