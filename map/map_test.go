package learn_map

import (
	"testing"
)

func TestRangeMap(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "测试map遍历的顺序",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RangeMap()
		})
	}
}

func TestInitMap(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "测试map初始化",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitMap()
		})
	}
}
