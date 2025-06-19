package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func httpGetter(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil { // check if http request error.
		if ctx.Err() != nil { // check if the error caused by context timeout.
			return "", ctx.Err()
		}
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func doLongRunningTask(ctx context.Context, jobID string) error {
	fmt.Println("starting: ", jobID)

	_, err := httpGetter(ctx, "https://httpbin.org/anything/service-1")
	if err != nil {
		return err
	}

	_, err = httpGetter(ctx, "https://httpbin.org/anything/service-2")
	if err != nil {
		return err
	}

	_, err = httpGetter(ctx, "https://httpbin.org/anything/service-3")
	if err != nil {
		return err
	}

	_, err = httpGetter(ctx, "https://httpbin.org/anything/service-4")
	if err != nil {
		return err
	}

	_, err = httpGetter(ctx, "https://httpbin.org/anything/service-5")
	if err != nil {
		return err
	}

	_, err = httpGetter(ctx, "https://httpbin.org/anything/service-6")
	if err != nil {
		return err
	}

	fmt.Println("finished: ", jobID)
	return nil
}

func main() {
	var (
		start   = time.Now()
		wg      sync.WaitGroup
		chanerr = make(chan error, 2) // use buffered channel as many goroutine number.
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := doLongRunningTask(ctx, "job-1"); err != nil {
			chanerr <- fmt.Errorf("job 1: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := doLongRunningTask(ctx, "job-2"); err != nil {
			chanerr <- fmt.Errorf("job 2: %v", err)
		}
	}()

	wg.Wait()
	close(chanerr) //  always close chanerr after all goroutine finished.

	for e := range chanerr {
		fmt.Println(e.Error())
	}

	fmt.Println("execution: ", time.Since(start))
}
