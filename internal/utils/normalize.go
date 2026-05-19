package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func Normalize(input string, fields ...string) string {
	for _, field := range fields {
		re := regexp.MustCompile(fmt.Sprintf(`(%q\s*:\s*(("[^"]+"|\d+|null|true|false)|\[[^\]]*\]|\{[^\}]*\}))`, field))
		input = re.ReplaceAllString(input, fmt.Sprintf(`%q: "placeholder"`, field))
	}
	return strings.TrimSpace(input)
}
