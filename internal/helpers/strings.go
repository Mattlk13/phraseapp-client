package helpers

import "strings"

// Contains returns true if str is an element of seq.
func Contains(seq []string, str string) bool {
	for _, elem := range seq {
		if str == elem {
			return true
		}
	}
	return false
}

// ContainsAnySub returns true if s contains at least one element in subs.
func ContainsAnySub(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}

	return false
}

// UniqueElements returns s with all duplicate elements removed.
func UniqueElements(s []string) []string {
	seen := map[string]bool{}
	clean := []string{}

	for _, elem := range s {
		if !seen[elem] {
			seen[elem] = true
			clean = append(clean, elem)
		}
	}

	return clean
}
