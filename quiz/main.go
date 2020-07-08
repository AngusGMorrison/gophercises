package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	q, a string
}

var (
	out         io.Writer = os.Stdout
	csvFilename string
	shuffle     bool
	timeLimit   int
)

func main() {
	flag.StringVar(&csvFilename, "csv", "problems.csv", "a csv file in the format 'question,answer'")
	flag.BoolVar(&shuffle, "shuffle", false, "shuffle questions")
	flag.IntVar(&timeLimit, "limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(csvFilename)
	if err != nil {
		exit(fmt.Sprintf("failed to open CSV file: %s", csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("failed to parse the provided CSV file"))
	}
	problems := parseLines(lines)
	if shuffle {
		shuffleProblems(problems)
	}

	correct, err := ask(problems, os.Stdin)
	if err != nil {
		exit(fmt.Sprintf("getting user input: %v", err))
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: formatAnswer(line[1]),
		}
	}
	return ret
}

func shuffleProblems(problems []problem) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(problems),
		func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
}

func ask(problems []problem, in io.Reader) (correct int, err error) {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanWords)
	answerCh := make(chan string)
	errorCh := make(chan error)

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)

		go func() {
			if scanner.Scan() {
				answerCh <- formatAnswer(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				errorCh <- err
			}
		}()

		select {
		case <-timer.C:
			fmt.Fprintln(out, "Times up!")
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		case err = <-errorCh:
			return
		}
	}
	return
}

func formatAnswer(answer string) string {
	return strings.ToLower(strings.TrimSpace(answer))
}

func exit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
