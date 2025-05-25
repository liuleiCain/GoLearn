package for_range

import "testing"

func TestBreakForSelect(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试for-select的break"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BreakForSelect()
		})
	}
}

func TestBreakForSelectLabel(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试label的break"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BreakForSelectLabel()
		})
	}
}

func TestBreakForLabel(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试多重循环label的break"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BreakForLabel()
		})
	}
}
