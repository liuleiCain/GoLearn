# Slice 切片模块

## 学习目标
深入理解Go语言切片的底层原理、扩容机制、内存模型和常见操作技巧。

## 核心概念

### 1. 切片底层结构
```go
type SliceHeader struct {
    Data uintptr  // 指向底层数组的指针
    Len  int      // 切片长度
    Cap  int      // 切片容量
}
```

### 2. 切片 vs 数组
| 特性 | 数组 | 切片 |
|------|------|------|
| 大小 | 固定 | 动态 |
| 值类型 | 是 | 引用类型 |
| 传递 | 复制整个数组 | 复制切片头 |
| 比较 | 可比较 | 不可比较(只能与nil比较) |

### 3. 扩容机制

#### 扩容触发条件
当使用 `append` 操作向切片添加元素时，如果切片的剩余容量不足以容纳新元素，就会触发扩容。

#### 扩容规则（Go 1.18+）

| 切片大小 | 扩容规则 | 计算方式 |
|---------|---------|--------|
| 小切片（cap < 256） | 容量翻倍 | newCap = oldCap * 2 |
| 大切片（cap >= 256） | 线性增长 | newCap = oldCap + oldCap/4 + 192 |

#### 详细计算示例

**小切片扩容示例：**
- 初始容量：10，添加1个元素 → 新容量：20（10*2）
- 初始容量：128，添加1个元素 → 新容量：256（128*2）
- 初始容量：255，添加1个元素 → 新容量：510（255*2）

**大切片扩容示例：**
- 初始容量：256，添加1个元素 → 新容量：256 + 64 + 192 = 512
- 初始容量：512，添加1个元素 → 新容量：512 + 128 + 192 = 832
- 初始容量：1000，添加1个元素 → 新容量：1000 + 250 + 192 = 1442

#### 扩容实现原理

1. **容量计算**：根据当前容量大小选择不同的扩容策略
2. **内存分配**：向内存管理器申请新的内存空间
3. **数据拷贝**：将原切片的元素拷贝到新的内存空间
4. **返回新切片**：返回指向新内存空间的切片

#### 扩容注意事项

1. **容量对齐**：最终分配的容量会进行内存对齐，确保是8的倍数
2. **容量上限**：扩容后的容量不会超过 `maxSliceCap`（通常为 2^64-1）
3. **性能影响**：频繁扩容会导致内存分配和数据拷贝，影响性能
4. **内存碎片**：扩容可能会产生内存碎片

#### 扩容源码解析（Go 1.23+）

以下是 Go 1.23+ runtime/slice.go 中的 `growslice` 和 `nextslicecap` 函数完整源码：

