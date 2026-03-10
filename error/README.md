# Error 错误处理模块

## 学习目标
深入理解Go语言错误处理机制，包括自定义错误类型、错误包装、错误判断、panic/recover等高级特性。

## 核心概念

### 1. error接口
Go语言内置的error接口：
```go
type error interface {
    Error() string
}
```

### 2. nil判断陷阱
当自定义error类型返回nil指针时，与nil比较可能得到意外的结果。

### 3. error接口结构
接口值由两部分组成：类型(type)和值(value)。只有两者都为nil时，接口才等于nil。

### 4. 错误包装 (Go 1.13+)
使用 `%w` 格式化动词包装错误，形成错误链。

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| MyErrorEqualsNil | 自定义error的nil判断陷阱 | 中级 |
| CustomError | 自定义错误类型 | 初级 |
| ErrorWrapping | 错误包装演示 | 中级 |
| ErrorUnwrap | 错误解包演示 | 中级 |
| ErrorIs | errors.Is判断错误 | 中级 |
| ErrorAs | errors.As转换错误类型 | 中级 |
| PanicRecover | panic和recover演示 | 中级 |
| MultiErrorDemo | 多错误处理 | 高级 |
| ErrorWithContextDemo | 带上下文的错误 | 高级 |
| DeferPanicOrder | defer、panic、recover执行顺序 | 中级 |

## nil判断陷阱详解

### 问题示例
```go
type MyError struct {
    error
}

func returnsError() error {
    var e *MyError = nil
    return e  // 返回nil指针
}

func main() {
    err := returnsError()
    fmt.Println(err == nil)  // false! 为什么？
}
```

### 原因分析
```
接口值 = (type, value)

var e error = nil
// type = nil, value = nil
// e == nil → true

var e *MyError = nil
// type = *MyError, value = nil
// 作为error返回时，接口值 = (*MyError, nil)
// e == nil → false (因为type不是nil)
```

### 图示说明
```
┌─────────────────────────────────────┐
│           error接口                  │
├─────────────┬───────────────────────┤
│    type     │        value          │
├─────────────┼───────────────────────┤
│    nil      │        nil            │ → == nil: true
├─────────────┼───────────────────────┤
│  *MyError   │        nil            │ → == nil: false
├─────────────┼───────────────────────┤
│  *MyError   │   &MyError{...}       │ → == nil: false
└─────────────┴───────────────────────┘
```

## 新增示例详解

### CustomError - 自定义错误类型
创建包含详细信息的自定义错误类型。
```go
type DivisionError struct {
    Dividend int
    Divisor  int
    Message  string
}

func (e *DivisionError) Error() string {
    return fmt.Sprintf("division error: %s (dividend=%d, divisor=%d)", 
        e.Message, e.Dividend, e.Divisor)
}
```

### ErrorWrapping - 错误包装
使用 `%w` 包装错误，形成错误链。
```go
// 底层错误
err := fmt.Errorf("open file failed")

// 中间层包装
err = fmt.Errorf("read config: %w", err)

// 上层再包装
err = fmt.Errorf("load application: %w", err)
```

### ErrorUnwrap - 错误解包
使用 `errors.Unwrap` 逐层解包错误。
```go
baseErr := errors.New("base error")
wrapped := fmt.Errorf("layer: %w", baseErr)

unwrapped := errors.Unwrap(wrapped) // 返回 baseErr
```

### ErrorIs - 错误判断
使用 `errors.Is` 判断错误类型（支持错误链）。
```go
var ErrNotFound = errors.New("resource not found")

err := fmt.Errorf("find user: %w", ErrNotFound)

if errors.Is(err, ErrNotFound) {
    // 处理资源未找到
}
```

### ErrorAs - 错误类型转换
使用 `errors.As` 将错误转换为特定类型。
```go
var netErr *NetworkError
if errors.As(err, &netErr) {
    fmt.Printf("网络错误: code=%d\n", netErr.Code)
}
```

### PanicRecover - panic处理
使用 `recover` 捕获 panic 并转换为错误。
```go
func safeExecute(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()
    fn()
    return nil
}
```

### MultiError - 多错误处理
批量操作时收集多个错误。
```go
type MultiError struct {
    Errors []error
}

func (e *MultiError) Add(err error) {
    if err != nil {
        e.Errors = append(e.Errors, err)
    }
}
```

### ErrorWithContext - 带上下文的错误
错误包含操作上下文信息。
```go
type ErrorWithContext struct {
    Op   string    // 操作类型
    Path string    // 文件路径
    Err  error     // 原始错误
    Time string    // 时间戳
}
```

## 常见陷阱

### 1. 返回nil指针
```go
// 错误: 返回nil指针导致接口不为nil
func doSomething() error {
    var err *MyError = nil
    return err  // 接口不为nil!
}

// 正确: 直接返回nil
func doSomething() error {
    return nil
}
```

### 2. 检查自定义error
```go
// 使用类型断言检查
if err, ok := err.(*MyError); ok && err != nil {
    // 处理MyError
}
```

### 3. 错误比较
```go
// 错误: 使用 == 比较错误
if err == someError { }

// 正确: 使用 errors.Is
if errors.Is(err, someError) { }
```

### 4. panic后的defer不执行
```go
defer fmt.Println("1")
panic("error")
defer fmt.Println("2")  // 不会执行
```

## 最佳实践

1. **直接返回nil**: 不要返回nil指针，直接返回nil
2. **使用errors.New**: 简单错误使用`errors.New()`
3. **使用fmt.Errorf**: 需要格式化时使用`fmt.Errorf()`
4. **自定义错误类型**: 实现`Error()`方法
5. **错误包装**: 使用`fmt.Errorf("context: %w", err)`包装错误
6. **错误判断**: 使用`errors.Is`判断错误值
7. **错误转换**: 使用`errors.As`转换错误类型
8. **哨兵错误**: 定义可导出的错误变量用于比较

## 标准库错误处理

```go
import "errors"

// 创建简单错误
err := errors.New("something went wrong")

// 创建格式化错误
err := fmt.Errorf("failed to open %s", filename)

// 包装错误 (Go 1.13+)
err := fmt.Errorf("context: %w", originalErr)

// 解包错误
unwrapped := errors.Unwrap(err)

// 判断错误类型
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    // 处理PathError
}

// 判断错误值
if errors.Is(err, os.ErrNotExist) {
    // 文件不存在
}
```

## 错误处理模式

### 1. 哨兵错误模式
```go
var ErrNotFound = errors.New("not found")

if errors.Is(err, ErrNotFound) {
    // 处理
}
```

### 2. 自定义错误类型模式
```go
type MyError struct {
    Code int
    Msg  string
}

func (e *MyError) Error() string {
    return fmt.Sprintf("code=%d, msg=%s", e.Code, e.Msg)
}
```

### 3. 错误包装模式
```go
func outer() error {
    err := inner()
    if err != nil {
        return fmt.Errorf("outer context: %w", err)
    }
    return nil
}
```

### 4. 多错误收集模式
```go
var errs []error
for _, item := range items {
    if err := process(item); err != nil {
        errs = append(errs, err)
    }
}
if len(errs) > 0 {
    return errors.Join(errs...) // Go 1.20+
}
```

## 运行测试

```bash
go test -v ./error/...
```

## 参考资料

- [Go语言错误处理](https://go.dev/blog/error-handling-and-go)
- [Go 1.13错误处理](https://go.dev/blog/go1.13-errors)
- [Go错误处理最佳实践](https://go.dev/doc/effective_go#errors)
