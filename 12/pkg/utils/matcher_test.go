package utils

import (
	"testing"
)

func TestCompileMatcher(t *testing.T) {
	tests := []struct {
		name            string
		pattern         string
		caseInsensitive bool
		fixedTemplate   bool
		input           string
		expected        bool
	}{
		{
			name:            "fixed case-sensitive match",
			pattern:         "hello",
			caseInsensitive: false,
			fixedTemplate:   true,
			input:           "hello world",
			expected:        true,
		},
		{
			name:            "fixed case-sensitive no match",
			pattern:         "hello",
			caseInsensitive: false,
			fixedTemplate:   true,
			input:           "Hello world",
			expected:        false,
		},
		{
			name:            "fixed case-insensitive match",
			pattern:         "hello",
			caseInsensitive: true,
			fixedTemplate:   true,
			input:           "Hello World",
			expected:        true,
		},
		{
			name:            "fixed case-insensitive no match",
			pattern:         "xyz",
			caseInsensitive: true,
			fixedTemplate:   true,
			input:           "Hello World",
			expected:        false,
		},
		{
			name:            "regex exact match",
			pattern:         "^start",
			caseInsensitive: false,
			fixedTemplate:   false,
			input:           "start of string",
			expected:        true,
		},
		{
			name:            "regex exact no match",
			pattern:         "^start",
			caseInsensitive: false,
			fixedTemplate:   false,
			input:           "begin of string",
			expected:        false,
		},
		{
			name:            "regex case-insensitive match",
			pattern:         "ERROR",
			caseInsensitive: true,
			fixedTemplate:   false,
			input:           "An error occurred",
			expected:        true,
		},
		{
			name:            "regex case-insensitive no match",
			pattern:         "CRITICAL",
			caseInsensitive: true,
			fixedTemplate:   false,
			input:           "An error occurred",
			expected:        false,
		},
		{
			name:            "regex with meta characters",
			pattern:         "a.b",
			caseInsensitive: false,
			fixedTemplate:   false,
			input:           "axb",
			expected:        true,
		},
		{
			name:            "regex meta no match",
			pattern:         "a.b",
			caseInsensitive: false,
			fixedTemplate:   false,
			input:           "ab",
			expected:        false,
		},
		{
			name:            "fixed template with regex meta",
			pattern:         "a.b",
			caseInsensitive: false,
			fixedTemplate:   true,
			input:           "a.b in text",
			expected:        true,
		},
		{
			name:            "fixed template with regex meta no match",
			pattern:         "a.b",
			caseInsensitive: false,
			fixedTemplate:   true,
			input:           "axb",
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := CompileMatcher(tt.pattern, tt.caseInsensitive, tt.fixedTemplate)
			result := matcher(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
