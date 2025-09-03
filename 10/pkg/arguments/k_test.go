package arguments

import "testing"

func TestColSort(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		col     int
		numeric bool
		want    []string
	}{
		{
			name:    "sort by first column alphabetically",
			input:   []string{"b\t3", "a\t2", "c\t1"},
			col:     0,
			numeric: false,
			want:    []string{"a\t2", "b\t3", "c\t1"},
		},
		{
			name:    "sort by second column alphabetically",
			input:   []string{"a\t3", "a\t1", "a\t2"},
			col:     1,
			numeric: false,
			want:    []string{"a\t1", "a\t2", "a\t3"},
		},
		{
			name:    "sort by first column numerically",
			input:   []string{"2\tname", "1\tname", "10\tname"},
			col:     0,
			numeric: true,
			want:    []string{"1\tname", "2\tname", "10\tname"},
		},
		{
			name:    "sort by second column numerically",
			input:   []string{"x\t10", "x\t2", "x\t1"},
			col:     1,
			numeric: true,
			want:    []string{"x\t1", "x\t2", "x\t10"},
		},
		{
			name:    "column out of range",
			input:   []string{"a", "b\tc", "d\t"},
			col:     5,
			numeric: false,
			want:    []string{"a", "b\tc", "d\t"},
		},
		{
			name:    "empty strings in column",
			input:   []string{"\tx", "a\ty", "\tz"},
			col:     0,
			numeric: false,
			want:    []string{"\tx", "\tz", "a\ty"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]string, len(tt.input))
			copy(got, tt.input)
			ColSort(&got, tt.col, tt.numeric)

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
