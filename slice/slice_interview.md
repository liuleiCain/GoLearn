# Slice 切片面试题集

本模块整理了 Go 语言 Slice 切片的常见面试题，包含问题和详细解答。

## 目录

- [基础概念面试题](#基础概念面试题)
- [扩容机制面试题](#扩容机制面试题)
- [底层原理面试题](#底层原理面试题)
- [操作技巧面试题](#操作技巧面试题)
- [常见陷阱面试题](#常见陷阱面试题)
- [源码解析面试题](#源码解析面试题)

## 基础概念面试题

### 1. Slice和数组的区别？

**答案：**

| 特性 | 数组 | 切片 |
|------|------|------|
| 大小 | 固定 | 动态 |
| 值类型 | 是（复制整个数组） | 引用类型（复制切片头） |
| 传递方式 | 值传递 | 引用传递 |
| 比较 | 可比较（逐元素） | 只能与nil比较 |
| 容量 | 无概念 | 有容量概念 |
| 声明方式 | [3]int | []int |

**示例：**
```go
// 数组 - 值类型
arr1 := [3]int{1, 2, 3}
arr2 := arr1          // 复制整个数组
arr2[0] = 100
fmt.Println(arr1[0])  // 输出: 1，未改变

// 切片 - 引用类型
s1 := []int{1, 2, 3}
s2 := s1              // 复制切片头（ptr, len, cap）
s2[0] = 100
fmt.Println(s1[0])    // 输出: 100，已改变
```

### 2. Slice底层结构是什么？

**答案：**

Slice 底层是一个结构体，包含三个字段：

```go
type slice struct {
    array unsafe.Pointer  // 指向底层数组的指针
    len   int            // 切片长度（元素个数）
    cap   int            // 切片容量（底层数组大小）
}
```

**图示：**
```
切片 s := arr[0:3]  // 从数组截取，len=3, cap=5
       ↓
┌─────────────┐
│ SliceHeader │
├─────────────┤
│ array:  ptr │ ──────────────────┐
│ len:   3    │                  │
│ cap:   5    │                  ↓
└─────────────┘          ┌───┬───┬───┬───┬───┐
                        │ 1 │ 2 │ 3 │ 4 │ 5 │
                        └───┴───┴───┴───┴───┘
                        ↑           ↑
                        │           │
                      s[0:3]    s[3:5]
                      len=3     cap=2
```

### 3. Slice的nil和空切片有什么区别？

**答案：**

| 类型 | 字面量 | 长度 | 容量 | 底层数组 |
|------|--------|------|------|----------|
| nil切片 | var s []int | 0 | 0 | nil |
| 空切片 | s := []int{} | 0 | 0 | 非nil，但长度为0 |
| 初始化切片 | s := make([]int, 0) | 0 | 0 | 非nil，空数组 |

**示例：**
```go
var nilSlice []int
emptySlice := []int{}
makeSlice := make([]int, 0)

fmt.Println(nilSlice == nil)      // true
fmt.Println(emptySlice == nil)    // false
fmt.Println(makeSlice == nil)     // false

fmt.Println(len(nilSlice), cap(nilSlice))   // 0, 0
fmt.Println(len(emptySlice), cap(emptySlice)) // 0, 0
fmt.Println(len(makeSlice), cap(makeSlice))   // 0, 0
```

### 4. Slice是如何实现动态扩容的？

**答案：**

Slice 本身是一个固定长度的结构体，真正的动态扩容是通过：
1. 当 `append` 操作超出容量时
2. 调用 `growslice` 函数分配新的底层数组
3. 将原数组的元素拷贝到新数组
4. 返回新的切片（新的array指针、len、cap）

**示例：**
```go
s := make([]int, 0, 2)  // len=0, cap=2
fmt.Printf("初始: len=%d, cap=%d\n", len(s), cap(s))

s = append(s, 1, 2)     // len=2, cap=2
fmt.Printf("添加2个: len=%d, cap=%d\n", len(s), cap(s))

s = append(s, 3)        // 触发扩容，cap变为4
fmt.Printf("添加3个: len=%d, cap=%d\n", len(s), cap(s))
```

## 扩容机制面试题

### 1. Slice的扩容规则是什么？

**答案：**（Go 1.18+）

| 条件 | 扩容规则 | 公式 |
|------|----------|------|
| 新容量 > 2倍旧容量 | 直接使用新容量 | newCap = newLen |
| 旧容量 < 256 | 翻倍增长 | newCap = oldCap * 2 |
| 旧容量 >= 256 | 平滑增长 | newCap = oldCap + (oldCap+768)/4 |

**源码（Go 1.23+）：**
```go
func nextslicecap(newLen, oldCap int) int {
    newcap := oldCap
    doublecap := newcap + newcap
    if newLen > doublecap {
        return newLen
    }

    const threshold = 256
    if oldCap < threshold {
        return doublecap
    }
    for {
        newcap += (newcap + 3*threshold) >> 2
        if uint(newcap) >= uint(newLen) {
            break
        }
    }

    if newcap <= 0 {
        return newLen
    }
    return newcap
}
```

### 2. Slice扩容后底层数组是如何变化的？

**答案：**

1. 分配一块新的内存空间（容量更大的数组）
2. 将原数组的所有元素拷贝到新数组
3. 新切片指向新的底层数组
4. 旧的底层数组如果没有引用会被GC回收

**示例：**
```go
s := []int{1, 2, 3}
fmt.Printf("s: len=%d, cap=%d, ptr=%p\n", len(s), cap(s), &s[0])

s = append(s, 4)
fmt.Printf("s: len=%d, cap=%d, ptr=%p\n", len(s), cap(s), &s[0])
// 输出: 容量翻倍，指针地址改变
```

### 3. 为什么Slice扩容要设计成翻倍+平滑增长？

**答案：**

1. **性能考量**：减少扩容次数，避免频繁内存分配
2. **内存考量**：避免一次性分配过大内存造成浪费
3. **平衡策略**：
   - 小切片（<256）：翻倍策略，因为容量小，即使浪费也不多
   - 大切片（>=256）：平滑增长（约1.25倍），避免内存浪费过大
4. **公式推导**：
   ```
   newcap = oldcap + (oldcap + 768) / 4
          ≈ oldcap * 1.25 + 192
   ```

### 4. Slice容量的内存对齐是什么？

**答案：**

内存对齐是CPU访问内存的优化策略，要求数据的地址是其大小的整数倍。

**为什么需要内存对齐？**

| 对齐 | 说明 |
|------|------|
| 性能 | CPU一次读取多个字节，未对齐需要多次读取 |
| 硬件 | 某些硬件平台只支持对齐访问 |
| 编译器 | 编译器自动插入padding |

**Go中的内存对齐规则：**

```go
// 编译器会自动插入padding
type Person struct {
    name string  // 8字节
    age  int     // 8字节
    id   byte    // 1字节 → 编译器会插入7字节padding
}
// 实际大小: 24字节，不是17字节
```

**Slice扩容时的内存对齐处理：**

```go
// 根据元素大小进行不同的对齐处理
switch {
case et.Size_ == 1:        // 1字节元素（如byte）
    capmem = roundupsize(uintptr(newcap), noscan)
    // 无需特殊对齐

case et.Size_ == goarch.PtrSize:  // 8字节（64位平台）
    capmem = roundupsize(uintptr(newcap)*goarch.PtrSize, noscan)
    // 按指针大小对齐

case isPowerOfTwo(et.Size_):       // 2的幂次大小（如2,4,8,16）
    capmem = roundupsize(uintptr(newcap)<<shift, noscan)
    // 使用位运算计算

default:                          // 其他大小（如结构体）
    capmem, overflow = math.MulUintptr(et.Size_, uintptr(newcap))
    capmem = roundupsize(capmem, noscan)
    // 按元素大小*数量计算后对齐
}
```

**实际影响示例：**

```go
// 假设初始容量256，按平滑增长公式计算：
// newcap = 256 + (256 + 768) / 4 = 512

// 但内存对齐后可能变成：
// 元素大小=8字节(int64) → 按8字节对齐 → 512 * 8 = 4096
// roundupsize(4096) → 可能对齐到 4096 或更大

// 示例
s := make([]int, 0)
for i := 0; i < 600; i++ {
    s = append(s, i)
    if i == 255 || i == 256 || i == 512 {
        fmt.Printf("len=%d, cap=%d\n", len(s), cap(s))
    }
}
// 可能输出：
// len=256, cap=256
// len=257, cap=512
// len=513, cap=848  ← 实际是按公式计算后再对齐
```

**图示：内存对齐过程**

```
原始需求：容量=513，元素大小=8字节
         ↓
计算内存：513 * 8 = 4104 字节
         ↓
roundupsize 对齐：4104 → 对齐到 4160（mspanSize）或其他合适大小
         ↓
最终容量：4160 / 8 = 520 个元素
```

**关键点：**
1. `roundupsize` 会将计算出的内存大小向上取整到合适的尺寸
2. 实际分配的容量可能比公式计算的稍大
3. 这是Go运行时内存管理器的行为，程序员无法控制
4. 内存对齐是为了提高CPU访问效率

## 底层原理面试题

### 1. Slice是如何在函数间传递的？

**答案：**

Slice 作为函数参数时，是**值传递**：
- 复制的是 SliceHeader 结构体（3个字段）
- 底层数组是共享的
- 修改元素会影响原切片
- 修改 len 或 cap 不会影响原切片

**示例：**
```go
func modifySlice(s []int) {
    s[0] = 100          // 修改底层数组，原切片受影响
    s = append(s, 4)    // 创建新切片，不影响原切片
    s[1] = 200          // 修改新切片，不影响原切片
}

func main() {
    original := []int{1, 2, 3}
    modifySlice(original)
    fmt.Println(original)  // 输出: [100 2 3]
}
```

### 2. Slice的append操作会发生什么？

**答案：**

```go
s := make([]int, 0, 3)
s = append(s, 1, 2, 3)  // len=3, cap=3

// 情况1: 超出容量，触发扩容
s = append(s, 4)         // 扩容，cap变为6

// 情况2: 未超出新容量，直接添加
s = append(s, 5)         // len=5, cap=6，直接添加
```

**关键点：**
1. `append` 返回新的切片
2. 如果容量足够，在原底层数组添加
3. 如果容量不够，扩容后添加

### 3. Slice共享底层数组的情况？

**答案：**

多个切片可能共享同一个底层数组：

```go
arr := [5]int{1, 2, 3, 4, 5}
s1 := arr[0:3]  // [1,2,3], ptr指向arr[0], len=3, cap=5
s2 := arr[2:5]  // [3,4,5], ptr指向arr[2], len=3, cap=3

// 修改s1会影响s2
s1[2] = 100
fmt.Println(s2[0])  // 输出: 100
```

**图示：**
```
arr:  [1] [2] [100] [4] [5]
        ↑           ↑
       s1[0:3]    s2[2:5]
```

### 4. Slice和数组的内存布局有什么区别？

**答案：**

```go
// 数组 - 连续内存块
arr := [4]int{1, 2, 3, 4}
// 内存布局: [1][2][3][4] (16字节，int=4字节)

// 切片 - 切片头 + 底层数组
s := []int{1, 2, 3, 4}
// 内存布局:
//   SliceHeader: [ptr][len][cap] (24字节，指针=8, int=8)
//   底层数组:     [1][2][3][4] (32字节)
```

## 操作技巧面试题

### 1. 如何正确删除Slice中的元素？

**答案：**

```go
// 删除索引i的元素
s = append(s[:i], s[i+1:]...)

// 删除区间[i:j]的元素
s = append(s[:i], s[j:]...)

// 删除第一个元素
s = s[1:]

// 删除最后一个元素
s = s[:len(s)-1]

// 清空切片
s = s[:0]  // 保留内存
s = nil    // 释放内存
```

### 2. 如何合并两个Slice？

**答案：**

```go
// 方法1: append
s1 := []int{1, 2}
s2 := []int{3, 4}
s3 := append(s1, s2...)

// 方法2: 使用copy
s3 := make([]int, len(s1)+len(s2))
copy(s3, s1)
copy(s3[len(s1):], s2)

// 方法3: slices.Concat (Go 1.22+)
s3 = slices.Concat(s1, s2)
```

### 3. 如何深拷贝一个Slice？

**答案：**

**重要：copy 的目标切片必须有足够的长度！**

```go
// copy 只复制 min(len(dst), len(src)) 个元素
var backup []int                    // len=0, cap=0, nil切片
original := []int{1, 2, 3, 4, 5}    // len=5

n := copy(backup, original)
fmt.Println(n)          // 输出: 0 ← 复制了0个元素！
fmt.Println(backup)      // 输出: [] ← 空的
```

| 目标切片 | len | 源切片 len | 实际复制 |
|----------|-----|-----------|---------|
| `var backup []int` | 0 | 5 | 0 |
| `make([]int, 5)` | 5 | 5 | 5 |
| `make([]int, 3)` | 3 | 5 | 3 |

**正确用法：**

```go
// 方法1: 使用copy（必须先make分配足够长度）
dst := make([]int, len(src))
copy(dst, src)

// 方法2: 使用append（无需预先make，append会自动扩容）
dst := append([]int(nil), src...)

// 方法3: 完整复制（带容量）
dst := make([]int, len(src), cap(src))
copy(dst, src)

// 对于复杂类型，需要手动拷贝每个字段
type Person struct {
    Name string
    Age  int
}
dst := make([]Person, len(src))
for i := range src {
    dst[i] = Person{
        Name: src[i].Name,  // string是值类型，直接拷贝
        Age:  src[i].Age,   // int是值类型，直接拷贝
    }
}
```

### 4. Slice的最佳实践有哪些？

**答案：**

1. **预分配容量**
   ```go
   // 不预分配
   s := make([]int, 0)
   for i := 0; i < 1000; i++ {
       s = append(s, i)  // 多次扩容
   }

   // 预分配
   s := make([]int, 0, 1000)
   for i := 0; i < 1000; i++ {
       s = append(s, i)  // 一次分配
   }
   ```

2. **避免内存泄漏**
   ```go
   large := make([]int, 1000000)
   small := large[:10]

   // 解决：使用copy
   small := make([]int, 10)
   copy(small, large[:10])
   ```

3. **注意append的返回值**
   ```go
   s := make([]int, 0, 3)
   append(s, 1, 2, 3)  // 错误：s不变
   s = append(s, 1, 2, 3)  // 正确：接收返回值
   ```

## 常见陷阱面试题

### 1. 循环中的Slice append有什么陷阱？

**答案：**

**Go 1.22之前的问题：**
```go
// 错误示例
s := make([][]int, 3)
for i := range s {
    s[i] = append(s[i], i)  // 每次append可能触发扩容
}
// 正确做法
s := make([][]int, 3)
for i := range s {
    s[i] = make([]int, 0, 1)  // 预分配
    s[i] = append(s[i], i)
}
```

**Go 1.22+已修复：** 循环变量现在会在每次迭代中重新创建。

### 2. Slice截取后会有什么潜在问题？

**答案：**

```go
large := make([]int, 1000000)
small := large[:10]  // small引用large的底层数组

// large不再使用，但small仍引用，large无法被GC
large = nil  // large被置为nil，但small仍引用底层数组
fmt.Println(small[0])  // 正常工作
```

**解决：**
```go
small := make([]int, 10)
copy(small, large[:10])  // 独立副本
large = nil
```

### 3. 为什么不要对nil Slice进行操作？

**答案：**

| 操作 | nil Slice | 空 Slice |
|------|-----------|----------|
| append | 正常工作 | 正常工作 |
| len | 0 | 0 |
| cap | 0 | 0 |
| range | 不panic | 不panic |
| 索引访问 | panic越界 | Panic越界 |

**但要注意：**
```go
var s []int
s = append(s, 1)  // 正常工作，s 变为 [1]，len=1, cap>=1
s[0] = 100        // 安全，append 后 cap 必定 >= 1

s2 := []int{}
s2 = append(s2, 1)  // 正常工作
s2[0] = 100         // 安全
```

### 4. Slice比较有什么坑？

**答案：**

```go
// 不能直接比较切片
s1 := []int{1, 2, 3}
s2 := []int{1, 2, 3}
// s1 == s2  // 编译错误

// 使用reflect.DeepEqual
fmt.Println(reflect.DeepEqual(s1, s2))  // true

// 使用slices.Equal (Go 1.21+)
fmt.Println(slices.Equal(s1, s2))  // true

// bytes包用于byte切片
b1 := []byte{1, 2, 3}
b2 := []byte{1, 2, 3}
fmt.Println(bytes.Equal(b1, b2))  // true
```

## 源码解析面试题

### 1. Slice的growslice函数主要逻辑是什么？

**答案：**

```go
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
    oldLen := newLen - num

    if newLen < 0 {
        panic(errorString("growslice: len out of range"))
    }

    if et.Size_ == 0 {
        return slice{unsafe.Pointer(&zerobase), newLen, newLen}
    }

    // 1. 计算新容量
    newcap := nextslicecap(newLen, oldCap)

    // 2. 内存对齐处理
    // ... 根据元素大小进行不同的对齐计算

    // 3. 分配新内存
    p := mallocgc(capmem, et, true)

    // 4. 拷贝数据
    memmove(p, oldPtr, lenmem)

    return slice{p, newLen, newcap}
}
```

### 2. nextslicecap函数为什么这样设计？

**答案：**

**设计考量：**

1. **快速路径**：新长度超过2倍容量时，直接使用新长度
   ```go
   if newLen > doublecap {
       return newLen
   }
   ```

2. **小切片翻倍**：容量<256时，翻倍增长
   ```go
   if oldCap < threshold {
       return doublecap
   }
   ```

3. **大切片平滑增长**：避免内存浪费
   ```go
   for {
       newcap += (newcap + 3*threshold) >> 2
       if uint(newcap) >= uint(newLen) {
           break
       }
   }
   ```

4. **溢出处理**：防止整数溢出
   ```go
   if newcap <= 0 {
       return newLen
   }
   ```

### 3. Slice的内存分配策略是什么？

**答案：**

1. **栈上分配**：小对象可能栈上分配
2. **堆上分配**：大对象或逃逸对象堆上分配
3. **内存复用**：释放后可能被其他对象复用
4. **对齐分配**：`roundupsize` 确保内存对齐

### 4. Slice和Make的区别是什么？

**答案：**

| 特性 | make([]T, len) | make([]T, len, cap) | make([]T, 0, cap) |
|------|----------------|---------------------|-------------------|
| 长度 | len | len | 0 |
| 容量 | len | cap | cap |
| 初始化 | 零值 | 零值 | 无 |
| 适用场景 | 已知长度 | 预分配 | 动态添加 |

```go
s1 := make([]int, 3)        // [0, 0, 0], len=3, cap=3
s2 := make([]int, 3, 10)   // [0, 0, 0], len=3, cap=10
s3 := make([]int, 0, 10)   // [], len=0, cap=10
```

## 相关资源

- [Go语言 Slice 官方文档](https://golang.org/doc/article/slices)
- [Go语言源码 runtime/slice.go](https://github.com/golang/go/blob/master/src/runtime/slice.go)
- [Effective Go - Slices](https://golang.org/doc/effective_go#slices)
