# Performance 性能优化模块

## 学习目标
掌握Go语言性能优化技巧，包括数据结构优化、字符串处理、并发优化、内存管理等。

## 核心概念

### 1. 预分配
提前分配容量避免动态扩容，提升性能。

### 2. 逃逸分析
Go编译器分析变量是否逃逸到堆，影响性能。

### 3. 对象复用
使用sync.Pool复用对象，减少GC压力。

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| SlicePreallocate | 切片预分配 | 初级 |
| MapPreallocate | map预分配 | 初级 |
| StringConcat | 字符串拼接方式对比 | 初级 |
| DeferPerformance | defer性能对比 | 中级 |
| InterfaceCast | 接口类型断言对比 | 中级 |
| LockPerformance | 锁性能对比 | 中级 |
| ValueVsPointer | 值 vs 指针接收者 | 中级 |
| EscapeAnalysis | 逃逸分析示例 | 中级 |
| SyncPool | sync.Pool对象池示例 | 中级 |
| BenchmarkTemplate | 基准测试模板 | 初级 |
| BestPractices | 最佳实践总结 | 初级 |

## 新增示例详解

### SlicePreallocate 切片预分配
```go
// 不预分配
s1 := make([]int, 0)
for i := 0; i < 1000000; i++ {
    s1 = append(s1, i)  // 多次扩容
}

// 预分配
s2 := make([]int, 0, 1000000)
for i := 0; i < 1000000; i++ {
    s2 = append(s2, i)  // 一次分配
}

// 速度提升: 2-10x
```

### MapPreallocate map预分配
```go
// 不预分配
m1 := make(map[int]int)
for i := 0; i < 100000; i++ {
    m1[i] = i  // 多次rehash
}

// 预分配
m2 := make(map[int]int, 100000)
for i := 0; i < 100000; i++ {
    m2[i] = i  // 一次分配
}
```

### StringConcat 字符串拼接方式对比
```go
// 使用+ (慢)
s1 := ""
for i := 0; i < 10000; i++ {
    s1 += "a"  // 每次创建新字符串
}

// 使用strings.Builder (快)
var builder strings.Builder
for i := 0; i < 10000; i++ {
    builder.WriteString("a")
}
s2 := builder.String()

// 预分配strings.Builder (最快)
var builder2 strings.Builder
builder2.Grow(10000)
for i := 0; i < 10000; i++ {
    builder2.WriteString("a")
}

// 速度提升: 100-1000x
```

### LockPerformance 锁性能对比
```go
// sync.Mutex
var mu sync.Mutex
counter1 := 0
mu.Lock()
counter1++
mu.Unlock()

// atomic (无锁，更快)
var counter2 int64
atomic.AddInt64(&counter2, 1)

// channel (最慢)
ch := make(chan struct{}, 1)
ch <- struct{}{}
counter3++
<-ch
```

### SyncPool 对象池示例
```go
type Buffer struct {
    data []byte
}

pool := &sync.Pool{
    New: func() interface{} {
        return &Buffer{data: make([]byte, 1024)}
    },
}

// 获取对象
buf1 := pool.Get().(*Buffer)
defer pool.Put(buf1)

// 使用对象
buf1.data = buf1.data[:0]
buf1.data = append(buf1.data, "hello"...)

// 归还对象
pool.Put(buf1)

// 再次获取（可能复用）
buf2 := pool.Get().(*Buffer)
```

## 性能优化最佳实践

### 1. 数据结构优化
- 预分配slice和map容量
- 使用值类型 vs 指针类型（小对象用值）
- 选择合适的数据结构

### 2. 字符串处理
- 使用strings.Builder拼接字符串
- 预分配strings.Builder容量
- 避免不必要的字符串转换

### 3. 并发优化
- 用atomic代替mutex（简单场景）
- 使用sync.Pool复用对象
- 避免过多的goroutine创建

### 4. 内存优化
- 减少堆分配（逃逸分析）
- 复用对象和缓冲区
- 避免内存泄漏

### 5. 一般原则
- 先测量，后优化
- 使用pprof分析
- 避免过早优化
- 保持代码可读性

## 逃逸分析

使用 `go build -gcflags '-m'` 查看逃逸分析结果：

```bash
go build -gcflags '-m' ./...
```

### 常见逃逸情况
1. 函数返回局部变量的指针
2. 将变量作为interface{}传递
3. 变量被闭包捕获
4. 变量大小超过栈限制

### 减少堆分配的方法
1. 避免不必要的指针
2. 复用对象（sync.Pool）
3. 预分配容量
4. 使用值而不是指针（小对象）

## 基准测试

Go基准测试规则：
1. 函数名以Benchmark开头
2. 参数为*testing.B
3. 使用b.N控制循环次数

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = 1 + 1
    }
}
```

运行基准测试：
```bash
go test -bench=.
go test -bench=. -benchmem
go test -bench=. -count=5
```

## 性能分析工具

### pprof
```go
import _ "net/http/pprof"

// 启动HTTP服务器
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

```bash
# CPU分析
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 内存分析
go tool pprof http://localhost:6060/debug/pprof/heap
```

### trace
```go
import "runtime/trace"

f, _ := os.Create("trace.out")
trace.Start(f)
defer trace.Stop()
```

```bash
go tool trace trace.out
```

## 运行测试

```bash
go test -v ./performance/...
```

## 参考资料

- [Go性能优化](https://go.dev/blog/pprof)
- [基准测试](https://pkg.go.dev/testing#hdr-Benchmarks)
- [sync.Pool](https://pkg.go.dev/sync#Pool)
