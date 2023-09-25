package main

import (
	"fmt"
	"time"
	"math/rand"
	. "QuizProject/Quiz"
)

// Main function of the program
func main(){
	fileName, shuffle, timeLimit := ReadArguments()

	// Open the file provided
	file, err := OpenFile(fileName)

	// If the file does not open
	if err != nil {
		// Exit the program with an error message
		Exit(fmt.Sprintf("Failed to open file: %s. Error: %s", fileName, err.Error()))
	}

	// Get the lines from the CSV file
	problems, err := ReadCSV(file)

	// If there was an error while reading the file
	if err != nil {
		Exit(err.Error())
	}

	// If shuffle flag was provided as true
	if shuffle{
		// Set a random seed
		rand.NewSource(time.Now().UnixNano())
		// Shuffle the array
		rand.Shuffle(len(problems), func(i, j int) { 
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	// Ask the problems from the file
	score, err := AskQuestions(problems, timeLimit)

	// If there was an error while asking the questions
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	// Create an additional string for the message to display
	// whether the user's done well or not
	addString := ""
	
	percentage := float64(score) / float64(len(problems))

	if percentage == 0{
		addString = "Poor job"
	} else if percentage > 0 && percentage <= 0.5{
		addString = "Try a little harder"
	} else if percentage > 0.5 && percentage <= 0.9{
		addString = "You've done alright"
	} else if percentage >= 0.9{
		addString = "Excellent job"
	}

	// Print out the result
	fmt.Printf(addString+" %d/%d\n", score, len(problems))
}