package helpers

// package helpers

// import (
// 	"strings"
// 	"testing"
// 	"time"

// 	"gotest.tools/assert"
// )

// func TestQuestions(t *testing.T) {
// 	// Test table with many test cases
// 	for i, testCase := range []struct{
// 		name string
// 		p problem
// 		d int
// 		a string
// 	}{
// 		{
// 			name: "Letters",
// 			p: problem{q: "a+b", a: "ab"},
// 			d: 2,
// 			a: "ab",
// 		},
// 		{
// 			name: "Digits",
// 			p: problem{q: "1+2+3", a: "6"},
// 			d: 2,
// 			a: "6",
// 		},
// 		{
// 			name: "Log",
// 			p: problem{q: "log(1)", a: "0"},
// 			d: 2,
// 			a: "0",
// 		},
// 		{
// 			name: "Case",
// 			p: problem{q: "What is the capital of Great Britain?", a: "London"},
// 			d: 2,
// 			a: "london",
// 		},
// 		{
// 			name: "Spaces",
// 			p: problem{q: "What is the capital of Great Britain?", a: "London"},
// 			d: 2,
// 			a: "  London",
// 		},
// 	}{
// 		// Create variables for answer and error
// 		var ans string
// 		var err error
		
// 		t.Run(testCase.name, func(t *testing.T) {
// 			// Create a timer for the specified duration
// 			timer := time.NewTimer(time.Duration(testCase.d) * time.Second).C
// 			// Create a channel to mock user input
// 			done := make(chan string)
// 			questionAnswered := make(chan bool)

// 			go func() {
// 				// Use the function
// 				ans, err = askASingleQuestion(testCase.p.q, testCase.p.a, i, timer, done)
// 				questionAnswered <- true
// 			}()

// 			// Send the supposed answer
// 			done <- testCase.a

// 			<-questionAnswered
// 			// Is an error occured - print in out
// 			if err != nil && ans == "" {
// 				t.Error(err)
// 			}

// 			// Assert the supposed and calculated answers
// 			assert.Equal(t, ans, strings.ToLower(testCase.p.a))
// 		})
// 	}
// }

// func TestReadCSV(t *testing.T) {
// 	// Test table as a map
// 	for name, testCase := range map [string]string{
// 		"Digits": "1+1,2\n2+1,3\n9+9,18\n",
// 		"Log and digits": "a+b,ab\nlog(1),0\n1+3+5,9\n",
// 		"String questions": "What is London,A city\n What is life,No Idea\n",
// 	} {
// 		t.Run(name, func(t *testing.T){
// 			// Create questions from the test table
// 			questions, err := ReadCSV(strings.NewReader(testCase))

// 			// If an error occured while reading the file
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			// Split the test case into lines
// 			lines := strings.Split(testCase, "\n")

// 			// For each line
// 			for i, line := range lines {
// 				if line == "" {
// 					continue
// 				}

// 				// Cut the answer from the test case
// 				answer := strings.Split(line, ",")[1]

// 				// Assert the supposed answer to the actual
// 				assert.Equal(t, questions[i].a, answer)
// 			}
// 		})
// 	}
// }