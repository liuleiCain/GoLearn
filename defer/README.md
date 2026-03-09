# Defer 延迟执行模块

## 学习目标
深入理解Go语言defer语句的执行顺序、参数捕获机制和与panic/recover的配合使用。

## 核心概念

### 1. defer基础
defer语句将函数调用推迟到外层函数返回之前执行。

### 2. 执行顺序
多个defer按LIFO（后进先出）顺序执行，类似栈结构。

### 3. 参数捕获
defer语句中的参数在声明时就已确定，而不是在执行时。

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| DeferFor | 循环中使用defer | 初级 |
| DeferForParams | defer参数捕获机制 | 中级 |
| DeferForReturn | defer与返回值的关系 | 中级 |
| DeferForReturn2 | 命名返回值与defer | 中级 |
| DeferForPanic | defer配合recover处理panic | 中级 |
| DeferForOldParams | defer保存旧值 | 初级 |

## defer执行规则

### 1. LIFO顺序
```go
defer fmt.Println(1)  // 最后执行
defer fmt.Println(2)  // 中间执行
defer fmt.Println(3)  // 最先执行
// 输出: 3, 2, 1
```

### 2. 参数预计算
```go
x := 1
defer fmt.Println(x)  // 输出1，不是2
x = 2
```

### 3. 闭包捕获引用
```go
x := 1
defer func() {
    fmt.Println(x)  // 输出2，捕获的是引用
}()
x = 2
```

### 4. 命名返回值
```go
func example() (x int) {
    x = 1
    defer func() {
        x++  // 可以修改命名返回值
    }()
    return x  // 返回2
}
```

## 常见陷阱

### 1. 循环中使用defer
```go
// 错误: defer在循环结束后才执行，可能导致资源延迟释放
for _, file := range files {
    f, _ := os.Open(file)
    defer f.Close()  // 所有文件在函数结束时才关闭
}

// 正确: 使用闭包立即释放
for _, file := range files {
    func() {
        f, _ := os.Open(file)
        defer f.Close()  // 每次迭代结束时关闭
    }()
}
```

### 2. defer与返回值
```go
func example() int {
    x := 1
    defer func() {
        x++  // 修改的是局部变量x，不影响返回值
    }()
    return x  // 返回1
}

func example2() (x int) {
    x = 1
    defer func() {
        x++  // 修改的是命名返回值
    }()
    return x  // 返回2
}
```

### 3. panic后的defer不执行
```go
defer fmt.Println("1")
panic("error")
defer fmt.Println("2")  // 不会执行
```

## 最佳实践

1. **资源释放**: 使用defer确保文件、锁等资源被释放
2. **错误处理**: 使用defer + recover处理panic
3. **避免循环defer**: 在循环中使用闭包包裹defer
4. **参数传递**: 需要当前值时，将变量作为参数传递给defer
5. **命名返回值**: 需要defer修改返回值时使用命名返回值
