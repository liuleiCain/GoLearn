# Switch Case 条件分支模块

## 学习目标
掌握Go语言switch语句的各种用法，包括fallthrough穿透、类型switch、无表达式switch等特性。

## 核心概念

### 1. switch基础
Go语言的switch比C/Java更灵活：
- 默认不需要break，匹配后自动退出
- case可以有多个条件
- 可以没有表达式，变成if-else链

### 2. fallthrough
使用`fallthrough`穿透到下一个case，**不判断条件直接执行**。

### 3. 类型switch
使用`.(type)`进行类型判断，只能用在switch语句中。

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| BasicSwitch | 基础switch用法 | 初级 |
| FallthroughDemo | fallthrough穿透行为 | 初级 |
| MultiCaseDemo | 多条件case | 初级 |
| NoExpressionSwitch | 无表达式switch | 初级 |
| TypeSwitch | 类型switch | 初级 |
| SwitchBreak | switch中的break | 中级 |
| SwitchInFor | for循环中的switch | 中级 |
| SwitchInitialization | switch初始化语句 | 初级 |
| SwitchPitfalls | 常见陷阱 | 中级 |

## switch语法详解

### 1. 基础语法
```go
switch 表达式 {
case 值1:
    // 执行代码
case 值2, 值3:  // 多个条件
    // 执行代码
default:
    // 默认执行
}
```

### 2. 无表达式switch
```go
switch {
case 条件1:
    // 执行代码
case 条件2:
    // 执行代码
default:
    // 默认执行
}
```

### 3. 类型switch
```go
switch v := x.(type) {
case int:
    fmt.Printf("int: %d\n", v)
case string:
    fmt.Printf("string: %s\n", v)
default:
    fmt.Printf("unknown: %T\n", v)
}
```

### 4. 初始化语句
```go
switch os := runtime.GOOS; os {
case "windows":
    // ...
default:
    // ...
}
```

## fallthrough行为

### 工作原理
```go
switch i {
case 0:
    fmt.Println("case 0")
    fallthrough  // 穿透到case 1，不判断i==1
case 1:
    fmt.Println("case 1")  // 即使i!=1也会执行
}
```

### 执行流程
```
i=0: case 0 -> fallthrough -> case 1 (不判断条件)
i=1: case 1 (直接匹配)
i=2: 无匹配case
```

## 常见陷阱

### 1. fallthrough不判断条件
```go
switch x {
case 1:
    fallthrough
case 2:
    // 即使x!=2也会执行！
}
```

### 2. fallthrough必须在case块末尾
```go
// 错误: fallthrough后不能有其他语句
case 1:
    fmt.Println("1")
    fallthrough
    fmt.Println("error")  // 编译错误
```

### 3. 无表达式switch的case必须是布尔表达式
```go
// 错误
switch {
case 1:  // 编译错误，不是布尔表达式
}

// 正确
switch {
case x == 1:  // 布尔表达式
}
```

## switch vs if-else

| 特性 | switch | if-else |
|------|--------|---------|
| 可读性 | 多条件时更好 | 条件复杂时更好 |
| 性能 | 编译器优化 | 无特殊优化 |
| 灵活性 | 固定模式 | 任意条件 |
| 类型判断 | 支持 | 需要类型断言 |

## 最佳实践

1. **优先使用switch**: 多个固定值判断时使用switch
2. **避免过度fallthrough**: 代码可读性下降
3. **使用初始化语句**: 限制变量作用域
4. **类型switch**: 处理interface{}时使用
5. **default放最后**: 即使Go不强制，也保持习惯

## 与其他语言对比

| 特性 | Go | C/Java |
|------|-----|--------|
| 默认break | ✓ | ✗ |
| fallthrough | 显式声明 | 默认行为 |
| 多条件case | ✓ | ✗ |
| 无表达式switch | ✓ | ✗ |
| 类型switch | ✓ | ✗ |
