# Closure 闭包模块

## 学习目标
理解Go语言闭包的工作原理、变量捕获机制和常见陷阱。

## 核心概念

### 1. 什么是闭包
闭包是一个函数值，它引用了其函数体之外的变量。该函数可以访问和修改这些变量。

```go
func adder() func(int) int {
    sum := 0
    return func(x int) int {
        sum += x  // 引用外部变量sum
        return sum
    }
}
```

### 2. 变量捕获
闭包捕获的是变量的引用，而不是值的副本：
- 闭包内修改会影响外部变量
- 外部修改会影响闭包内的值

### 3. 闭包的作用
- 状态封装（类似私有变量）
- 函数工厂
- 回调和延迟执行
- 函数式编程模式

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| BasicClosure | 基础闭包概念 | 初级 |
| ClosureFactory | 闭包工厂模式 | 初级 |
| LoopTrap | 循环变量捕获陷阱 | 中级 |
| GoroutineTrap | Goroutine中的闭包陷阱 | 中级 |
| ClosureState | 闭包状态封装 | 中级 |
| ClosureRecursion | 闭包递归 | 中级 |
| ClosurePerformance | 闭包性能考量 | 中级 |
| ClosureDefer | defer中的闭包 | 初级 |
| ClosureTimer | 闭包实现计时器 | 初级 |

## 常见陷阱

### 1. 循环变量捕获（经典陷阱）
```go
// 问题代码
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // 所有goroutine可能打印相同的值
    }()
}

// 解决方案1: 传递参数
for i := 0; i < 3; i++ {
    go func(n int) {
        fmt.Println(n)
    }(i)
}

// 解决方案2: 创建局部变量
for i := 0; i < 3; i++ {
    i := i  // 创建新的局部变量
    go func() {
        fmt.Println(i)
    }()
}
```

### 2. defer中的闭包
```go
// 问题: defer捕获的是最终值
x := 1
defer func() {
    fmt.Println(x)  // 打印100
}()
x = 100

// 解决: 传递参数
x := 1
defer func(v int) {
    fmt.Println(v)  // 打印1
}(x)
x = 100
```

## 闭包应用场景

### 1. 函数工厂
```go
func multiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}

double := multiplier(2)
triple := multiplier(3)
```

### 2. 状态封装
```go
func counter() func() int {
    count := 0  // 私有变量
    return func() int {
        count++
        return count
    }
}
```

### 3. 计时器
```go
func timer(name string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s: %v\n", name, time.Since(start))
    }
}

defer timer("operation")()
```

## Go 1.22的变化

Go 1.22修复了for循环变量捕获问题：
```go
// Go 1.22之前: 需要手动创建局部变量
// Go 1.22之后: 每次迭代自动创建新变量
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // Go 1.22+ 正确打印 0, 1, 2
    }()
}
```

## 最佳实践

1. **注意变量捕获**: 理解闭包捕获的是引用
2. **传递参数**: 在循环中使用闭包时，将变量作为参数传递
3. **使用go vet**: 检查循环变量捕获问题
4. **考虑性能**: 频繁创建闭包可能有性能开销
5. **升级Go版本**: Go 1.22+自动修复循环变量问题
