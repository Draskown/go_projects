package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main(){
	fileName := flag.String("csv", "problems.csv", "a csv file in format of 'question,answer'")
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz")
	flag.Parse()

	log.Print(*shuffle)

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

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)

		var answer string

		fmt.Scanln(&answer)

		if answer == p.a {
			correct++
		}
	}

	fmt.Printf("%d/%d", correct, len(problems))
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