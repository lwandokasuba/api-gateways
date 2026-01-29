package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type response struct {
	Service string `json:"service"`
	Path    string `json:"path"`
	Status  string `json:"status"`
}

func main() {
	serviceName := getenv("SERVICE_NAME", "service")
	port := getenv("SERVICE_PORT", "8080")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, response{Service: serviceName, Path: r.URL.Path, Status: "ok"})
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, response{Service: serviceName, Path: r.URL.Path, Status: "forced error"})
	})

	mux.HandleFunc("/delay", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(750 * time.Millisecond)
		writeJSON(w, http.StatusOK, response{Service: serviceName, Path: r.URL.Path, Status: "delayed"})
	})

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           logRequests(mux, serviceName),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("%s listening on :%s", serviceName, port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func logRequests(next http.Handler, serviceName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", serviceName, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
