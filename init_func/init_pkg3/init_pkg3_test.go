package init_pkg1

import "testing"

func TestInitPkgSort(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试init方法依赖执行顺序"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitPkgSort()
		})
	}
}
