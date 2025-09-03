package arguments

import "testing"

func TestReverse(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "normal strings",
			input: []string{"a", "b", "c"},
			want:  []string{"c", "b", "a"},
		},
		{
			name:  "two elements",
			input: []string{"first", "last"},
			want:  []string{"last", "first"},
		},
		{
			name:  "single element",
			input: []string{"only"},
			want:  []string{"only"},
		},
		{
			name:  "empty slice",
			input: []string{},
			want:  []string{},
		},
		{
			name:  "palindrome",
			input: []string{"x", "y", "x"},
			want:  []string{"x", "y", "x"},
		},
		{
			name:  "mixed length strings",
			input: []string{"short", "very long string", "tiny"},
			want:  []string{"tiny", "very long string", "short"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]string, len(tt.input))
			copy(got, tt.input)
			Reverse(got)

			if len(got) != len(tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("expected %v, got %v", tt.want, got)
				}
			}
		})
	}
}
