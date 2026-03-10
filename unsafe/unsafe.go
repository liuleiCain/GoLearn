package unsafe

import (
	"fmt"
	"unsafe"
)

// BasicPointer 基础指针操作
func BasicPointer() {
	fmt.Println("=== 基础指针操作 ===")

	var x int = 42
	fmt.Printf("x的值: %d\n", x)
	fmt.Printf("x的地址: %p\n", &x)

	// unsafe.Pointer 可以转换任何类型的指针
	ptr := unsafe.Pointer(&x)
	fmt.Printf("unsafe.Pointer: %v\n", ptr)

	// 转换回 *int
	intPtr := (*int)(ptr)
	fmt.Printf("通过指针读取: %d\n", *intPtr)

	// 修改值
	*intPtr = 100
	fmt.Printf("修改后x的值: %d\n", x)
}

// SizeOf 计算大小
func SizeOf() {
	fmt.Println("=== SizeOf 计算大小 ===")

	var i int
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	var f64 float64
	var s string
	var b bool
	var p *int
	var arr [5]int
	var sl []int
	var m map[int]int

	fmt.Printf("int:         %d bytes\n", unsafe.Sizeof(i))
	fmt.Printf("int8:        %d bytes\n", unsafe.Sizeof(i8))
	fmt.Printf("int16:       %d bytes\n", unsafe.Sizeof(i16))
	fmt.Printf("int32:       %d bytes\n", unsafe.Sizeof(i32))
	fmt.Printf("int64:       %d bytes\n", unsafe.Sizeof(i64))
	fmt.Printf("float64:     %d bytes\n", unsafe.Sizeof(f64))
	fmt.Printf("string:      %d bytes (指针+长度)\n", unsafe.Sizeof(s))
	fmt.Printf("bool:        %d bytes\n", unsafe.Sizeof(b))
	fmt.Printf("*int:        %d bytes\n", unsafe.Sizeof(p))
	fmt.Printf("[5]int:      %d bytes\n", unsafe.Sizeof(arr))
	fmt.Printf("[]int:       %d bytes (指针+长度+容量)\n", unsafe.Sizeof(sl))
	fmt.Printf("map[int]int: %d bytes\n", unsafe.Sizeof(m))
}

// AlignOf 计算对齐
func AlignOf() {
	fmt.Println("=== AlignOf 计算对齐 ===")

	var i int
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	var f64 float64
	var s string
	var p *int

	fmt.Printf("int:         %d 字节对齐\n", unsafe.Alignof(i))
	fmt.Printf("int8:        %d 字节对齐\n", unsafe.Alignof(i8))
	fmt.Printf("int16:       %d 字节对齐\n", unsafe.Alignof(i16))
	fmt.Printf("int32:       %d 字节对齐\n", unsafe.Alignof(i32))
	fmt.Printf("int64:       %d 字节对齐\n", unsafe.Alignof(i64))
	fmt.Printf("float64:     %d 字节对齐\n", unsafe.Alignof(f64))
	fmt.Printf("string:      %d 字节对齐\n", unsafe.Alignof(s))
	fmt.Printf("*int:        %d 字节对齐\n", unsafe.Alignof(p))
}

// OffsetOf 计算字段偏移
func OffsetOf() {
	fmt.Println("=== OffsetOf 计算字段偏移 ===")

	type StructWithPadding struct {
		A int8   // 1字节
		B int    // 8字节
		C int16  // 2字节
		D string // 16字节
	}

	var s StructWithPadding
	fmt.Printf("结构体整体大小: %d bytes\n", unsafe.Sizeof(s))
	fmt.Println()

	fmt.Println("字段偏移:")
	fmt.Printf("A: 偏移 %d bytes, 大小 %d bytes\n", unsafe.Offsetof(s.A), unsafe.Sizeof(s.A))
	fmt.Printf("B: 偏移 %d bytes, 大小 %d bytes\n", unsafe.Offsetof(s.B), unsafe.Sizeof(s.B))
	fmt.Printf("C: 偏移 %d bytes, 大小 %d bytes\n", unsafe.Offsetof(s.C), unsafe.Sizeof(s.C))
	fmt.Printf("D: 偏移 %d bytes, 大小 %d bytes\n", unsafe.Offsetof(s.D), unsafe.Sizeof(s.D))

	fmt.Println()
	fmt.Println("注意: 有填充字节(padding)用于对齐!")
}

