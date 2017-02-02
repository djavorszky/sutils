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
	occurrences, err := FindIgnoreCase(haystack, needle)
	if err != nil {
		return 0, err
	}

	return len(occurrences), nil
}

// CountCaseSensitive searches an io.Reader for a given string in a case sensitive way.
// It returns the number of occurrences it found, or an error if something went wrong.
func CountCaseSensitive(haystack io.Reader, needle string) (count int, err error) {
	occurrences, err := FindCaseSensitive(haystack, needle)
	if err != nil {
		return 0, err
	}

	return len(occurrences), nil
}

// FindIgnoreCase searches an io.Reader for a given string in a case-insensitive way.
// It returns the line numbers where it such strings found, or an error if something went wrong.
func FindIgnoreCase(haystack io.Reader, needle string) (occurrences []int, err error) {
	lines := 0
	scanner := bufio.NewScanner(haystack)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		if IContains(scanner.Text(), needle) {
			occurrences = append(occurrences, lines)
		}
		lines++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return occurrences, nil
}

// FindCaseSensitive searches an io.Reader for a given string in a case sensitive way.
// It returns the line numbers where it such strings found, or an error if something went wrong.
func FindCaseSensitive(haystack io.Reader, needle string) (occurrences []int, err error) {
	lines := 0
	scanner := bufio.NewScanner(haystack)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), needle) {
			occurrences = append(occurrences, lines)
		}
		lines++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return occurrences, nil
}
