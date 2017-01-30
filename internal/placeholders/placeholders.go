package placeholders

import "regexp"

var Regexp = regexp.MustCompile("<(locale_name|tag|locale_code)>")

func ContainsAnyPlaceholders(s string) bool {
	return Regexp.MatchString(s)
}