```go
// runtime/slice.go

// growslice allocates new backing store for a slice.
//
// arguments:
//
//      oldPtr = pointer to the slice's backing array
//      newLen = new length (= oldLen + num)
//      oldCap = original slice's capacity.
//         num = number of elements being added
//          et = element type
//
// return values:
//
//      newPtr = pointer to the new backing store
//      newLen = same value as the argument
//      newCap = capacity of the new backing store
//
// Requires that uint(newLen) > uint(oldCap).
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
        oldLen := newLen - num

        if newLen < 0 {
                panic(errorString("growslice: len out of range"))
        }

        if et.Size_ == 0 {
                return slice{unsafe.Pointer(&zerobase), newLen, newLen}
        }

        newcap := nextslicecap(newLen, oldCap)

        var overflow bool
        var lenmem, newlenmem, capmem uintptr
        noscan := !et.Pointers()
        switch {
        case et.Size_ == 1:
                lenmem = uintptr(oldLen)
                newlenmem = uintptr(newLen)
                capmem = roundupsize(uintptr(newcap), noscan)
                overflow = uintptr(newcap) > maxAlloc
                newcap = int(capmem)
        case et.Size_ == goarch.PtrSize:
                lenmem = uintptr(oldLen) * goarch.PtrSize
                newlenmem = uintptr(newLen) * goarch.PtrSize
                capmem = roundupsize(uintptr(newcap)*goarch.PtrSize, noscan)
                overflow = uintptr(newcap) > maxAlloc/goarch.PtrSize
                newcap = int(capmem / goarch.PtrSize)
        case isPowerOfTwo(et.Size_):
                var shift uintptr
                if goarch.PtrSize == 8 {
                        shift = uintptr(sys.TrailingZeros64(uint64(et.Size_))) & 63
                } else {
                        shift = uintptr(sys.TrailingZeros32(uint32(et.Size_))) & 31
                }
                lenmem = uintptr(oldLen) << shift
                newlenmem = uintptr(newLen) << shift
                capmem = roundupsize(uintptr(newcap)<<shift, noscan)
                overflow = uintptr(newcap) > (maxAlloc >> shift)
                newcap = int(capmem >> shift)
                capmem = uintptr(newcap) << shift
        default:
                lenmem = uintptr(oldLen) * et.Size_
                newlenmem = uintptr(newLen) * et.Size_
                capmem, overflow = math.MulUintptr(et.Size_, uintptr(newcap))
                capmem = roundupsize(capmem, noscan)
                newcap = int(capmem / et.Size_)
                capmem = uintptr(newcap) * et.Size_
        }

        if overflow || capmem > maxAlloc {
                panic(errorString("growslice: len out of range"))
        }

        var p unsafe.Pointer
        if !et.Pointers() {
                p = mallocgc(capmem, nil, false)
                memclrNoHeapPointers(add(p, newlenmem), capmem-newlenmem)
        } else {
                p = mallocgc(capmem, et, true)
                if lenmem > 0 && writeBarrier.enabled {
                        bulkBarrierPreWriteSrcOnly(uintptr(p), uintptr(oldPtr), lenmem-et.Size_+et.PtrBytes, et)
                }
        }
        memmove(p, oldPtr, lenmem)

        return slice{p, newLen, newcap}
}

// nextslicecap computes the next appropriate slice length.
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
                // Transition from growing 2x for small slices
                // to growing 1.25x for large slices. This formula
                // gives a smooth-ish transition between the two.
                newcap += (newcap + 3*threshold) >> 2

                // We need to check `newcap >= newLen` and whether `newcap` overflowed.
                // newLen is guaranteed to be larger than zero, hence
                // when newcap overflows then `uint(newcap) > uint(newLen)`.
                // This allows to check for both with the same comparison.
                if uint(newcap) >= uint(newLen) {
                        break
                }
        }

        // Set newcap to the requested cap when
        // the newcap calculation overflowed.
        if newcap <= 0 {
                return newLen
        }
        return newcap
}
```

#### 扩容算法说明

1. **快速路径**：如果新长度大于当前容量的两倍，直接使用新长度作为新容量
2. **小切片（< 256）**：容量翻倍
3. **大切片（≥ 256）**：使用平滑增长公式 `newcap += (newcap + 768) >> 2`
4. **内存对齐**：最终容量会根据元素大小进行内存对齐

#### 平滑增长公式详解

公式 `newcap += (newcap + 3*threshold) / 4` 可以展开为：
```
newcap = newcap + (newcap + 768) / 4
       = newcap * 1.25 + 192
```

这意味着：
- 基础增长率：1.25x（即 25%）
- 额外固定增量：192

#### 扩容示例

| 初始容量 | 添加元素后 | 新容量 | 增长率 | 说明 |
|---------|-----------|-------|--------|------|
| 0 | 1 | 1 | - | 空切片特殊处理 |
| 1 | 2 | 2 | 2.00x | 翻倍 |
| 2 | 3 | 4 | 2.00x | 翻倍 |
| 4 | 5 | 8 | 2.00x | 翻倍 |
| 8 | 9 | 16 | 2.00x | 翻倍 |
| 16 | 17 | 32 | 2.00x | 翻倍 |
| 32 | 33 | 64 | 2.00x | 翻倍 |
| 64 | 65 | 128 | 2.00x | 翻倍 |
| 128 | 129 | 256 | 2.00x | 翻倍 |
| 256 | 257 | 512 | 2.00x | 刚好在阈值，还是翻倍 |
| 512 | 513 | 848 | 1.66x | 开始平滑增长 |
| 848 | 849 | 1280 | 1.51x | 增长率逐渐趋近1.25x |
| 1024 | 1025 | 1696 | 1.66x | 平滑增长 |

