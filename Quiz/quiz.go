package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Main function of the program
func main(){
	// Create custom flags of a file name, shuffle and time limit
	fileName := flag.String("csv", "problems.csv", "a csv file in format of 'question,answer'")
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz")
	timeLimit := flag.Int("time", 30, "time limit for the quiz")

	// Parse the flags
	flag.Parse()

	// Open the file provided
	file, err := os.Open(*fileName)

	// If the file does not open
	if err != nil {
		// Exit the program with an error message
		exit(fmt.Sprintf("Failed to open file: %s.", *fileName))
	}

	// Else - create a reader for the file
	r := csv.NewReader(file)

	// Read all lines from the file
	lines, err := r.ReadAll()

	// If the reader failed to read the file
	if err != nil{
		// Exit the program with an error message
		exit("Failed to parse the file.")
	}

	// Parse the read lines into a problem structure
	problems := parseLines(lines)

	// If shuffle flag was provided as true
	if *shuffle{
		// Set a random seed
		rand.Seed(time.Now().UnixNano())
		// Shuffle the array
		rand.Shuffle(len(problems), func(i, j int) { 
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	// Create a timer with specified time in a flag
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	// Correct count is set to zero
	correct := 0

	// Go through all of the problems
	for i, p := range problems {
		// Print the current problem
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		// Create a channel for a goroutine to listen for answers
		answerC := make(chan string)
		
		// Create a goroutine to listen for answers
		go func(){
			// Wait for a user to print their answer
			var answer string
			fmt.Scanln(&answer)
			// Send the read answer to the channel
			answerC <- answer
		}()

		// Handle the channels
		select {
		// If timer has sent a message (that it was done)
		case <-timer.C:
			// Exit the program with the result
			fmt.Printf("Your time has expired.\n%d/%d\n", correct, len(problems))
			return
		// If the answer channel has some message in it
		case answer := <-answerC:
			// Check if the received answer is correct
			if answer == p.a {
				// If it is - add it to the counter
				correct++
			}
		}
	}

	// Create an additional string for the message to display
	// whether the user's done well or not
	addString := ""
	if correct == 0{
		addString = "Poor job"
	} else if correct > 0 && correct < 5{
		addString = "Try a little harder"
	} else if correct >= 5 && correct <8{
		addString = "You've done alright"
	} else if correct >= 8{
		addString = "Excellent job"
	}

	// Print out the result
	fmt.Printf(addString+" %d/%d", correct, len(problems))
}

// Custom structure of a problem
// where q is a question, a - an answer
type problem struct{
	q string
	a string
}

// Parses the read lines from csv file
// into an array of custom structured problems
func parseLines(input [][]string) []problem{
	// Make a new slice for the output
	res := make([]problem, len(input))

	// Read through the each line
	for i, line := range input {
		// Initialize a structure
		res[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	// Return the output
	return res
}

// Exits the program with a given message
func exit(msg string){
	fmt.Println(msg)
	os.Exit(1)
}