package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/w8rbt/bloom"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
)

var fp = .001                                                                      // False Positive Rate (1 in 1,000)
var n = 517238891.0                                                                // Hash list cardinality
var m = math.Ceil((n * math.Log(fp)) / math.Log(1.0/math.Pow(2.0, math.Log(2.0)))) // Number of bits in the filter
var k = uint(10)                                                                   // Number of hash functions
var filter = bloom.New(uint(m), k)

var hex = "0123456789ABCDEF"

func hexOnly(hash string) bool {
	for _, c := range hash {
		if !strings.Contains(hex, string(c)) {
			return false
		}
	}
	return true
}

// check - check the bloom filter for the hash
// return 200 if found, 400 on bad request or 404 if not found
// to the client
func check(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Access-Control-Allow-Methods", "HEAD")
	w.Header().Set("Content-Security-Policy", "default-src 'self';")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "same-origin")

	vars := mux.Vars(r)
	hash := strings.ToUpper(vars["hash"])

	if len(hash) != 16 || !hexOnly(hash) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	} else if filter.Test([]byte(hash)) {
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

// index - have a blank index page (Docker Container)
func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Security-Policy", "default-src 'self';")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "same-origin")
}

func main() {
    // Pick option 1. or 2. to read in the filter.
    // Option 1. is much faster. Download a prebuilt filter here:
    // https://drive.google.com/open?id=1TTpxHrqgp8T7GvLlyR9ooEN1KSJzJb5o

	// Option 1. Get the filter from a local file (Jack)
	f, err := os.Open("/home/jack/filter/hibp3.filter")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Option 2. Get the filter from a URL (Docker Container)
	//url := "https://url.to.the.filter/hibp3.filter"
	//resp, err := http.Get(url)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer resp.Body.Close()

    // Uncomment depending on choice 1. or 2. from above
	bytesRead, err := filter.ReadFrom(f)
	// bytesRead, err := filter.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bytes read from filter: %d\n", bytesRead)

	// create router & routes
	router := mux.NewRouter()
	router.HandleFunc("/hashes/sha1/{hash}", check).Methods("GET")
	router.HandleFunc("/", index).Methods("GET")

	// server routes over TLS
	log.Fatal(http.ListenAndServeTLS("0.0.0.0:9379",
		"/tmp/cert.pem",
		"/tmp/privkey.pem",
		// Don't use ProxyHeaders unless running behind a proxy or a LB.
		handlers.CombinedLoggingHandler(os.Stdout, handlers.ProxyHeaders(router))))
}
