# Pointer 指针模块

## 学习目标
理解Go语言指针的基本概念、值接收者与指针接收者的区别、逃逸分析和性能优化。

## 核心概念

### 1. 指针基础
```go
x := 42
p := &x    // p是指向x的指针
*p = 100   // 通过指针修改x的值
```

### 2. 指针特点
- Go指针不支持算术运算（除非使用unsafe包）
- 自动垃圾回收，无悬空指针问题
- nil指针访问会panic

### 3. 值类型 vs 引用类型
| 类型 | 传递方式 | 示例 |
|------|----------|------|
| 值类型 | 复制 | int, float, struct, array |
| 引用类型 | 引用 | slice, map, channel, pointer |

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| BasicPointer | 基础指针操作 | 初级 |
| PointerArithmetic | unsafe指针运算 | 高级 |
| ReceiverCompare | 值接收者vs指针接收者 | 中级 |
| NilPointer | nil指针处理 | 初级 |
| PointerToPointer | 多级指针 | 中级 |
| newAndMake | new和make的区别 | 初级 |
| EscapeAnalysis | 逃逸分析 | 高级 |
| PointerPerformance | 指针性能考量 | 中级 |
| PointerSafety | 指针安全性 | 中级 |

## 值接收者 vs 指针接收者

### 值接收者
```go
func (p Person) SetName(name string) {
    p.Name = name  // 修改的是副本
}
```

### 指针接收者
```go
func (p *Person) SetName(name string) {
    p.Name = name  // 修改的是原值
}
```

### 选择原则
| 场景 | 推荐 |
|------|------|
| 需要修改接收者 | 指针接收者 |
| 结构体较大 | 指针接收者 |
| 一致性要求 | 统一使用指针接收者 |
| 小结构体、只读 | 值接收者 |

## new vs make

| 函数 | 用途 | 返回值 | 适用类型 |
|------|------|--------|----------|
| new | 分配内存 | 指针 | 所有类型 |
| make | 初始化 | 值 | slice, map, channel |

```go
p := new(int)        // 返回 *int，值为0
s := make([]int, 5)  // 返回 []int，长度为5
```

## 逃逸分析

### 什么是逃逸
变量从栈"逃逸"到堆，由Go编译器自动决定。栈分配更快，但生命周期受限；堆分配需要GC管理。

### 查看逃逸分析
```bash
go build -gcflags="-m"      # 基本逃逸分析
go build -gcflags="-m -m"   # 详细逃逸分析
```

### 逃逸场景详解

#### 1. 返回局部变量的指针
```go
// 逃逸到堆
func heapAlloc() *int {
    x := 42
    return &x  // x逃逸，因为指针返回到函数外
}

// 栈分配
func stackAlloc() int {
    x := 42
    return x  // x不逃逸，值复制返回
}
```

#### 2. 闭包捕获
```go
func closureEscape() func() int {
    x := 100
    return func() int {
        return x  // x逃逸，被闭包捕获
    }
}
```

#### 3. 接口类型
```go
func interfaceEscape() interface{} {
    x := 200
    return x  // x逃逸，转换为interface{}
}
```

#### 4. Channel发送
```go
func channelEscape(ch chan *int) {
    x := 42
    ch <- &x  // x逃逸，指针发送到channel
}
```

#### 5. 大对象
```go
func largeObject() {
    // 超过32KB的对象通常在堆上分配
    s := make([]byte, 33*1024)
    _ = s
}
```

### 逃逸分析结果示例
```bash
$ go build -gcflags="-m" pointer.go
./pointer.go:10:6: moved to heap: x  # 逃逸到堆
./pointer.go:18:6: stack object: x   # 栈分配
```

### 性能影响
| 分配位置 | 优点 | 缺点 |
|----------|------|------|
| 栈 | 快速、无GC压力 | 生命周期受限 |
| 堆 | 生命周期灵活 | GC压力、分配较慢 |

## 性能考量

性能分为两个方面：**时间性能**和**内存性能**

### 时间性能测试

#### 小结构体 (<=32字节)
```go
type SmallStruct struct {
    a, b int  // 16字节
}

// 值传递和指针传递性能相近
func process(s SmallStruct) { }
func processPtr(s *SmallStruct) { }
```

**测试结果:**
```
小结构体(16字节) - 值传递:   160ms (1.60 ns/op)
小结构体(16字节) - 指针传递: 159ms (1.59 ns/op)
结论: 性能相近，优先值传递
```

#### 大结构体 (>32字节)
```go
type LargeStruct struct {
    data [10000]int  // 80KB
}

// 指针传递避免复制
func process(s *LargeStruct) { }
```

**测试结果:**
```
大结构体(80KB) - 值传递:   13ms (1.29 μs/op)
大结构体(80KB) - 指针传递: 0ms  (0.00 μs/op)
结论: 指针传递快很多 (避免复制80KB)
```

### 内存性能测试

#### 切片存储内存对比
存储1000个大结构体(80KB):
```
值类型切片:   80KB × 1000 = 78.1 MB
指针类型切片: 8字节 × 1000 = 8 KB
内存节省:     90%
```

#### 函数调用内存分配
```
返回值类型:   0 bytes (栈分配)
返回指针类型: 8 MB (堆分配，100次调用)
结论: 频繁创建时，值返回避免堆分配
```

#### 结构体内存对齐
```go
type Aligned1 struct {
    a bool    // 1字节 + 7字节填充
    b int64   // 8字节
    c bool    // 1字节 + 7字节填充
}  // 总共24字节

type Aligned2 struct {
    b int64   // 8字节
    a bool    // 1字节
    c bool    // 1字节 + 6字节填充
}  // 总共16字节
```

### 选择建议

| 场景 | 推荐方式 | 原因 |
|------|---------|------|
| 小结构体(<=32字节) | 值传递 | 复制快，避免指针解引用 |
| 大结构体(>32字节) | 指针传递 | 避免复制开销 |
| 存储大量元素 | 指针切片 | 节省内存 |
| 频繁创建返回 | 值返回 | 避免堆分配 |
| 共享数据 | 指针 | 避免重复存储 |

## 常见陷阱

### 1. nil指针访问
```go
var p *Person
fmt.Println(p.Name)  // panic: nil pointer dereference
```

### 2. 返回局部变量指针（不是陷阱！）
```go
func newPerson() *Person {
    p := Person{Name: "张三"}
    return &p  // Go会自动将p分配到堆上
}
```

### 3. 循环变量指针
```go
// 问题代码
var pointers []*int
for i := 0; i < 3; i++ {
    pointers = append(pointers, &i)  // 所有指针指向同一个地址
}

// 解决方案
for i := 0; i < 3; i++ {
    i := i  // 创建新变量
    pointers = append(pointers, &i)
}
```

## 最佳实践

1. **优先使用值类型**: 小结构体使用值传递
2. **一致性**: 同一类型的所有方法使用同一种接收者
3. **避免过度使用指针**: 不是所有地方都需要指针
4. **注意nil检查**: 解引用前检查nil
5. **慎用unsafe**: 只在必要时使用unsafe包

## 性能考量

### 值传递
- 优点：无堆分配，缓存友好
- 缺点：大结构体复制开销

### 指针传递
- 优点：避免复制，适合大结构体
- 缺点：可能导致堆分配，缓存不友好

### 建议
- 小结构体（< 3个字段）：值传递
- 大结构体：指针传递
- 需要修改：指针传递
