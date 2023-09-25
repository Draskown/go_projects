package quiz

import (
	"os"
	"io"
	"fmt"
	"time"
	"flag"
	"bufio"
	"strings"
	"encoding/csv"
)

// Opens the file from given directory
func OpenFile(name string) (*os.File, error) {
	return os.Open(name)
}

// Reads argumetns from a command line
func ReadArguments() (string, bool, int) {
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
func ReadCSV(file io.Reader) ([]problem, error) {
	// Read all lines from the file
	lines, err := csv.NewReader(file).ReadAll()

	// If the reader failed to read the file
	if err != nil{
		// Exit the program with an error message
		Exit(fmt.Sprintf("Failed to parse the file. Error %s", err.Error()))
	}

	// Count the lines of the file
	numOfLines := len(lines)
	// If there are zero lines
	if numOfLines == 0{
		// Return an error
		return nil, fmt.Errorf("file is empty")
	}

	// Parse the read lines into a problem structure
	problems := ParseLines(lines)

	// If all is good - return the read lines
	return problems, nil
}

// Custom structure of a problem
// where q is a question, a - an answer
type problem struct{
	q string
	a string
}

// Parses the read lines from csv file
// into an array of custom structured problems
func ParseLines(input [][]string) []problem{
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
func AskQuestions(problems []problem, timelimit int) (int, error){
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
		if err != nil && ans == ""{
			// Stop asking
			return correct, err
		} else if err != nil{
			// Print that the answer was wrong
			fmt.Println(err.Error())
		}

		// Otherwise, inrease the amount of correct answers
		if strings.EqualFold(ans, p.a){
			correct++
		}
	}

	// Return the amount of correct answers
	return correct, nil
}

// Asks a single question from the all of the read
func askASingleQuestion(question, answer string, index int, timer <-chan time.Time, input <-chan string) (string, error){
	// Print the current problem
	fmt.Printf("Problem #%d: %s = \n", index+1, question)

	for {
		select {
		// If timer has sent a message (that it was done)
		case <- timer:
			return "", fmt.Errorf("Your time has expired")
		
		// If the input channel has some message in it
		case ans := <-input:
			// Return the lowercase answer with no spaces and newline symbols
			return strings.ToLower(strings.TrimSpace(strings.TrimRight(ans, "\r\n"))), nil
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

		// If an error occured
		if err != nil{
			Exit(fmt.Sprintf("Error while reading the input, %s", err.Error()))
		}

		// If not - put the result into the channel
		inp <- result
	}
}

// Exits the program with a given message
func Exit(msg string){
	// Print the message
	fmt.Println(msg)

	// Exit the program
	os.Exit(1)
}