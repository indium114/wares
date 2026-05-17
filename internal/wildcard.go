package internal

import (
	"regexp"
	"strings"
)

func wildcardMatch(pattern, s string) bool {
	// escape regex special chars first
	regex := regexp.QuoteMeta(pattern)

	// Replace escaped \* with .*
	regex = strings.ReplaceAll(regex, `\*`, ".*")

	// Match entire string
	regex = "^" + regex + "$"

	matched, err := regexp.MatchString(regex, s)
	if err != nil {
		return false
	}

	return matched
}

func filterMatches(pattern string, values []string) []string {
	var matches []string

	for _, v := range values {
		if wildcardMatch(pattern, v) {
			matches = append(matches, v)
		}
	}

	return matches
}
