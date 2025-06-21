package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const numberOfWorker = 5

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

func consumer(ctx context.Context, workerID int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for email := range jobs {
		fmt.Println("worker:", workerID, "processing", email)

		// no need to check cancellation using select.
		// because consumer and httpGetter is already context aware.
		// if parents context canceled, and this function automatically cancel.
		url := "https://httpbin.org/anything/" + email
		_, err := httpGetter(ctx, url)
		if err != nil {
			fmt.Println("error processing:", email)
			continue
		}

		fmt.Println(email, "processed successfully")
	}
}

func producer(ctx context.Context, jobs chan<- string, emails []string) {
	defer close(jobs)

	for _, email := range emails {
		// check cancellation before publish data to a channel.
		select {
		case <-ctx.Done():
			return
		case jobs <- email:
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()
	emails := []string{
		"jason.miller93@example.com",
		"emily.sanders42@example.org",
		"kevin.torres87@example.net",
		"sophia.wilson18@example.com",
		"liam.johnson55@example.org",
		"olivia.brown23@example.net",
		"ethan.martinez67@example.com",
		"ava.garcia12@example.org",
		"noah.anderson99@example.net",
		"isabella.thomas31@example.com",
		"mason.jackson73@example.org",
		"mia.white88@example.net",
		"logan.harris04@example.com",
		"amelia.martin29@example.org",
		"lucas.thompson16@example.net",
		"charlotte.moore65@example.com",
		"elijah.walker77@example.org",
		"harper.hall39@example.net",
		"james.lewis52@example.com",
		"abigail.young83@example.org",
		"benjamin.allen90@example.net",
		"ella.king08@example.com",
		"sebastian.wright46@example.org",
		"avery.scott19@example.net",
		"henry.green35@example.com",
		"scarlett.baker58@example.org",
		"jack.adams27@example.net",
		"grace.nelson64@example.com",
		"alexander.hill11@example.org",
		"lily.rivera70@example.net",
	}

	// single channel for producer and consumer.
	// make it buffered channel because it has multiple consumers.
	// buffer size is equal to number of worker.
	jobs := make(chan string, numberOfWorker)

	// producer sending email data thru channel.
	go producer(ctx, jobs, emails)

	var wg sync.WaitGroup

	for i := 0; i < numberOfWorker; i++ {
		wg.Add(1)
		go consumer(ctx, i, jobs, &wg)
	}

	wg.Wait()
	fmt.Println("\nexecution: ", time.Since(start))
}
