package learn_defer_test

import (
	learndefer "go-learn/defer"
	"testing"
)

func TestDeferFor(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试for循环里面的defer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learndefer.DeferFor()
		})
	}
}

func TestDeferForParams(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试参数在defer中的传值"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learndefer.DeferForParams()
		})
	}
}

func TestDeferForReturn(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"测试参数在defer中的传值", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := learndefer.DeferForReturn(); got != tt.want {
				t.Errorf("DeferForReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeferForReturn2(t *testing.T) {
	tests := []struct {
		name  string
		wantX int
	}{
		{"测试参数在defer中的传值", 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotX := learndefer.DeferForReturn2(); gotX != tt.wantX {
				t.Errorf("DeferForReturn2() = %v, want %v", gotX, tt.wantX)
			}
		})
	}
}

func TestDeferForPanic(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试defer中的panic"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learndefer.DeferForPanic()
		})
	}
}

func TestDeferForOldParams(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试defer还原旧值的方法"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learndefer.DeferForOldParams()
		})
	}
}
