package sutils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
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
		result              []int
		expectedLen         int
		expectedLineNumbers = []int{4, 5}

		buf = bytes.NewBufferString(testString)
	)

	result, _ = FindWith(IContains, buf, []string{"my chamber"})
	expectedLen = 2

	if len(result) != expectedLen {
		t.Fatalf("Mismatch. Expected count=%d, got result=%d", expectedLen, len(result))
	}

	if result[0] != expectedLineNumbers[0] || result[1] != expectedLineNumbers[1] {
		t.Errorf("Error. Found match at %d and %d, but should've been at %d and %d", result[0], result[1], expectedLineNumbers[0], expectedLineNumbers[1])
	}

	buf = bytes.NewBufferString(testString)

	result, _ = FindWith(IContains, buf, []string{"MY chamber"})
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

	result, _ = FindWith(strings.Contains, buf, []string{"my chamber"})
	expectedLen = 2

	if len(result) != expectedLen {
		t.Fatalf("Mismatch. Expected count=%d, got result=%d", expectedLen, len(result))
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

	result, _ = FindWith(strings.HasPrefix, buf, []string{"Once"})
	expectedLen = 1

	if len(result) != expectedLen {
		t.Fatalf("Mismatch. Expected count=%d, got result=%d", expectedLen, len(result))
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

func TestFindWith(t *testing.T) {
	tests := []struct {
		Haystack string
		Needle   []string
		Expected []int
	}{
		{"looking for this\nbut not for that\n", []string{""}, []int{}},
		{"looking for this\nbut not for that\n", []string{"Madness"}, []int{}},
		{"looking for this\nbut not for that\n", []string{"this"}, []int{1}},
		{"looking for this\nbut not for that\n", []string{"that"}, []int{2}},
		{"looking for this\nbut not for that\n", []string{"for"}, []int{1, 2}},
	}

	for _, test := range tests {
		found, err := FindWith(strings.Contains, bytes.NewBufferString(test.Haystack), test.Needle)
		if err != nil {
			t.Errorf("FindWith(strings.Contains, %q, %q) errored out: %v", test.Haystack, test.Needle, err)
		}

		if !reflect.DeepEqual(test.Expected, found) {
			t.Errorf("FindWith(strings.Contains, %q, %q) result mismatch. Expected %#v, got %#v", test.Haystack, test.Needle, test.Expected, found)
		}
	}
}

func TestOccursWith(t *testing.T) {
	tests := []struct {
		Haystack string
		Needles  []string
		Expected []int
		Method   func(string, string) bool
	}{
		{"looking for this\nbut not for that\n", []string{""}, []int{}, strings.Contains},
		{"looking for this\nbut not for that\n", []string{"this"}, []int{1}, strings.Contains},
		{"looking for this\nbut not for that\n", []string{"that"}, []int{2}, strings.Contains},
		{"looking for this\nbut not for that\n", []string{"for"}, []int{1, 2}, strings.Contains},

		{"looking for this\nbut not for that", []string{""}, []int{}, strings.Contains},
		{"looking for this\nbut not for that", []string{"this"}, []int{1}, strings.Contains},
		{"looking for this\nbut not for that", []string{"that"}, []int{2}, strings.Contains},
		{"looking for this\nbut not for that", []string{"for"}, []int{1, 2}, strings.Contains},

		{"looking for this\r\nbut not for that\r\n", []string{""}, []int{}, strings.Contains},
		{"looking for this\r\nbut not for that\r\n", []string{"this"}, []int{1}, strings.Contains},
		{"looking for this\r\nbut not for that\r\n", []string{"that"}, []int{2}, strings.Contains},
		{"looking for this\r\nbut not for that\r\n", []string{"for"}, []int{1, 2}, strings.Contains},

		{"looking for this\r\nbut not for that", []string{""}, []int{}, strings.Contains},
		{"looking for this\r\nbut not for that", []string{"this"}, []int{1}, strings.Contains},
		{"looking for this\r\nbut not for that", []string{"that"}, []int{2}, strings.Contains},
		{"looking for this\r\nbut not for that", []string{"for"}, []int{1, 2}, strings.Contains},

		{"looking for this\nbut not for that\n", []string{""}, []int{}, strings.HasPrefix},
		{"looking for this\nbut not for that\n", []string{"looking"}, []int{1}, strings.HasPrefix},
		{"looking for this\nbut not for that\n", []string{"but"}, []int{2}, strings.HasPrefix},
		{"looking for this\nlooking for that\n", []string{"looking"}, []int{1, 2}, strings.HasPrefix},
	}

	for _, test := range tests {
		found, err := FindWith(test.Method, bytes.NewBufferString(test.Haystack), test.Needles)
		if err != nil {
			t.Errorf("FindWith(test.Method, %q, %q) errored out: %v", test.Haystack, test.Needles, err)
		}

		if !reflect.DeepEqual(test.Expected, found) {
			t.Errorf("FindWith(strings.Contains, %q, %q) result mismatch. Expected %#v, got %#v", test.Haystack, test.Needles, test.Expected, found)
		}
	}
}

func TestOccursWithInFile(t *testing.T) {
	tests := []struct {
		Needles  []string
		Expected []int
		Method   func(string, string) bool
	}{
		{[]string{""}, []int{}, strings.Contains},
		{[]string{"this"}, []int{1}, strings.Contains},
		{[]string{"that"}, []int{2}, strings.Contains},
		{[]string{"for"}, []int{1, 2}, strings.Contains},

		{[]string{""}, []int{}, strings.Contains},
		{[]string{"this"}, []int{1}, strings.Contains},
		{[]string{"that"}, []int{2}, strings.Contains},
		{[]string{"for"}, []int{1, 2}, strings.Contains},

		{[]string{""}, []int{}, strings.Contains},
		{[]string{"this"}, []int{1}, strings.Contains},
		{[]string{"that"}, []int{2}, strings.Contains},
		{[]string{"for"}, []int{1, 2}, strings.Contains},

		{[]string{""}, []int{}, strings.HasPrefix},
		{[]string{"looking"}, []int{1, 2}, strings.HasPrefix},
		{[]string{"but"}, []int{}, strings.HasPrefix},
		{[]string{"looking"}, []int{1, 2}, strings.HasPrefix},
	}

	tmpfile, err := ioutil.TempFile("./", "test")
	if err != nil {
		t.Errorf("failed creating temp file: %v", err)
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	fmt.Fprintf(tmpfile, "looking for this\n")
	fmt.Fprintf(tmpfile, "looking for that\n")

	reader := ioutil.NopCloser(tmpfile)
	for _, test := range tests {
		found, err := FindWith(test.Method, reader, test.Needles)
		if err != nil {
			t.Errorf("FindWith(test.Method, tempfile, %q) errored out: %v", test.Needles, err)
		}

		if !reflect.DeepEqual(test.Expected, found) {
			t.Errorf("FindWith(test.Method, tmpfile, %q) result mismatch. Expected %#v, got %#v", test.Needles, test.Expected, found)
		}

		tmpfile.Seek(0, 0)
	}
}

func TestCopyLines(t *testing.T) {
	tests := []struct {
		From     string
		Lines    []int
		Expected string
	}{
		{"LineOne\nLineTwo\nLineThree\n", []int{1, 2, 3}, "LineOne\nLineTwo\nLineThree\n"},
		{"LineOne\nLineTwo\nLineThree\n", []int{1}, "LineOne\n"},
		{"LineOne\nLineTwo\nLineThree\n", []int{2}, "LineTwo\n"},
		{"LineOne\nLineTwo\nLineThree\n", []int{3}, "LineThree\n"},
		{"LineOne\nLineTwo\nLineThree\n", []int{1, 2}, "LineOne\nLineTwo\n"},
		{"LineOne\nLineTwo\nLineThree\n", []int{1, 3}, "LineOne\nLineThree\n"},
		{"LineOne\nLineTwo\nLineThree\n", []int{}, ""},

		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{1, 2, 3}, "LineOne\nLineTwo\nLineThree\n"},
		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{1}, "LineOne\n"},
		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{2}, "LineTwo\n"},
		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{3}, "LineThree\n"},
		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{1, 2}, "LineOne\nLineTwo\n"},
		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{1, 3}, "LineOne\nLineThree\n"},
		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{}, ""},

		{"LineOne\r\nLineTwo\r\nLineThree", []int{1, 2, 3}, "LineOne\nLineTwo\nLineThree\n"},
		{"LineOne\r\nLineTwo\r\nLineThree", []int{1}, "LineOne\n"},
		{"LineOne\r\nLineTwo\r\nLineThree", []int{2}, "LineTwo\n"},
		{"LineOne\r\nLineTwo\r\nLineThree", []int{3}, "LineThree\n"},
		{"LineOne\r\nLineTwo\r\nLineThree", []int{1, 2}, "LineOne\nLineTwo\n"},
		{"LineOne\r\nLineTwo\r\nLineThree", []int{1, 3}, "LineOne\nLineThree\n"},
		{"LineOne\r\nLineTwo\r\nLineThree", []int{}, ""},
	}

	for _, test := range tests {
		var to bytes.Buffer

		err := CopyLines(bytes.NewBufferString(test.From), test.Lines, &to)
		if err != nil {
			t.Errorf("CopyLines(%q, %v, %v) errored out: %v", test.From, test.Lines, to, err)
		}

		read, _ := ioutil.ReadAll(&to)

		if string(read) != test.Expected {
			t.Errorf("CopyLines(%q, %v, ..) result mismatch. Expected %q, got %q", test.From, test.Lines, test.Expected, string(read))
		}
	}
}

