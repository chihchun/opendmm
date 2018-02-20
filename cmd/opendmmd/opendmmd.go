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
)

func searchHandler(timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		metach := opendmm.Search(q)
		select {
		case meta, ok := <-metach:
			if ok {
				metajson, _ := json.MarshalIndent(meta, "", "  ")
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
	timeout := flag.Duration("timeout", 30*time.Second, "Timeout of Search API")
	flag.Parse()

	http.HandleFunc("/search", searchHandler(*timeout))
	http.HandleFunc("/guess", guessHandler)
	glog.Fatal(http.ListenAndServe(":"+port, nil))
}
