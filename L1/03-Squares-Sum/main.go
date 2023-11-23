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

// Uses a Wait Group and a Mutex to calculate
// sum of an array's squares
func basic(arr []int) {
	// Initialise the sum to be calculated
	var sum int

	// Init Waitgroup for syncing the goroutines
	var wg sync.WaitGroup

	// Init mutex to lock sum variable
	var mut sync.Mutex

	// Add a group element for each
	// of the array element
	wg.Add(len(arr))

	// Iterate through each array element
	for _, v := range arr {
		// Run a goroutine for every element
		//
		// Can be done w/o the pointers to the WaitGroup, Mutex and sum
		// because they are initialised in the global scope,
		// but the more practical way is to
		go func(val int, wg *sync.WaitGroup, m *sync.Mutex, s *int) {
			// Lock all the global variables
			mut.Lock()
			// Unlock them after finishing the goroutine
			defer mut.Unlock()

			// Calculate the square of the element
			// Using pointer for the value to change
			*s += val * val

			// Return that the WaitGroup element
			// has finished its work
			wg.Done()
		}(v, &wg, &mut, &sum)
	}

	// Wait for all of the Wait Group elelements to be Done
	wg.Wait()

	fmt.Printf("Sum of the squares by basic approach: %d\n", sum)
}

// Uses a channel calculate
// sum of an array's squares
func idiomatic(arr []int) {
	// Initialise sum channel and variable
	// for alternative calculation
	sumCh := make(chan int, 1)
	var sumAlt int

	// Iterate through each array element
	for _, v := range arr {
		// Or, more idiomatic way is to calculate the sum through channels:
		//
		// No values can be locked, and by syncronising the channel input
		// with output, no WaitGroup is needed
		go func(val int, ch chan int) {
			ch <- val * val
		}(v, sumCh)

		sumAlt += <-sumCh
	}
	fmt.Printf("Sum of the squares by idiomatic approach: %d\n", sumAlt)
}
