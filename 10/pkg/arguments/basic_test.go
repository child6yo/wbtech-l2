package arguments

import "testing"

func TestBasicSort(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "normal strings",
			input: []string{"banana", "apple", "cherry"},
			want:  []string{"apple", "banana", "cherry"},
		},
		{
			name:  "already sorted",
			input: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "reverse order",
			input: []string{"z", "y", "x"},
			want:  []string{"x", "y", "z"},
		},
		{
			name:  "single element",
			input: []string{"single"},
			want:  []string{"single"},
		},
		{
			name:  "empty slice",
			input: []string{},
			want:  []string{},
		},
		{
			name:  "duplicate strings",
			input: []string{"cat", "dog", "cat", "apple"},
			want:  []string{"apple", "cat", "cat", "dog"},
		},
		{
			name:  "mixed case",
			input: []string{"Zebra", "apple", "Banana"},
			want:  []string{"Banana", "Zebra", "apple"},
		},
		{
			name:  "with empty strings",
			input: []string{"", "a", "", "b"},
			want:  []string{"", "", "a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]string, len(tt.input))
			copy(got, tt.input)
			BasicSort(&got)

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