// TypeConversion 类型转换
func TypeConversion() {
	fmt.Println("=== 类型转换 ===")

	type IntStruct struct {
		Value int
	}

	type FloatStruct struct {
		Value float64
	}

	var i IntStruct
	i.Value = 42
	fmt.Printf("原始值: %+v\n", i)

	// 内存重新解释 - 危险但强大
	fPtr := (*FloatStruct)(unsafe.Pointer(&i))
	fmt.Printf("内存重新解释为FloatStruct: %+v\n", *fPtr)

	// 修改float值会影响int值
	fPtr.Value = 3.14
	fmt.Printf("修改float后, int值变为: %+v\n", i)
}

// SliceManipulation 切片操作
func SliceManipulation() {
	fmt.Println("=== 切片操作 ===")

	// slice结构: [pointer, len, cap]
	type sliceHeader struct {
		Data unsafe.Pointer
		Len  int
		Cap  int
	}

	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("原始切片: %v\n", s)

	// 获取slice内部结构
	header := (*sliceHeader)(unsafe.Pointer(&s))
	fmt.Printf("slice header: Data=%v, Len=%d, Cap=%d\n", header.Data, header.Len, header.Cap)

	// 直接通过指针访问元素
	fmt.Println("通过指针访问元素:")
	for i := 0; i < len(s); i++ {
		elemPtr := (*int)(unsafe.Pointer(uintptr(header.Data) + uintptr(i)*unsafe.Sizeof(0)))
		fmt.Printf("  s[%d] = %d\n", i, *elemPtr)
	}

	// 修改元素
	*(*int)(unsafe.Pointer(uintptr(header.Data) + uintptr(0)*unsafe.Sizeof(0))) = 100
	fmt.Printf("修改后: %v\n", s)
}

// StringManipulation 字符串操作
func StringManipulation() {
	fmt.Println("=== 字符串操作 ===")

	// string结构: [pointer, len]
	type stringHeader struct {
		Data unsafe.Pointer
		Len  int
	}

	s := "Hello World"
	fmt.Printf("原始字符串: %s\n", s)

	header := (*stringHeader)(unsafe.Pointer(&s))
	fmt.Printf("string header: Data=%v, Len=%d\n", header.Data, header.Len)

	// 字符串是不可变的，但可以通过unsafe转换为[]byte修改
	// 注意：这会破坏字符串不变性，可能导致问题
	fmt.Println("\n将字符串转换为[]byte (共享内存):")
	bytes := *(*[]byte)(unsafe.Pointer(&struct {
		string
		Cap int
	}{s, len(s)}))
	fmt.Printf("[]byte: %v\n", bytes)

	// 警告：修改会影响原字符串！
	// bytes[0] = 'h' // 可能导致崩溃，因为字符串在只读内存中
}

// ArrayToSlice 数组转切片
func ArrayToSlice() {
	fmt.Println("=== 数组转切片 ===")

	arr := [5]int{1, 2, 3, 4, 5}
	fmt.Printf("数组: %v\n", arr)

	// 方式1: 标准方式
	sl1 := arr[:]
	fmt.Printf("标准方式切片: %v\n", sl1)

	// 方式2: 使用unsafe
	type sliceHeader struct {
		Data unsafe.Pointer
		Len  int
		Cap  int
	}

	sl2 := *(*[]int)(unsafe.Pointer(&sliceHeader{
		Data: unsafe.Pointer(&arr[0]),
		Len:  len(arr),
		Cap:  len(arr),
	}))
	fmt.Printf("unsafe方式切片: %v\n", sl2)

	// 修改切片会影响数组
	sl2[0] = 100
	fmt.Printf("修改后数组: %v\n", arr)
}

