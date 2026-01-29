package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var (
		baseURL  = flag.String("base", "http://localhost:49000", "Kong proxy base URL")
		routes   = flag.String("routes", "/a,/b", "Comma-separated route paths")
		workers  = flag.Int("workers", 20, "Number of concurrent workers")
		rps      = flag.Int("rps", 200, "Approx requests per second total")
		duration = flag.Duration("duration", 30*time.Second, "How long to run")
	)
	flag.Parse()

	paths := splitComma(*routes)
	if len(paths) == 0 {
		log.Fatal("no routes provided")
	}

	client := &http.Client{Timeout: 3 * time.Second}
	var total uint64
	var errors uint64

	stop := time.After(*duration)
	interval := time.Second / time.Duration(max(1, *rps))
	if interval < time.Millisecond {
		interval = time.Millisecond
	}

	log.Printf("loadgen starting: base=%s routes=%v workers=%d rps=%d duration=%s", *baseURL, paths, *workers, *rps, *duration)

	jobs := make(chan string, 1024)
	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for path := range jobs {
				url := *baseURL + path
				resp, err := client.Get(url)
				if err != nil {
					atomic.AddUint64(&errors, 1)
					continue
				}
				_ = resp.Body.Close()
				atomic.AddUint64(&total, 1)
			}
		}(i)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			path := paths[rand.Intn(len(paths))]
			select {
			case jobs <- path:
			default:
				// drop if workers are saturated
			}
		}
	}()

	<-stop
	close(jobs)
	wg.Wait()

	log.Printf("done: total=%d errors=%d", total, errors)
}

func splitComma(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if i > start {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	return out
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
