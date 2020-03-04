package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	flagBindTo = flag.String("bind-to", ":8000", "bind to ip:port")
)

func main() {
	flag.Parse()

	router := mux.NewRouter()

	router.HandleFunc("/register/backend/{name}/{port}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ping(vars["name"], vars["port"])
	})

	router.HandleFunc("/{name}/{file:.*}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		url := getURL(vars["name"], vars["file"])

		if url == nil {
			http.NotFound(w, r)
			return
		}

		response := fetch(*url)
		if response.err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		for k, v := range response.headers {
			for _, vv := range v {
				w.Header().Set(k, vv)
			}
		}
		w.Header().Set("Content-Lenght", fmt.Sprintf("%d", len(response.body)))
		w.Write(response.body)
	})

	http.ListenAndServe(*flagBindTo, router)
}

type Response struct {
	body    []byte
	headers http.Header
	code    int
	err     error
}

func fetch(url string) *Response {

	log.Printf("fetching: %v", url)

	hc := http.Client{Timeout: 10 * time.Second}

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "iptv/1.0")

	response, err := hc.Do(request)
	if err != nil {
		return &Response{
			err: err,
		}
	}

	if response.StatusCode/100 != 2 {
		return &Response{
			err:  fmt.Errorf("Invalid status code: %v", response.StatusCode),
			code: response.StatusCode,
		}
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &Response{
			err:  err,
			code: response.StatusCode,
		}
	}
	return &Response{
		body:    b,
		headers: response.Header.Clone(),
		code:    response.StatusCode,
	}
}
