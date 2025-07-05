package learn_error_test

import (
	learnerror "go-learn/error"
	"testing"
)

func TestMyErrorEqualsNil(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试自定义error结构体空值"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learnerror.MyErrorEqualsNil()
		})
	}
}
