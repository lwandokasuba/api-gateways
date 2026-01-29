package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type target struct {
	name string
	path string
}

func main() {
	baseURL := "http://localhost"
	services := []string{"alpha.localhost", "beta.localhost", "gamma.localhost"}
	paths := []string{"/", "/health", "/error", "/missing"}

	client := &http.Client{Timeout: 5 * time.Second}

	var wg sync.WaitGroup
	for _, service := range services {
		for _, path := range paths {
			wg.Add(1)
			go func(service, path string) {
				defer wg.Done()
				call(client, baseURL, service, path)
			}(service, path)
		}
	}

	// Add a few delayed calls to see non-instant responses.
	for _, service := range services {
		wg.Add(1)
		go func(service string) {
			defer wg.Done()
			call(client, baseURL, service, "/delay")
		}(service)
	}

	wg.Wait()
}

func call(client *http.Client, baseURL, host, path string) {
	req, err := http.NewRequest(http.MethodGet, baseURL+path, nil)
	if err != nil {
		fmt.Printf("%s %s -> request error: %v\n", host, path, err)
		return
	}
	req.Host = host

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s %s -> transport error: %v\n", host, path, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("%s %s -> %d %s\n", host, path, resp.StatusCode, truncate(body, 120))
}

func truncate(b []byte, limit int) string {
	if len(b) <= limit {
		return string(b)
	}
	return string(b[:limit]) + "..."
}
