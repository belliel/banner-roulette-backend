package hash

import (
	"testing"
	"unicode/utf8"
)

func TestRandHash_Hash(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			"first",
			32,
		},
		{
			"second",
			32,
		},
		{
			"third",
			32,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RandHash{}
			if got := r.Hash(); utf8.RuneCountInString(got) != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
