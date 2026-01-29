package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type resp struct {
	Service string `json:"service"`
	Method  string `json:"method"`
	Path    string `json:"path"`
	Time    string `json:"time"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		payload := resp{
			Service: "svc-b",
			Method:  r.Method,
			Path:    r.URL.Path,
			Time:    time.Now().UTC().Format(time.RFC3339),
		}
		_ = json.NewEncoder(w).Encode(payload)
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}
	log.Printf("svc-b listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
