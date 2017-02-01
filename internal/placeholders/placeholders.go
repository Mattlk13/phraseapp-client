package placeholders

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/phrase/phraseapp-client/internal/helpers"
)

var Regexp = regexp.MustCompile("<(locale_name|tag|locale_code)>")

func ContainsAnyPlaceholders(s string) bool {
	return Regexp.MatchString(s)
}

func ToGlobbingPattern(s string) string {
	return Regexp.ReplaceAllString(s, "*")
}

// Resolve matches s against pattern and maps placeholders in pattern to
// substrings of s.
// Resolve handles '*' wildcards in the pattern, but will return an error
// if the pattern contains '**'.
func Resolve(s, pattern string) (map[string]string, error) {
	if strings.Contains(pattern, "**") {
		return map[string]string{}, fmt.Errorf("'**' wildcard not allowed in pattern")
	}

	placeholders := Regexp.FindAllString(pattern, -1)
	if len(placeholders) <= 0 {
		return map[string]string{}, nil
	}

	patternRE := regexp.QuoteMeta(pattern)
	patternRE = strings.Replace(patternRE, "\\*", ".*", -1)

	for _, placeholder := range helpers.UniqueElements(placeholders) {
		placeholder = regexp.QuoteMeta(placeholder)
		placeholderRE := fmt.Sprintf("(?P%s[^%s]+)", placeholder, string(filepath.Separator)) // build named subexpression (capturing group) from placeholder
		patternRE = strings.Replace(patternRE, placeholder, placeholderRE, -1)
	}

	patternRegex, err := regexp.Compile(patternRE)
	if err != nil {
		return nil, err
	}

	matchNames := patternRegex.SubexpNames()
	matches := patternRegex.FindStringSubmatch(s)

	if len(matches) < len(placeholders)+1 || matches[0] != s {
		return nil, fmt.Errorf("string %q does not match pattern %q", s, patternRE)
	}

	// drop first element, which is the entire string s wich match name ""
	matches, matchNames = matches[1:], matchNames[1:]

	values := map[string]string{}
	for i, match := range matches {
		placeholder := matchNames[i]
		if value, ok := values[placeholder]; ok {
			if match != value {
				return nil, fmt.Errorf("string %q does not match pattern %q: placeholder %q is used twice with different values", s, patternRE, placeholder)
			}
		}

		values[placeholder] = match
	}

	return values, nil
}
