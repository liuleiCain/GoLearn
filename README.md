# Go语言学习教程

## 项目介绍

这是一个全面的Go语言学习教程项目，通过实际代码示例帮助新手和老手掌握Go语言的核心特性和高级技巧。每个模块都包含详细的代码注释、测试用例和README文档。

## 学习路径

### 初级阶段 - 基础语法

| 模块 | 内容 | 难度 | 文档 |
|------|------|------|------|
| [switch_case](./switch_case) | switch语句和fallthrough | ⭐ | [README](./switch_case/README.md) |
| [map](./map) | map初始化、遍历顺序 | ⭐ | [README](./map/README.md) |
| [json](./json) | JSON序列化与反序列化 | ⭐ | [README](./json/README.md) |
| [slice](./slice) | 切片原理、扩容机制、append、copy | ⭐⭐ | [README](./slice/README.md) |
| [interface](./interface) | 接口定义、多态、空接口、类型断言 | ⭐⭐ | [README](./interface/README.md) |

### 中级阶段 - 并发编程

| 模块 | 内容 | 难度 | 文档 |
|------|------|------|------|
| [goroutine](./goroutine) | 协程基础、WaitGroup、sync.Once | ⭐⭐ | [README](./goroutine/README.md) |
| [channel](./channel) | 无缓冲/有缓冲通道、select、关闭通道 | ⭐⭐ | [README](./channel/README.md) |
| [context](./context) | 超时控制、取消信号、值传递 | ⭐⭐ | [README](./context/README.md) |
| [concurrency_patterns](./concurrency_patterns) | Worker Pool、Pipeline、Fan-out/Fan-in | ⭐⭐⭐ | [README](./concurrency_patterns/README.md) |

### 中级阶段 - 语言特性

| 模块 | 内容 | 难度 | 文档 |
|------|------|------|------|
| [defer](./defer) | defer执行顺序、闭包捕获、panic处理 | ⭐⭐ | [README](./defer/README.md) |
| [for_range](./for_range) | for-range遍历各种数据结构 | ⭐⭐ | [README](./for_range/README.md) |
| [break](./break) | break标签、select跳出循环 | ⭐⭐ | [README](./break/README.md) |
| [error](./error) | 自定义error、nil判断陷阱 | ⭐⭐ | [README](./error/README.md) |
| [receiver](./receiver) | 值接收者vs指针接收者 | ⭐⭐ | [README](./receiver/README.md) |
| [closure](./closure) | 闭包原理、变量捕获、循环陷阱 | ⭐⭐ | [README](./closure/README.md) |
| [pointer](./pointer) | 指针基础、逃逸分析、性能考量 | ⭐⭐ | [README](./pointer/README.md) |
| [init_func](./init_func) | init函数执行顺序 | ⭐⭐ | [README](./init_func/README.md) |

### 高级阶段 - 现代特性

| 模块 | 内容 | 难度 | 文档 |
|------|------|------|------|
| [generics](./generics) | 泛型函数、泛型类型、类型约束 | ⭐⭐⭐ | [README](./generics/README.md) |

## 快速开始

### 环境要求
- Go 1.23.6 或更高版本

### 运行测试
```bash
# 运行所有测试
go test ./...

# 运行特定模块测试
go test ./goroutine/...
go test ./channel/...
```

### 运行示例
```go
package main

import (
    "fmt"
    "go-learn/goroutine"
)

func main() {
    goroutine.BasicGoroutine()
    goroutine.WaitGroupDemo()
}
```

## 模块详解

### 并发编程核心

#### Goroutine（协程）
Go语言的轻量级线程，由Go运行时调度。
- 启动方式：`go func()`
- 同步机制：`sync.WaitGroup`
- 单次执行：`sync.Once`

#### Channel（通道）
goroutine之间的通信机制。
- 无缓冲：同步通信
- 有缓冲：异步通信
- select：多路复用

#### Context（上下文）
控制goroutine的生命周期。
- 超时控制：`WithTimeout`
- 取消信号：`WithCancel`
- 值传递：`WithValue`

### 语言特性深入

#### Defer（延迟执行）
确保函数退出时执行清理操作。
- LIFO顺序执行
- 参数预计算
- 配合recover处理panic

#### Interface（接口）
Go语言的多态实现。
- 隐式实现
- 空接口`any`
- 类型断言

#### Generics（泛型）
Go 1.18+引入的类型参数。
- 泛型函数
- 泛型类型
- 类型约束

## 常见陷阱

### 1. 循环变量捕获
```go
// 错误
for i := 0; i < 3; i++ {
    go func() { fmt.Println(i) }()  // 可能打印相同的值
}

// 正确
for i := 0; i < 3; i++ {
    go func(n int) { fmt.Println(n) }(i)
}
```

### 2. 接口nil判断
```go
var p *Person  // nil
var s Shape = p  // s != nil!
```

### 3. 切片内存泄漏
```go
// 小切片引用大数组
large := make([]int, 1000000)
small := large[:10]  // 仍引用整个大数组

// 解决方案
small := make([]int, 10)
copy(small, large[:10])
```

## 最佳实践

1. **错误处理**: 不要忽略错误，及时处理
2. **并发安全**: 使用channel或sync包保护共享资源
3. **资源管理**: 使用defer确保资源释放
4. **接口设计**: 保持接口小而专注
5. **代码风格**: 遵循Go官方代码规范

## 项目结构

```
go-learn/
├── goroutine/          # 协程模块
├── channel/            # 通道模块
├── context/            # 上下文模块
├── interface/          # 接口模块
├── slice/              # 切片模块
├── map/                # map模块
├── defer/              # defer模块
├── for_range/          # for-range模块
├── break/              # break模块
├── error/              # error模块
├── receiver/           # 方法接收者模块
├── closure/            # 闭包模块
├── pointer/            # 指针模块
├── generics/           # 泛型模块
├── concurrency_patterns/  # 并发模式模块
├── switch_case/        # switch模块
├── json/               # JSON模块
├── init_func/          # init函数模块
├── go.mod
└── README.md
```

## 参考资源

- [Go官方文档](https://golang.org/doc/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go语言圣经](https://gopl-zh.github.io/)

## 贡献指南

欢迎提交Issue和Pull Request来完善这个教程项目。

## 许可证

[MIT License](./LICENSE)
