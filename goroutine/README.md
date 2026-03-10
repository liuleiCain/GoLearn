# Goroutine 协程模块

## 学习目标
深入理解Go语言协程的基本概念、启动方式、同步机制、并发模式和最佳实践。

## 核心概念

### 1. Goroutine基础
- Goroutine是Go语言轻量级线程的实现
- 使用`go`关键字启动一个goroutine
- Goroutine的调度由Go运行时管理

### 2. sync.WaitGroup
- 用于等待一组goroutine完成
- `Add(delta int)`: 增加计数器
- `Done()`: 减少计数器(通常在defer中调用)
- `Wait()`: 阻塞直到计数器归零

### 3. sync.Once
- 确保某个操作只执行一次
- 常用于单例模式、初始化等场景

### 4. runtime包
- `runtime.NumGoroutine()`: 获取当前goroutine数量
- `runtime.Gosched()`: 让出CPU时间片
- `runtime.GOMAXPROCS(n)`: 设置使用的CPU核心数

### 5. Context包
- `context.WithCancel`: 可取消的context
- `context.WithTimeout`: 带超时的context
- `context.WithValue`: 携带值的context

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| BasicGoroutine | 基础goroutine启动方式 | 初级 |
| WaitGroupDemo | WaitGroup同步等待 | 初级 |
| GoroutineLeakDemo | Goroutine泄漏问题演示 | 中级 |
| GoroutineNumberDemo | 监控goroutine数量 | 中级 |
| GoschedDemo | runtime.Gosched使用 | 中级 |
| GOMAXPROCSDemo | CPU核心数设置 | 中级 |
| OnceDemo | sync.Once单次执行 | 初级 |
| MutexDemo | 互斥锁保护共享资源 | 初级 |
| RWMutexDemo | 读写锁使用 | 中级 |
| AtomicDemo | 原子操作 | 中级 |
| ContextCancel | Context取消goroutine | 中级 |
| ContextTimeout | Context超时控制 | 中级 |
| ContextValue | Context传递值 | 中级 |
| WorkerPool | 工作池模式 | 中级 |
| Pipeline | 管道模式 | 中级 |
| FanOutFanIn | 扇出扇入模式 | 高级 |
| TimeoutPattern | 超时模式 | 中级 |
| GracefulShutdown | 优雅关闭 | 高级 |
| GoroutineLocal | Goroutine本地存储模拟 | 中级 |
| SelectDemo | Select多路复用 | 初级 |
| NonBlockingSelect | 非阻塞Select | 中级 |
| TimerDemo | 定时器使用 | 初级 |
| TickerDemo | 周期执行 | 初级 |

## 新增示例详解

### MutexDemo - 互斥锁
保护共享资源，防止竞态条件。
```go
var mu sync.Mutex
mu.Lock()
counter++
mu.Unlock()
```

### RWMutexDemo - 读写锁
读多写少场景下提高性能。
```go
var rwmu sync.RWMutex
rwmu.Lock()    // 写锁
rwmu.RLock()   // 读锁
```

### AtomicDemo - 原子操作
无锁并发操作。
```go
atomic.AddInt64(&counter, 1)
atomic.CompareAndSwapInt64(&value, 100, 200)
atomic.LoadInt64(&value)
atomic.StoreInt64(&value, 300)
```

### ContextCancel - 取消goroutine
```go
ctx, cancel := context.WithCancel(context.Background())
go func() {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            // 工作
        }
    }
}()
cancel() // 取消
```

### ContextTimeout - 超时控制
```go
ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
defer cancel()
select {
case <-result:
    // 完成
case <-ctx.Done():
    // 超时
}
```

### WorkerPool - 工作池模式
```go
jobs := make(chan int, 100)
results := make(chan int, 100)

// 启动固定数量的worker
for w := 1; w <= 3; w++ {
    go worker(w, jobs, results)
}

// 发送任务
go func() {
    for j := 1; j <= 10; j++ {
        jobs <- j
    }
    close(jobs)
}()
```

