package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	Dest = ""
)

func main() {
	var src, dst string
	flag.Parse()
	args := flag.Args()
	if len(args) >= 1 {
		dst = args[0]
	} else {
		dst = "127.0.0.1:8080"
	}
	if len(args) == 2 {
		src = args[1]
	} else {
		src = ":80"
	}
	Dest = dst
	log.Printf("Starting...")
	http.HandleFunc("/", configFunc)
	log.Fatal(http.ListenAndServe(src, nil))
}

func configFunc(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(Dest)
	r.Host = u.Host
	pass, _ := u.User.Password()
	r.SetBasicAuth(u.User.Username(), pass)
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &myTransport{}
	proxy.ServeHTTP(w, r)
}

type myTransport struct{}

func (t *myTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)
	body, err := httputil.DumpResponse(response, true)
	if err != nil {
		return nil, err
	}

	log.Print(string(body))

	return response, err
}
