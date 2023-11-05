package helpers

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Reads arguments from the command line
// And checks them for formatting
func ReadArguments() (bool, bool, string, string) {
	if len(os.Args) < 3 {
		panic("Too few arguments!\nUsage: 'go run main.go <input> <output>'")
	}

	b := flag.Bool("b", false, "Numbers that precede the main sequence")
	a := flag.Bool("a", false, "Numbers that follow the main sequence")

	flag.Parse()

	if len(os.Args[1]) != 4{
		panic("Wrong input size! Must be 4 letters")
	}
	if len(os.Args[2]) != 4{
		panic("Wrong output size! Must be 4 letters")
	}
	if !checkStrings(os.Args[1], os.Args[2]){
		panic("Wrong string input format!")
	}

	fmt.Println(*b, *a)
	os.Exit(1)

	return *b, *a, os.Args[1], os.Args[2]
}

// Checks the input strings for correct formatting
func checkStrings(inStrings ...string) (bool) {
	ref := "rgby"

	for _, s := range inStrings {
		for _, c := range s {
			if !strings.Contains(ref, string(c)){
				return false
			}
		}
	}

	return true
}

// Checks the input numbers for correct formatting
func checkNumber(num string) (bool) {
	ref := "1234"

	for _, c := range num {
		tempChar := string(c)
		if !strings.Contains(ref, tempChar){
			return false
		}

		ref = strings.ReplaceAll(ref, tempChar, "")
	}

	return true
}

// Reads user input for sequence numbers
func ReadSequenceInput(arr *[]string, msg string) () {
	var input string

	for {
		fmt.Printf("Input %s sequence of numbers (print 's' to stop)\n", msg)
		fmt.Scanln(&input)
		
		if input == "s"{
			break
		}

		if len(input) != 4 {
			fmt.Println("Wrong number length. Try again")
			continue
		}
		
		if !checkNumber(input){
			fmt.Println("Wrong number format. Try again")
			continue
		}

		*arr = append(*arr, input)
	}
}

// Resolves the decoding of provided sequences
func Decode(b []string, m []string, a []string, in string, out string) string {
	var answer string
	var temp string

	if len(b) > 0{
		for _, code := range b {
			temp = ""
			for _, c := range code {
				temp += string(in[int(c - '0') - 1])
			}
			in = temp
		}
	}
	
	if len(a) > 0{
		for i := len(b) - 1; i >= 0; i-- {
			temp = ""
			for _, c := range b[i] {
				temp += string(out[int(c - '0') - 1])
			}
			out = temp
		}
	}

	for _, code := range m {		
		for  i, c := range code {
			if in[int(c - '0') - 1] != out[i] {
				break
			}
			answer += string(c)
		}

		if len(answer) == len(in){
			break
		}
	}

	return answer;
}