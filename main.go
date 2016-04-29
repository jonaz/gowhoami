package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8080", "Port to listen on")
}

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	log.Println("Starting server on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RemoteAddr: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "Host: %s\n", r.Host)
	fmt.Fprintf(w, "Protocol: %s\n", r.Proto)

	keys := make([]string, 0, len(r.Header))
	for key := range r.Header {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "Headers:\n")
	for _, v := range keys {
		fmt.Fprintf(w, "%s:%s\n", v, r.Header.Get(v))
	}

	log.Printf("| %s | %s | %s", r.RemoteAddr, r.Method, r.RequestURI)
}
