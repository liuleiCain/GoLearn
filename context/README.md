# Context 上下文模块

## 学习目标
掌握Go语言context包的核心功能，理解如何在并发编程中进行超时控制、取消信号传播和请求范围值传递。

## 核心概念

### 1. Context接口
```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

### 2. 创建Context
| 函数 | 说明 |
|------|------|
| `context.Background()` | 返回一个非nil的空Context，作为根Context |
| `context.TODO()` | 当不确定使用哪个Context时使用 |

### 3. 派生Context
| 函数 | 说明 |
|------|------|
| `WithCancel(parent)` | 返回一个可取消的Context |
| `WithTimeout(parent, duration)` | 返回一个带超时的Context |
| `WithDeadline(parent, time)` | 返回一个带截止时间的Context |
| `WithValue(parent, key, value)` | 返回一个携带键值对的Context |

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| BackgroundDemo | Background和TODO的区别 | 初级 |
| WithCancelDemo | 手动取消Context | 初级 |
| WithTimeoutDemo | 超时自动取消 | 初级 |
| WithDeadlineDemo | 指定截止时间 | 初级 |
| WithValueDemo | 请求范围值传递 | 初级 |
| ChainContextDemo | 链式Context派生 | 中级 |
| PropagationDemo | Context在调用链中传播 | 中级 |
| ErrDemo | Context错误处理 | 中级 |

## 使用规则

### 1. 不要将Context存储在结构体中
```go
// 错误示例
type Handler struct {
    ctx context.Context
}

// 正确示例 - 作为函数参数传递
func (h *Handler) Handle(ctx context.Context) {
}
```

### 2. Context作为函数第一个参数
```go
func ProcessRequest(ctx context.Context, id int) error {
}
```

### 3. 不要传递nil Context
```go
// 错误示例
ProcessRequest(nil, 1)

// 正确示例
ProcessRequest(context.Background(), 1)
```

### 4. Context.Value仅用于请求范围值
- 不要用于传递可选参数
- 不要用于传递业务数据
- 适合传递trace ID、认证信息等

## 常见陷阱

### 1. 忘记调用cancel
```go
// 错误示例 - 可能导致goroutine泄漏
ctx, _ := context.WithTimeout(context.Background(), time.Second)

// 正确示例 - 始终调用cancel
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
```

### 2. Context取消后继续使用
```go
ctx, cancel := context.WithCancel(context.Background())
cancel()

// 取消后ctx.Done()会立即返回
select {
case <-ctx.Done():
    fmt.Println("已取消")  // 立即执行
}
```

## 最佳实践

1. 始终使用`defer cancel()`确保资源释放
2. Context作为函数第一个参数传递
3. 使用自定义类型作为Value的key，避免冲突
4. 在长时间运行的操作中定期检查`ctx.Done()`
5. 使用`context.WithoutCancel`(Go 1.21+)在取消后继续某些操作
