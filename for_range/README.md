# For-Range 循环遍历模块

## 学习目标
深入理解Go语言for-range语句的遍历机制、变量捕获和Go 1.22的变化。

## 核心概念

### 1. for-range语法
```go
for index, value := range collection {
    // index: 索引/键
    // value: 值的副本
}
```

### 2. 遍历类型
- 数组/切片: `for i, v := range slice`
- map: `for k, v := range map`
- 字符串: `for i, r := range string`
- channel: `for v := range ch`
- 整数 (Go 1.22+): `for i := range n`

### 3. 迭代变量特性
- 迭代变量在每次循环中重用
- value是值的副本，不是引用

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| AnalysisBasic | 基础遍历、字符串遍历 | 初级 |
| AnalysisArr | 数组遍历特性 | 初级 |
| AnalysisSlice | 切片遍历与修改 | 初级 |
| AnalysisMap | map遍历顺序随机性 | 初级 |
| AnalysisRangeSliceIndex | goroutine中的闭包陷阱 | 中级 |
| AnalysisRangeStruct | 结构体遍历与指针 | 中级 |

## 遍历行为详解

### 1. 数组遍历
```go
arr := [3]int{1, 2, 3}
for i, v := range arr {
    // v是arr[i]的副本
    // 修改v不影响arr
}
```

### 2. 切片遍历
```go
slice := []int{1, 2, 3}
for i, v := range slice {
    slice[i] = v + 1  // 可以通过索引修改
}
```

### 3. map遍历顺序
```go
// map遍历顺序是随机的！
m := map[string]int{"a": 1, "b": 2, "c": 3}
for k, v := range m {
    // 每次遍历顺序可能不同
}
```

### 4. 字符串遍历
```go
s := "你好世界"
for i, r := range s {
    // i: 字节位置
    // r: rune(Unicode码点)
}
```

## Go 1.22的变化

### 循环变量作用域
```go
// Go 1.22之前: 循环变量共享，需要手动创建副本
for i, v := range slice {
    go func() {
        fmt.Println(i, v)  // 可能打印相同的值
    }()
}

// Go 1.22之后: 每次迭代创建新变量
for i, v := range slice {
    go func() {
        fmt.Println(i, v)  // 正确打印每个值
    }()
}
```

### 整数遍历
```go
// Go 1.22+支持
for i := range 5 {
    fmt.Println(i)  // 0, 1, 2, 3, 4
}
```

## 常见陷阱

### 1. 循环变量捕获 (Go 1.22之前)
```go
// 问题代码
var funcs []func()
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() {
        fmt.Println(i)  // 都打印3
    })
}

// 解决方案1: 传递参数
for i := 0; i < 3; i++ {
    funcs = append(funcs, func(n int) func() {
        return func() { fmt.Println(n) }
    }(i))
}

// 解决方案2: 创建局部变量
for i := 0; i < 3; i++ {
    i := i
    funcs = append(funcs, func() {
        fmt.Println(i)
    })
}
```

### 2. 遍历时修改map
```go
m := map[string]int{"a": 1, "b": 2}
for k := range m {
    delete(m, k)  // 安全
    m["c"] = 3    // 可能导致不确定行为
}
```

### 3. value是副本
```go
type Person struct { Name string }
people := []Person{{"Alice"}, {"Bob"}}
for _, p := range people {
    p.Name = "Changed"  // 不影响原切片
}
```

## 最佳实践

1. **Go 1.22+**: 升级Go版本解决循环变量问题
2. **需要修改**: 使用索引访问`slice[i]`
3. **大结构体**: 使用指针切片`[]*Struct`
4. **map顺序**: 如需有序，先排序key
5. **goroutine**: 注意闭包捕获问题
