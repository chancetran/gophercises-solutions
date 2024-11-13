/*
Solution to the Quiz Game exercise.

Citations:
https://pkg.go.dev/encoding/csv
https://gobyexample.com/structs
https://gobyexample.com/line-filters
https://tip.golang.org/doc/comment
*/

package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

type Problem struct {
	Question string
	Answer   string
}

// Load data into memory from the specified filepath.
func loadData(filepath *string) [][]string {

	file, err := os.Open(*filepath)
	if err != nil {
		log.Fatal("An error occurred while reading ", *filepath)
	}

	r := csv.NewReader(file)
	data, err := r.ReadAll()
	if err != nil {
		log.Fatal("An error occured while parsing ", *filepath)
	}

	return data
}

// Process the data into a Problem array.
func processData(data [][]string) []Problem {

	processedData := make([]Problem, len(data))
	for index, row := range data {
		processedData[index] = Problem{row[0], row[1]}
	}

	return processedData
}

// Quiz a user using problems from a CSV file.
// Keep track of the number of questions the user gets right and display it at the end.
// Regardless if the user got a problem right, provide the next problem immediately after.
func quiz(data []Problem) {

	score := 0
	for index, problem := range data {

		fmt.Printf("Question %d: %s?\n", index+1, problem.Question)

		r := bufio.NewScanner(os.Stdin)
		r.Scan()
		response := r.Text()

		if response == problem.Answer {
			score++
		}

	}

	fmt.Printf("Your score: %d\n", score)
}

func main() {

	filepath := flag.String("csv", "problems.csv", "a CSV file in the following format: 'question,answer'")

	data := loadData(filepath)
	processedData := processData(data)

	quiz(processedData)

	os.Exit(0)
}