### Pipeline - 管道模式
```go
// 阶段1: 生成数据
nums := generator(done, 1, 2, 3, 4, 5)
// 阶段2: 处理数据
squared := square(done, nums)
// 阶段3: 消费数据
consumer(done, squared)
```

### FanOutFanIn - 扇出扇入模式
```go
// 扇出: 多个worker读取同一输入
input := producer(1, 2, 3, 4, 5)
c1 := worker("Worker1", input)
c2 := worker("Worker2", input)
c3 := worker("Worker3", input)

// 扇入: 合并多个输出
for result := range merger(c1, c2, c3) {
    fmt.Println(result)
}
```

### GracefulShutdown - 优雅关闭
```go
ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
defer cancel()

// 启动服务
for i := 1; i <= 3; i++ {
    go service(ctx, i)
}

// 等待所有服务关闭
wg.Wait()
```

### SelectDemo - 多路复用
```go
select {
case msg1 := <-ch1:
    fmt.Println("ch1:", msg1)
case msg2 := <-ch2:
    fmt.Println("ch2:", msg2)
case <-time.After(100 * time.Millisecond):
    fmt.Println("超时")
}
```

### NonBlockingSelect - 非阻塞操作
```go
select {
case ch <- 1:
    fmt.Println("发送成功")
default:
    fmt.Println("发送失败")
}
```

### TimerDemo - 定时器
```go
timer := time.NewTimer(200 * time.Millisecond)
<-timer.C
fmt.Println("定时器触发")

timer.Stop()   // 停止
timer.Reset(100 * time.Millisecond)  // 重置
```

### TickerDemo - 周期执行
```go
ticker := time.NewTicker(100 * time.Millisecond)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        fmt.Println("Ticker触发")
    case <-done:
        return
    }
}
```

## 常见陷阱

### 1. 循环变量捕获
```go
// 错误示例 - 所有goroutine可能打印相同的值
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i) // 捕获的是变量i的引用
    }()
}

// 正确示例 - 将i作为参数传递
for i := 0; i < 3; i++ {
    go func(n int) {
        fmt.Println(n)
    }(i)
}
```

### 2. Goroutine泄漏
- 如果goroutine永远无法结束，会导致内存泄漏
- 确保每个goroutine都有退出的条件

### 3. 竞态条件
```go
// 错误: 多个goroutine同时写入
var counter int
for i := 0; i < 1000; i++ {
    go func() { counter++ }()
}

// 正确: 使用锁或原子操作
var mu sync.Mutex
mu.Lock()
counter++
mu.Unlock()
```

### 4. 死锁
```go
// 错误: channel无缓冲且无接收者
ch := make(chan int)
ch <- 1  // 死锁

// 正确: 使用缓冲或确保有接收者
ch := make(chan int, 1)
ch <- 1
```

## 并发模式总结

### 1. Worker Pool（工作池）
固定数量的worker处理任务，适合CPU密集型任务。

### 2. Pipeline（管道）
数据流经多个处理阶段，适合数据转换场景。

### 3. Fan-out/Fan-in（扇出扇入）
多个worker并行处理，结果合并，适合IO密集型任务。

### 4. Timeout（超时）
防止goroutine长时间阻塞。

### 5. Graceful Shutdown（优雅关闭）
确保资源正确释放。

## 最佳实践

1. 使用`go vet`检查循环变量捕获问题
2. 使用`sync.WaitGroup`等待goroutine完成
3. 监控goroutine数量，防止泄漏
4. 使用`context`包实现goroutine的取消机制
5. 使用锁保护共享资源
6. 优先使用channel通信，而不是共享内存
7. 为长时间运行的goroutine设置退出条件
8. 使用`runtime.NumGoroutine()`监控goroutine数量

## 运行测试

```bash
go test -v ./goroutine/...
```

## 参考资料

- [Go语言并发编程](https://go.dev/blog/pipelines)
- [Go Concurrency Patterns](https://go.dev/talks/2012/concurrency.slide)
- [Context包文档](https://pkg.go.dev/context)
