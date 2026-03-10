package performance

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// SlicePreallocate 切片预分配
func SlicePreallocate() {
	fmt.Println("=== 切片预分配对比 ===")

	const count = 1000000

	// 不预分配
	start := time.Now()
	s1 := make([]int, 0)
	for i := 0; i < count; i++ {
		s1 = append(s1, i)
	}
	duration1 := time.Since(start)
	fmt.Printf("不预分配: %v\n", duration1)

	// 预分配
	start = time.Now()
	s2 := make([]int, 0, count)
	for i := 0; i < count; i++ {
		s2 = append(s2, i)
	}
	duration2 := time.Since(start)
	fmt.Printf("预分配: %v\n", duration2)

	speedup := float64(duration1) / float64(duration2)
	fmt.Printf("速度提升: %.2fx\n", speedup)
}

// MapPreallocate map预分配
func MapPreallocate() {
	fmt.Println("=== Map预分配对比 ===")

	const count = 100000

	// 不预分配
	start := time.Now()
	m1 := make(map[int]int)
	for i := 0; i < count; i++ {
		m1[i] = i
	}
	duration1 := time.Since(start)
	fmt.Printf("不预分配: %v\n", duration1)

	// 预分配
	start = time.Now()
	m2 := make(map[int]int, count)
	for i := 0; i < count; i++ {
		m2[i] = i
	}
	duration2 := time.Since(start)
	fmt.Printf("预分配: %v\n", duration2)

	speedup := float64(duration1) / float64(duration2)
	fmt.Printf("速度提升: %.2fx\n", speedup)
}

// StringConcat 字符串拼接方式对比
func StringConcat() {
	fmt.Println("=== 字符串拼接对比 ===")

	const count = 10000

	// 使用+拼接
	start := time.Now()
	s1 := ""
	for i := 0; i < count; i++ {
		s1 += "a"
	}
	duration1 := time.Since(start)
	fmt.Printf("使用+: %v\n", duration1)

	// 使用strings.Builder
	start = time.Now()
	var builder strings.Builder
	for i := 0; i < count; i++ {
		builder.WriteString("a")
	}
	s2 := builder.String()
	duration2 := time.Since(start)
	fmt.Printf("使用strings.Builder: %v\n", duration2)

	// 预分配的strings.Builder
	start = time.Now()
	var builder2 strings.Builder
	builder2.Grow(count)
	for i := 0; i < count; i++ {
		builder2.WriteString("a")
	}
	s3 := builder2.String()
	duration3 := time.Since(start)
	fmt.Printf("预分配strings.Builder: %v\n", duration3)

	speedup1 := float64(duration1) / float64(duration2)
	speedup2 := float64(duration1) / float64(duration3)
	fmt.Printf("strings.Builder速度提升: %.2fx\n", speedup1)
	fmt.Printf("预分配strings.Builder速度提升: %.2fx\n", speedup2)

	_ = s1
	_ = s2
	_ = s3
}

// DeferPerformance defer性能对比
func DeferPerformance() {
	fmt.Println("=== Defer性能对比 ===")

	const count = 1000000

	// 不使用defer
	start := time.Now()
	for i := 0; i < count; i++ {
		func() {
			x := 1
			_ = x
		}()
	}
	duration1 := time.Since(start)
	fmt.Printf("不使用defer: %v\n", duration1)

	// 使用defer
	start = time.Now()
	for i := 0; i < count; i++ {
		func() {
			var x int
			defer func() { _ = x }()
			x = 1
		}()
	}
	duration2 := time.Since(start)
	fmt.Printf("使用defer: %v\n", duration2)

	fmt.Printf("性能差异: %.2f%%\n", float64(duration2-duration1)/float64(duration1)*100)
	fmt.Println("注意: Go 1.14+ defer性能已大幅优化")
}

// InterfaceCast 接口类型断言对比
func InterfaceCast() {
	fmt.Println("=== 接口类型断言对比 ===")

	const count = 1000000

	var iface interface{} = 42

	// 使用类型断言
	start := time.Now()
	sum1 := 0
	for i := 0; i < count; i++ {
		if v, ok := iface.(int); ok {
			sum1 += v
		}
	}
	duration1 := time.Since(start)
	fmt.Printf("类型断言: %v\n", duration1)

	// 使用类型switch
	start = time.Now()
	sum2 := 0
	for i := 0; i < count; i++ {
		switch v := iface.(type) {
		case int:
			sum2 += v
		}
	}
	duration2 := time.Since(start)
	fmt.Printf("类型switch: %v\n", duration2)

	_ = sum1
	_ = sum2
}

// LockPerformance 锁性能对比
func LockPerformance() {
	fmt.Println("=== 锁性能对比 ===")

	const count = 1000000
	const goroutines = 4

	// 使用sync.Mutex
	var mu sync.Mutex
	counter1 := 0
	start := time.Now()
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < count/goroutines; i++ {
				mu.Lock()
				counter1++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	duration1 := time.Since(start)
	fmt.Printf("sync.Mutex: %v, counter=%d\n", duration1, counter1)

	// 使用atomic
	var counter2 int64
	start = time.Now()
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < count/goroutines; i++ {
				atomic.AddInt64(&counter2, 1)
			}
		}()
	}
	wg.Wait()
	duration2 := time.Since(start)
	fmt.Printf("atomic: %v, counter=%d\n", duration2, counter2)

	// 使用channel
	ch := make(chan struct{}, 1)
	counter3 := 0
	start = time.Now()
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < count/goroutines; i++ {
				ch <- struct{}{}
				counter3++
				<-ch
			}
		}()
	}
	wg.Wait()
	duration3 := time.Since(start)
	fmt.Printf("channel: %v, counter=%d\n", duration3, counter3)
}

