package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	list := []string{"this", "is", "data", "to", "stream"}
	w.Header().Add("Content-Type", "text/event-stream")
	for i := 0; i < len(list); i += 1 {
		fmt.Fprintf(w, "data: %s\n\n", list[i])
		w.(http.Flusher).Flush()
		time.Sleep(2 * time.Second)
	}
}
