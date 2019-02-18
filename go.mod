module github.com/libredmm/opendmm

require (
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/andybalholm/cascadia v1.0.0
	github.com/benbjohnson/phantomjs v0.0.0-20181211182228-6499a20f5cd6
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set v1.7.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db
	github.com/heroku/x v0.0.0-20181102215100-85e5aa5e6aa1
	github.com/junzh0u/httpx v0.0.0-20190216221210-07a21eee7e61
	github.com/junzh0u/ioutilx v0.0.0-20180827205417-117bc779d863
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/syndtr/goleveldb v0.0.0-20190203031304-2f17a3356c66
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd
	golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4 // indirect
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a // indirect
	golang.org/x/text v0.3.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

// +heroku goVersion go1.11
// +heroku install ./cmd/...
