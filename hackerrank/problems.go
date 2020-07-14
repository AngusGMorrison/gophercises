package hackerrank

import (
	"unicode"
)

func countCamelCaseWords(str string) int {
	count := 1
	for _, r := range str {
		if r >= 'A' && r <= 'Z' {
			count++
		}
	}
	return count
}

func caesarCipher(n int, str string, rot int) string {
	buf := make([]rune, n)
	for i, r := range str {
		if !unicode.IsLetter(r) {
			buf[i] = r
			continue
		}

		var baseline int
		if unicode.IsUpper(r) {
			baseline = int('A')
		} else {
			baseline = int('a')
		}
		rotated := (int(r)+rot-baseline)%26 + baseline
		buf[i] = rune(rotated)
	}
	return string(buf)
}
