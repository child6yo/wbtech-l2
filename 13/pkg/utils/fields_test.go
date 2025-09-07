package utils

import (
	"reflect"
	"testing"
)

func TestParseFields(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[int]struct{}
		hasError bool
	}{
		{
			name:     "single field",
			input:    "1",
			expected: map[int]struct{}{1: {}},
			hasError: false,
		},
		{
			name:     "multiple fields",
			input:    "1,2,3",
			expected: map[int]struct{}{1: {}, 2: {}, 3: {}},
			hasError: false,
		},
		{
			name:     "range",
			input:    "5-7",
			expected: map[int]struct{}{5: {}, 6: {}, 7: {}},
			hasError: false,
		},
		{
			name:     "mixed fields and range",
			input:    "1,3,5-7,10",
			expected: map[int]struct{}{1: {}, 3: {}, 5: {}, 6: {}, 7: {}, 10: {}},
			hasError: false,
		},
		{
			name:     "reversed range",
			input:    "7-5",
			expected: nil,
			hasError: true,
		},
		{
			name:     "negative number",
			input:    "-1",
			expected: nil,
			hasError: true,
		},
		{
			name:     "negative in range",
			input:    "1--3",
			expected: nil,
			hasError: true,
		},
		{
			name:     "zero field",
			input:    "0",
			expected: nil,
			hasError: true,
		},
		{
			name:     "range with zero",
			input:    "0-2",
			expected: nil,
			hasError: true,
		},
		{
			name:     "malformed range: too many parts",
			input:    "1-2-3",
			expected: nil,
			hasError: true,
		},
		{
			name:     "non-integer input",
			input:    "abc",
			expected: nil,
			hasError: true,
		},
		{
			name:     "non-integer in range",
			input:    "1-abc",
			expected: nil,
			hasError: true,
		},
		{
			name:     "trailing comma",
			input:    "1,2,",
			expected: map[int]struct{}{1: {}, 2: {}},
			hasError: false,
		},
		{
			name:     "leading comma",
			input:    ",1,2",
			expected: map[int]struct{}{1: {}, 2: {}},
			hasError: false,
		},
		{
			name:     "spaces around numbers",
			input:    "  1 , 2 , 3  ",
			expected: map[int]struct{}{1: {}, 2: {}, 3: {}},
			hasError: false,
		},
		{
			name:     "spaces in range",
			input:    "5 - 7",
			expected: map[int]struct{}{5: {}, 6: {}, 7: {}},
			hasError: false,
		},
		{
			name:     "duplicate fields",
			input:    "1,1,2",
			expected: map[int]struct{}{1: {}, 2: {}},
			hasError: false,
		},
		{
			name:     "empty input",
			input:    "",
			expected: map[int]struct{}{},
			hasError: false,
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: map[int]struct{}{},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseFields(tt.input)

			if tt.hasError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
