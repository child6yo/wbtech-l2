package arguments

import "testing"

func TestNumSort(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "normal numbers",
			input: []string{"3", "1", "2"},
			want:  []string{"1", "2", "3"},
		},
		{
			name:  "already sorted",
			input: []string{"1", "2", "3"},
			want:  []string{"1", "2", "3"},
		},
		{
			name:  "reverse order",
			input: []string{"5", "4", "3", "2", "1"},
			want:  []string{"1", "2", "3", "4", "5"},
		},
		{
			name:  "single element",
			input: []string{"42"},
			want:  []string{"42"},
		},
		{
			name:  "empty slice",
			input: []string{},
			want:  []string{},
		},
		{
			name:  "duplicate numbers",
			input: []string{"2", "1", "2", "3", "1"},
			want:  []string{"1", "1", "2", "2", "3"},
		},
		{
			name:  "large numbers",
			input: []string{"100", "20", "3"},
			want:  []string{"3", "20", "100"},
		},
		{
			name:  "zero and positive",
			input: []string{"10", "0", "5"},
			want:  []string{"0", "5", "10"},
		},
		{
			name:  "invalid numbers only",
			input: []string{"xyz", "abc", "def"},
			want:  []string{"abc", "def", "xyz"},
		},
		{
			name:  "mixed valid and invalid",
			input: []string{"3", "invalid", "1", "bad", "2"},
			want:  []string{"1", "2", "3", "bad", "invalid"},
		},
		{
			name:  "invalid before valid",
			input: []string{"xyz", "1", "abc", "2"},
			want:  []string{"1", "2", "abc", "xyz"},
		},
		{
			name:  "numbers and symbols",
			input: []string{"5", "@@@", "1", "!!!", "3"},
			want:  []string{"1", "3", "5", "!!!", "@@@"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]string, len(tt.input))
			copy(got, tt.input)
			NumSort(&got)

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
