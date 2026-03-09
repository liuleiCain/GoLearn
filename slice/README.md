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
Go 1.18+的扩容规则：
- 小切片（cap < 256）: 翻倍
- 大切片（cap >= 256）: cap + cap/4 + 192

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

1. **预估容量**: 使用`make([]T, 0, cap)`预分配容量，避免频繁扩容
2. **避免内存泄漏**: 从大切片截取小切片时使用copy
3. **注意共享**: 多个切片可能共享底层数组，修改要小心
4. **使用copy**: 需要独立副本时使用copy而不是切片表达式
5. **检查边界**: 访问切片前检查长度，避免panic
