package sutils

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

// IContains returns true if the haystack contains the needle.
// It searches in a case-insensitive way
func IContains(haystack, needle string) bool {
	// short out
	if haystack == "" || needle == "" {
		return false
	}

	re := regexp.MustCompile("(?i)" + needle)
	m := re.FindString(haystack)

	if m == "" {
		return false
	}

	return true
}

// Present checks whether all of its parameters are non-empty.
func Present(reqFields ...string) bool {
	for _, field := range reqFields {
		if field == "" {
			return false
		}
	}

	return true
}

// CountIgnoreCase searches an io.Reader for a given string in a case-insensitive way.
// It returns the number of occurrences it found, or an error if something went wrong.
func CountIgnoreCase(haystack io.Reader, needle string) (count int, err error) {
	scanner := bufio.NewScanner(haystack)

	for scanner.Scan() {
		if IContains(scanner.Text(), needle) {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

// CountCaseSensitive searches an io.Reader for a given string in a case sensitive way.
// It returns the number of occurrences it found, or an error if something went wrong.
func CountCaseSensitive(haystack io.Reader, needle string) (count int, err error) {
	scanner := bufio.NewScanner(haystack)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), needle) {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}
