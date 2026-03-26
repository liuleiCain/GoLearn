package internals

import (
	"fmt"
	"unsafe"
)

// ============================================================
// Slice 底层实现原理
// ============================================================

// SliceHeader 切片头部结构（与runtime.sliceHeader一致）
// 这是切片在内存中的真实布局
type SliceHeader struct {
	Data uintptr // 指向底层数组的指针
	Len  int     // 切片长度
	Cap  int     // 切片容量
}

// DemoSliceStructure 演示切片的内存结构
func DemoSliceStructure() {
	fmt.Println("=== 切片内存结构 ===")

	// 创建一个切片
	s := make([]int, 3, 5)
	s[0] = 1
	s[1] = 2
	s[2] = 3

	// 通过反射获取切片头部
	header := (*SliceHeader)(unsafe.Pointer(&s))

	fmt.Printf("切片: %v\n", s)
	fmt.Printf("Data指针: 0x%x\n", header.Data)
	fmt.Printf("Len: %d\n", header.Len)
	fmt.Printf("Cap: %d\n", header.Cap)

	// 内存布局图示
	fmt.Println("\n内存布局:")
	fmt.Println("+--------+--------+--------+")
	fmt.Println("|  Data  |  Len   |  Cap   |")
	fmt.Printf("| 0x%x |   %d    |   %d    |\n", header.Data, header.Len, header.Cap)
	fmt.Println("+--------+--------+--------+")
	fmt.Println("         |")
	fmt.Println("         v")
	fmt.Println("+----+----+----+----+----+")
	fmt.Println("|  1 |  2 |  3 |    |    |")
	fmt.Println("+----+----+----+----+----+")
	fmt.Println("  [0]  [1]  [2]  [3]  [4]")
	fmt.Println("  <- Len=3 ->")
	fmt.Println("  <------ Cap=5 ------>")
}

// DemoSliceSharing 演示切片共享底层数组
func DemoSliceSharing() {
	fmt.Println("\n=== 切片共享底层数组 ===")

	// 原始切片
	original := make([]int, 5, 10)
	for i := 0; i < 5; i++ {
		original[i] = i + 1
	}

	// 创建子切片
	s1 := original[1:3] // [2, 3]
	s2 := original[2:5] // [3, 4, 5]

	// 获取各切片的头部信息
	h0 := (*SliceHeader)(unsafe.Pointer(&original))
	h1 := (*SliceHeader)(unsafe.Pointer(&s1))
	h2 := (*SliceHeader)(unsafe.Pointer(&s2))

	fmt.Printf("original: Data=0x%x, Len=%d, Cap=%d\n", h0.Data, h0.Len, h0.Cap)
	fmt.Printf("s1:       Data=0x%x, Len=%d, Cap=%d\n", h1.Data, h1.Len, h1.Cap)
	fmt.Printf("s2:       Data=0x%x, Len=%d, Cap=%d\n", h2.Data, h2.Len, h2.Cap)

	// 验证共享：Data指针相差的偏移量
	fmt.Printf("\n指针偏移: s1.Data - original.Data = %d (等于1 * 8 = 8字节)\n", h1.Data-h0.Data)
	fmt.Printf("指针偏移: s2.Data - original.Data = %d (等于2 * 8 = 16字节)\n", h2.Data-h0.Data)

	// 修改s1会影响original和s2
	s1[1] = 100
	fmt.Printf("\n修改s1[1]=100后:\n")
	fmt.Printf("original: %v\n", original)
	fmt.Printf("s1: %v\n", s1)
	fmt.Printf("s2: %v\n", s2)
}

// DemoSliceGrowth 演示切片扩容机制
func DemoSliceGrowth() {
	fmt.Println("\n=== 切片扩容机制 ===")

	// Go 1.18+ 扩容算法
	fmt.Println("扩容规则 (Go 1.18+):")
	fmt.Println("1. 如果新容量 > 2倍旧容量 -> 使用新容量")
	fmt.Println("2. 如果旧容量 < 256 -> 新容量 = 2倍旧容量")
	fmt.Println("3. 如果旧容量 >= 256 -> 新容量 = 旧容量 + (旧容量+3*256)/4")
	fmt.Println("4. 最后进行内存对齐")

	s := make([]int, 0)
	oldCap := 0

	fmt.Println("\n扩容过程:")
	for i := 0; i < 20; i++ {
		s = append(s, i)
		if cap(s) != oldCap {
			growthRate := float64(cap(s)) / float64(oldCap)
			if oldCap == 0 {
				growthRate = 0
			}
			fmt.Printf("Len=%2d, Cap: %3d -> %3d (增长率: %.2fx)\n",
				len(s), oldCap, cap(s), growthRate)
			oldCap = cap(s)
		}
	}
}

// DemoSliceCopy 演示copy函数的实现原理
func DemoSliceCopy() {
	fmt.Println("\n=== copy函数原理 ===")

	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)

	// copy会复制min(len(src), len(dst))个元素
	n := copy(dst, src)

	fmt.Printf("src: %v (len=%d)\n", src, len(src))
	fmt.Printf("dst: %v (len=%d)\n", dst, len(dst))
	fmt.Printf("复制了 %d 个元素\n", n)

	// copy的本质是memmove
	fmt.Println("\ncopy内部实现（伪代码）:")
	fmt.Println("func copy(dst, src []T) int {")
	fmt.Println("    n := min(len(dst), len(src))")
	fmt.Println("    memmove(dst.Data, src.Data, n*sizeof(T))")
	fmt.Println("    return n")
	fmt.Println("}")
}

