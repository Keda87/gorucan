package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func httpGetter(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Inject random failure (30% chance)
	if rand.Float32() < 0.3 {
		return "", fmt.Errorf("simulated random failure")
	}

	return string(bytes), nil
}

func main() {
	var (
		wg           sync.WaitGroup
		start        = time.Now()
		successCount int32
		failureCount int32
		urls         = []string{
			"https://httpbin.org/anything/hi",
			"https://httpbin.org/anything/hello",
			"https://httpbin.org/anything/page",
			"https://httpbin.org/anything/example",
			"https://httpbin.org/anything/birthday",
			"https://httpbin.org/anything/greeter",
		}
	)

	for i := 0; i < len(urls); i++ {
		url := urls[i]

		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			res, err := httpGetter(u)
			if err != nil {
				fmt.Println("FAILED: ", u)
				atomic.AddInt32(&failureCount, 1)
			} else {
				fmt.Println("SUCCESS: ", u, " with data: ", res)
				atomic.AddInt32(&successCount, 1)
			}
		}(url)
	}
	wg.Wait()

	fmt.Println("total url: ", len(urls))
	fmt.Println("total success: ", successCount)
	fmt.Println("total failure: ", failureCount)
	fmt.Println("execution: ", time.Since(start))
}
