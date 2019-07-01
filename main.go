package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Question struct {
	question string
	answer   string
}

func readQuestions(fileName string) []Question {
	csvFile, _ := os.Open(fileName)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var questions []Question

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		questions = append(questions, Question{question: line[0], answer: line[1]})
	}

	return questions
}

func askQuestion(q Question, reader *bufio.Reader, timer *time.Timer, c chan string) (bool, error) {
	fmt.Printf("%s=", q.question)
	go func() {
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(text, "\n")
		c <- text
	}()

	select {
	case <-timer.C:
		return false, errors.New("time's up")
	case msg := <-c:
		return msg == q.answer, nil
	}
}

func printScore(score int) {
	fmt.Printf("You scored %d points\n", score)
}

func main() {
	fileName := flag.String("file", "problems.csv", "Name of the csv file")
	timeLimit := flag.Int("time", 5, "quiz time limit in seconds")
	flag.Parse()
	questions := readQuestions(*fileName)

	stdinReader := bufio.NewReader(os.Stdin)
	score := 0

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	ch := make(chan string)
	for _, q := range questions {
		isAnswerCorrect, err := askQuestion(q, stdinReader, timer, ch)
		if err != nil {
			printScore(score)
			return
		}
		if isAnswerCorrect {
			score++
		}
	}
	printScore(score)
}