#### Go 1.17 vs Go 1.18+ 对比

| 旧容量 | Go 1.17 | Go 1.18+ | 差异 |
|--------|---------|----------|------|
| 512 | 1024 | 848 | -17.2% |
| 1024 | 1280 | 1696 | +32.5% |
| 2048 | 2560 | 3408 | +33.1% |

Go 1.18+ 的改进：
- 阈值从 1024 降低到 256，更早进入平滑增长
- 中小切片内存效率提升 17-40%
- 增长曲线更平滑，避免突变

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| SliceStructure | 切片底层结构 | 初级 |
| SliceCreate | 切片创建方式 | 初级 |
| AppendDemo | append操作和扩容 | 初级 |
| GrowthRule | 扩容规则详解 | 中级 |
| CopyDemo | copy函数使用 | 初级 |
| SliceExpression | 切片表达式 | 初级 |
| SliceTricks | 切片常用技巧 | 中级 |
| SliceMemory | 内存泄漏问题 | 高级 |
| SlicePointer | 切片指针分析 | 高级 |

## 切片表达式

### 简单切片表达式
```go
arr := [5]int{1, 2, 3, 4, 5}
arr[low:high]   // low <= index < high
arr[:3]         // 等价于 arr[0:3]
arr[2:]         // 等价于 arr[2:len(arr)]
arr[:]          // 等价于 arr[0:len(arr)]
```

### 完整切片表达式
```go
arr[low:high:max]  // 限制容量为 max-low
arr[1:3:4]         // len=2, cap=3
```

## 切片技巧速查

### 删除元素
```go
// 删除索引i
s = append(s[:i], s[i+1:]...)

// 删除索引i到j
s = append(s[:i], s[j:]...)
```

### 插入元素
```go
// 在索引i插入x
s = append(s[:i], append([]int{x}, s[i:]...)...)
```

### 清空切片
```go
// 方法1: 设为nil
s = nil

// 方法2: 保留内存
s = s[:0]
```

### 复制切片
```go
// 浅拷贝
newSlice := make([]int, len(s))
copy(newSlice, s)

// 或使用append
newSlice := append([]int(nil), s...)
```

## 常见陷阱

### 1. 内存泄漏
```go
// 问题: 小切片引用大数组
var largeSlice = make([]int, 1000000)
var smallSlice = largeSlice[:10]  // 仍引用整个大数组

// 解决: 使用copy
var smallSlice = make([]int, 10)
copy(smallSlice, largeSlice[:10])
```

### 2. 共享底层数组
```go
arr := [5]int{1, 2, 3, 4, 5}
s1 := arr[:3]
s2 := arr[3:]
s1 = append(s1, 100)  // 可能修改s2的数据!
```

### 3. append返回新切片
```go
// 错误: append可能返回新的切片
append(s, 1)  // s没有改变

// 正确
s = append(s, 1)
```

## 最佳实践

### 1. 预估容量：预分配避免扩容

**错误做法：**
```go
// 每次append都可能触发扩容
s := make([]int, 0)
for i := 0; i < 10000; i++ {
    s = append(s, i)  // 频繁扩容
}
```

**正确做法：**
```go
// 预分配容量，一次分配
s := make([]int, 0, 10000)
for i := 0; i < 10000; i++ {
    s = append(s, i)  // 无需扩容
}
```

