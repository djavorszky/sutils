package sutils

import "testing"
import "bytes"

type testcase struct {
	Value    string
	Search   string
	Expected bool
}

func TestIContains(t *testing.T) {
	pairs := []testcase{
		// check exact matches with varying cases
		testcase{"hello", "hello", true},
		testcase{"hello", "HELLO", true},
		testcase{"HELLO", "hello", true},
		testcase{"heLLo", "hello", true},
		testcase{"hello", "HeLLo", true},
		// check with space
		testcase{"HELLO there", "hello", true},
		testcase{"heLLo there", "lo th", true},
		testcase{"hello there", "LO TH", true},
		// check submatches with varying cases
		testcase{"hello", "el", true},
		testcase{"hello", "EL", true},
		testcase{"hello", "eL", true},
		// check false hits
		testcase{"hello", "oh", false},
		testcase{"hello", "OH", false},
		testcase{"hello", "oH", false},
		testcase{"this is", "ss", false},
		testcase{"this is", "sS", false},
		testcase{"this is", "SS", false},
		testcase{"", "", false},
		testcase{"hello", "", false},
		testcase{"", "hello", false},
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

	result, _ = CountIgnoreCase(buf, "my")
	expected = 2

	if result != expected {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expected, result)
	}

	buf = bytes.NewBufferString(testString)

	result, _ = CountIgnoreCase(buf, "MY")
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

	result, _ = CountCaseSensitive(buf, "my")
	expected = 2

	if result != expected {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expected, result)
	}

	buf = bytes.NewBufferString(testString)

	result, _ = CountCaseSensitive(buf, "MY")
	expected = 0

	if result != expected {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expected, result)
	}
}

func TestFindIgnoreCase(t *testing.T) {
	var (
		result      []int
		expectedLen int

		expectedLineNumbers = []int{3, 4}

		buf = bytes.NewBufferString(testString)
	)

	result, _ = FindIgnoreCase(buf, "my")
	expectedLen = 2

	if len(result) != expectedLen {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expectedLen, len(result))
	}

	if result[0] != expectedLineNumbers[0] || result[1] != expectedLineNumbers[1] {
		t.Errorf("Error. Found match at %d and %d, but should've been at %d and %d", result[0], result[1], expectedLineNumbers[0], expectedLineNumbers[1])
	}

	buf = bytes.NewBufferString(testString)

	result, _ = FindIgnoreCase(buf, "MY")
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

		expectedLineNumbers = []int{3, 4}

		buf = bytes.NewBufferString(testString)
	)

	result, _ = FindCaseSensitive(buf, "my")
	expectedLen = 2

	if len(result) != expectedLen {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expectedLen, len(result))
	}

	if result[0] != expectedLineNumbers[0] || result[1] != expectedLineNumbers[1] {
		t.Errorf("Error. Found match at %d and %d, but should've been at %d and %d", result[0], result[1], expectedLineNumbers[0], expectedLineNumbers[1])
	}

	buf = bytes.NewBufferString(testString)

	result, _ = FindCaseSensitive(buf, "MY")
	expectedLen = 0

	if len(result) != expectedLen {
		t.Errorf("Mismatch. Expected count=%d, got result=%d", expectedLen, len(result))
	}
}

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
