# Go GC垃圾回收详解

本文档详细介绍 Go 语言的垃圾回收机制，包括三色标记算法、写屏障、GC调优等内容。

---

## 一、GC概述

### 1.1 为什么需要GC？

```
手动内存管理的问题：
1. 忘记释放 -> 内存泄漏
2. 重复释放 -> 程序崩溃
3. 悬空指针 -> 数据损坏
4. 开发效率低

GC的优势：
1. 自动管理内存
2. 避免常见内存错误
3. 提高开发效率
4. 安全性更高
```

### 1.2 Go GC的特点

```
1. 并发标记-清除
   - 标记与用户代码并发执行
   - 低延迟

2. 三色标记算法
   - 增量式标记
   - 减少STW时间

3. 写屏障
   - 保证并发标记正确性
   - 混合写屏障（Go 1.8+）

4. 非分代
   - 没有分代GC
   - 但通过逃逸分析优化

5. 非紧凑
   - 不进行内存整理
   - 依赖内存分配器
```

---

## 二、三色标记算法

### 2.1 三色抽象

```
白色（White）：
- 未被标记的对象
- 可能是垃圾对象

灰色（Gray）：
- 已被标记，但其引用的对象未被标记
- 需要扫描其引用

黑色（Black）：
- 已被标记，且其引用的对象也已被标记
- 不需要再扫描

不变式：
1. 强三色不变式：黑色对象不会直接引用白色对象
2. 弱三色不变式：黑色对象引用的白色对象，必须被灰色对象保护
```

### 2.2 标记过程

```
初始状态：
所有对象都是白色

标记开始：
1. 将根对象（全局变量、栈上变量）标记为灰色

标记阶段：
while 存在灰色对象 {
    1. 选择一个灰色对象
    2. 将其引用的白色对象标记为灰色
    3. 将该灰色对象标记为黑色
}

标记结束：
- 黑色对象：存活
- 白色对象：垃圾

清除阶段：
- 回收所有白色对象
```

### 2.3 标记流程图

```
初始：
[白] -> [白] -> [白] -> [白]
  |
  v
[灰] (根对象)

标记中：
[黑] -> [灰] -> [白] -> [白]
         |
         v
       [灰]

标记完成：
[黑] -> [黑] -> [黑] -> [白]
                        |
                        v
                     回收
```

---

## 三、写屏障

### 3.1 为什么需要写屏障？

```
问题：并发标记时，用户代码可能修改对象引用

场景：
1. 黑色对象A引用白色对象B
2. 灰色对象C原本引用B
3. 用户代码：C.B = nil（断开引用）
4. 结果：B变成垃圾，但A已经标记为黑色

问题：
- A是黑色，不会再扫描
- B是白色，会被回收
- 但A引用了B，导致悬空指针

解决：写屏障
```

### 3.2 Dijkstra写屏障

```go
// 伪代码
func writePointer(slot *unsafe.Pointer, ptr unsafe.Pointer) {
    // 将新引用的对象标记为灰色
    shade(ptr)
    // 更新引用
    *slot = ptr
}

// 问题：需要所有写操作都经过写屏障
// 性能开销较大
```

### 3.3 Yuasa写屏障

```go
// 伪代码
func writePointer(slot *unsafe.Pointer, ptr unsafe.Pointer) {
    // 将旧引用的对象标记为灰色
    shade(*slot)
    // 更新引用
    *slot = ptr
}

// 问题：可能保留一些垃圾对象
```

### 3.4 混合写屏障（Go 1.8+）

```go
// 伪代码
func writePointer(slot *unsafe.Pointer, ptr unsafe.Pointer) {
    shade(*slot)  // 标记旧对象
    shade(ptr)    // 标记新对象
    *slot = ptr
}

// 优化：栈不需要写屏障
// 1. GC开始时，扫描栈，将所有引用标记为黑色
// 2. 栈上的写操作不需要写屏障
// 3. 只有堆上的写操作需要写屏障
```

---

## 四、GC流程详解

### 4.1 GC触发条件

