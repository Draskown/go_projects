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

func main(){
	fileName := flag.String("csv", "problems.csv", "a csv file in format of 'question,answer'")
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz")
	timeLimit := flag.Int("time", 30, "time limit for the quiz")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open file: %s.", *fileName))
	}
	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil{
		exit("Failed to parse the file.")
	}
	problems := parseLines(lines)

	if *shuffle{
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) { 
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		answerC := make(chan string)
		
		go func(){
			var answer string
			fmt.Scanln(&answer)
			answerC <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("Your time has expired. %d/%d", correct, len(problems))
			return
		case answer := <-answerC:
			if answer == p.a {
				correct++
			}
		}
	}

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

	fmt.Printf(addString+" %d/%d", correct, len(problems))
}

type problem struct{
	q string
	a string
}

func parseLines(input [][]string) []problem{
	res := make([]problem, len(input))

	for i, line := range input {
		res[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return res
}

func exit(msg string){
	fmt.Println(msg)
	os.Exit(1)
}