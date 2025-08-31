package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	testCases := []struct {
		input   string
		wantErr bool
		output  string
	}{
		{
			input:   "a4bc2d5e",
			wantErr: false,
			output:  "aaaabccddddde",
		},
		{
			input:   "abcd",
			wantErr: false,
			output:  "abcd",
		},
		{
			input:   "45",
			wantErr: true,
			output:  "",
		},
		{
			input:   "",
			wantErr: false,
			output:  "",
		},
		{
			input:   `qwe\4\5`,
			wantErr: false,
			output:  "qwe45",
		},
		{
			input:   `qwe\45`,
			wantErr: false,
			output:  "qwe44444",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			res, err := unpackString(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got none")
				}
				return
			}
			if err != nil && !tc.wantErr {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if res != tc.output {
				t.Errorf("want %q, got %q", tc.output, res)
			}
		})
	}
}
