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

func TestAskScoring(t *testing.T) {
	problems := []problem{{"1+2", "3"}, {"10+4", "14"}, {"9+6", "15"}}
	tests := []struct {
		inputPs      []problem
		answerReader io.Reader
		wantScore    int
	}{
		{
			inputPs:      problems,
			answerReader: strings.NewReader("3 14 15"),
			wantScore:    3,
		},
		{
			inputPs:      problems,
			answerReader: strings.NewReader("4 20 0"),
			wantScore:    0,
		},
		{
			inputPs:      problems,
			answerReader: strings.NewReader("3 20 0"),
			wantScore:    1,
		},
		{
			inputPs:      []problem{},
			answerReader: strings.NewReader("3 14 15"),
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
