package main

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Versão com regex
func normalize(input string) string {
	s := strings.ToUpper(input)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	re := regexp.MustCompile(`[^A-Z0-9 ]+`)
	result = re.ReplaceAllString(result, "")
	return result
}

// Versão mais performática (strings.Map)
func normalizeFast(s string) string {
	s = norm.NFD.String(s)
	return strings.Map(func(r rune) rune {
		switch {
		case unicode.Is(unicode.Mn, r): // remove acento
			return -1
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			return unicode.ToUpper(r)
		case unicode.IsSpace(r):
			return ' '
		default:
			return -1
		}
	}, s)
}
