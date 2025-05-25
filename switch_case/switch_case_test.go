package switch_case

import "testing"

func TestSwitchCase(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试SwitchCase"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SwitchCase()
		})
	}
}

func TestSwitchCaseMulti(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试SwitchCase多个条件执行同样操作"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SwitchCaseMulti()
		})
	}
}
