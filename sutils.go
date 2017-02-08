package sutils

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/icrowley/fake"
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
// It returns the line numbers where it found such strings, or an error if something went wrong.
func CountIgnoreCase(haystack io.Reader, needle string) (int, error) {
	occurrences, err := FindIgnoreCase(haystack, needle)
	if err != nil {
		return 0, err
	}

	return len(occurrences), nil
}

// CountCaseSensitive searches an io.Reader for a given string in a case sensitive way.
// It returns the line numbers where it found such strings, or an error if something went wrong.
func CountCaseSensitive(haystack io.Reader, needle string) (int, error) {
	occurrences, err := FindCaseSensitive(haystack, needle)
	if err != nil {
		return 0, err
	}

	return len(occurrences), nil
}

// FindIgnoreCase searches an io.Reader for a given string in a case-insensitive way.
// It returns the line numbers where it found such strings, or an error if something went wrong.
func FindIgnoreCase(haystack io.Reader, needle string) (occurrences []int, err error) {
	lines := 1
	reader := bufio.NewReader(haystack)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		if IContains(line, needle) {
			occurrences = append(occurrences, lines)
		}

		lines++
	}
	return occurrences, nil
}

// FindCaseSensitive searches an io.Reader for a given string in a case sensitive way.
// It returns the line numbers where it found such strings, or an error if something went wrong.
func FindCaseSensitive(haystack io.Reader, needle string) (occurrences []int, err error) {
	lines := 1
	reader := bufio.NewReader(haystack)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		if strings.Contains(line, needle) {
			occurrences = append(occurrences, lines)
		}

		lines++
	}

	return occurrences, nil
}

// FindStartsWith searches an io.Reader for all lines that start with a given string in a case sensitive way.
// It returns the line numbers where it found such strings, or an error if something went wrong.
func FindStartsWith(haystack io.Reader, needle string) (occurrences []int, err error) {
	lines := 1
	reader := bufio.NewReader(haystack)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		if strings.HasPrefix(line, needle) {
			occurrences = append(occurrences, lines)
		}

		lines++
	}

	return occurrences, nil
}

// TrimNL trims the newline from the end of the string.
func TrimNL(s string) string {
	s = strings.TrimSuffix(s, "\n")
	s = strings.TrimSuffix(s, "\r")

	return s
}

func FindWith(find func(string, string) bool, haystack io.Reader, needle string) (occurrences []int, err error) {

	lines := 1
	reader := bufio.NewReader(haystack)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		if find(line, needle) {
			occurrences = append(occurrences, lines)
		}

		lines++
	}

	return occurrences, nil
}

// RandUserName returns a random username that can be used for databases
func RandUserName() string {
	res := strings.Split(fake.ProductName(), " ")[:1]

	return strings.ToLower(strings.Join(res, "_"))
}

// RandDBName returns a random name that can be used as a name for a database
func RandDBName() string {
	return strings.ToLower(fake.Brand())
}

// RandPassword returns a random password that can be used for databases
func RandPassword() string {
	return RandUserName()
}
