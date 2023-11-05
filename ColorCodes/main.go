package main

import (
	hlp "ColorCodes/Helpers"
	"fmt"
)

func main() {
	seq_before := []string{}
	seq_after := []string{}
	seq_main := []string{}

	isSeqBefore, isSeqAfter, input, output := hlp.ReadArguments()

	if isSeqBefore{
		hlp.ReadSequenceInput(&seq_before, "the before")
	}

	hlp.ReadSequenceInput(&seq_main, "main")

	if isSeqAfter{
		hlp.ReadSequenceInput(&seq_after, "the after")
	}

	answer := hlp.Decode(seq_before, seq_main, seq_after, input, output)

	if len(answer) == 0{
		fmt.Println("The're is no answer to these sequences")
	}

	fmt.Printf("The answer is %s\n", answer)
}