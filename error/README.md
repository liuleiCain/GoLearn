# Error 错误处理模块

## 学习目标
理解Go语言错误处理机制，特别是自定义error类型和nil判断的陷阱。

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

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| MyErrorEqualsNil | 自定义error的nil判断陷阱 | 中级 |

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

## 最佳实践

1. **直接返回nil**: 不要返回nil指针，直接返回nil
2. **使用errors.New**: 简单错误使用`errors.New()`
3. **使用fmt.Errorf**: 需要格式化时使用`fmt.Errorf()`
4. **自定义错误类型**: 实现`Error()`方法
5. **错误包装**: 使用`fmt.Errorf("context: %w", err)`包装错误

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