// ValueVsPointer 值 vs 指针接收者
func ValueVsPointer() {
	fmt.Println("=== 值 vs 指针接收者 ===")

	type Data struct {
		values [100]int
	}

	const count = 100000

	// 值接收者
	var d1 Data
	start := time.Now()
	for i := 0; i < count; i++ {
		d1.values[0] = i
	}
	duration1 := time.Since(start)
	fmt.Printf("值操作: %v\n", duration1)

	// 指针接收者
	var d2 Data
	pd2 := &d2
	start = time.Now()
	for i := 0; i < count; i++ {
		pd2.values[0] = i
	}
	duration2 := time.Since(start)
	fmt.Printf("指针操作: %v\n", duration2)

	fmt.Println("注意: 小结构体差异不大，大结构体指针优势明显")
}

// EscapeAnalysis 逃逸分析示例
func EscapeAnalysis() {
	fmt.Println("=== 逃逸分析 ===")
	fmt.Println("使用 go build -gcflags '-m' 查看逃逸分析结果")
	fmt.Println()
	fmt.Println("常见逃逸情况:")
	fmt.Println("1. 函数返回局部变量的指针")
	fmt.Println("2. 将变量作为interface{}传递")
	fmt.Println("3. 变量被闭包捕获")
	fmt.Println("4. 变量大小超过栈限制")
	fmt.Println()
	fmt.Println("减少堆分配的方法:")
	fmt.Println("1. 避免不必要的指针")
	fmt.Println("2. 复用对象（sync.Pool）")
	fmt.Println("3. 预分配容量")
	fmt.Println("4. 使用值而不是指针（小对象）")
}

// SyncPool 对象池示例
func SyncPool() {
	fmt.Println("=== sync.Pool 对象池 ===")

	type Buffer struct {
		data []byte
	}

	// 创建对象池
	pool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("创建新Buffer")
			return &Buffer{data: make([]byte, 1024)}
		},
	}

	// 获取对象
	buf1 := pool.Get().(*Buffer)
	fmt.Printf("获取buf1: %v\n", buf1)

	// 使用对象
	buf1.data = buf1.data[:0]
	buf1.data = append(buf1.data, "hello"...)

	// 归还对象
	pool.Put(buf1)
	fmt.Println("归还buf1到pool")

	// 再次获取（可能复用）
	buf2 := pool.Get().(*Buffer)
	fmt.Printf("获取buf2: %v (可能是复用的)\n", buf2)

	// 归还
	pool.Put(buf2)

	fmt.Println()
	fmt.Println("sync.Pool适用场景:")
	fmt.Println("1. 频繁创建销毁的临时对象")
	fmt.Println("2. 对象创建成本高")
	fmt.Println("3. 对象可以被复用")
	fmt.Println()
	fmt.Println("注意:")
	fmt.Println("1. Pool中的对象可能被GC回收")
	fmt.Println("2. 不能依赖Pool中的对象状态")
	fmt.Println("3. 适用于高并发场景")
}

// BenchmarkTemplate 基准测试模板
func BenchmarkTemplate() {
	fmt.Println("=== 基准测试 ===")
	fmt.Println()
	fmt.Println("Go基准测试规则:")
	fmt.Println("1. 函数名以Benchmark开头")
	fmt.Println("2. 参数为*testing.B")
	fmt.Println("3. 使用b.N控制循环次数")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println(`func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = 1 + 1
    }
}`)
	fmt.Println()
	fmt.Println("运行基准测试:")
	fmt.Println("  go test -bench=.")
	fmt.Println("  go test -bench=. -benchmem")
	fmt.Println("  go test -bench=. -count=5")
}

// BestPractices 最佳实践总结
func BestPractices() {
	fmt.Println("=== 性能优化最佳实践 ===")
	fmt.Println()
	fmt.Println("1. 数据结构优化:")
	fmt.Println("   - 预分配slice和map容量")
	fmt.Println("   - 使用值类型 vs 指针类型")
	fmt.Println("   - 选择合适的数据结构")
	fmt.Println()
	fmt.Println("2. 字符串处理:")
	fmt.Println("   - 使用strings.Builder拼接字符串")
	fmt.Println("   - 预分配strings.Builder容量")
	fmt.Println("   - 避免不必要的字符串转换")
	fmt.Println()
	fmt.Println("3. 并发优化:")
	fmt.Println("   - 用atomic代替mutex（简单场景）")
	fmt.Println("   - 使用sync.Pool复用对象")
	fmt.Println("   - 避免过多的goroutine创建")
	fmt.Println()
	fmt.Println("4. 内存优化:")
	fmt.Println("   - 减少堆分配（逃逸分析）")
	fmt.Println("   - 复用对象和缓冲区")
	fmt.Println("   - 避免内存泄漏")
	fmt.Println()
	fmt.Println("5. 一般原则:")
	fmt.Println("   - 先测量，后优化")
	fmt.Println("   - 使用pprof分析")
	fmt.Println("   - 避免过早优化")
	fmt.Println("   - 保持代码可读性")
}