func TestCopyWithoutLines(t *testing.T) {
	tests := []struct {
		From     string
		Lines    []int
		Expected string
	}{
		{"LineOne\nLineTwo\nLineThree\n", []int{}, "LineOne\nLineTwo\nLineThree\n"},
		{"LineOne\nLineTwo\nLineThree\n", []int{2, 3}, "LineOne\n"},
		{"LineOne\nLineTwo\nLineThree\n", []int{1, 3}, "LineTwo\n"},
		{"LineOne\nLineTwo\nLineThree\n", []int{1, 2}, "LineThree\n"},
		{"LineOne\nLineTwo\nLineThree\n", []int{3}, "LineOne\nLineTwo\n"},
		{"LineOne\nLineTwo\nLineThree\n", []int{2}, "LineOne\nLineThree\n"},
		{"LineOne\nLineTwo\nLineThree\n", []int{}, "LineOne\nLineTwo\nLineThree\n"},

		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{}, "LineOne\nLineTwo\nLineThree\n"},
		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{2, 3}, "LineOne\n"},
		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{1, 3}, "LineTwo\n"},
		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{1, 2}, "LineThree\n"},
		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{3}, "LineOne\nLineTwo\n"},
		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{2}, "LineOne\nLineThree\n"},
		{"LineOne\r\nLineTwo\r\nLineThree\r\n", []int{}, "LineOne\nLineTwo\nLineThree\n"},

		{"LineOne\r\nLineTwo\r\nLineThree", []int{}, "LineOne\nLineTwo\nLineThree\n"},
		{"LineOne\r\nLineTwo\r\nLineThree", []int{2, 3}, "LineOne\n"},
		{"LineOne\r\nLineTwo\r\nLineThree", []int{1, 3}, "LineTwo\n"},
		{"LineOne\r\nLineTwo\r\nLineThree", []int{1, 2}, "LineThree\n"},
		{"LineOne\r\nLineTwo\r\nLineThree", []int{3}, "LineOne\nLineTwo\n"},
		{"LineOne\r\nLineTwo\r\nLineThree", []int{2}, "LineOne\nLineThree\n"},
		{"LineOne\r\nLineTwo\r\nLineThree", []int{}, "LineOne\nLineTwo\nLineThree\n"},
	}

	for _, test := range tests {
		var to bytes.Buffer

		err := CopyWithoutLines(bytes.NewBufferString(test.From), test.Lines, &to)
		if err != nil {
			t.Errorf("CopyWithoutLines(%q, %v, %v) errored out: %v", test.From, test.Lines, to, err)
		}

		read, _ := ioutil.ReadAll(&to)

		if string(read) != test.Expected {
			t.Errorf("CopyWithoutLines(%q, %v, ..) result mismatch. Expected %q, got %q", test.From, test.Lines, test.Expected, string(read))
		}
	}
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
		FindWith(strings.Contains, bytes.NewBufferString(testString), []string{"weary"})
	}
}

func BenchmarkFindWithStringsHasPrefix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindWith(strings.HasPrefix, bytes.NewBufferString(testString), []string{"weary"})
	}
}

func BenchmarkFindWithStringsContainsAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindWith(strings.ContainsAny, bytes.NewBufferString(testString), []string{"weary"})
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
