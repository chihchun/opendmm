package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/libredmm/opendmm"
	"github.com/syndtr/goleveldb/leveldb"
)

func searchHandler(timeout time.Duration, db *leveldb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		metajson, err := db.Get([]byte(q), nil)
		if err == nil {
			glog.Infof("Hit cache: %s", q)
			fmt.Fprintf(w, string(metajson))
			return
		}

		metach := opendmm.Search(q)
		select {
		case meta, ok := <-metach:
			if ok {
				metajson, _ := json.MarshalIndent(meta, "", "  ")
				db.Put([]byte(q), metajson, nil)
				fmt.Fprintf(w, string(metajson))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		case <-time.After(timeout):
			w.WriteHeader(http.StatusGatewayTimeout)
		}
	}
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	keywords := opendmm.Guess(q)
	if keywords.Cardinality() == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		json, _ := json.MarshalIndent(keywords, "", "  ")
		fmt.Fprintf(w, string(json))
	}
}

func main() {
	port := os.Getenv("PORT")

	flag.Set("stderrthreshold", "FATAL")
	timeout := flag.Duration("timeout", 60*time.Second, "Timeout of Search API")
	flag.Parse()

	db, err := leveldb.OpenFile("opendmmd.leveldb", nil)
	if err != nil {
		glog.Fatal("Failed to open DB")
	}
	defer db.Close()

	http.HandleFunc("/search", searchHandler(*timeout, db))
	http.HandleFunc("/guess", guessHandler)
	glog.Fatal(http.ListenAndServe(":"+port, nil))
}
