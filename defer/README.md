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
| DeferTiming | defer执行时机演示 | 初级 |
| DeferStackTrace | defer打印堆栈信息 | 中级 |
| DeferMutexUnlock | defer解锁互斥锁 | 中级 |
| DeferFileOperation | defer文件操作模拟 | 中级 |
| DeferMethod | defer调用方法 | 初级 |
| DeferPerformance | defer性能对比 | 高级 |
| DeferNilFunction | defer nil函数演示 | 中级 |
| DeferArguments | defer参数计算时机 | 中级 |
| DeferNamedReturnMultiple | 多命名返回值与defer | 高级 |

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

## 新增示例详解

### DeferTiming - 执行时机
演示defer在外层函数返回前执行，包括return语句赋值之后。
```go
func deferTimingHelper() int {
    defer fmt.Println("defer执行: 在return之后")
    return 42  // defer在此之后执行
}
```

### DeferMutexUnlock - 互斥锁解锁
这是defer最常见的使用场景之一，确保锁一定会被释放。
```go
mu.Lock()
defer mu.Unlock()  // 无论是否发生panic，都会解锁
```

### DeferFileOperation - 文件操作
模拟文件操作中defer的使用，确保资源释放。
```go
f, err := os.Open("file.txt")
if err != nil {
    return err
}
defer f.Close()  // 确保文件关闭
```

### DeferMethod - 方法调用
defer可以调用方法，方法接收者会被正确捕获。
```go
p := &person{name: "张三"}
defer p.sayHello()  // defer可以调用方法
```

### DeferPerformance - 性能对比
演示defer的性能影响，在现代Go版本中defer性能已大幅优化。
```go
// 不使用defer
func withoutDefer() {
    x := 1
    x++
}

// 使用defer
func withDefer() {
    x := 1
    defer func() { _ = x }()
    x++
}
```

### DeferNilFunction - nil函数panic
defer调用nil函数会触发panic。
```go
var fn func()
defer fn()  // fn是nil，会panic
```

### DeferArguments - 参数计算时机
参数在defer声明时就已经计算，闭包捕获的是引用。
```go
x := 1
defer fmt.Printf("x=%d\n", x)  // 输出1，参数预计算
defer func() {
    fmt.Printf("x=%d\n", x)  // 输出当前值，闭包捕获引用
}()
```

### DeferNamedReturnMultiple - 多命名返回值
defer可以修改多个命名返回值。
```go
func example() (a int, b string, c bool) {
    a, b, c = 10, "hello", true
    defer func() {
        a *= 2    // 修改命名返回值
        b += " world"
        c = false
    }()
    return  // 返回 20, "hello world", false
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

### 4. defer nil函数
```go
var fn func()
defer fn()  // 会panic: nil function
```

## 最佳实践

1. **资源释放**: 使用defer确保文件、锁等资源被释放
2. **错误处理**: 使用defer + recover处理panic
3. **避免循环defer**: 在循环中使用闭包包裹defer
4. **参数传递**: 需要当前值时，将变量作为参数传递给defer
5. **命名返回值**: 需要defer修改返回值时使用命名返回值
6. **性能考虑**: 在性能关键路径上谨慎使用defer（现代Go已优化）
7. **方法调用**: defer可以调用方法，注意接收者的捕获

## 运行测试

```bash
go test -v ./defer/...
```

## 参考资料

- [Go语言defer详解](https://go.dev/blog/defer-panic-and-recover)
- [Go defer性能优化](https://go.dev/doc/go1.14#runtime)
