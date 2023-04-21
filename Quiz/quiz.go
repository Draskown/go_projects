package main

import (
	"bufio"
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
	fileName, shuffle, timeLimit := readArguments()

	// Open the file provided
	file, err := openFile(fileName)

	// If the file does not open
	if err != nil {
		// Exit the program with an error message
		exit(fmt.Sprintf("Failed to open file: %s. Error: %s", fileName, err.Error()))
	}

	// Get the lines from the CSV file
	lines, err := readCSV(file)

	// If the reader failed to read the file
	if err != nil{
		// Exit the program with an error message
		exit(fmt.Sprintf("Failed to parse the file. Error %s", err.Error()))
	}

	// Parse the read lines into a problem structure
	problems := parseLines(lines)

	// If shuffle flag was provided as true
	if shuffle{
		// Set a random seed
		rand.Seed(time.Now().UnixNano())
		// Shuffle the array
		rand.Shuffle(len(problems), func(i, j int) { 
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	// Ask the problems from the file
	score, err := askQuestions(problems, timeLimit)

	// If there was an error while asking the questions
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	// Create an additional string for the message to display
	// whether the user's done well or not
	addString := ""
	if score == 0{
		addString = "Poor job"
	} else if score > 0 && score < 5{
		addString = "Try a little harder"
	} else if score >= 5 && score <8{
		addString = "You've done alright"
	} else if score >= 8{
		addString = "Excellent job"
	}

	// Print out the result
	fmt.Printf(addString+" %d/%d\n", score, len(problems))
}

// Opens the file from given directory
func openFile(name string) (*os.File, error) {
	return os.Open(name)
}

// Reads argumetns from a command line
func readArguments() (string, bool, int) {
	// Create custom flags of a file name, shuffle and time limit
	filename := flag.String("csv", "problems.csv", "a csv file in format of 'question,answer'")
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz")
	timelimit := flag.Int("time", 30, "time limit for the quiz")

	// Parse the flags
	flag.Parse()

	// Return the read flags
	return *filename, *shuffle, *timelimit
}

// Reads the csv file from opened file
func readCSV(file *os.File) ([][]string, error) {
	// Read all lines from the file
	lines, err := csv.NewReader(file).ReadAll()

	// If there is an error while reading the file
	if err != nil{
		// Close the file
		file.Close()
		// Return nothing and an error
		return nil, err
	}

	// Count the lines of the file
	numOfLines := len(lines)
	// If there are zero lines
	if numOfLines == 0{
		// Return an error
		return nil, fmt.Errorf("File is empty")
	}

	// If all is good - return the read lines
	return lines, nil
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

// Asks the user all of the questions
func askQuestions(problems []problem, timelimit int) (int, error){
	// Create a timer with specified time in a flag
	timer := time.NewTimer(time.Duration(timelimit) * time.Second)

	// Correct count is set to zero
	correct := 0

	input := make(chan string)

	go getInput(input)

	// Go through all of the problems
	for i, p := range problems {
		// Asks for the each question
		ans, err := askASingleQuestion(p.q, p.a, i, timer.C, input)

		// If a question returns an error
		if err != nil && ans == -1{
			// Stop asking
			return correct, err
		} else if err != nil{
			// Print that the answer was wrong
			fmt.Println(err.Error())
		}

		// Otherwise, inrease the amount of correct answers
		correct+=ans
	}

	// Return the amount of correct answers
	return correct, nil
}

// Asks a single question from the all of the read
func askASingleQuestion(question, answer string, index int, timer <- chan time.Time, input <-chan string) (int, error){
	// Print the current problem
	fmt.Printf("Problem #%d: %s = \n", index+1, question)

	for {
		select {
		// If timer has sent a message (that it was done)
		case <- timer:
			return -1, fmt.Errorf("Your time has expired!")
		
		// If the input channel has some message in it
		case ans := <-input:
			// If the input euqals to answer - return 1 correct
			// Else - return 0 correct
			if strings.Compare(strings.Trim(strings.ToLower(ans), "\n"), answer) == 0{
				return 1, nil
			} else {
				return 0, fmt.Errorf("Wrong answer!")
			}
		}
	}
}

// Gets the user input for the answer
func getInput(inp chan string){
	// Infinite loop to listen for the answer
	for {
		// Get the answer
		reader := bufio.NewReader(os.Stdin)

		// Read the string
		result, err := reader.ReadString('\n')

		fmt.Println(result)
		
		// If an error occured
		if err != nil{
			exit(fmt.Sprintf("Error while reading the input, %s", err.Error()))
		}

		// If not - put the result into the channel
		inp <- result
	}
}

// Exits the program with a given message
func exit(msg string){
	// Print the message
	fmt.Println(msg)

	// Exit the program
	os.Exit(1)
}