package main

import (
	"fmt"
	"sync"
	"time"
)

func greeter1Second(name string) {
	fmt.Println("Hi: ", name)
	time.Sleep(time.Second)
}

func executeInSerial() {
	greeter1Second("Alex")
	greeter1Second("Kosim")
	greeter1Second("John")
	greeter1Second("Diana")
}

func executeInConcurrent() {
	var wg sync.WaitGroup

	wg.Add(4) // number of task.
	go func() {
		defer wg.Done()
		greeter1Second("Alex")
	}()

	go func() {
		defer wg.Done()
		greeter1Second("Kosim")
	}()

	go func() {
		defer wg.Done()
		greeter1Second("John")
	}()

	go func() {
		defer wg.Done()
		greeter1Second("Dian")
	}()

	// block main goroutine until all goroutine above finished.
	wg.Wait()
}

func main() {
	serialStart := time.Now()
	executeInSerial()
	fmt.Println("Serial execution: ", time.Since(serialStart))
	fmt.Println()

	concurrentStart := time.Now()
	executeInConcurrent()
	fmt.Println("Concurrent execution: ", time.Since(concurrentStart))
}
