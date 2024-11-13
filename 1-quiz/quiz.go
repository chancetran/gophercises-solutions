/*
Solution to the Quiz Game exercise.

Citations:
https://pkg.go.dev/encoding/csv
https://gobyexample.com/structs
https://gobyexample.com/line-filters
https://tip.golang.org/doc/comment
https://go.dev/tour/concurrency/1
https://go.dev/tour/concurrency/2
https://go.dev/tour/concurrency/5
https://gobyexample.com/timeouts
https://stackoverflow.com/questions/12264789/shuffle-array-in-go
*/

package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/exp/rand"
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
// If the user passes the 'shuffle' flag as `true`, then the data is randomly permutated prior to processing.
func processData(data [][]string, shuffle *bool) []Problem {

	if *shuffle {
		for i := range data {
			j := rand.Intn(i + 1)
			data[i], data[j] = data[j], data[i]
		}
	}

	processedData := make([]Problem, len(data))
	for index, row := range data {
		processedData[index] = Problem{row[0], row[1]}
	}

	return processedData
}

// Quiz a user using problems from a CSV file.
// Keep track of the number of questions the user gets right and display it at the end.
// Regardless if the user got a problem right, provide the next problem immediately after.
// The user is asked to confirm when they're ready to begin the quiz, triggering the timer to start.
// If the quiz is not completed within the time limit, the quiz is terminated even if the user is inputting an answer.
func quiz(data []Problem, time_limit *int) {

	// Confirm that the user is ready to take the quiz
	fmt.Printf("To begin the quiz, please press 'Enter':\n")

	r := bufio.NewScanner(os.Stdin)
	r.Scan()

	if len(r.Text()) > 0 {
		log.Fatal("Text detected prior to pressing 'Enter'; aborting....")
	}

	score := 0
	timer := time.NewTimer(time.Duration(*time_limit) * time.Second)

QuizLoop:
	for index, problem := range data {

		answerCh := make(chan string)
		go func() {
			fmt.Printf("Question %d: %s?\n", index+1, problem.Question)

			r.Scan()
			answerCh <- r.Text()
		}()

		select {
		case <-timer.C:
			fmt.Print("\nTime's up!\n")
			break QuizLoop
		case answer := <-answerCh:
			if answer == problem.Answer {
				score++
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", score, len(data))
}

func main() {

	filepath := flag.String("csv", "problems.csv", "a CSV file in the following format: 'question,answer'")
	shuffle := flag.Bool("shuffle", false, "shuffle the data...")
	time_limit := flag.Int("time_limit", 30, "the time limit of the quiz (in seconds)")

	flag.Parse()

	data := loadData(filepath)
	processedData := processData(data, shuffle)

	quiz(processedData, time_limit)

	os.Exit(0)
}
