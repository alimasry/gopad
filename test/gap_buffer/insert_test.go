package gapbuffer_test

import (
	"testing"

	"github.com/alimasry/gopad/internal/pkg/gapbuffer"
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
		gapBuffer := gapbuffer.NewGapBuffer(1)
		for _, insert := range test.inserts {
			gapBuffer.Insert(insert)
		}
		output := gapBuffer.String()
		if output != test.expected {
			t.Errorf("Incorrect output %s, expected %s", output, test.expected)
		}
	}
}
