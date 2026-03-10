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
| RangeChannel | 遍历channel | 初级 |
| RangeStringDetail | 字符串遍历细节 | 中级 |
| RangeWithBreak | break和continue用法 | 初级 |
| RangeNilCollection | 遍历nil集合 | 中级 |
| RangeCopyValue | range值是副本 | 中级 |
| RangeTwoValues | range返回值用法 | 初级 |
| RangeModifyDuringIteration | 遍历时修改集合 | 中级 |
| RangePerformance | range性能考量 | 高级 |

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

### 5. channel遍历
```go
ch := make(chan int)
go func() {
    ch <- 1
    ch <- 2
    close(ch) // 必须关闭，否则range会死锁
}()
for v := range ch {
    fmt.Println(v)
}
```

## 新增示例详解

### RangeChannel - 遍历Channel
使用range遍历channel，必须在发送方关闭channel。
```go
ch := make(chan int, 5)
go func() {
    for i := 1; i <= 5; i++ {
        ch <- i
    }
    close(ch) // 必须关闭
}()
for v := range ch {
    fmt.Println(v)
}
```

### RangeStringDetail - 字符串遍历细节
range遍历字符串按rune遍历，普通for循环按字节遍历。
```go
s := "Hello世界"
// range: 按rune遍历，索引是字节位置
for i, r := range s {
    fmt.Printf("位置: %d, rune: %c\n", i, r)
}
// 普通for: 按字节遍历
for i := 0; i < len(s); i++ {
    fmt.Printf("字节: %x\n", s[i])
}
```

### RangeWithBreak - break和continue
```go
// break提前退出
for _, v := range nums {
    if condition {
        break
    }
}

// continue跳过
for _, v := range nums {
    if condition {
        continue
    }
}

// 标签break退出外层循环
outer:
    for i := 0; i < 10; i++ {
        for j := 0; j < 10; j++ {
            if condition {
                break outer
            }
        }
    }
```

### RangeNilCollection - 遍历nil集合
```go
// nil切片和nil map: 安全，不执行循环体
var nilSlice []int
for i, v := range nilSlice { } // 不执行

var nilMap map[string]int
for k, v := range nilMap { } // 不执行

// nil channel: 永久阻塞！
var nilChan chan int
for v := range nilChan { } // 死锁
```

### RangeCopyValue - 值是副本
```go
items := []Item{{1}, {2}, {3}}

// 修改副本不影响原切片
for _, item := range items {
    item.Value = 10 // 无效
}

// 通过索引修改
for i := range items {
    items[i].Value = 10 // 有效
}

// 使用指针切片
ptrItems := []*Item{{1}, {2}, {3}}
for _, item := range ptrItems {
    item.Value = 10 // 有效
}
```

### RangeTwoValues - 返回值用法
```go
// 只获取key/索引
for k := range map { }
for i := range slice { }

// 只获取value
for _, v := range map { }
for _, v := range slice { }

// 获取两者
for k, v := range map { }
for i, v := range slice { }
```

### RangeModifyDuringIteration - 遍历时修改
```go
// 切片：追加不影响当前遍历
s := []int{1, 2, 3}
for i, v := range s {
    s = append(s, v+10) // 不影响当前遍历
}

// Map：修改已有key的值
m := map[string]int{"a": 1, "b": 2}
for k, v := range m {
    m["b"] = 100 // 可能影响后续读取
}

// Map：删除元素
for k := range m {
    delete(m, k) // 安全
}
```

### RangePerformance - 性能考量
```go
// 大结构体切片：使用索引避免拷贝
type BigStruct struct { data [100]int }
bigSlice := make([]BigStruct, 10000)

// 值拷贝开销大
for _, v := range bigSlice { }

// 索引访问无拷贝
for i := range bigSlice {
    _ = bigSlice[i].data[0]
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

### 4. nil channel死锁
```go
var ch chan int
for v := range ch { }  // 永久阻塞！
```

## 最佳实践

1. **Go 1.22+**: 升级Go版本解决循环变量问题
2. **需要修改**: 使用索引访问`slice[i]`
3. **大结构体**: 使用指针切片`[]*Struct`
4. **map顺序**: 如需有序，先排序key
5. **goroutine**: 注意闭包捕获问题
6. **channel**: 使用range遍历必须关闭channel
7. **nil集合**: nil切片和map遍历安全，nil channel会阻塞

## 运行测试

```bash
go test -v ./for_range/...
```

## 参考资料

- [Go语言for-range](https://go.dev/blog/range)
- [Go 1.22 Loop Changes](https://go.dev/doc/go1.22#language)
