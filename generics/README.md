# Generics 泛型模块

## 学习目标
掌握Go 1.18+引入的泛型特性，理解类型参数、类型约束和泛型类型的使用。

## 核心概念

### 1. 泛型函数
```go
func Print[T any](value T) {
    fmt.Println(value)
}

Print(42)
Print("hello")
```

### 2. 类型参数
- `[T any]` - T可以是任何类型
- `[T comparable]` - T必须支持 == 和 !=
- `[T cmp.Ordered]` - T必须支持 < > <= >=

### 3. 泛型类型
```go
type Stack[T any] struct {
    elements []T
}
```

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| BasicGeneric | 基础泛型函数 | 初级 |
| GenericSlice | 泛型切片操作(map/filter) | 初级 |
| TypeConstraint | 类型约束详解 | 中级 |
| GenericSum | 自定义约束类型 | 中级 |
| GenericType | 泛型类型(栈实现) | 中级 |
| MultiTypeParam | 多类型参数 | 中级 |
| GenericMethod | 泛型方法 | 中级 |
| GenericInterface | 泛型接口 | 高级 |

## 类型约束

### 内置约束

| 约束 | 说明 | 支持的操作 |
|------|------|------------|
| `any` | 任意类型 | 无 |
| `comparable` | 可比较类型 | == != |
| `cmp.Ordered` | 有序类型 | < > <= >= |

### 自定义约束
```go
type Number interface {
    int | int64 | float64 | float32
}

func Sum[T Number](s []T) T {
    // ...
}
```

### 约束语法
```go
type Constraint interface {
    int | float64 | string  // 类型联合
    ~int                     // 底层类型是int
    Method()                 // 方法约束
}
```

## 泛型 vs 接口

### 使用泛型的场景
- 需要类型安全
- 需要保留原始类型信息
- 性能敏感（避免类型断言）
- 容器类型（栈、队列、链表）

### 使用接口的场景
- 需要运行时多态
- 类型不确定
- 需要动态类型检查

## 常见陷阱

### 1. 不能在方法中使用类型参数
```go
// 错误: 方法不能有类型参数
func (s *Stack[T]) Map[U any](f func(T) U) []U {
    // ...
}

// 正确: 使用泛型函数
func MapStack[T, U any](s *Stack[T], f func(T) U) []U {
    // ...
}
```

### 2. 类型推断失败
```go
// 有时需要显式指定类型参数
func Process[T any](v T) { }
Process(42)           // 类型推断成功
Process[int](42)      // 显式指定
```

### 3. 约束不满足
```go
type Number interface {
    int | float64
}

func Add[T Number](a, b T) T {
    return a + b  // 编译错误: + 操作符不在约束中
}

// 解决方案: 使用cmp.Ordered或自定义约束
```

## 最佳实践

1. **优先使用具体类型**: 如果不需要泛型，不要强制使用
2. **保持简单**: 复杂的泛型代码难以理解和维护
3. **合理使用约束**: 选择最小必要的约束
4. **文档化类型参数**: 解释每个类型参数的用途
5. **避免过度抽象**: 泛型不是解决所有问题的方案

## Go泛型发展历史

| 版本 | 特性 |
|------|------|
| Go 1.18 | 泛型正式发布 |
| Go 1.20 | 类型推断改进 |
| Go 1.21 | slices、maps泛型包 |
| Go 1.22 | 循环变量作用域修复 |

## 标准库泛型包

Go 1.21+提供了泛型工具包：

```go
import "slices"
import "maps"

slices.Sort(numbers)
slices.Contains(numbers, 42)
maps.Keys(m)
maps.Values(m)
```
