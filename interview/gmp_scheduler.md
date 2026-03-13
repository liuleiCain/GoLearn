# Go GMP调度器详解

本文档详细介绍 Go 语言的 GMP 调度模型，这是理解 Go 并发编程的核心。

---

## 一、调度器概述

### 1.1 为什么需要调度器？

Go 语言的设计目标是支持高并发，如果使用操作系统线程：

```
问题：
1. 线程创建成本高（1MB+栈空间）
2. 线程切换成本高（进入内核态）
3. 线程数量有限（几千个就很多了）
4. 线程间通信复杂（共享内存+锁）

Go的解决方案：
1. 用户态调度（不进入内核）
2. 轻量级协程（2KB栈）
3. M:N调度模型
4. Channel通信
```

### 1.2 调度器演进

```
GM模型（Go 1.0之前）：
- 全局G队列
- 所有M从全局队列获取G
- 问题：全局锁竞争严重

GMP模型（Go 1.1+）：
- 引入P（Processor）
- 每个P有本地G队列
- 减少锁竞争
```

---

## 二、GMP模型详解

### 2.1 G - Goroutine

G 是 Goroutine 的缩写，代表一个协程。

```go
// G的结构（简化版）
type g struct {
    stack       stack   // 栈内存 [lo, hi)
    stackguard0 uintptr // 栈溢出检查
    stackguard1 uintptr // 栈溢出检查（Cgo）

    _panic    *_panic // panic链表
    _defer    *_defer // defer链表
    m         *m      // 当前绑定的M
    sched     gobuf   // 保存的调度信息
    
    atomicstatus uint32 // G的状态
    goid         int64  // G的唯一ID
    
    waitsince    int64  // 等待开始时间
    waitreason   waitReason // 等待原因
    
    preempt      bool   // 抢占标志
    lockedm      *m     // 锁定的M
}

// G的状态
const (
    _Gidle       = iota // 0: 刚分配，未初始化
    _Grunnable          // 1: 在运行队列中，等待运行
    _Grunning           // 2: 正在执行
    _Gsyscall           // 3: 正在执行系统调用
    _Gwaiting           // 4: 被阻塞（IO、channel等）
    _Gdead              // 6: 已经退出或正在退出
    _Gcopystack         // 8: 栈正在复制
    _Gpreempted         // 9: 被抢占
)
```

### 2.2 M - Machine

M 是 Machine 的缩写，代表一个操作系统线程。

```go
// M的结构（简化版）
type m struct {
    g0      *g     // 用于执行调度的特殊G
    curg    *g     // 当前正在运行的G
    p       *p     // 绑定的P
    nextp   *p     // 即将绑定的P
    oldp    *p     // 系统调用前绑定的P
    
    spinning bool   // 是否正在寻找可运行的G
    blocked  bool   // 是否被阻塞
    
    park     note   // 休眠/唤醒
    alllink  *m     // 所有M的链表
    schedlink muintptr // 调度链表
    
    mOS              // 操作系统相关字段
}

// M的最大数量默认是10000
// 可通过 runtime/debug.SetMaxThreads 设置
```

### 2.3 P - Processor

P 是 Processor 的缩写，代表一个逻辑处理器。

```go
// P的结构（简化版）
type p struct {
    id          int32
    status      uint32 // P的状态
    
    m           muintptr // 绑定的M
    
    // 本地G队列（无锁）
    runqhead uint32
    runqtail uint32
    runq     [256]guintptr
    runnext  guintptr // 优先运行的G
    
    // 空闲G列表
    gFree struct {
        gList
        n int32
    }
    
    // GC相关
    gcAssistTime     int64
    gcBgMarkWorker   guintptr
    
    // 定时器
    timers           *timersBucket
}

// P的状态
const (
    _Pidle    = iota // 0: 空闲
    _Prunning         // 1: 运行中
    _Psyscall         // 2: 系统调用中
    _Pgcstop          // 3: GC停止
    _Pdead            // 4: 已停止
)
```

---

## 三、调度流程

### 3.1 调度器初始化

```
程序启动流程：
1. runtime.main 调用 schedinit
2. 创建初始的M和P
3. 创建main goroutine
4. 开始调度

schedinit:
- 设置GOMAXPROCS（P的数量）
- 初始化全局调度器
- 创建初始的P
```

### 3.2 调度循环

```go
// 调度循环（简化版）
func schedule() {
    // 1. 尝试获取G
    gp, inheritTime, tryWakeP := findRunnable()
    
    // 2. 执行G
    execute(gp, inheritTime)
}

// 获取可运行的G
func findRunnable() (gp *g, inheritTime, tryWakeP bool) {
    // 1. 检查是否有需要运行的GC worker
    // 2. 从本地队列获取
    // 3. 从全局队列获取
    // 4. 从网络轮询器获取
    // 5. 从其他P窃取
    // 6. 从其他P窃取GC worker
    // 7. 如果找不到，休眠
}
```

### 3.3 获取G的优先级

```
findRunnable 的查找顺序：

1. 每隔61次调度，检查全局队列
   - 防止全局队列饥饿

2. 本地队列的runnext
   - 最新创建的G优先执行

3. 本地队列
   - 无锁访问，效率最高

4. 全局队列
   - 需要加锁

5. 网络轮询器
   - netpoll

6. 工作窃取
   - 从其他P的本地队列窃取

7. 休眠
   - 没有G可运行时，M会休眠
```

---

## 四、工作窃取（Work Stealing）

### 4.1 窃取算法