```
1. 内存分配触发
   - 堆内存达到阈值
   - 阈值 = 上次GC后的堆大小 * (1 + GOGC/100)
   - 默认GOGC=100，即堆大小翻倍时触发

2. 手动触发
   - runtime.GC()

3. 定时触发
   - 2分钟未触发GC
```

### 4.2 GC阶段

```
阶段1：标记准备（STW）
- 停止所有P
- 开启写屏障
- 扫描根对象（栈、全局变量）
- 创建标记任务
- 恢复P运行

阶段2：并发标记
- 标记工作协程并发执行
- 用户代码并发执行
- 写屏障记录引用变更
- 辅助标记（防止用户分配过快）

阶段3：标记终止（STW）
- 停止所有P
- 关闭写屏障
- 清理工作
- 统计GC信息
- 恢复P运行

阶段4：并发清除
- 异步回收白色对象
- 用户代码并发执行
```

### 4.3 GC时间线

```
时间线：
|----STW----|----并发标记----|----STW----|----并发清除----|
   ~100μs       ~ms级          ~100μs       ~ms级

STW时间：
- Go 1.5: ~1ms
- Go 1.8+: ~100μs
- 现代Go: ~10-100μs
```

---

## 五、GC调优

### 5.1 GOGC参数

```go
// 设置GOGC
// GOGC=100: 堆大小翻倍时触发GC（默认）
// GOGC=200: 堆大小增长200%时触发GC
// GOGC=off: 禁用自动GC

// 通过环境变量
// GOGC=200 ./program

// 通过代码
import "runtime/debug"
debug.SetGCPercent(200)
```

### 5.2 内存限制（Go 1.19+）

```go
import "runtime/debug"

// 设置软内存限制
debug.SetMemoryLimit(1 * 1024 * 1024 * 1024) // 1GB

// 当内存接近限制时，GC会更积极地运行
```

### 5.3 GC监控

```go
import (
    "runtime"
    "runtime/debug"
)

func monitorGC() {
    var stats debug.GCStats
    debug.ReadGCStats(&stats)
    
    fmt.Printf("GC次数: %d\n", stats.NumGC)
    fmt.Printf("暂停总时间: %v\n", stats.PauseTotal)
    fmt.Printf("上次暂停: %v\n", stats.Pause[0])
}

func monitorMem() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("堆分配: %d MB\n", m.HeapAlloc/1024/1024)
    fmt.Printf("系统内存: %d MB\n", m.Sys/1024/1024)
    fmt.Printf("GC次数: %d\n", m.NumGC)
}
```

### 5.4 GC追踪

```bash
# 启用GC追踪
GODEBUG=gctrace=1 ./program

# 输出示例
gc 1 @0.003s 0%: 0.016+0.12+0.003 ms clock, 0.13+0.099/0.12/0.018+0.026 ms cpu, 4->4->0 MB, 5 MB goal, 8 P

# 字段含义：
# gc 1: 第1次GC
# @0.003s: 程序启动后0.003秒
# 0%: GC占用的CPU时间百分比
# 0.016+0.12+0.003 ms: STW标记准备 + 并发标记 + STW标记终止
# 4->4->0 MB: GC前堆大小 -> GC后堆大小 -> 存活对象大小
# 5 MB goal: 下次GC的目标堆大小
# 8 P: P的数量
```

---

## 六、内存泄漏与排查

### 6.1 常见内存泄漏

```go
// 1. 全局变量持有引用
var globalData []byte

func leak1() {
    globalData = make([]byte, 1<<30) // 1GB
}

// 2. 未关闭的资源
func leak2() {
    f, _ := os.Open("file.txt")
    // 忘记 f.Close()
}

// 3. 无限增长的map
var cache = make(map[string][]byte)

func leak3(key string, data []byte) {
    cache[key] = data // 只增不减
}

// 4. goroutine泄漏
func leak4() {
    ch := make(chan int)
    go func() {
        <-ch // 永远阻塞
    }()
}

// 5. 闭包捕获
func leak5() {
    var data []byte
    runtime.SetFinalizer(&data, func(d *[]byte) {
        // 闭包捕获data，导致无法回收
    })
}
```

### 6.2 排查工具

