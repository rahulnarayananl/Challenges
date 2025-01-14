package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", LogRequest)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func LogRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request from %s \n", r.RemoteAddr)
	fmt.Printf("%s %s %s \n", r.Method, r.Pattern, r.Proto)
	fmt.Printf("Host : %s \n", r.Host)
	for h, v := range r.Header {
		fmt.Printf("%s: %s\n", h, v)
	}
	fmt.Println()
}
