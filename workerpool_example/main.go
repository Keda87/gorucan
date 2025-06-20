package main

import (
	"fmt"
	"strings"
	"time"
)

// jobs <-chan string   : it's a channel for receiving data only (reading channel).
// result chan<- string : it's a channel for producer (publish data).
func workerStringUpperCase(id int, jobs <-chan string, result chan<- string) {
	for job := range jobs {
		fmt.Println("Worker: ", id, " started job")
		fmt.Println("Processing to convert: ", job)
		time.Sleep(time.Second)
		fmt.Println("Worker: ", id, " finished job")
		result <- strings.ToUpper(job)
	}
}

func main() {
	var (
		words = []string{
			"testing",
			"goal",
			"journey",
			"acceptance",
			"anger",
			"flexible",
			"terror",
			"max",
		}

		start        = time.Now()
		numberOfJobs = len(words)

		// channel capacity is the same capacity with the array length. (buffered channel)
		jobs   = make(chan string, numberOfJobs)
		result = make(chan string, numberOfJobs)
	)

	// spawn 3 workers to standby processing incoming task.
	for i := 1; i <= 3; i++ {
		go workerStringUpperCase(i, jobs, result)
	}

	// send or publish all word to jobs channel to be process.
	for _, word := range words {
		jobs <- word
	}

	// close jobs channel after all words are published.
	// so jobs channel no longer accept incoming data, then
	// just waiting all the word processed by the worker.
	close(jobs)

	// printing all worker result based on given words.
	for i := 0; i < len(words); i++ {
		<-result
	}
	fmt.Println("execution: ", time.Since(start))
}
