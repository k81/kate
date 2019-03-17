package utils

import (
	"strings"
	"unicode"
)

// ToSnake converts a given string to snake case
func ToSnake(s string) string {
	var result string
	var words []string
	var lastPos int
	rs := []rune(s)

	for i := 0; i < len(rs); i++ {
		if i > 0 {
			switch {
			case unicode.IsUpper(rs[i]):
				words = append(words, s[lastPos:i])
				lastPos = i
			case rs[i] == '-':
				fallthrough
			case rs[i] == '.':
				words = append(words, s[lastPos:i])
				i++
				lastPos = i
			}
		}
	}

	// append the last word
	if s[lastPos:] != "" {
		words = append(words, s[lastPos:])
	}

	for k, word := range words {
		if k > 0 {
			result += "_"
		}

		result += strings.ToLower(word)
	}

	return result
}

// ToCamel returns a string converted from snake case to uppercase
func ToCamel(s string) string {
	return snakeToCamel(s, true)
}

// ToCamelLower returns a string converted from snake case to lowercase
func ToCamelLower(s string) string {
	return snakeToCamel(s, false)
}

func snakeToCamel(s string, upperCase bool) string {
	var result string

	words := strings.Split(s, "_")

	for i, word := range words {
		if (upperCase || i > 0) && len(word) > 0 {
			w := []rune(word)
			w[0] = unicode.ToUpper(w[0])
			result += string(w)
		} else {
			result += word
		}
	}

	return result
}
