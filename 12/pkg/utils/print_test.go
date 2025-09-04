package utils

import (
	"bufio"
	"bytes"
	"testing"
)

func TestPrintResults(t *testing.T) {
	tests := []struct {
		name            string
		lines           []string
		matches         []int
		beforeContext   int
		afterContext    int
		showLineNumbers bool
		expectedOutput  string
	}{
		{
			name:            "no context, no line numbers",
			lines:           []string{"apple", "banana", "cherry", "date"},
			matches:         []int{1, 3},
			beforeContext:   0,
			afterContext:    0,
			showLineNumbers: false,
			expectedOutput:  "banana\ndate\n",
		},
		{
			name:            "no context, with line numbers",
			lines:           []string{"apple", "banana", "cherry", "date"},
			matches:         []int{1, 3},
			beforeContext:   0,
			afterContext:    0,
			showLineNumbers: true,
			expectedOutput:  "2:banana\n4:date\n",
		},
		{
			name:            "with before context",
			lines:           []string{"one", "two", "three", "four"},
			matches:         []int{2},
			beforeContext:   1,
			afterContext:    0,
			showLineNumbers: false,
			expectedOutput:  "two\nthree\n",
		},
		{
			name:            "with after context",
			lines:           []string{"one", "two", "three", "four"},
			matches:         []int{1},
			beforeContext:   0,
			afterContext:    1,
			showLineNumbers: false,
			expectedOutput:  "two\nthree\n",
		},
		{
			name:            "with context around",
			lines:           []string{"a", "b", "c", "d", "e"},
			matches:         []int{2},
			beforeContext:   1,
			afterContext:    1,
			showLineNumbers: false,
			expectedOutput:  "b\nc\nd\n",
		},
		{
			name:            "overlapping matches with context",
			lines:           []string{"a", "b", "c", "d", "e"},
			matches:         []int{1, 2},
			beforeContext:   1,
			afterContext:    1,
			showLineNumbers: false,
			expectedOutput:  "a\nb\nc\nd\n",
		},
		{
			name:            "context with line numbers",
			lines:           []string{"first", "second", "third", "fourth"},
			matches:         []int{2},
			beforeContext:   1,
			afterContext:    1,
			showLineNumbers: true,
			expectedOutput:  "2:second\n3:third\n4:fourth\n",
		},
		{
			name:            "match at start with before context",
			lines:           []string{"a", "b", "c"},
			matches:         []int{0},
			beforeContext:   2,
			afterContext:    1,
			showLineNumbers: false,
			expectedOutput:  "a\nb\n",
		},
		{
			name:            "match at end with after context",
			lines:           []string{"a", "b", "c"},
			matches:         []int{2},
			beforeContext:   1,
			afterContext:    2,
			showLineNumbers: false,
			expectedOutput:  "b\nc\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			PrintResults(writer, tt.lines, tt.matches, tt.beforeContext, tt.afterContext, tt.showLineNumbers)

			writer.Flush()
			got := buf.String()
			if got != tt.expectedOutput {
				t.Errorf("expected %q, got %q", tt.expectedOutput, got)
			}
		})
	}
}

func TestPrintOnlyNumLines(t *testing.T) {
	tests := []struct {
		name           string
		matches        []int
		expectedOutput string
	}{
		{
			name:           "three matches",
			matches:        []int{0, 2, 4},
			expectedOutput: "3\n",
		},
		{
			name:           "one match",
			matches:        []int{5},
			expectedOutput: "1\n",
		},
		{
			name:           "no matches",
			matches:        []int{},
			expectedOutput: "0\n",
		},
		{
			name:           "five matches",
			matches:        []int{1, 2, 3, 4, 5},
			expectedOutput: "5\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			PrintOnlyNumLines(writer, tt.matches)

			writer.Flush()
			got := buf.String()
			if got != tt.expectedOutput {
				t.Errorf("expected %q, got %q", tt.expectedOutput, got)
			}
		})
	}
}