```bash
# 1. pprof内存分析
go tool pprof http://localhost:6060/debug/pprof/heap

# 2. 查看goroutine
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 3. 查看allocs
go tool pprof http://localhost:6060/debug/pprof/allocs

# 4. 生成堆转储
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

### 6.3 pprof使用示例

```go
import (
    "net/http"
    _ "net/http/pprof"
)

func main() {
    // 启动pprof服务
    go func() {
        http.ListenAndServe(":6060", nil)
    }()
    
    // 主程序
    // ...
}
```

```
# 在pprof交互模式
(pprof) top10    # 显示内存占用前10
(pprof) list funcName  # 查看函数详情
(pprof) web      # 生成可视化图表
(pprof) png > out.png  # 保存图表
```

---

## 七、GC最佳实践

### 7.1 减少堆分配

```go
// 1. 使用值类型
type Point struct {
    X, Y int
}

// 好：值类型，分配在栈上
func processPoint(p Point) { }

// 避免：指针类型，可能逃逸到堆
func processPoint(p *Point) { }

// 2. 预分配切片
// 不好
var s []int
for i := 0; i < 1000; i++ {
    s = append(s, i)
}

// 好
s := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    s = append(s, i)
}

// 3. 使用sync.Pool
var bufPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func process() {
    buf := bufPool.Get().([]byte)
    defer bufPool.Put(buf)
    // 使用buf
}
```

### 7.2 避免逃逸

```go
// 查看逃逸分析
// go build -gcflags="-m" main.go

// 逃逸场景1：返回局部变量指针
func escape1() *int {
    x := 1
    return &x // 逃逸
}

// 逃逸场景2：接口转换
func escape2() interface{} {
    x := 1
    return x // 逃逸
}

// 逃逸场景3：闭包捕获
func escape3() func() int {
    x := 1
    return func() int {
        return x // 逃逸
    }
}

// 逃逸场景4：大对象
func escape4() {
    x := make([]byte, 64*1024) // 大于32KB，逃逸
}
```

### 7.3 合理设置GOGC

```
场景1：内存充足，追求低延迟
GOGC=50 或更低
- GC更频繁
- STW时间更短
- 堆内存更小

场景2：内存紧张，追求低内存
GOGC=200 或更高
- GC更少
- 堆内存更大
- STW时间可能更长

场景3：批量处理
GOGC=off
- 禁用自动GC
- 手动在合适时机触发
```

---

## 八、GC相关API

### 8.1 runtime包

```go
import "runtime"

// 手动触发GC
runtime.GC()

// 获取内存统计
var m runtime.MemStats
runtime.ReadMemStats(&m)

// 设置CPU核心数
runtime.GOMAXPROCS(n)

// 让出CPU
runtime.Gosched()

// 获取goroutine数量
n := runtime.NumGoroutine()
```

### 8.2 runtime/debug包

```go
import "runtime/debug"

// 设置GOGC
old := debug.SetGCPercent(50)

// 设置内存限制
debug.SetMemoryLimit(1 << 30) // 1GB

// 获取GC统计
var stats debug.GCStats
debug.ReadGCStats(&stats)

// 设置最大线程数
debug.SetMaxThreads(1000)

// 设置栈大小限制
debug.SetMaxStack(1 << 20) // 1MB
```

---

## 九、总结

### 9.1 GC关键指标

```
- STW时间：~10-100μs
- GC触发：堆大小翻倍（GOGC=100）
- 标记算法：三色并发标记
- 写屏障：混合写屏障
- 内存整理：无
```

### 9.2 调优建议

```
1. 减少堆分配
   - 使用值类型
   - 预分配
   - sync.Pool

2. 避免内存泄漏
   - 及时释放引用
   - 关闭资源
   - 控制goroutine

3. 合理设置参数
   - GOGC
   - SetMemoryLimit
   - GOMAXPROCS

4. 监控GC
   - gctrace
   - pprof
   - runtime.ReadMemStats
```

### 9.3 调试命令

```bash
# GC追踪
GODEBUG=gctrace=1 ./program

# 详细GC追踪
GODEBUG=gctrace=1,gcpacertrace=1 ./program

# 内存分配追踪
GODEBUG=allocfreetrace=1 ./program

# 禁用GC
GODEBUG=gcstop=1 ./program
```
