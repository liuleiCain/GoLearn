# Channel 通道模块

## 学习目标
掌握Go语言通道的核心概念、使用方法和常见并发模式。

## 核心概念

### 1. 通道类型
- **无缓冲通道**: `make(chan T)` - 同步通信，发送和接收必须同时准备好
- **有缓冲通道**: `make(chan T, n)` - 异步通信，可以缓存n个值

### 2. 通道操作
- **发送**: `ch <- value`
- **接收**: `value := <-ch` 或 `value, ok := <-ch`
- **关闭**: `close(ch)` - 只应由发送者关闭

### 3. Select语句
- 多路复用，同时监听多个通道
- 随机选择一个就绪的case执行
- 配合`default`实现非阻塞操作

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| UnbufferedChannel | 无缓冲通道基础 | 初级 |
| BufferedChannel | 有缓冲通道使用 | 初级 |
| ChannelDirection | 单向通道(只发送/只接收) | 初级 |
| CloseChannel | 通道关闭和range遍历 | 初级 |
| SelectDemo | Select多路复用 | 中级 |
| SelectTimeout | Select配合超时 | 中级 |
| SelectNonBlocking | 非阻塞Select | 中级 |
| NilChannel | nil通道的特殊行为 | 高级 |
| WorkerPool | Worker Pool并发模式 | 中级 |
| PingPong | Ping-Pong通信模式 | 中级 |

## 通道规则

### 1. 关闭通道的规则
- 只有发送者应该关闭通道
- 接收者不应关闭通道
- 关闭已关闭的通道会panic
- 向已关闭的通道发送会panic
- 从已关闭的通道接收会返回零值和false

### 2. nil通道的行为
| 操作 | 行为 |
|------|------|
| 从nil通道接收 | 永久阻塞 |
| 向nil通道发送 | 永久阻塞 |
| 关闭nil通道 | panic |
| select中的nil通道case | 永远被忽略 |

## 常见陷阱

### 1. 死锁
```go
// 死锁示例
ch := make(chan int)
ch <- 1  // 永久阻塞，没有接收者

// 正确做法
ch := make(chan int, 1)
ch <- 1  // 有缓冲，不会阻塞
```

### 2. 向已关闭通道发送
```go
ch := make(chan int)
close(ch)
ch <- 1  // panic: send on closed channel
```

### 3. 关闭已关闭的通道
```go
ch := make(chan int)
close(ch)
close(ch)  // panic: close of closed channel
```

## 最佳实践

1. 使用`defer close(ch)`确保通道被关闭
2. 使用`for range`遍历通道，自动检测关闭
3. 使用`select`配合`time.After`实现超时
4. 使用缓冲通道提高性能，但要注意容量规划
5. 在select中使用nil通道来禁用某些case
