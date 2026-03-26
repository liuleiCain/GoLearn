# Go 底层设计与实现

本模块深入讲解 Go 语言核心数据结构的底层实现原理，涵盖面试高频考点。

## 目录

- [Slice 底层实现](#slice-底层实现)
- [Map 底层实现](#map-底层实现)
- [Channel 底层实现](#channel-底层实现)
- [Interface 底层实现](#interface-底层实现)
- [String 底层实现](#string-底层实现)
- [运行测试](#运行测试)

---

## Slice 底层实现

### 核心结构

```go
type SliceHeader struct {
    Data uintptr  // 底层数组指针
    Len  int      // 长度
    Cap  int      // 容量
}
```

### 面试高频问题

#### 1. 切片扩容机制（Go 1.18+）

```
扩容规则：
1. 新容量 > 2倍旧容量 -> 使用新容量
2. 旧容量 < 256 -> 新容量 = 2倍旧容量
3. 旧容量 >= 256 -> 新容量 = 旧容量 + (旧容量+3*256)/4
4. 最后进行内存对齐
```

#### 2. 切片共享底层数组

```go
original := make([]int, 5, 10)
s1 := original[1:3]  // 共享底层数组
s2 := original[2:5]  // 共享底层数组
// 修改s1会影响original和s2
```

#### 3. nil切片 vs 空切片

| 类型 | Data指针 | Len | Cap | == nil | JSON |
|------|----------|-----|-----|--------|------|
| nil切片 | 0 | 0 | 0 | true | null |
| 空切片 | 非nil | 0 | 0 | false | [] |

#### 4. 切片内存泄漏

```go
// 问题：子切片引用大数组
largeSlice := make([]int, 1000000)
smallSlice := largeSlice[:10]  // 整个大数组无法GC

// 解决：复制一份
correctSlice := make([]int, 10)
copy(correctSlice, largeSlice[:10])
largeSlice = nil  // 释放引用
```

---

## Map 底层实现

### 核心结构

```go
type hmap struct {
    count     int           // 元素个数
    B         uint8         // bucket数量 = 2^B
    hash0     uint32        // hash种子
    buckets   unsafe.Pointer // bucket数组
    oldbuckets unsafe.Pointer // 扩容时的旧bucket
}

type bmap struct {
    tophash [8]uint8  // hash高8位
    // 后面是8个key + 8个value + overflow指针
}
```

### 面试高频问题

#### 1. Map查找过程

```
1. 计算key的hash值
2. 取低B位定位bucket
3. 取高8位在tophash中查找
4. 匹配则比较完整key
5. 不匹配则查找overflow bucket
```

#### 2. Map扩容条件

```
1. 负载因子 > 6.5 -> 翻倍扩容
2. overflow bucket过多 -> 等量扩容（整理）
```

#### 3. Map并发问题

```go
// map不是并发安全的！
// 解决方案：
// 1. sync.Mutex
// 2. sync.RWMutex
// 3. sync.Map
```

#### 4. Map Key类型限制

| 可用 | 不可用 |
|------|--------|
| int, float, string, bool | slice |
| 指针 *T | map |
| 数组 [N]T | func |
| 结构体（字段可比较） | 包含不可比较字段的结构体 |

---

## Channel 底层实现

### 核心结构

```go
type hchan struct {
    qcount   uint           // 队列元素数
    dataqsiz uint           // 缓冲区大小
    buf      unsafe.Pointer // 循环队列
    sendx    uint           // 发送索引
    recvx    uint           // 接收索引
    recvq    waitq          // 接收等待队列
    sendq    waitq          // 发送等待队列
    lock     mutex          // 互斥锁
}
```

### 面试高频问题

#### 1. Channel发送流程

```
1. 获取锁
2. 检查是否已关闭 -> panic
3. 检查recvq有等待 -> 直接发送
4. 检查缓冲区有空间 -> 放入缓冲区
5. 否则 -> 阻塞，加入sendq
```

#### 2. Channel接收流程

```
1. 获取锁
2. 检查sendq有等待 -> 直接接收
3. 检查缓冲区有数据 -> 从缓冲区接收
4. 检查已关闭 -> 返回零值
5. 否则 -> 阻塞，加入recvq
```

#### 3. 关闭Channel的行为

| 操作 | 结果 |
|------|------|
| 发送 | panic |
| 接收 | 返回零值 + false |
| 再次关闭 | panic |

#### 4. Nil Channel行为

| 操作 | 结果 |
|------|------|
| 发送 | 永久阻塞 |
| 接收 | 永久阻塞 |
| 关闭 | panic |

---

## Interface 底层实现

### 核心结构

```go
// 空接口
type eface struct {
    _type *_type         // 类型信息
    data  unsafe.Pointer // 数据指针
}

// 非空接口
type iface struct {
    tab  *itab           // 接口表
    data unsafe.Pointer  // 数据指针
}

type itab struct {
    inter *interfacetype // 接口类型
    _type *_type         // 具体类型
    fun   [1]uintptr     // 方法表
}
```

### 面试高频问题

#### 1. nil接口 vs nil值接口

```go
var nilInterface interface{}  // eface{_type: nil, data: nil}
var p *int = nil
var nilValueInterface interface{} = p  // eface{_type: *int, data: nil}

nilInterface == nil        // true
nilValueInterface == nil   // false!
```

#### 2. 类型断言实现

```
1. 检查eface._type是否等于目标类型
2. 相等则返回data指针
3. 不等则返回零值和false
```

#### 3. 接口装箱与逃逸

```go
x := 42
var i interface{} = x  // x逃逸到堆
// 原因：接口需要存储值的指针
```

#### 4. 接口比较

```go
// 可比较：类型相同且可比较
i1 == i2  // 调用_type.equal比较data

// 不可比较的类型会导致panic
var s1 interface{} = []int{1, 2}
var s2 interface{} = []int{1, 2}
// s1 == s2  // panic!
```

---

## String 底层实现

### 核心结构

```go
type StringHeader struct {
    Data uintptr  // 字节数组指针
    Len  int      // 字节长度
}
```

### 面试高频问题

#### 1. String不可变性

```go
s := "hello"
// s[0] = 'H'  // 编译错误

// 修改方法
b := []byte(s)
b[0] = 'H'
s = string(b)
```

#### 2. String与[]byte转换

```go
// 标准转换（复制）
b := []byte(s)
s := string(b)

// 零拷贝转换（unsafe）
func bytesToString(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}
```

#### 3. UTF-8编码

```
编码规则：
0xxxxxxx: 1字节 (ASCII)
110xxxxx 10xxxxxx: 2字节
1110xxxx 10xxxxxx 10xxxxxx: 3字节
11110xxx 10xxxxxx 10xxxxxx 10xxxxxx: 4字节
```

#### 4. 子字符串内存泄漏

```go
s := make([]byte, 1<<20)  // 1MB
sub := string(s[:10])      // 整个1MB无法GC

// 解决：复制一份
sub := string([]byte(s[:10]))
```

---

## 运行测试

```bash
# 运行所有测试
go test -v ./internals

# 运行特定测试
go test -v ./internals -run TestDemoSlice

# 查看逃逸分析
go build -gcflags="-m" ./internals
```

## 使用示例

```go
package main

import "github.com/yourusername/GoLearn/internals"

func main() {
    // Slice底层
    internals.DemoSliceStructure()
    internals.DemoSliceGrowth()
    
    // Map底层
    internals.DemoMapStructure()
    internals.DemoMapGrow()
    
    // Channel底层
    internals.DemoChannelStructure()
    internals.DemoChannelSend()
    
    // Interface底层
    internals.DemoInterfaceStructure()
    internals.DemoInterfaceNil()
    
    // String底层
    internals.DemoStringStructure()
    internals.DemoStringUTF8()
}
```

## 面试要点总结

### Slice
- 结构：Data + Len + Cap
- 扩容：2倍增长（<256），1.25倍增长（>=256）
- 共享底层数组，注意内存泄漏

### Map
- 结构：hmap + bmap
- 查找：hash -> bucket -> tophash -> key
- 扩容：负载因子>6.5 或 overflow过多
- 非并发安全

### Channel
- 结构：hchan + 循环队列 + 等待队列
- 发送/接收：锁 + 条件检查 + 阻塞
- 关闭：唤醒所有等待者

### Interface
- 结构：eface（空接口）/ iface（非空接口）
- nil接口 vs nil值接口
- 装箱导致逃逸

### String
- 结构：Data + Len
- 不可变
- UTF-8编码
- 子字符串可能内存泄漏

## 相关资源

- [Go语言设计与实现](https://draveness.me/golang/)
- [Go源码](https://github.com/golang/go)
- [Go内存模型](https://golang.org/ref/mem)
