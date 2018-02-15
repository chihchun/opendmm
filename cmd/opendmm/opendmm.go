package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/libredmm/opendmm"
)

func search(query string, timeout time.Duration) {
	metach := opendmm.Search(query)
	select {
	case meta, ok := <-metach:
		if ok {
			metajson, _ := json.MarshalIndent(meta, "", "  ")
			fmt.Println(string(metajson))
		} else {
			glog.Exit("Not found")
		}
	case <-time.After(timeout):
		glog.Fatal("Timeout")
	}
}

func main() {
	flag.Set("stderrthreshold", "FATAL")
	timeout := flag.Duration("timeout", 30*time.Second, "timeout")
	flag.Parse()

	switch flag.Arg(0) {
	case "search":
		search(flag.Arg(1), *timeout)

	case "guess":
		for keyword := range opendmm.Guess(flag.Arg(1)).Iter() {
			fmt.Println(keyword)
		}

	default:
		search(flag.Arg(0), *timeout)
	}
}
