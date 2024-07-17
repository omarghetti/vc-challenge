package util

import "testing"

func TestIsStopword(t *testing.T) {
	tests := []struct {
		name string
		word string
		want bool
	}{
		{
			name: "stopword",
			word: "between",
			want: true,
		},
		{
			name: "not stopword",
			word: "hello",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStopword(tt.word); got != tt.want {
				t.Errorf("IsStopword() = %v, want %v", got, tt.want)
			}
		})
	}
}
