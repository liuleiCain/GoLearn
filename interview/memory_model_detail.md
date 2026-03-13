# Go 内存模型详解

本文档详细介绍 Go 语言的内存模型，包括 Happens-Before 关系、数据竞争、同步原语等内容。

---

## 一、内存模型概述

### 1.1 什么是内存模型？

```
内存模型定义：
- 在一个goroutine中对变量的写入
- 如何能被另一个goroutine中对该变量的读取观察到

核心问题：
- 多个goroutine并发执行
- 编译器和CPU可能重排序指令
- 如何保证正确的可见性？
```

### 1.2 为什么需要内存模型？

```go
// 问题代码
var x, y int

func goroutine1() {
    x = 1  // A
    y = 2  // B
}

func goroutine2() {
    if y == 2 {  // C
        print(x) // D: x可能是0或1
    }
}

// 可能的执行顺序：
// 1. A -> B -> C -> D (x=1)
// 2. B -> C -> D -> A (x=0)
// 3. 编译器/CPU重排序后的其他顺序
```

---

## 二、Happens-Before关系

### 2.1 定义

```
如果操作A happens-before 操作B：
- A的结果对B可见
- A在B之前执行

传递性：
- A happens-before B
- B happens-before C
- 则 A happens-before C
```

### 2.2 Happens-Before规则

```
规则1：goroutine创建
- go语句 happens-before 新goroutine开始执行

规则2：goroutine销毁
- goroutine的退出不保证 happens-before 任何操作
- 需要使用同步机制确保可见性

规则3：channel发送
- 向channel发送 happens-before 从该channel接收完成

规则4：channel关闭
- 关闭channel happens-before 从该channel接收到零值

规则5：无缓冲channel
- 从无缓冲channel接收 happens-before 发送完成

规则6：有缓冲channel
- 第k次接收 happens-before 第(k+C)次发送完成
- C是channel容量

规则7：锁
- 对于sync.Mutex或sync.RWMutex
- 第n次Unlock happens-before 第n+1次Lock

规则8：Once
- once.Do(f)中的f() happens-before once.Do(f)返回

规则9：WaitGroup
- Add happens-before Wait返回
- Done happens-before Wait返回

规则10：atomic
- 原子操作遵循顺序一致性
```

### 2.3 Happens-Before图示

```
goroutine创建：
main goroutine          new goroutine
    |                        |
    v                        v
go func()  --------->  func()开始执行


channel通信：
sender goroutine         receiver goroutine
    |                        |
    v                        v
ch <- x    --------->  x := <-ch
(发送开始)              (接收完成)


无缓冲channel：
receiver goroutine       sender goroutine
    |                        |
    v                        v
<-ch       --------->  ch <-
(接收开始)              (发送完成)


Mutex锁：
goroutine1              goroutine2
    |                        |
    v                        v
Unlock()   --------->  Lock()
(释放锁)                (获取锁)
```

---

## 三、数据竞争

### 3.1 什么是数据竞争？

```
数据竞争定义：
两个goroutine同时访问同一变量，且满足：
1. 至少有一个是写操作
2. 没有同步机制

后果：
- 未定义行为
- 数据损坏
- 程序崩溃
```

### 3.2 数据竞争示例

```go
// 示例1：简单的数据竞争
var counter int

func race1() {
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++ // 数据竞争！
        }()
    }
    wg.Wait()
    fmt.Println(counter) // 结果不确定
}

// 示例2：看似安全的代码
var done bool
var data string

func race2() {
    go func() {
        data = "hello"
        done = true
    }()
    
    for !done {
        // 可能永远看不到done变为true
        // 因为编译器可能优化掉循环
    }
    fmt.Println(data)
}
```

### 3.3 检测数据竞争

```bash
# 使用race检测器
go run -race main.go
go test -race ./...

# 输出示例
==================
WARNING: DATA RACE
Write at 0x000001234567 by goroutine 7:
  main.main.func1()
      /path/main.go:10 +0x42

Previous read at 0x000001234567 by goroutine 6:
  main.main.func2()
      /path/main.go:15 +0x3a
==================
```

---

## 四、同步原语详解

### 4.1 Mutex（互斥锁）

