package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main(){
	fileName := flag.String("csv", "problems.csv", "a csv file in format of 'question,answer'")
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