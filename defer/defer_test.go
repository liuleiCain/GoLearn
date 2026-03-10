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

func TestDeferTiming(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试defer执行时机"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learndefer.DeferTiming()
		})
	}
}

func TestDeferStackTrace(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试defer打印堆栈信息"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learndefer.DeferStackTrace()
		})
	}
}

func TestDeferMutexUnlock(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试defer解锁互斥锁"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learndefer.DeferMutexUnlock()
		})
	}
}

func TestDeferFileOperation(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试defer文件操作"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learndefer.DeferFileOperation()
		})
	}
}

func TestDeferMethod(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试defer调用方法"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learndefer.DeferMethod()
		})
	}
}

func TestDeferPerformance(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试defer性能对比"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learndefer.DeferPerformance()
		})
	}
}

func TestDeferNilFunction(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试defer nil函数"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learndefer.DeferNilFunction()
		})
	}
}

func TestDeferArguments(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试defer参数计算时机"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learndefer.DeferArguments()
		})
	}
}

func TestDeferNamedReturnMultiple(t *testing.T) {
	tests := []struct {
		name    string
		wantA   int
		wantB   string
		wantC   bool
	}{
		{"测试多命名返回值与defer", 20, "hello world", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotB, gotC := learndefer.DeferNamedReturnMultiple()
			if gotA != tt.wantA {
				t.Errorf("DeferNamedReturnMultiple() a = %v, want %v", gotA, tt.wantA)
			}
			if gotB != tt.wantB {
				t.Errorf("DeferNamedReturnMultiple() b = %v, want %v", gotB, tt.wantB)
			}
			if gotC != tt.wantC {
				t.Errorf("DeferNamedReturnMultiple() c = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
