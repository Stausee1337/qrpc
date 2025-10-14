package source

import (
	"strings"
	"unicode"
)

func ToPascalCase(s string) string {
	separators := []string{"_", "-", " "}
	for _, sep := range separators {
		s = strings.ReplaceAll(s, sep, " ")
	}
	parts := strings.Fields(s)

	for i, p := range parts {
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, "")
}
