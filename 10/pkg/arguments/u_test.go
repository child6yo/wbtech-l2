package arguments

import (
	"bufio"
	"strings"
	"testing"
)

func TestPrintUniqueOnly(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		output string
	}{
		{
			name:   "no duplicates",
			input:  []string{"a", "b", "c"},
			output: "a\nb\nc\n",
		},
		{
			name:   "all duplicates",
			input:  []string{"x", "x", "x"},
			output: "x\n",
		},
		{
			name:   "mixed duplicates",
			input:  []string{"cat", "dog", "cat", "bird", "dog", "cat"},
			output: "cat\ndog\nbird\n",
		},
		{
			name:   "empty input",
			input:  []string{},
			output: "",
		},
		{
			name:   "single element",
			input:  []string{"hello"},
			output: "hello\n",
		},
		{
			name:   "empty strings",
			input:  []string{"", "a", "", "b", ""},
			output: "\na\nb\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			writer := bufio.NewWriter(&sb)

			PrintUniqueOnly(writer, tt.input)

			writer.Flush()

			got := sb.String()
			if got != tt.output {
				t.Fatalf("expected %q, got %q", tt.output, got)
			}
		})
	}
}
