package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Read user input for N seconds to send data
	var N int
	fmt.Println("Enter N seconds: ")
	_, err := fmt.Scanln(&N)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Create a channel for the data
	ch := make(chan int)
	// And run the sender goroutine
	go sendData(ch, N)

	// Print values as long as the channel is open
	for val := range ch {
		fmt.Println(val)
	}

	fmt.Println("Finished main")
}

// Send data for the time provided by user input
func sendData(c chan int, duration int) {
	// Create timer
	timer := time.NewTimer(time.Duration(duration) * time.Second)

	// Randomise the sent values
	randomiser := rand.New(
		rand.NewSource(
			time.Now().UnixNano(),
		),
	)

	// Sent values in the eternal loop
	for {
		select {
		// If the timer has finished counting
		case <-timer.C:
			// Close the channel and leave sender
			fmt.Println("Finish sending")
			close(c)
			return
		default:
			// Elsewise, send data to the channel
			c <- randomiser.Intn(1000)
			time.Sleep(333 * time.Millisecond)
		}
	}
}
