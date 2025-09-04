package utils

import (
	"regexp"
	"strings"
	"testing"
)

func TestFindMatchesIndexes(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		matcher  Matcher
		invert   bool
		expected []int
	}{
		{
			name:     "no matches",
			lines:    []string{"apple", "banana", "cherry"},
			matcher:  func(s string) bool { return strings.Contains(s, "xyz") },
			invert:   false,
			expected: []int{},
		},
		{
			name:     "one match",
			lines:    []string{"apple", "banana", "cherry"},
			matcher:  func(s string) bool { return strings.Contains(s, "banana") },
			invert:   false,
			expected: []int{1},
		},
		{
			name:     "multiple matches",
			lines:    []string{"test", "Test", "best", "rest"},
			matcher:  func(s string) bool { return strings.Contains(strings.ToLower(s), "test") },
			invert:   false,
			expected: []int{0, 1},
		},
		{
			name:     "all lines match",
			lines:    []string{"a", "aa", "aaa"},
			matcher:  func(s string) bool { return strings.Contains(s, "a") },
			invert:   false,
			expected: []int{0, 1, 2},
		},
		{
			name:     "invert: no matches become all",
			lines:    []string{"apple", "banana", "cherry"},
			matcher:  func(s string) bool { return strings.Contains(s, "xyz") },
			invert:   true,
			expected: []int{0, 1, 2},
		},
		{
			name:     "invert: one match excluded",
			lines:    []string{"apple", "banana", "cherry"},
			matcher:  func(s string) bool { return strings.Contains(s, "banana") },
			invert:   true,
			expected: []int{0, 2},
		},
		{
			name:     "invert: multiple matches excluded",
			lines:    []string{"test", "Test", "best", "rest"},
			matcher:  func(s string) bool { return strings.Contains(strings.ToLower(s), "test") },
			invert:   true,
			expected: []int{2, 3},
		},
		{
			name:     "invert: all matches excluded",
			lines:    []string{"a", "aa", "aaa"},
			matcher:  func(s string) bool { return strings.Contains(s, "a") },
			invert:   true,
			expected: []int{},
		},
		{
			name:  "regex match",
			lines: []string{"error", "warning", "info", "critical error"},
			matcher: func(s string) bool {
				re := regexp.MustCompile(`error`)
				return re.MatchString(s)
			},
			invert:   false,
			expected: []int{0, 3},
		},
		{
			name:  "invert regex",
			lines: []string{"error", "warning", "info", "critical error"},
			matcher: func(s string) bool {
				re := regexp.MustCompile(`error`)
				return re.MatchString(s)
			},
			invert:   true,
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMatchesIndexes(tt.lines, tt.matcher, tt.invert)
			if len(got) != len(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Fatalf("expected %v, got %v", tt.expected, got)
				}
			}
		})
	}
}