// DemoSliceAppend 演示append函数的实现原理
func DemoSliceAppend() {
	fmt.Println("\n=== append函数原理 ===")

	s := make([]int, 2, 4)
	s[0] = 1
	s[1] = 2

	h := (*SliceHeader)(unsafe.Pointer(&s))
	oldData := h.Data

	fmt.Printf("追加前: %v, len=%d, cap=%d\n", s, len(s), cap(s))

	// 追加元素（容量足够，不扩容）
	s = append(s, 3)
	fmt.Printf("追加后: %v, len=%d, cap=%d\n", s, len(s), cap(s))
	fmt.Printf("Data指针变化: %v (未扩容，指针不变)\n", h.Data == oldData)

	// 继续追加（容量不足，需要扩容）
	s = append(s, 4)
	s = append(s, 5) // 触发扩容
	fmt.Printf("扩容后: %v, len=%d, cap=%d\n", s, len(s), cap(s))
	fmt.Printf("Data指针变化: %v (扩容后，指针改变)\n", h.Data != oldData)

	// append内部实现
	fmt.Println("\nappend内部实现（伪代码）:")
	fmt.Println("func append(s []T, elems ...T) []T {")
	fmt.Println("    if len(s) + len(elems) > cap(s) {")
	fmt.Println("        newCap := growslice(cap(s), len(s)+len(elems))")
	fmt.Println("        newData := mallocgc(newCap * sizeof(T))")
	fmt.Println("        memmove(newData, s.Data, len(s)*sizeof(T))")
	fmt.Println("        s.Data = newData")
	fmt.Println("        s.Cap = newCap")
	fmt.Println("    }")
	fmt.Println("    memmove(s.Data+len(s), elems, len(elems)*sizeof(T))")
	fmt.Println("    s.Len += len(elems)")
	fmt.Println("    return s")
	fmt.Println("}")
}

// DemoSliceExpression 演示切片表达式
func DemoSliceExpression() {
	fmt.Println("\n=== 切片表达式 ===")

	arr := [5]int{1, 2, 3, 4, 5}
	fmt.Printf("原数组: %v\n", arr)

	// 简单切片表达式 s[low:high]
	s1 := arr[1:4] // [2, 3, 4]
	fmt.Printf("arr[1:4] = %v\n", s1)

	// 完整切片表达式 s[low:high:max]
	// max限制新切片的容量
	s2 := arr[1:4:4] // len=3, cap=3
	s3 := arr[1:4:5] // len=3, cap=4

	h2 := (*SliceHeader)(unsafe.Pointer(&s2))
	h3 := (*SliceHeader)(unsafe.Pointer(&s3))

	fmt.Printf("arr[1:4:4] -> len=%d, cap=%d\n", h2.Len, h2.Cap)
	fmt.Printf("arr[1:4:5] -> len=%d, cap=%d\n", h3.Len, h3.Cap)

	fmt.Println("\n切片表达式规则:")
	fmt.Println("s[low:high]      -> len = high-low, cap = len(s)-low")
	fmt.Println("s[low:high:max]  -> len = high-low, cap = max-low")
	fmt.Println("注意: 0 <= low <= high <= max <= cap(s)")
}

// DemoSliceNil 演示nil切片和空切片的区别
func DemoSliceNil() {
	fmt.Println("\n=== nil切片 vs 空切片 ===")

	// nil切片
	var nilSlice []int

	// 空切片
	emptySlice := make([]int, 0)
	// 或者
	emptySlice2 := []int{}

	fmt.Printf("nil切片: %v, len=%d, cap=%d, nil=%v\n",
		nilSlice, len(nilSlice), cap(nilSlice), nilSlice == nil)
	fmt.Printf("空切片1: %v, len=%d, cap=%d, nil=%v\n",
		emptySlice, len(emptySlice), cap(emptySlice), emptySlice == nil)
	fmt.Printf("空切片2: %v, len=%d, cap=%d, nil=%v\n",
		emptySlice2, len(emptySlice2), cap(emptySlice2), emptySlice2 == nil)

	// 查看Data指针
	hNil := (*SliceHeader)(unsafe.Pointer(&nilSlice))
	hEmpty := (*SliceHeader)(unsafe.Pointer(&emptySlice))

	fmt.Printf("\nnil切片Data指针: 0x%x (值为0)\n", hNil.Data)
	fmt.Printf("空切片Data指针: 0x%x (指向一个非nil地址)\n", hEmpty.Data)

	fmt.Println("\nJSON序列化差异:")
	fmt.Println("nil切片 -> null")
	fmt.Println("空切片  -> []")
}

// DemoSliceMemoryLeak 演示切片内存泄漏问题
func DemoSliceMemoryLeak() {
	fmt.Println("\n=== 切片内存泄漏问题 ===")

	// 场景1：子切片引用大数组
	fmt.Println("场景1：子切片引用大数组")
	largeSlice := make([]int, 1000000)
	for i := range largeSlice {
		largeSlice[i] = i
	}

	// 只需要前10个元素，但整个大数组无法被GC
	smallSlice := largeSlice[:10]
	fmt.Printf("smallSlice: len=%d, cap=%d\n", len(smallSlice), cap(smallSlice))
	fmt.Println("问题: 虽然只需要10个元素，但整个100万元素的数组无法被回收")

	// 正确做法：复制一份
	correctSlice := make([]int, 10)
	copy(correctSlice, largeSlice[:10])
	largeSlice = nil // 释放大数组的引用
	fmt.Printf("correctSlice: len=%d, cap=%d\n", len(correctSlice), cap(correctSlice))

	// 场景2：切片作为函数参数
	fmt.Println("\n场景2：切片作为函数参数")
	fmt.Println("切片作为参数传递时，只复制切片头（24字节）")
	fmt.Println("底层数组不会被复制，多个切片共享同一数组")
}
