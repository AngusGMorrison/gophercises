package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type problem struct {
	q, a string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format 'question,answer'")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("failed to open CSV file: %s", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("failed to parse the provided CSV file"))
	}
	problems := parseLines(lines)
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
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func ask(problems []problem, in io.Reader) (correct int, err error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanWords)
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		if scanner.Scan() && scanner.Text() == p.a {
			correct++
		}
		if err := scanner.Err(); err != nil {
			return correct, err
		}
	}
	return
}

func exit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
