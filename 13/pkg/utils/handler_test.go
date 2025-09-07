package utils

import (
	"testing"
)

func TestHandleString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter string
		fields    map[int]struct{}
		separated bool
		expected  string
		success   bool
	}{
		{
			name:      "select multiple fields by index",
			input:     "a,b,c,d",
			delimiter: ",",
			fields:    map[int]struct{}{1: {}, 3: {}},
			separated: false,
			expected:  "a,c",
			success:   true,
		},
		{
			name:      "select all fields in order",
			input:     "x\ty\tz",
			delimiter: "\t",
			fields:    map[int]struct{}{1: {}, 2: {}, 3: {}},
			separated: false,
			expected:  "x\ty\tz",
			success:   true,
		},
		{
			name:      "no fields selected",
			input:     "one,two,three",
			delimiter: ",",
			fields:    map[int]struct{}{},
			separated: false,
			expected:  "",
			success:   true,
		},
		{
			name:      "fields is nil",
			input:     "a,b,c",
			delimiter: ",",
			fields:    nil,
			separated: false,
			expected:  "",
			success:   true,
		},
		{
			name:      "some indices out of range",
			input:     "p1,p2,p3",
			delimiter: ",",
			fields:    map[int]struct{}{2: {}, 5: {}},
			separated: false,
			expected:  "p2",
			success:   true,
		},
		{
			name:      "single part but separated=true",
			input:     "alone",
			delimiter: "|",
			fields:    map[int]struct{}{1: {}},
			separated: true,
			expected:  "",
			success:   false,
		},
		{
			name:      "two parts and separated=true",
			input:     "first|second",
			delimiter: "|",
			fields:    map[int]struct{}{2: {}},
			separated: true,
			expected:  "second",
			success:   true,
		},
		{
			name:      "multiple delimiters and separated=true",
			input:     "a:b:c:d",
			delimiter: ":",
			fields:    map[int]struct{}{1: {}, 4: {}},
			separated: true,
			expected:  "a:d",
			success:   true,
		},
		{
			name:      "empty string input",
			input:     "",
			delimiter: ",",
			fields:    map[int]struct{}{1: {}},
			separated: false,
			expected:  "",
			success:   true,
		},
		{
			name:      "delimiter not found",
			input:     "nodelem",
			delimiter: ",",
			fields:    map[int]struct{}{1: {}},
			separated: false,
			expected:  "nodelem",
			success:   true,
		},
		{
			name:      "delimiter not found but separated=true",
			input:     "single",
			delimiter: ",",
			fields:    map[int]struct{}{1: {}},
			separated: true,
			expected:  "",
			success:   false,
		},
		{
			name:      "consecutive delimiters",
			input:     "a,,b,",
			delimiter: ",",
			fields:    map[int]struct{}{1: {}, 2: {}, 4: {}},
			separated: false,
			expected:  "a,,",
			success:   true,
		},
		{
			name:      "non contiguous valid indices",
			input:     "f1,f2,f3",
			delimiter: ",",
			fields:    map[int]struct{}{1: {}, 100: {}},
			separated: false,
			expected:  "f1",
			success:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := HandleString(tt.input, tt.delimiter, tt.fields, tt.separated)
			if ok != tt.success {
				t.Fatalf("expected success=%v, got %v", tt.success, ok)
			}
			if !tt.success {
				return
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
