package gc

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

// GCStats 展示如何获取 GC 统计信息
func GCStats() {
	var stats debug.GCStats
	debug.ReadGCStats(&stats)

	fmt.Println("GC统计信息:")
	fmt.Println("  GC次数:", stats.NumGC)
	fmt.Println("  上次GC时间:", stats.LastGC)
	fmt.Println("  暂停总时间:", stats.PauseTotal)
	if len(stats.Pause) > 0 {
		fmt.Println("  最近暂停时间:", stats.Pause[0])
	}
}

// ManualGC 展示如何手动触发 GC
func ManualGC() {
	fmt.Println("手动触发 GC 前...")
	runtime.GC()
	fmt.Println("手动触发 GC 完成")
}

// MemoryAllocation 展示内存分配与 GC 的关系
func MemoryAllocation() {
	fmt.Println("开始内存分配测试...")

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("初始分配内存:", m.Alloc, "bytes")

	// 分配大量内存
	var data [][]byte
	for i := 0; i < 1000; i++ {
		data = append(data, make([]byte, 1024*10)) // 每次分配 10KB
	}

	runtime.ReadMemStats(&m)
	fmt.Println("分配后内存:", m.Alloc, "bytes")

	// 释放引用
	data = nil

	runtime.GC()

	runtime.ReadMemStats(&m)
	fmt.Println("GC后内存:", m.Alloc, "bytes")
}

// GCPercent 展示如何调整 GC 触发频率
func GCPercent() {
	original := debug.SetGCPercent(100)
	fmt.Println("原始 GOGC 值:", original)

	newValue := debug.SetGCPercent(200)
	fmt.Println("设置 GOGC 为:", newValue)

	// 恢复原始值
	debug.SetGCPercent(original)
	fmt.Println("恢复 GOGC 为:", original)
}

// Finalizer 展示如何使用终结器
func Finalizer() {
	type Resource struct {
		id   int
		name string
	}

	createResource := func(id int, name string) *Resource {
		r := &Resource{id: id, name: name}
		runtime.SetFinalizer(r, func(r *Resource) {
			fmt.Println("资源", r.name, "被 GC 回收")
		})
		return r
	}

	r := createResource(1, "test-resource")
	fmt.Println("创建资源:", r.name)

	// 释放引用
	r = nil

	// 触发 GC
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
}

// HeapProfile 展示如何获取堆内存分析
func HeapProfile() {
	fmt.Println("堆内存分析示例...")

	// 分配一些内存
	var data [][]byte
	for i := 0; i < 100; i++ {
		data = append(data, make([]byte, 1024*100))
	}

	// 获取堆 profile
	debug.WriteHeapDump(0)
	fmt.Println("堆 profile 已写入 (仅在调试模式下有效)")

	data = nil
}

// ConcurrentGC 展示并发场景下的 GC 行为
func ConcurrentGC() {
	var wg sync.WaitGroup

	// 启动多个 goroutine 分配内存
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var localData [][]byte
			for j := 0; j < 100; j++ {
				localData = append(localData, make([]byte, 1024))
				if j%20 == 0 {
					runtime.Gosched()
				}
			}
		}(i)
	}

	// 在主 goroutine 中监控 GC
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Println("当前内存分配:", m.Alloc, "bytes, GC次数:", m.NumGC)
		}
	}()

	wg.Wait()
	time.Sleep(200 * time.Millisecond)
}

// MemoryLeakDetection 展示如何检测内存泄漏
func MemoryLeakDetection() {
	fmt.Println("内存泄漏检测示例...")

	var leakyData [][]byte
	allocate := func() {
		leakyData = append(leakyData, make([]byte, 1024*1024)) // 1MB
	}

	// 模拟泄漏
	for i := 0; i < 5; i++ {
		allocate()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Println("分配", i+1, "次后:", m.Alloc, "bytes")
	}

	// 正常释放
	fmt.Println("释放引用...")
	leakyData = nil
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("GC后内存:", m.Alloc, "bytes")
}

// GCControl 展示如何控制 GC 行为
func GCControl() {
	fmt.Println("GC 控制示例...")

	// 禁用 GC
	debug.SetGCPercent(-1)
	fmt.Println("GC 已禁用")

	// 分配内存
	var data [][]byte
	for i := 0; i < 100; i++ {
		data = append(data, make([]byte, 1024*10))
	}

	// 重新启用 GC
	debug.SetGCPercent(100)
	fmt.Println("GC 已重新启用")

	// 手动触发 GC
	runtime.GC()

	data = nil
}

// EscapeAnalysis 展示逃逸分析的效果
func EscapeAnalysis() {
	fmt.Println("逃逸分析示例...")

	// 不逃逸的情况
	noEscape := func() {
		x := 42
		fmt.Println("栈上分配:", x)
	}

	// 逃逸的情况
	escape := func() *int {
		x := 42
		return &x // 逃逸到堆
	}

	noEscape()
	ptr := escape()
	fmt.Println("堆上分配:", *ptr)

	// 检查内存分配
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("当前堆分配:", m.HeapAlloc, "bytes")
}
