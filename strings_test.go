package sutils

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestIContains(t *testing.T) {
	pairs := []struct {
		Value    string
		Search   string
		Expected bool
	}{
		// check exact matches with varying cases
		{"hello", "hello", true},
		{"hello", "HELLO", true},
		{"HELLO", "hello", true},
		{"heLLo", "hello", true},
		{"hello", "HeLLo", true},
		// check with space
		{"HELLO there", "hello", true},
		{"heLLo there", "lo th", true},
		{"hello there", "LO TH", true},
		// check submatches with varying cases
		{"hello", "el", true},
		{"hello", "EL", true},
		{"hello", "eL", true},
		// check false hits
		{"hello", "oh", false},
		{"hello", "OH", false},
		{"hello", "oH", false},
		{"this is", "ss", false},
		{"this is", "sS", false},
		{"this is", "SS", false},
		{"", "", false},
		{"hello", "", false},
		{"", "hello", false},
	}

	for _, c := range pairs {
		if res := IContains(c.Value, c.Search); res != c.Expected {
			t.Errorf("Mismatch. Expected: %v, got: %v for testcase (%q, %q)", c.Expected, res, c.Value, c.Search)
		}
	}
}

var (
	testString = `Once upon a midnight dreary, while I pondered, weak and weary,
Over many a quaint and curious volume of forgotten lore—
    While I nodded, nearly napping, suddenly there came a tapping,
As of some one gently rapping, rapping at my chamber door.
“’Tis some visitor,” I muttered, “tapping at my chamber door—
            Only this and nothing more.”`
)

func TestCountIgnoreCase(t *testing.T) {
	var (
		result   int
		expected int

		buf = bytes.NewBufferString(testString)
	)

	result, _ = CountIgnoreCase(buf, "my chamber")
	expected = 2

	if result != expected {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expected, result)
	}

	buf = bytes.NewBufferString(testString)

	result, _ = CountIgnoreCase(buf, "MY chamber")
	if result != expected {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expected, result)
	}
}

func TestCountCaseSensitive(t *testing.T) {
	var (
		result   int
		expected int

		buf = bytes.NewBufferString(testString)
	)

	result, _ = CountCaseSensitive(buf, "my chamber")
	expected = 2

	if result != expected {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expected, result)
	}

	buf = bytes.NewBufferString(testString)

	result, _ = CountCaseSensitive(buf, "MY chamber")
	expected = 0

	if result != expected {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expected, result)
	}
}

func TestFindIgnoreCase(t *testing.T) {
	var (
		result      []int
		expectedLen int

		expectedLineNumbers = []int{4, 5}

		buf = bytes.NewBufferString(testString)
	)

	result, _ = FindWith(IContains, buf, "my chamber")
	expectedLen = 2

	if len(result) != expectedLen {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expectedLen, len(result))
	}

	if result[0] != expectedLineNumbers[0] || result[1] != expectedLineNumbers[1] {
		t.Errorf("Error. Found match at %d and %d, but should've been at %d and %d", result[0], result[1], expectedLineNumbers[0], expectedLineNumbers[1])
	}

	buf = bytes.NewBufferString(testString)

	result, _ = FindWith(IContains, buf, "MY chamber")
	expectedLen = 2

	if len(result) != expectedLen {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expectedLen, len(result))
	}

	if result[0] != expectedLineNumbers[0] || result[1] != expectedLineNumbers[1] {
		t.Errorf("Error. Found match at %d and %d, but should've been at %d and %d", result[0], result[1], expectedLineNumbers[0], expectedLineNumbers[1])
	}
}

func TestFindCaseSensitive(t *testing.T) {
	var (
		result      []int
		expectedLen int

		expectedLineNumbers = []int{4, 5}

		buf = bytes.NewBufferString(testString)
	)

	result, _ = FindWith(strings.Contains, buf, "my chamber")
	expectedLen = 2

	if len(result) != expectedLen {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expectedLen, len(result))
	}

	if result[0] != expectedLineNumbers[0] || result[1] != expectedLineNumbers[1] {
		t.Errorf("Error. Found match at %d and %d, but should've been at %d and %d", result[0], result[1], expectedLineNumbers[0], expectedLineNumbers[1])
	}

	buf = bytes.NewBufferString(testString)

	result, _ = FindCaseSensitive(buf, "MY chamber")
	expectedLen = 0

	if len(result) != expectedLen {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expectedLen, len(result))
	}
}

func TestFindWithPrefix(t *testing.T) {
	var (
		result      []int
		expectedLen int

		expectedLineNumbers = []int{4, 5}

		buf = bytes.NewBufferString(testString)
	)

	result, _ = FindWith(strings.HasPrefix, buf, "Once")
	expectedLen = 1

	if len(result) != expectedLen {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expectedLen, len(result))
	}

	if result[0] != 1 {
		t.Errorf("Error. Found match at %d and %d, but should've been at %d and %d", result[0], result[1], expectedLineNumbers[0], expectedLineNumbers[1])
	}

	buf = bytes.NewBufferString(testString)

	result, _ = FindCaseSensitive(buf, "once")
	expectedLen = 0

	if len(result) != expectedLen {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expectedLen, len(result))
	}
}

func TestTrimNL(t *testing.T) {
	const (
		input        = "Nothing\n"
		inputWindows = "Nothing\r\n"
		trimInput    = "Nothing"

		noInput        = "\n"
		noInputWindows = "\r\n"
		trimNoInput    = ""
	)

	for _, in := range []string{input, inputWindows} {
		trimmed := TrimNL(in)

		if msg, ok := expect(trimInput, trimmed); !ok {
			t.Error(msg)
		}
	}

	for _, in := range []string{noInput, noInputWindows} {
		trimmed := TrimNL(in)

		if msg, ok := expect(trimNoInput, trimmed); !ok {
			t.Error(msg)
		}
	}

}

func TestOccurs(t *testing.T) {

}

/*
============== Benchmarks ==============
*/

func BenchmarkIContainsFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IContains("This is a rather long line and I'm curious whether that thing is in there or not.", "hat")
	}
}

func BenchmarkIContainsNotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IContains("This is a rather long line and I'm curious whether that thing is in there or not.", "moo")
	}
}

func BenchmarkIContainsShorted(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IContains("This is a rather long line and I'm curious whether that thing is in there or not.", "")
	}
}

func BenchmarkFindIgnoreCase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindIgnoreCase(bytes.NewBufferString(testString), "my")
	}
}

func BenchmarkFindCaseSensitive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindCaseSensitive(bytes.NewBufferString(testString), "my")
	}
}

func BenchmarkFindWithStringsContains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindWith(strings.Contains, bytes.NewBufferString(testString), "weary")
	}
}

func BenchmarkFindWithStringsHasPrefix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindWith(strings.HasPrefix, bytes.NewBufferString(testString), "weary")
	}
}

func BenchmarkFindWithStringsContainsAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindWith(strings.ContainsAny, bytes.NewBufferString(testString), "weary")
	}
}

/*
============== Utils ==============
*/

func expect(expected, actual string) (string, bool) {
	if actual != expected {
		return fmt.Sprintf("Mismatch: Expected %q, got %q", expected, actual), false
	}

	return "", true
}