// ZeroCopyConversion 零拷贝转换
func ZeroCopyConversion() {
	fmt.Println("=== 零拷贝转换 ===")

	// string <-> []byte 零拷贝转换
	s := "Hello Go"
	fmt.Printf("原始string: %s\n", s)

	// string -> []byte 零拷贝
	bytes := unsafe.Slice(unsafe.StringData(s), len(s))
	fmt.Printf("零拷贝[]byte: %v\n", bytes)

	// []byte -> string 零拷贝
	s2 := unsafe.String(&bytes[0], len(bytes))
	fmt.Printf("零拷贝转回string: %s\n", s2)

	fmt.Println("\n注意：零拷贝共享内存，修改一个会影响另一个!")
	fmt.Println("警告：修改string内存可能导致程序崩溃!")
}

// PointerArithmetic 指针算术
func PointerArithmetic() {
	fmt.Println("=== 指针算术 ===")

	arr := []int{10, 20, 30, 40, 50}
	fmt.Printf("数组: %v\n", arr)

	// 获取第一个元素的指针
	basePtr := unsafe.Pointer(&arr[0])

	fmt.Println("通过指针算术访问元素:")
	for i := 0; i < len(arr); i++ {
		// 计算偏移: base + i * size
		offset := uintptr(i) * unsafe.Sizeof(arr[0])
		elemPtr := (*int)(unsafe.Pointer(uintptr(basePtr) + offset))
		fmt.Printf("  索引 %d: %d\n", i, *elemPtr)
	}
}

// StructFieldAccess 结构体字段直接访问
func StructFieldAccess() {
	fmt.Println("=== 结构体字段直接访问 ===")

	type Person struct {
		Name string
		Age  int
		City string
	}

	p := Person{Name: "Alice", Age: 30, City: "Beijing"}
	fmt.Printf("原始: %+v\n", p)

	// 直接通过偏移访问字段
	namePtr := (*string)(unsafe.Pointer(&p))
	fmt.Printf("Name (偏移0): %s\n", *namePtr)

	ageOffset := unsafe.Offsetof(p.Age)
	agePtr := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + ageOffset))
	fmt.Printf("Age (偏移%d): %d\n", ageOffset, *agePtr)

	cityOffset := unsafe.Offsetof(p.City)
	cityPtr := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + cityOffset))
	fmt.Printf("City (偏移%d): %s\n", cityOffset, *cityPtr)

	// 直接修改
	*agePtr = 35
	fmt.Printf("修改后: %+v\n", p)
}

// UnionLike 类Union操作
func UnionLike() {
	fmt.Println("=== 类Union操作 ===")

	// 同一块内存可以解释为不同类型
	type Union struct {
		data [8]byte
	}

	u := Union{}

	// 解释为int64
	intPtr := (*int64)(unsafe.Pointer(&u.data[0]))
	*intPtr = 0x123456789ABCDEF0
	fmt.Printf("作为int64: 0x%X\n", *intPtr)

	// 解释为float64
	floatPtr := (*float64)(unsafe.Pointer(&u.data[0]))
	fmt.Printf("作为float64: %g\n", *floatPtr)

	// 查看字节
	fmt.Printf("字节数据: ")
	for i := 0; i < 8; i++ {
		fmt.Printf("%02X ", u.data[i])
	}
	fmt.Println()
}

// PerformanceComparison 性能对比
func PerformanceComparison() {
	fmt.Println("=== unsafe vs 普通操作对比 ===")

	arr := make([]int, 1000000)
	for i := range arr {
		arr[i] = i
	}

	// 注意：这只是演示，实际性能影响取决于具体场景
	fmt.Println("说明:")
	fmt.Println("- unsafe操作跳过了Go的安全检查")
	fmt.Println("- 可能更快，但也更危险")
	fmt.Println("- 除非绝对必要，否则不推荐使用")
	fmt.Println()
	fmt.Println("unsafe包的用途:")
	fmt.Println("1. 与C代码交互")
	fmt.Println("2. 零拷贝类型转换")
	fmt.Println("3. 性能优化（极端情况）")
	fmt.Println("4. 底层内存操作")
}
