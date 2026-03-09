package json

import "testing"

func TestUnmarshalArr(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UnmarshalArr()
		})
	}
}
