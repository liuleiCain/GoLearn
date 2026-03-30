package interview

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

// ============================================================
// 运行时与底层设计示例
// ============================================================

// AllocatorStatsDemo 观测分配统计信息，理解分配器的行为特征。
func AllocatorStatsDemo() {
	fmt.Println("\n=== 内存分配器观测 ===")

	var before, after, afterGC runtime.MemStats
	runtime.ReadMemStats(&before)

	const smallN = 10000
	small := make([][]byte, smallN)
	for i := 0; i < smallN; i++ {
		size := 32 + (i % 128)
		small[i] = make([]byte, size)
	}

	big := make([][]byte, 64)
	for i := 0; i < len(big); i++ {
		big[i] = make([]byte, 256*1024)
	}

	runtime.ReadMemStats(&after)
	fmt.Printf("Mallocs: %d -> %d (+%d)\n", before.Mallocs, after.Mallocs, after.Mallocs-before.Mallocs)
	fmt.Printf("HeapObjects: %d -> %d (+%d)\n", before.HeapObjects, after.HeapObjects, after.HeapObjects-before.HeapObjects)
	fmt.Printf("HeapAlloc: %d -> %d bytes\n", before.HeapAlloc, after.HeapAlloc)

	runtime.KeepAlive(small)
	runtime.KeepAlive(big)

	small = nil
	big = nil
	runtime.GC()
	runtime.ReadMemStats(&afterGC)
	fmt.Printf("GC后 HeapObjects: %d\n", afterGC.HeapObjects)
}

// StackGrowthDemo 观测栈地址变化，理解栈拷贝的可能性。
func StackGrowthDemo() {
	fmt.Println("\n=== 栈增长与拷贝观测 ===")

	const depth = 32
	addrs := make([]uintptr, 0, depth+1)

	var walk func(int)
	walk = func(n int) {
		var local [256]byte
		local[0] = byte(n) // 防止优化
		addrs = append(addrs, uintptr(unsafe.Pointer(&local)))
		if n > 0 {
			walk(n - 1)
		}
	}
	walk(depth)

	if len(addrs) > 0 {
		fmt.Printf("采样深度: %d\n", len(addrs))
		fmt.Printf("栈地址范围: 0x%x -> 0x%x\n", addrs[0], addrs[len(addrs)-1])
	}

	monotonic := true
	for i := 1; i < len(addrs); i++ {
		if addrs[i] >= addrs[i-1] {
			monotonic = false
			break
		}
	}
	fmt.Printf("地址是否单调递减: %v (非单调可能发生栈拷贝)\n", monotonic)

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Printf("StackInuse=%d bytes, StackSys=%d bytes\n", ms.StackInuse, ms.StackSys)
}

// SyncPoolDemo 展示 sync.Pool 的基本行为：对象可能被 GC 清理。
func SyncPoolDemo() {
	fmt.Println("\n=== sync.Pool 行为 ===")

	var pool sync.Pool
	pool.New = func() any {
		return make([]byte, 1024)
	}

	buf := pool.Get().([]byte)
	pool.Put(buf)

	runtime.GC()

	buf2 := pool.Get().([]byte)
	same := len(buf) > 0 && len(buf2) > 0 && unsafe.Pointer(&buf[0]) == unsafe.Pointer(&buf2[0])
	fmt.Printf("buf2 len=%d cap=%d, 是否复用=%v\n", len(buf2), cap(buf2), same)
}
