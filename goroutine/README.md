# Goroutine 协程模块

## 学习目标
理解Go语言协程的基本概念、启动方式、同步机制和最佳实践。

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

## 最佳实践

1. 使用`go vet`检查循环变量捕获问题
2. 使用`sync.WaitGroup`等待goroutine完成
3. 监控goroutine数量，防止泄漏
4. 使用`context`包实现goroutine的取消机制