```go
type Mutex struct {
    state int32  // 锁状态
    sema  uint32 // 信号量
}

// 使用示例
var (
    mu      sync.Mutex
    counter int
)

func safeIncrement() {
    mu.Lock()
    defer mu.Unlock()
    counter++
}

// 内部实现（简化）
func (m *Mutex) Lock() {
    // 快速路径：CAS尝试获取锁
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        return
    }
    // 慢速路径：自旋或阻塞
    m.lockSlow()
}
```

### 4.2 RWMutex（读写锁）

```go
type RWMutex struct {
    w           Mutex  // 写锁
    writerSem   uint32 // 写等待信号量
    readerSem   uint32 // 读等待信号量
    readerCount int32  // 读计数
    readerWait  int32  // 写等待的读者数
}

// 使用示例
var (
    rwMu sync.RWMutex
    data map[string]string
)

func read(key string) string {
    rwMu.RLock()
    defer rwMu.RUnlock()
    return data[key]
}

func write(key, value string) {
    rwMu.Lock()
    defer rwMu.Unlock()
    data[key] = value
}
```

### 4.3 WaitGroup

```go
type WaitGroup struct {
    noCopy noCopy
    state1 [3]uint32 // state和sema
}

// 使用示例
func waitGroupExample() {
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            // 工作
        }(i)
    }
    
    wg.Wait()
}

// Happens-Before保证：
// Add() happens-before Wait()返回
// Done() happens-before Wait()返回
```

### 4.4 Once

```go
type Once struct {
    done uint32
    m    Mutex
}

// 使用示例
var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}

// Happens-Before保证：
// f() happens-before Do()返回
```

### 4.5 Cond（条件变量）

```go
type Cond struct {
    noCopy noCopy
    L Locker
    notify list.List
}

// 使用示例
var (
    mu    sync.Mutex
    cond  = sync.NewCond(&mu)
    ready bool
)

func wait() {
    mu.Lock()
    for !ready {
        cond.Wait() // 释放锁并等待
    }
    mu.Unlock()
}

func signal() {
    mu.Lock()
    ready = true
    cond.Signal() // 唤醒一个等待者
    mu.Unlock()
}
```

### 4.6 atomic（原子操作）

```go
// 原子操作类型
var (
    counter int64
    value   atomic.Value
)

// Add: 原子加
atomic.AddInt64(&counter, 1)

// CompareAndSwap: 比较并交换
if atomic.CompareAndSwapInt64(&counter, 0, 1) {
    // 成功
}

// Swap: 交换
old := atomic.SwapInt64(&counter, 10)

// Load: 原子读取
v := atomic.LoadInt64(&counter)

// Store: 原子存储
atomic.StoreInt64(&counter, 100)

// Value: 存储任意类型
value.Store("hello")
s := value.Load().(string)
```

---

## 五、Channel同步详解

### 5.1 无缓冲Channel

```go
// 无缓冲channel：同步通信
ch := make(chan int)

// 发送者
go func() {
    ch <- 1  // 阻塞直到接收
}()

// 接收者
<-ch  // 阻塞直到发送

// Happens-Before：
// 接收开始 happens-before 发送完成
// 发送开始 happens-before 接收完成
```

### 5.2 有缓冲Channel

```go
// 有缓冲channel：异步通信
ch := make(chan int, 3)

// 发送者：缓冲区满时阻塞
ch <- 1
ch <- 2
ch <- 3
// ch <- 4  // 阻塞

// 接收者：缓冲区空时阻塞
<-ch  // 1
<-ch  // 2
<-ch  // 3
// <-ch  // 阻塞

// Happens-Before：
// 第k次接收 happens-before 第(k+C)次发送完成
// C是容量
```

### 5.3 Channel关闭

```go
ch := make(chan int, 3)

go func() {
    ch <- 1
    ch <- 2
    close(ch) // 关闭
}()

// 从已关闭的channel读取
v, ok := <-ch  // 1, true
v, ok = <-ch   // 2, true
v, ok = <-ch   // 0, false（零值，已关闭）

// range遍历
for v := range ch {
    fmt.Println(v) // 1, 2
}
```

---

## 六、内存顺序与重排序

