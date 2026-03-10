# Go 内存模型模块

本模块详细介绍 Go 语言的内存模型及其并发编程的核心概念。

## 目录

- [内存模型简介](#内存模型简介)
- [示例函数](#示例函数)
- [使用说明](#使用说明)
- [关键概念](#关键概念)

## 内存模型简介

Go 内存模型定义了在一个 goroutine 中对变量的写入，如何能被另一个 goroutine 中对该变量的读取观察到。理解内存模型对于编写正确的并发程序至关重要。

## 示例函数

### 1. `DataRace()` - 数据竞争问题
展示在没有同步机制的情况下，多个 goroutine 同时读写共享变量会导致的数据竞争问题。

### 2. `AtomicOperation()` - 原子操作
展示如何使用 `sync/atomic` 包提供的原子操作来避免数据竞争。

### 3. `MutexLock()` - 互斥锁
展示如何使用 `sync.Mutex` 互斥锁来保护共享资源的访问。

### 4. `HappensBefore()` - Happens-Before 关系
展示 Go 内存模型中 happens-before 关系的重要性，确保操作的可见性。

### 5. `ChannelSynchronization()` - 通道同步
展示如何使用通道（channel）进行 goroutine 之间的同步，这是 Go 中推荐的同步方式。

### 6. `OnceInitialization()` - 一次性初始化
展示如何使用 `sync.Once` 确保某个函数只被执行一次，常用于单例模式。

### 7. `WaitGroupSynchronization()` - WaitGroup 同步
展示如何使用 `sync.WaitGroup` 等待一组 goroutine 完成。

### 8. `MemoryOrdering()` - 内存重排序
展示编译器和 CPU 可能进行的指令重排序，以及为什么需要内存屏障。

### 9. `RWMutexExample()` - 读写锁
展示如何使用 `sync.RWMutex` 读写锁，允许多个读者同时访问，但写操作互斥。

### 10. `AtomicValue()` - 原子值
展示如何使用 `atomic.Value` 安全地读写任意类型的值。

## 使用说明

### 基本使用

```go
package main

import "github.com/yourusername/GoLearn/memory_model"

func main() {
    memorymodel.DataRace()
    memorymodel.AtomicOperation()
    memorymodel.ChannelSynchronization()
}
```

### 检测数据竞争

使用 `-race` 标志运行程序来检测数据竞争：

```bash
go run -race your_program.go
go test -race ./memory_model
```

### 运行测试

```bash
cd memory_model
go test -v
```

## 关键概念

### 1. Happens-Before 关系

Happens-Before 是 Go 内存模型的核心概念，它定义了操作之间的顺序关系：

- **初始化顺序**：包级别的变量初始化 happens-before 该包中任何函数的执行
- **goroutine 创建**：`go` 语句 happens-before 新 goroutine 中的任何操作
- **channel 通信**：
  - 向 channel 发送 happens-before 从该 channel 接收完成
  - 关闭 channel happens-before 从该 channel 接收到零值
  - 对于无缓冲 channel，接收 happens-before 发送完成
- **锁**：
  - 对于 `sync.Mutex` 或 `sync.RWMutex`，第 n 次 `Unlock()` happens-before 第 n+1 次 `Lock()` 返回
  - 对于 `sync.RWMutex`，任何 `RLock()` happens-before 对应的 `RUnlock()`
- **Once**：`sync.Once.Do(f)` 中对 `f` 的调用 happens-before 任何 `Once.Do(f)` 调用返回

### 2. 数据竞争

当两个 goroutine 同时访问同一个变量，且至少有一个访问是写入操作时，就会发生数据竞争。使用 `-race` 标志可以检测数据竞争。

### 3. 同步原语选择

- **通道**：优先使用，符合 Go 的设计哲学
- **互斥锁**：适用于简单的临界区保护
- **读写锁**：适用于读多写少的场景
- **原子操作**：适用于简单的数值操作
- **WaitGroup**：适用于等待多个 goroutine 完成
- **Once**：适用于一次性初始化

## 注意事项

1. **避免数据竞争**：始终使用同步机制保护共享变量
2. **优先使用通道**：Go 推荐使用通道进行 goroutine 间通信
3. **使用 race 检测器**：在测试时始终使用 `-race` 标志
4. **理解 happens-before**：确保你的并发程序满足正确的 happens-before 关系
5. **不要过度使用原子操作**：原子操作难以正确使用，优先考虑更高层次的同步原语

## 相关文档

- [Go 内存模型规范](https://golang.org/ref/mem)
- [Go 并发编程](https://golang.org/doc/effective_go#concurrency)
- [Go 数据竞争检测](https://golang.org/doc/articles/race_detector.html)
