package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	// Create a map to fill in by goroutines
	initMap := make(map[string]int)

	// Slice to take data from
	arr := []string{
		"People",
		"Animals",
		"Aircraft",
		"Cars",
		"Boats",
		"Bots",
		"Boxers",
		"Thai",
		"Chinese",
		"Russians",
		"English",
		"Germans",
	}

	// Use two approaches of filling data
	basic(initMap, arr)
	idiomatic(initMap, arr)
}

// Uses Wait Group and Mutex for data synchronisation
// and keeping the map from being blocked
func basic(initMap map[string]int, arr []string) {
	// Init a Wait Group for data synchronisation
	wg := sync.WaitGroup{}
	// Add a number of goroutines that corresponds
	// to the slice's length
	wg.Add(len(arr))

	// Init a mutex to lock data from being accessed
	// while being operated at the same time
	mut := sync.Mutex{}

	// Loop over slice's elements
	for _, v := range arr {
		// And call a goroutine for each
		go func(goMap map[string]int, s string, goMut *sync.Mutex, g *sync.WaitGroup) {
			// Lock the memory
			goMut.Lock()
			// Executed on goroutine exit, unlocks memory and sends
			// a signal that the goroutine has finished
			defer goMut.Unlock()
			defer g.Done()

			// Set values to the map
			goMap[s] = rand.Intn(1000)
		}(initMap, v, &mut, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Printf("Map by basic approach: %v\n", initMap)
}

// Uses channels to operate data
func idiomatic(initMap map[string]int, arr []string) {
	// Create two channels for a map:
	// one for keys and one for values
	chI := make(chan int)
	chS := make(chan string)

	// Loop over slice's elements
	for _, v := range arr {
		// And call a goroutine for each
		go func(s string, cI chan int, cS chan string) {
			// Send data to the channels
			cS <- s
			cI <- rand.Intn(1000)
		}(v, chI, chS)

		// Recieve data from the channels
		// and write them into the map
		initMap[<-chS] = <-chI
	}

	fmt.Printf("Map by idiomatic approach: %v\n", initMap)
}
