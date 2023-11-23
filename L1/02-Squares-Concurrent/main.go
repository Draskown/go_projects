package main

import (
	"fmt"
	"sync"
)

func main() {
	// Initialise the array with preset values
	initArray := []int{2, 4, 6, 8, 10}

	// Perform calculations with two approaches
	basic(initArray)
	idiomatic(initArray)

}

// Uses Wait Group to synchronise
func basic(arr []int) {
	// Init Waitgroup for syncing the goroutines
	var wg sync.WaitGroup

	// Add a group element for each
	// of the array element
	wg.Add(len(arr))

	// Iterate through each array element
	for _, v := range arr {

		// Run a goroutine for every element
		go func(val int, wg *sync.WaitGroup) {
			// Calculate the square of the element
			fmt.Printf("Square of %d by basic approach is %d\n", val, val*val)

			// Return that the WaitGroup element
			// has finished its work
			wg.Done()
		}(v, &wg)
	}

	// Wait for all of the Wait Group elelements to be Done
	wg.Wait()
}

// Uses a channel to synchronise between main and goroutines
func idiomatic(arr []int) {
	// Initialise channel for idiomatic approach
	ch := make(chan int, 1)

	// More idiomatic approach is through a channel
	// No need for a WaitGroup as channel's inputs
	// are synchronised with outputs
	for _, v := range arr {
		go func(val int, c chan int) {
			c <- v * v
		}(v, ch)

		fmt.Printf("Square of %d by idiomatic approach is %d\n", v, <-ch)
	}
}
