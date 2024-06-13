package gapbuffer

import (
	"testing"
)

func TestGapBufferInsert(t *testing.T) {
	type insertTest struct {
		inserts  []string
		expected string
	}

	var insertTests = []insertTest{
		{[]string{"welcome"}, "welcome"},
		{[]string{"hello", ", world!"}, "hello, world!"},
	}

	for _, test := range insertTests {
		gapBuffer := NewGapBuffer(1)
		for _, insert := range test.inserts {
			gapBuffer.Insert(insert)
		}
		output := gapBuffer.String()
		if output != test.expected {
			t.Errorf("Incorrect output %s, expected %s", output, test.expected)
		}
	}
}

func TestGapBufferMoveCursor(t *testing.T) {
	type moveAndInsert struct {
		position int
		insert   string
	}

	type moveCursorTest struct {
		moves    []moveAndInsert
		expected string
	}

	var moveCursorTests = []moveCursorTest{
		{[]moveAndInsert{{0, "world!"}, {0, "Hello, "}}, "Hello, world!"},
		{[]moveAndInsert{{0, "Hoyu?"}, {2, "w are "}, {9, "o"}}, "How are you?"},
	}

	for _, test := range moveCursorTests {
		gapBuffer := NewGapBuffer(1)
		for _, move := range test.moves {
			gapBuffer.MoveCursor(move.position)
			gapBuffer.Insert(move.insert)
		}
		output := gapBuffer.String()
		if output != test.expected {
			t.Errorf("Incorrect output %s, expected %s", output, test.expected)
		}
	}
}
