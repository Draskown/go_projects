package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Set up a new random seed from current time
	randomiser := rand.New(
		rand.NewSource(
			time.Now().UnixNano(),
		),
	)

	// Prompt for the N of workers
	fmt.Println("Enter N of workers: ")
	var N int

	// Scan user input into the variable
	_, err := fmt.Scanln(&N)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Add N items to the waitgroup
	// waiting for them to finish
	wg := sync.WaitGroup{}
	wg.Add(N)

	// Make a data channel
	dataChan := make(chan int, N)

	// Init N workers
	for i := 0; i < N; i++ {
		go worker(dataChan, i, &wg)
	}

	// Channel waiting for SIGTERM and SIGINT signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// Eternal loop for sending data
	// to the data channel
	for {
		select {
		// If quit has a signal
		case <-quit:
			// Close the channel to stop
			// all goroutines
			close(dataChan)
			wg.Wait()
			os.Exit(0)
		// If no signal was sent
		default:
			// Send data to the channel
			dataChan <- randomiser.Intn(1000)
		}
	}
}

// Start a eternal loop for reading from the
// data channel
func worker(c chan int, num int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// If channel is not yet closed,
		// print it out and wait a second
		if val, ok := <-c; ok {
			fmt.Printf("Worker %d: %d\n", num, val)
			time.Sleep(1 * time.Second)
		} else {
			// Otherwise exit the func
			fmt.Printf("Worker %d has been stopped\n", num)
			return
		}
	}
}