```
当P的本地队列为空时：

1. 随机选择一个P
2. 窃取其本地队列的一半G
3. 如果失败，继续选择下一个P
4. 尝试最多4次

窃取过程：
P1的队列: [G1, G2, G3, G4, G5, G6]
P2窃取后:
P1的队列: [G1, G2, G3]
P2的队列: [G4, G5, G6]
```

### 4.2 窃取示意图

```
初始状态：
P0: [G1, G2, G3]  -> M0 执行
P1: []            -> M1 空闲

窃取后：
P0: [G1]          -> M0 执行
P1: [G2, G3]      -> M1 执行
```

---

## 五、系统调用处理

### 5.1 阻塞式系统调用

```go
// 当G执行阻塞式系统调用时：

1. M进入系统调用状态
2. P与M解绑（handoffp）
3. P寻找新的M或创建新M
4. 系统调用返回后：
   - 尝试重新绑定原来的P
   - 如果失败，放入全局队列
   - M进入休眠或销毁
```

### 5.2 非阻塞式系统调用

```go
// 网络IO使用netpoll（非阻塞）

1. G发起网络请求
2. 使用非阻塞IO
3. 如果没有数据，G进入等待状态
4. M继续执行其他G
5. 数据就绪后，G被唤醒
```

### 5.3 系统调用流程图

```
G执行系统调用:
    |
    v
M进入syscall状态
    |
    v
P与M解绑
    |
    +---> P寻找新M或创建新M
    |
    v
系统调用返回
    |
    +---> 尝试重新绑定P
    |         |
    |         +---> 成功：继续执行
    |         |
    |         +---> 失败：G放入全局队列
    |
    v
M休眠或销毁
```

---

## 六、抢占式调度

### 6.1 基于协作的抢占（Go 1.13之前）

```
问题：
- G长时间占用CPU不释放
- 只在函数调用时检查抢占标志

场景：
for {
    // 没有函数调用，不会被抢占
}
```

### 6.2 基于信号的抢占（Go 1.14+）

```go
// 抢占流程：

1. 后台线程监控所有G
2. 如果G运行超过10ms
3. 发送SIGURG信号
4. G的信号处理函数保存上下文
5. 调度器选择其他G运行

// 抢占检查点
- 函数调用入口
- 循环回跳
- 信号处理
```

### 6.3 抢占时机

```
触发抢占的条件：
1. G运行时间超过10ms
2. GC需要停止所有G
3. GC需要扫描栈

安全的抢占点：
1. 函数调用
2. 循环回跳
3. 任何可以安全保存上下文的位置
```

---

## 七、调度器调优

### 7.1 GOMAXPROCS

```go
// 设置P的数量
runtime.GOMAXPROCS(n)

// 默认值：CPU核心数
// 建议：
// - CPU密集型：设置为CPU核心数
// - IO密集型：可以适当增加

// 查看当前值
n := runtime.GOMAXPROCS(0)
```

### 7.2 调度器追踪

```bash
# 查看调度器状态
GODEBUG=schedtrace=1000 ./program

# 输出示例
SCHED 1000ms: gomaxprocs=8 idleprocs=0 threads=10 spinningthreads=1 idlethreads=3 runqueue=0 [0 0 0 0 0 0 0 0]

# 字段含义：
# gomaxprocs: P的数量
# idleprocs: 空闲P的数量
# threads: M的数量
# spinningthreads: 正在寻找G的M数量
# idlethreads: 空闲M的数量
# runqueue: 全局队列中G的数量
# [0 0 ...]: 每个P本地队列中G的数量
```

### 7.3 性能分析

```go
import (
    "runtime"
    "runtime/pprof"
)

// 获取goroutine数量
n := runtime.NumGoroutine()

// 获取goroutine堆栈
buf := make([]byte, 1<<20)
n := runtime.Stack(buf, true)

// pprof分析
pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
```

---

## 八、常见问题

### 8.1 Goroutine泄漏

```go
// 泄漏场景1：无限循环
func leak1() {
    go func() {
        for {
            // 没有退出条件
        }
    }()
}

// 泄漏场景2：channel阻塞
func leak2() {
    ch := make(chan int)
    go func() {
        ch <- 1 // 没有接收者
    }()
}

// 解决方案：使用context
func noLeak(ctx context.Context) {
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
}
```

### 8.2 调度延迟

```go
// 问题：大量G导致调度延迟

// 解决方案1：控制并发数
sem := make(chan struct{}, 100)
for i := 0; i < 10000; i++ {
    sem <- struct{}{}
    go func() {
        defer func() { <-sem }()
        // 工作
    }()
}

// 解决方案2：使用Worker Pool
```

---

## 九、总结

### 9.1 GMP模型优势

```
1. 轻量级
   - G初始栈仅2KB
   - 创建成本低

2. 高效调度
   - 用户态调度
   - 无锁本地队列
   - 工作窃取

3. 可扩展
   - 支持百万级G
   - 动态调整M数量

4. 公平性
   - 抢占式调度
   - 防止G饥饿
```

### 9.2 关键数据

```
- G初始栈：2KB
- G最大栈：1GB（64位）
- P数量：默认CPU核心数
- M最大数量：10000
- 调度切换：~200ns
- 抢占时间：10ms
```

### 9.3 调试命令

```bash
# 查看调度器状态
GODEBUG=schedtrace=1000 ./program

# 查看详细调度信息
GODEBUG=schedtrace=1000,scheddetail=1 ./program

# 禁用抢占
GODEBUG=asyncpreemptoff=1 ./program
```
