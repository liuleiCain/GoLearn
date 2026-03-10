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

func TestCustomError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试自定义错误类型"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learnerror.CustomError()
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试错误包装"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learnerror.ErrorWrapping()
		})
	}
}

func TestErrorUnwrap(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试错误解包"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learnerror.ErrorUnwrap()
		})
	}
}

func TestErrorIs(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试errors.Is判断错误"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learnerror.ErrorIs()
		})
	}
}

func TestErrorAs(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试errors.As转换错误类型"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learnerror.ErrorAs()
		})
	}
}

func TestPanicRecover(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试Panic和Recover"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learnerror.PanicRecover()
		})
	}
}

func TestMultiErrorDemo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试多错误处理"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learnerror.MultiErrorDemo()
		})
	}
}

func TestErrorWithContextDemo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试带上下文的错误"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learnerror.ErrorWithContextDemo()
		})
	}
}

func TestDeferPanicOrder(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试Defer、Panic、Recover执行顺序"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learnerror.DeferPanicOrder()
		})
	}
}
