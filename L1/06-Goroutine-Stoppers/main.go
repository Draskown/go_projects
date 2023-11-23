package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gopkg.in/tomb.v2"
)

func main() {
	// stopBySendingChannel()
	// stopByClosingChannel()
	// stopByContext()
	stopByTomb()
}

// Uses an additional channel to close the goroutine
// if a values has been sent to it
func stopBySendingChannel() {
	// Create a Wait Group for synchronisation
	wg := sync.WaitGroup{}
	// Add one goroutine to the Wait Group
	wg.Add(4)

	// Create the channel that will stop the goroutine
	quit := make(chan bool)
	// Run the goroutine
	go func(q chan bool, g *sync.WaitGroup) {
		defer g.Done()

		// Eternal loop to mimic some calculations
		for {
			select {
			// When signal to the quit channel has been recieved
			case <-q:
				// Stop the goroutine
				fmt.Println("Goroutine has been stopped")
				return
			default:
				// Elsewise, continue operating
				fmt.Println("Goroutine is running")
				time.Sleep(1 * time.Second)
			}
		}
	}(quit, &wg)

	// Operate the goroutine for 5 seconds
	time.Sleep(5 * time.Second)
	// And then shut it down by sending a signal
	quit <- true

	// Wait for the goroutine to finish
	wg.Wait()
	fmt.Println("End of main")
}

// Uses a main data channel to close the goroutine after
// sending has been completed
func stopByClosingChannel() {
	// Create main data channel
	ch := make(chan int)

	// Create a Wait Group for synchronisation
	wg := sync.WaitGroup{}
	// Add one goroutine to the Wait Group
	wg.Add(1)

	// Run the goroutine
	go func(c chan int, g *sync.WaitGroup) {
		defer g.Done()

		// Read data from the main channel
		// as long as there is any
		for v := range c {
			fmt.Println(v)
		}
		// If no data - exit the goroutine
		fmt.Println("Goroutine has been stopped")
		return
	}(ch, &wg)

	// Send some data through the main data channel
	for _, v := range []int{1, 2, 3, 4, 5} {
		ch <- v
		time.Sleep(333 * time.Millisecond)
	}

	// Close the channel after all data has been sent
	close(ch)
	// Wait for the channel to close
	wg.Wait()
	fmt.Println("End of main")
}

// Uses in-built context package to send the signal
// to the goroutine
//
// Advantage of the context compared to the closing
// by channel signal is that multiple goroutines can
// receive this signal as they operate a copy of the
// main context structure instead of reading
// from one channel once for each goroutine
func stopByContext() {
	// Create Wait Group for synchronisation
	wg := sync.WaitGroup{}
	// Add one goroutine to the Wait Group
	wg.Add(1)
	// Create a context with cancel option to send signals
	ctx, cancel := context.WithCancel(context.Background())

	// Run the goroutine
	go func(ctx context.Context, g *sync.WaitGroup) {
		defer g.Done()

		// Eternal loop to mimic some calculations
		for {
			select {
			// If context has received cancel signal
			case <-ctx.Done():
				// Stop the goroutine
				fmt.Println("Goroutine has been stopped")
				return
			default:
				// Elsewise, do some work
				fmt.Println("Goroutine is running")
				time.Sleep(1 * time.Second)
			}
		}
	}(ctx, &wg)

	// Wait for some time before sending the sinal
	// to the context
	time.Sleep(5 * time.Second)
	cancel()

	// Wait fot the goroutine to finish
	wg.Wait()
	fmt.Println("This is the end of the main func")
}

// Uses a side tomb package to handle
// and terminate goroutines
//
// Combines both synchronisation and multiple
// goroutines operations, so no additional channel
// nor a Wait Group is needed
//
// The disadvantage of this approach is that
// the goroutine handled by Tomb can be only
// func() error, so no variables or channels
// can be passed to it and there is no way to
// share data between those goroutines.
//
// Although the tomb functions can be encapsulated
// into a function of a greater scope to resolve this
func stopByTomb() {
	// Create a new tomb instance
	var tb tomb.Tomb

	// Define a function that will be
	// passed to the tomb instance
	f := func() error {
		for {
			select {
			// When tomb receives a Kill signal
			case <-tb.Dying():
				// Exit the goroutine
				fmt.Println("Goroutine has been stopped")
				return nil
			default:
				// Elsewise, perform some work
				fmt.Println("Goroutine is running")
				time.Sleep(1 * time.Second)
			}
		}
	}

	// Add a function to the tomb
	// that will be launched in parallel
	tb.Go(f)

	// Wait for some work to be done
	time.Sleep(5 * time.Second)
	// Send a Kill signal to the tomb
	tb.Kill(fmt.Errorf("End of main"))

	// Wait for the goroutine to finish
	err := tb.Wait()
	fmt.Println(err.Error())
}