**性能对比：**
```go
// 不预分配
start := time.Now()
s1 := make([]int, 0)
for i := 0; i < 1000000; i++ {
    s1 = append(s1, i)
}
fmt.Printf("不预分配: %v\n", time.Since(start))

// 预分配
start = time.Now()
s2 := make([]int, 0, 1000000)
for i := 0; i < 1000000; i++ {
    s2 = append(s2, i)
}
fmt.Printf("预分配: %v\n", time.Since(start))
```

### 2. 避免内存泄漏：使用copy代替切片截取

**错误做法：**
```go
largeSlice := make([]int, 1000000)
smallSlice := largeSlice[:10]  // 仍引用整个大数组，GC无法回收largeSlice
```

**正确做法：**
```go
largeSlice := make([]int, 1000000)
smallSlice := make([]int, 10)
copy(smallSlice, largeSlice[:10])  // 独立副本，largeSlice可被回收
largeSlice = nil  // 显式释放
```

**图示：**
```
错误做法：
┌────────────────────────────────────┐
│ largeSlice (1000000元素)           │
├──────────────┬─────────────────────┤
│ smallSlice引用│  无法回收的部分      │
└──────────────┴─────────────────────┘

正确做法：
┌────────────────────────────────────┐
│ largeSlice (1000000元素)  ──[GC]──>│ 释放
└────────────────────────────────────┘
┌──────────────┐
│ smallSlice   │  独立副本
│ (10元素)     │
└──────────────┘
```

### 3. 注意共享：多个切片共享底层数组

**错误做法：**
```go
arr := [5]int{1, 2, 3, 4, 5}
s1 := arr[:3]  // [1,2,3]
s2 := arr[3:]   // [4,5]

s1 = append(s1, 100)  // 修改了s2的数据！
fmt.Println(s2[0])    // 输出: 100
```

**正确做法：**
```go
arr := [5]int{1, 2, 3, 4, 5}

// 方案1: 使用copy创建独立副本
s1 := make([]int, 3)
copy(s1, arr[:3])
s1 = append(s1, 100)
fmt.Println(arr[:3])  // [1,2,3] 不受影响

// 方案2: 使用完整切片表达式限制容量
s1 = arr[:3:3]  // cap=3，不与arr共享
s1 = append(s1, 100)  // 触发扩容，不影响arr
fmt.Println(arr)  // [1,2,3,4,5] 不受影响
```

### 4. 使用copy：创建独立副本

**错误做法：**
```go
original := []int{1, 2, 3, 4, 5}
backup := original  // 只是复制了切片头，共享底层数组
backup[0] = 100
fmt.Println(original[0])  // 输出: 100，也被改了
```

**正确做法：**
```go
original := []int{1, 2, 3, 4, 5}

// 方案1: 使用copy
backup := make([]int, len(original))
copy(backup, original)
backup[0] = 100
fmt.Println(original[0])  // 输出: 1，未受影响

// 方案2: 使用append
backup = append([]int(nil), original...)

// 方案3: 完整复制（保留容量）
backup = make([]int, len(original), cap(original))
copy(backup, original)
```

### 5. 检查边界：避免越界访问

**错误做法：**
```go
s := []int{1, 2, 3}
fmt.Println(s[5])  // panic: index out of range
```

**正确做法：**
```go
s := []int{1, 2, 3}

// 方案1: 先检查长度
if len(s) > 5 {
    fmt.Println(s[5])
}

// 方案2: 使用循环遍历
for i, v := range s {
    fmt.Printf("s[%d] = %d\n", i, v)
}

// 方案3: 使用slices包（Go 1.21+）
if ok := slices.Contains(s, target); ok {
    // ...
}
```

### 6. 善用切片表达式：控制容量

**完整切片表达式可以限制容量：**
```go
arr := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

s1 := arr[2:5]      // len=3, cap=8（指向arr[2]）
s2 := arr[2:5:5]    // len=3, cap=3（容量限制为3）

s1 = append(s1, 100)  // 不影响原数组
s2 = append(s2, 100)  // 触发扩容，不影响原数组
```

**使用场景：**
- 防止意外修改共享底层数组
- 明确切片的容量边界
- 提高代码可读性
