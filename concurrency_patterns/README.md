# Concurrency Patterns 并发模式模块

## 学习目标
掌握Go语言中常用的并发设计模式，理解如何构建高效、可靠的并发系统。

## 核心概念

### 1. 并发 vs 并行
- **并发**: 多个任务交替执行，逻辑上的同时
- **并行**: 多个任务同时执行，物理上的同时

### 2. Go并发原语
- goroutine: 轻量级协程
- channel: 通信机制
- select: 多路复用
- context: 上下文控制
- sync包: 同步原语

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| WorkerPoolPattern | Worker Pool工作池模式 | 中级 |
| PipelinePattern | Pipeline流水线模式 | 中级 |
| FanOutFanIn | Fan-out/Fan-in扇出扇入模式 | 中级 |
| TimeoutPattern | 超时控制模式 | 初级 |
| CancellationPattern | 取消信号模式 | 初级 |
| SemaphorePattern | 信号量限流模式 | 中级 |
| RateLimitingPattern | 速率限制模式 | 中级 |
| GracefulShutdown | 优雅关闭模式 | 中级 |

## 并发模式详解

### 1. Worker Pool（工作池）
```
┌─────────┐     ┌─────────┐
│  Jobs   │────▶│ Worker1 │────┐
└─────────┘     ├─────────┤    │
                │ Worker2 │    ▼
                ├─────────┤ ┌─────────┐
                │ Worker3 │ │ Results │
                └─────────┘ └─────────┘
```
- 固定数量的worker处理任务
- 避免创建过多goroutine
- 控制并发度

### 2. Pipeline（流水线）
```
┌──────────┐     ┌─────────┐     ┌─────────┐
│ Generator│────▶│ Square  │────▶│ Double  │
└──────────┘     └─────────┘     └─────────┘
```
- 数据经过多个阶段处理
- 每个阶段是一个独立的goroutine
- 通过channel连接

### 3. Fan-out/Fan-in（扇出/扇入）
```
              ┌─────────┐
         ┌───▶│ Worker1 │───┐
┌────────┤    ├─────────┤   │    ┌────────┐
│ Input  │───▶│ Worker2 │───┼───▶│ Output │
└────────┘    ├─────────┤   │    └────────┘
         └───▶│ Worker3 │───┘
              └─────────┘
```
- Fan-out: 多个worker读取同一channel
- Fan-in: 多个channel合并到一个channel

### 4. 超时控制
```go
select {
case result := <-ch:
    // 处理结果
case <-time.After(timeout):
    // 超时处理
}
```

### 5. 取消信号
```go
ctx, cancel := context.WithCancel(context.Background())
// 发送取消信号
cancel()
// worker检查
select {
case <-ctx.Done():
    return
}
```

### 6. 信号量限流
```go
sem := make(chan struct{}, maxConcurrent)

// 获取信号量
sem <- struct{}{}
// 释放信号量
<-sem
```

### 7. 速率限制
```go
limiter := time.Tick(interval)
for req := range requests {
    <-limiter  // 等待下一个时间窗口
    process(req)
}
```

## 最佳实践

### 1. 资源管理
- 使用Worker Pool限制并发数
- 避免goroutine泄漏
- 及时关闭channel

### 2. 错误处理
- 使用error channel传递错误
- 实现超时机制
- 正确处理panic

### 3. 优雅关闭
- 使用context传递取消信号
- 等待所有goroutine完成
- 清理资源

### 4. 性能考量
- 合理设置并发数
- 避免过度竞争
- 使用缓冲channel提高吞吐量

## 常见陷阱

### 1. Goroutine泄漏
```go
// 错误: goroutine永远阻塞
go func() {
    ch <- value  // 没有接收者
}()
```

### 2. 关闭已关闭的channel
```go
close(ch)  // 第一次关闭
close(ch)  // panic!
```

### 3. 向已关闭的channel发送
```go
close(ch)
ch <- value  // panic!
```

## 模式选择指南

| 场景 | 推荐模式 |
|------|----------|
| 批量任务处理 | Worker Pool |
| 数据流水线处理 | Pipeline |
| 高并发读取 | Fan-out/Fan-in |
| API调用 | 超时+取消 |
| 限流保护 | 速率限制 |
| 资源限制 | 信号量 |
