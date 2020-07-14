package hackerrank

import (
	"fmt"
	"testing"
)

func TestCountCamelCaseWords(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"camel", 1},
		{"camelCase", 2},
		{"camelCaseCamel", 3},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			got := countCamelCaseWords(test.in)
			if got != test.want {
				t.Fatalf("got %d, want %d", got, test.want)
			}
		})

	}
}

func TestCaesarCipher(t *testing.T) {
	tests := []struct {
		n    int
		in   string
		rot  int
		want string
	}{
		{1, "a", 1, "b"},
		{1, "a", 26, "a"},
		{1, "A", 1, "B"},
		{6, "aBcDyZ", 1, "bCdEzA"},
		{8, "a!BcD_yZ", 2, "c!DeF_aB"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("caesarCipher(%d, %s, %d)", test.n, test.in, test.rot), func(t *testing.T) {
			got := caesarCipher(test.n, test.in, test.rot)
			if got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}

func BenchmarkCountCamelCaseWords(b *testing.B) {
	for i := 0; i < b.N; i++ {
		countCamelCaseWords("camelCaseCamel")
	}
}
