package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "svc-a ok\nmethod=%s\npath=%s\ntime=%s\n", r.Method, r.URL.Path, time.Now().UTC().Format(time.RFC3339))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}
	log.Printf("svc-a listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
