package main

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestParseLines(t *testing.T) {
	tests := []struct {
		in        [][]string
		wantLen   int
		wantProbs []problem
	}{
		{
			in:        [][]string{},
			wantLen:   0,
			wantProbs: nil,
		},
		{
			in:        [][]string{{"1+2", "3"}},
			wantLen:   1,
			wantProbs: []problem{{"1+2", "3"}},
		},
		{
			in:        [][]string{{"1+2", "3"}, {"10+4", "14"}},
			wantLen:   2,
			wantProbs: []problem{{"1+2", "3"}, {"10+4", "14"}},
		},
	}

	for _, test := range tests {
		problems := parseLines(test.in)
		if len(problems) != test.wantLen {
			t.Errorf("parseLines(%v): output has length %d, want %d",
				test.in, len(problems), test.wantLen)
			continue
		}
		for i, p := range problems {
			wantProb := test.wantProbs[i]
			if p.q != wantProb.q {
				t.Errorf("parseLines(%v): problem at index %d has question %q, want %q",
					test.in, i, p.q, wantProb.q)
			}
			if p.a != wantProb.a {
				t.Errorf("parseLines(%v): problem at index %d has question %q, want %q",
					test.in, i, p.a, wantProb.a)
			}
		}
	}
}

func TestShuffleProblems(t *testing.T) {
	problems := []problem{{"1+2", "3"}, {"10+4", "14"}, {"9+6", "15"}, {"a + a", "2a"}}
	shuffled := make([]problem, 4)
	copy(shuffled, problems)
	shuffleProblems(shuffled)
	if problems[0] == shuffled[0] &&
		problems[1] == shuffled[1] &&
		problems[2] == shuffled[2] &&
		problems[3] == shuffled[3] {
		t.Errorf("problems array was not shuffled")
	}
}

func TestAskScoring(t *testing.T) {
	timeLimit = 30
	problems := []problem{{"1+2", "3"}, {"10+4", "14"}, {"9+6", "15"}, {"a + a", "2a"}}
	tests := []struct {
		inputPs      []problem
		answerReader io.Reader
		wantScore    int
	}{
		{
			inputPs:      problems,
			answerReader: strings.NewReader("3 14 15 2a"),
			wantScore:    4,
		},
		{
			inputPs:      problems,
			answerReader: strings.NewReader("3 14 15 2a"),
			wantScore:    4,
		},
		{
			inputPs:      problems,
			answerReader: strings.NewReader("3 14 15 2A"),
			wantScore:    4,
		},
		{
			inputPs:      problems,
			answerReader: strings.NewReader("4 20 0 b"),
			wantScore:    0,
		},
		{
			inputPs:      problems,
			answerReader: strings.NewReader("3 20 0 b"),
			wantScore:    1,
		},
		{
			inputPs:      []problem{},
			answerReader: strings.NewReader("3 14 15 2a"),
			wantScore:    0,
		},
	}

	for _, test := range tests {
		score, _ := ask(test.inputPs, test.answerReader)
		if score != test.wantScore {
			t.Errorf("ask(%v, answerReader): scored %d, want %d",
				test.inputPs, score, test.wantScore)
		}
	}
}

func TestAskTimeOut(t *testing.T) {
	timeLimit = 1
	testProblems := []problem{{"1+2", "3"}}
	done := make(chan struct{})
	go func() {
		ask(testProblems, os.Stdin)
		done <- struct{}{}
	}()

	select {
	case <-time.After(2 * time.Second):
		t.Errorf("ask(%v, os.Stdin) did not time out", testProblems)
	case <-done:
	}
}