### 6.1 编译器重排序

```go
// 原始代码
x = 1
y = 2

// 编译器可能重排序为
y = 2
x = 1

// 原因：优化性能
```

### 6.2 CPU重排序

```
CPU可能重排序内存操作：
- Store-Load重排序
- Store-Store重排序
- Load-Load重排序
- Load-Store重排序

x86架构：
- 只允许Store-Load重排序

ARM架构：
- 允许所有类型的重排序
```

### 6.3 内存屏障

```
内存屏障类型：
1. LoadLoad：禁止Load-Load重排序
2. StoreStore：禁止Store-Store重排序
3. LoadStore：禁止Load-Store重排序
4. StoreLoad：禁止Store-Load重排序

Go中的内存屏障：
- atomic操作隐含内存屏障
- channel操作隐含内存屏障
- 锁操作隐含内存屏障
```

---

## 七、常见陷阱与最佳实践

### 7.1 常见陷阱

```go
// 陷阱1：双重检查锁定错误
var instance *Singleton
var mu sync.Mutex

func GetInstance() *Singleton {
    if instance == nil {  // 第一次检查：数据竞争！
        mu.Lock()
        defer mu.Unlock()
        if instance == nil {
            instance = &Singleton{}
        }
    }
    return instance
}

// 正确做法：使用sync.Once

// 陷阱2：错误的时间戳检查
var start time.Time

func startTimer() {
    start = time.Now()  // 可能在goroutine中
}

func elapsed() time.Duration {
    return time.Since(start)  // 可能看到零值
}

// 正确做法：使用channel或锁同步

// 陷阱3：map并发访问
var m = make(map[string]int)

func write() {
    m["key"] = 1  // 并发写入
}

func read() {
    _ = m["key"]  // 并发读取
}

// 正确做法：使用sync.Mutex或sync.Map
```

### 7.2 最佳实践

```go
// 1. 使用channel而非共享内存
// 不要
var data int
var mu sync.Mutex

func write() {
    mu.Lock()
    data = 1
    mu.Unlock()
}

// 推荐
ch := make(chan int)
go func() { ch <- 1 }()
data := <-ch

// 2. 明确所有权
// 数据的"所有者"负责访问
type Queue struct {
    ch chan int
}

func (q *Queue) Put(v int) {
    q.ch <- v
}

func (q *Queue) Get() int {
    return <-q.ch
}

// 3. 使用context取消
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            // 工作
        }
    }
}

// 4. 使用atomic进行简单计数
var counter int64

func increment() {
    atomic.AddInt64(&counter, 1)
}
```

---

## 八、调试与检测

### 8.1 Race检测器

```bash
# 运行时检测
go run -race main.go
go test -race ./...

# 原理：
# 1. 插入检测代码
# 2. 记录内存访问
# 3. 检测并发冲突
```

### 8.2 ThreadSanitizer

```bash
# Go的race检测器基于ThreadSanitizer
# 检测：
# - 数据竞争
# - 死锁（部分）
# - 原子性违规

# 限制：
# - 只检测实际执行的代码路径
# - 有性能开销（5-10x）
# - 内存开销（5-10x）
```

### 8.3 静态分析

```bash
# 使用go vet
go vet ./...

# 使用staticcheck
staticcheck ./...

# 使用golangci-lint
golangci-lint run
```

---

## 九、总结

### 9.1 核心概念

```
1. Happens-Before关系
   - 定义操作之间的可见性顺序
   - 是正确同步的基础

2. 数据竞争
   - 未定义行为
   - 必须避免

3. 同步原语
   - channel：推荐
   - mutex：简单场景
   - atomic：性能敏感

4. 内存屏障
   - 保证内存顺序
   - 隐含在同步操作中
```

### 9.2 设计原则

```
1. 不要通过共享内存通信，通过通信共享内存

2. 明确数据的所有权

3. 使用race检测器

4. 保持简单

5. 避免过度优化
```

### 9.3 调试命令

```bash
# race检测
go run -race main.go

# 静态分析
go vet ./...

# 查看汇编
go build -gcflags="-S" main.go

# 查看逃逸分析
go build -gcflags="-m" main.go
```
