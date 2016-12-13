package main

import (
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/gorilla/handlers"
)

var pages = []string{
	"/section-a",
	"/section-a/subsection-a-b",
	"/section-a/",

	"/section-b",
	"/section-b/subsection-b-b",
	"/section-b/",

	"/",
}

func main() {
	logfd, err := os.Create("access.log")
	if err != nil {
		log.Fatal(err)
	}
	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	loggedRouter := handlers.LoggingHandler(logfd, r)

	ts := httptest.NewServer(loggedRouter)
	for {
		i := rand.Intn(len(pages))
		target := pages[i]

		res, err := http.Get(ts.URL + target)
		if err != nil {
			log.Fatal(err)
		}

		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()

		// s := time.Duration(rand.Intn(1500)) * time.Millisecond
		// s := time.Duration(rand.Intn(50)) * time.Millisecond
		s := time.Millisecond
		time.Sleep(s)
	}
}
