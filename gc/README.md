# Go 垃圾回收（GC）模块

本模块详细介绍 Go 语言的垃圾回收机制及其使用方法。

## 目录

- [GC 简介](#gc-简介)
- [示例函数](#示例函数)
- [使用说明](#使用说明)

## GC 简介

Go 语言使用三色标记-清除（Tri-color Mark-and-Sweep）算法进行垃圾回收，具有以下特点：

- **并发标记**：GC 标记阶段与用户程序并发执行
- **写屏障**：使用写屏障技术确保并发标记的正确性
- **分代回收**：虽然 Go 没有显式的分代，但通过逃逸分析优化了年轻代对象的回收
- **STW（Stop-The-World）**：主要在标记开始和标记结束时短暂停止

## 示例函数

### 1. `GCStats()` - 获取 GC 统计信息
展示如何获取 GC 的详细统计信息，包括 GC 次数、暂停时间等。

### 2. `ManualGC()` - 手动触发 GC
展示如何使用 `runtime.GC()` 手动触发垃圾回收。

### 3. `MemoryAllocation()` - 内存分配与 GC
展示内存分配与 GC 回收的关系，观察分配大量内存后 GC 的回收效果。

### 4. `GCPercent()` - 调整 GC 触发频率
展示如何使用 `debug.SetGCPercent()` 调整 GC 的触发阈值（GOGC 环境变量）。

### 5. `Finalizer()` - 使用终结器
展示如何使用 `runtime.SetFinalizer()` 在对象被 GC 回收前执行清理操作。

### 6. `HeapProfile()` - 堆内存分析
展示如何获取堆内存分析信息，用于内存泄漏检测。

### 7. `ConcurrentGC()` - 并发场景下的 GC
展示在多 goroutine 并发分配内存时的 GC 行为。

### 8. `MemoryLeakDetection()` - 内存泄漏检测
展示如何通过观察内存分配来检测潜在的内存泄漏。

### 9. `GCControl()` - GC 行为控制
展示如何禁用和重新启用 GC，以及手动控制 GC 的执行。

### 10. `EscapeAnalysis()` - 逃逸分析
展示逃逸分析的效果，对比栈上分配和堆上分配的区别。

## 使用说明

### 基本使用

```go
package main

import "github.com/yourusername/GoLearn/gc"

func main() {
    gc.GCStats()
    gc.ManualGC()
    gc.MemoryAllocation()
}
```

### 运行测试

```bash
cd gc
go test -v
```

### 调整 GC 参数

通过环境变量 GOGC 调整 GC 触发频率：

```bash
GOGC=200 go run your_program.go  # GC 阈值提高到 200%
GOGC=-1 go run your_program.go   # 禁用自动 GC
```

### 内存分析

使用 Go 自带的工具进行内存分析：

```bash
go test -bench=. -memprofile mem.prof
go tool pprof mem.prof
```

## 注意事项

1. **不要过度依赖 Finalizer**：Finalizer 的执行时机不确定，不应作为主要的资源释放方式
2. **避免频繁手动触发 GC**：手动触发 GC 会影响性能，应让 Go 运行时自动管理
3. **合理设置 GOGC**：根据应用的内存使用模式调整 GOGC 值
4. **关注逃逸分析**：通过 `go build -gcflags="-m"` 查看逃逸分析结果，优化内存分配

## 相关文档

- [Go 语言内存管理](https://golang.org/doc/gc)
- [Go 垃圾回收指南](https://github.com/golang/go/wiki/DesignDocuments)
