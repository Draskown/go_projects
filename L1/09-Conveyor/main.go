package main

import (
	"fmt"
	"time"
)

func main() {
	// Create conveyor's channels
	// for sending X values
	// and X^2
	chX := make(chan int)
	chXTwo := make(chan int)

	// Run X sender
	go sendX(chX)
	// Run X^2 sender
	go printXPowTwo(chX, chXTwo)

	// X^2 values receiver
	for v := range chXTwo {
		fmt.Printf("X^2: %d\n", v)
	}

	fmt.Println("Finished main")
}

// Sends X values from a slice into the c channel
func sendX(c chan int) {
	defer close(c)

	// Slice of data to send
	arr := []int{12, 11, 27, 13, 54, 52, 62221, 790}

	// Loop through every slice elements
	for _, v := range arr {
		// Send as X
		c <- v
		time.Sleep(333 * time.Millisecond)
	}

	fmt.Println("Finished X conveyor")
}

// Receives X values from the cx channel
// and sends X^2 values into the cxtwo channel
func printXPowTwo(cx chan int, cxtwo chan int) {
	defer close(cxtwo)

	// Recieve X data
	for v := range cx {
		// And send X^2 to the other channel
		cxtwo <- v * v
	}
	fmt.Println("Finished X^2 conveyor")
}
