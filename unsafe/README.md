# Unsafe 底层内存操作模块

## 学习目标
理解Go语言unsafe包的使用场景和潜在风险，掌握底层内存操作技术。

## 核心概念

### 1. unsafe.Pointer
可以转换任何类型的指针，是类型安全和不安全操作之间的桥梁。

```go
unsafe.Pointer(&x)  // 任何指针都可以转为unsafe.Pointer
(*int)(ptr)           // unsafe.Pointer可以转为任何指针
```

### 2. 三个基本函数
```go
unsafe.Sizeof(x)      // 计算大小（字节）
unsafe.Alignof(x)     // 计算对齐
unsafe.Offsetof(x.f)  // 计算字段偏移
```

### 3. 警告
unsafe包绕过了Go的类型安全检查，使用需谨慎！

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| BasicPointer | 基础指针操作 | 初级 |
| SizeOf | 计算大小 | 初级 |
| AlignOf | 计算对齐 | 初级 |
| OffsetOf | 计算字段偏移 | 中级 |
| TypeConversion | 类型转换 | 中级 |
| SliceManipulation | 切片操作 | 中级 |
| StringManipulation | 字符串操作 | 中级 |
| ArrayToSlice | 数组转切片 | 中级 |
| ZeroCopyConversion | 零拷贝转换 | 高级 |
| PointerArithmetic | 指针算术 | 高级 |
| StructFieldAccess | 结构体字段直接访问 | 中级 |
| UnionLike | 类Union操作 | 高级 |
| PerformanceComparison | 性能对比 | 中级 |

## 新增示例详解

### BasicPointer 基础指针操作
```go
var x int = 42
ptr := unsafe.Pointer(&x)
intPtr := (*int)(ptr)
*intPtr = 100  // 修改值
```

### SizeOf 计算大小
```go
fmt.Printf("int: %d bytes\n", unsafe.Sizeof(i))
fmt.Printf("string: %d bytes (指针+长度)\n", unsafe.Sizeof(s))
fmt.Printf("[]int: %d bytes (指针+长度+容量)\n", unsafe.Sizeof(sl))
```

### OffsetOf 计算字段偏移
```go
type StructWithPadding struct {
    A int8
    B int
    C int16
}

var s StructWithPadding
fmt.Printf("A: 偏移 %d bytes\n", unsafe.Offsetof(s.A))
fmt.Printf("B: 偏移 %d bytes\n", unsafe.Offsetof(s.B))
fmt.Printf("注意: 有填充字节(padding)用于对齐!")
```

### TypeConversion 类型转换
```go
type IntStruct struct { Value int }
type FloatStruct struct { Value float64 }

var i IntStruct
i.Value = 42

// 内存重新解释
fPtr := (*FloatStruct)(unsafe.Pointer(&i))
fPtr.Value = 3.14  // 修改float会影响int
```

### SliceManipulation 切片操作
```go
type sliceHeader struct {
    Data unsafe.Pointer
    Len  int
    Cap  int
}

header := (*sliceHeader)(unsafe.Pointer(&s))

// 直接通过指针访问元素
elemPtr := (*int)(unsafe.Pointer(
    uintptr(header.Data) + uintptr(i)*unsafe.Sizeof(0),
))
```

### ZeroCopyConversion 零拷贝转换
```go
// string -> []byte 零拷贝
bytes := unsafe.Slice(unsafe.StringData(s), len(s))

// []byte -> string 零拷贝
s2 := unsafe.String(&bytes[0], len(bytes))

fmt.Println("注意：零拷贝共享内存，修改一个会影响另一个!")
fmt.Println("警告：修改string内存可能导致程序崩溃!")
```

### PointerArithmetic 指针算术
```go
arr := []int{10, 20, 30}
basePtr := unsafe.Pointer(&arr[0])

for i := 0; i < len(arr); i++ {
    offset := uintptr(i) * unsafe.Sizeof(arr[0])
    elemPtr := (*int)(unsafe.Pointer(
        uintptr(basePtr) + offset,
    ))
    fmt.Printf("索引 %d: %d\n", i, *elemPtr)
}
```

### StructFieldAccess 结构体字段直接访问
```go
type Person struct {
    Name string
    Age  int
}

p := Person{Name: "Alice", Age: 30}

// 直接通过偏移访问字段
namePtr := (*string)(unsafe.Pointer(&p))
fmt.Printf("Name: %s\n", *namePtr)

ageOffset := unsafe.Offsetof(p.Age)
agePtr := (*int)(unsafe.Pointer(
    uintptr(unsafe.Pointer(&p)) + ageOffset,
))
fmt.Printf("Age: %d\n", *agePtr)
```

### UnionLike 类Union操作
```go
type Union struct {
    data [8]byte
}

u := Union{}

// 解释为int64
intPtr := (*int64)(unsafe.Pointer(&u.data[0]))
*intPtr = 0x123456789ABCDEF0

// 解释为float64
floatPtr := (*float64)(unsafe.Pointer(&u.data[0]))
```

## 常见陷阱

### 1. 内存安全
```go
// 错误: 超出边界访问
slice := []int{1, 2, 3}
ptr := unsafe.Pointer(&slice[0])
badPtr := unsafe.Pointer(uintptr(ptr) + 3*unsafe.Sizeof(0))
*(*int)(badPtr) = 4  // 可能崩溃

// 正确: 只访问有效内存
```

### 2. 字符串不变性
```go
// 错误: 修改字符串内存
s := "hello"
bytes := unsafe.Slice(unsafe.StringData(s), len(s))
bytes[0] = 'H'  // 可能崩溃，字符串在只读内存

// 正确: 先拷贝再修改
```

### 3. 垃圾回收
```go
// 错误: 保存指向临时对象的指针
var ptr unsafe.Pointer
func bad() {
    x := 42
    ptr = unsafe.Pointer(&x)
}
bad()
// 此时x可能被GC回收，ptr是悬垂指针
```

## unsafe包的用途

### 1. 与C代码交互
```go
// cgo中常用
// #include <stdio.h>
import "C"
```

### 2. 零拷贝类型转换
```go
// string <-> []byte 零拷贝
bytes := unsafe.Slice(unsafe.StringData(s), len(s))
```

### 3. 性能优化（极端情况）
```go
// 跳过边界检查等安全检查
// 仅在性能测试证明值得时使用
```

### 4. 底层内存操作
```go
// 直接操作内存布局
type sliceHeader struct {
    Data unsafe.Pointer
    Len  int
    Cap  int
}
```

## 最佳实践

1. **尽量避免**: 优先使用类型安全的方式
2. **隔离使用**: 将unsafe代码限制在小范围内
3. **充分测试**: 测试不同平台和架构
4. **添加注释**: 详细说明为什么使用unsafe
5. **考虑替代方案**: 代码生成、接口等
6. **安全检查**: 使用go vet和race detector

## 运行测试

```bash
go test -v ./unsafe/...
```

## 参考资料

- [unsafe包文档](https://pkg.go.dev/unsafe)
- [Go内存模型](https://go.dev/ref/mem)
