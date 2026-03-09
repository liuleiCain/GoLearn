package slice

import (
	"fmt"
	"reflect"
	"unsafe"
)

func SliceStructure() {
	fmt.Println("=== 切片结构示例 ===")

	arr := [5]int{1, 2, 3, 4, 5}
	slice := arr[1:4]

	fmt.Printf("数组: %v, 长度: %d\n", arr, len(arr))
	fmt.Printf("切片: %v, 长度: %d, 容量: %d\n", slice, len(slice), cap(slice))

	fmt.Println("\n切片底层结构:")
	fmt.Println("- 指针: 指向底层数组的起始位置")
	fmt.Println("- 长度(len): 切片中元素的数量")
	fmt.Println("- 容量(cap): 从起始位置到数组末尾的元素数量")
}

func SliceCreate() {
	fmt.Println("=== 切片创建方式示例 ===")

	s1 := []int{1, 2, 3}
	fmt.Printf("字面量创建: %v, len=%d, cap=%d\n", s1, len(s1), cap(s1))

	s2 := make([]int, 3)
	fmt.Printf("make(len): %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))

	s3 := make([]int, 3, 5)
	fmt.Printf("make(len,cap): %v, len=%d, cap=%d\n", s3, len(s3), cap(s3))

	arr := [5]int{1, 2, 3, 4, 5}
	s4 := arr[1:3]
	fmt.Printf("数组切片: %v, len=%d, cap=%d\n", s4, len(s4), cap(s4))

	var s5 []int
	fmt.Printf("nil切片: %v, len=%d, cap=%d, == nil: %v\n", s5, len(s5), cap(s5), s5 == nil)
}

func AppendDemo() {
	fmt.Println("=== Append扩容示例 ===")

	s := make([]int, 0, 2)
	fmt.Printf("初始: %v, len=%d, cap=%d\n", s, len(s), cap(s))

	for i := 1; i <= 10; i++ {
		oldCap := cap(s)
		s = append(s, i)
		if cap(s) != oldCap {
			fmt.Printf("添加%d后: %v, len=%d, cap=%d (扩容: %d -> %d)\n",
				i, s, len(s), cap(s), oldCap, cap(s))
		}
	}
}

func GrowthRule() {
	fmt.Println("=== 扩容规则示例 ===")

	fmt.Println("Go 1.18之前的扩容规则:")
	fmt.Println("- 容量 < 1024: 翻倍")
	fmt.Println("- 容量 >= 1024: 增加25%")

	fmt.Println("\nGo 1.18之后的扩容规则:")
	fmt.Println("- 更平滑的增长曲线")
	fmt.Println("- 避免小切片扩容后容量过大")

	testGrowth := func(initialCap int) {
		s := make([]int, 0, initialCap)
		fmt.Printf("初始: %v, len=%d, cap=%d\n", s, len(s), cap(s))
		for i := 0; i < 2100; i++ {
			oldCap := cap(s)
			s = append(s, i)
			if cap(s) != oldCap {
				fmt.Printf("第%d次遍历时触发扩容 初始cap=%d, 扩容: %d -> %d (增长%.1f%%)\n",
					i, initialCap, oldCap, cap(s), float64(cap(s)-oldCap)/float64(oldCap)*100)
				break
			}
		}
	}

	testGrowth(1)
	testGrowth(2)
	testGrowth(4)
	testGrowth(8)
	testGrowth(512)
	testGrowth(800)
	testGrowth(1024)
	testGrowth(2048)
}

func CopyDemo() {
	fmt.Println("=== Copy函数示例 ===")

	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)

	n := copy(dst, src)
	fmt.Printf("源: %v, 目标: %v, 复制数量: %d\n", src, dst, n)

	dst2 := make([]int, 10)
	n2 := copy(dst2, src)
	fmt.Printf("源: %v, 目标: %v, 复制数量: %d\n", src, dst2, n2)

	fmt.Println("\n切片共享底层数组问题:")
	arr := [5]int{1, 2, 3, 4, 5}
	s1 := arr[:3]
	s2 := arr[2:]
	fmt.Printf("s1: %v, s2: %v\n", s1, s2)
	s1[2] = 100
	fmt.Printf("修改s1[2]=100后: s1: %v, s2: %v\n", s1, s2)

	fmt.Println("\n使用copy避免共享:")
	s3 := make([]int, 3)
	copy(s3, arr[:3])
	s3[2] = 200
	fmt.Printf("s3: %v, arr: %v\n", s3, arr)
}

func SliceExpression() {
	fmt.Println("=== 切片表达式示例 ===")

	arr := [5]int{1, 2, 3, 4, 5}
	fmt.Printf("原数组: %v\n", arr)

	fmt.Printf("arr[1:3]: %v\n", arr[1:3])
	fmt.Printf("arr[:3]: %v\n", arr[:3])
	fmt.Printf("arr[2:]: %v\n", arr[2:])
	fmt.Printf("arr[:]: %v\n", arr[:])
	s1 := arr[1:3]
	fmt.Printf("arr[1:3]: %v, len=%d, cap=%d\n", s1, len(s1), cap(s1))
	fmt.Println("\n完整切片表达式 arr[low:high:max]:")
	s := arr[2:3:4]
	fmt.Printf("arr[1:3:4]: %v, len=%d, cap=%d\n", s, len(s), cap(s))
}

func SliceTricks() {
	fmt.Println("=== 切片技巧示例 ===")

	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("原切片: %v\n", s)

	s = append(s[:2], s[3:]...)
	fmt.Printf("删除索引2: %v\n", s)

	s = []int{1, 2, 3, 4, 5}
	s = append(s, 6)
	fmt.Printf("追加元素: %v\n", s)

	s = []int{1, 2, 3, 4, 5}
	s = append([]int{0}, s...)
	fmt.Printf("头部插入: %v\n", s)

	s = []int{1, 2, 3, 4, 5}
	s = append(s[:2], append([]int{100}, s[2:]...)...)
	fmt.Printf("索引2插入: %v\n", s)

	s = []int{1, 2, 3, 4, 5}
	for i := range s {
		s[i] = 0
	}
	fmt.Printf("清空切片: %v\n", s)
}

func SliceMemory() {
	fmt.Println("=== 切片内存泄漏示例 ===")

	fmt.Println("场景: 切片保留对大数组的引用")

	largeSlice := make([]int, 1000000)
	for i := range largeSlice {
		largeSlice[i] = i
	}

	smallSlice := largeSlice[:10]
	fmt.Printf("小切片: len=%d, cap=%d\n", len(smallSlice), cap(smallSlice))
	fmt.Println("问题: smallSlice仍引用整个大数组!")

	fmt.Println("\n解决方案: 使用copy")
	newSlice := make([]int, 10)
	copy(newSlice, largeSlice[:10])
	fmt.Printf("新切片: len=%d, cap=%d\n", len(newSlice), cap(newSlice))
	fmt.Println("新切片不再引用大数组")
}

func GetSliceHeader(s []int) (ptr uintptr, len_, cap_ int) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	return header.Data, int(header.Len), int(header.Cap)
}

func SlicePointer() {
	fmt.Println("=== 切片指针示例 ===")

	s := []int{1, 2, 3, 4, 5}
	ptr, len_, cap_ := GetSliceHeader(s)
	fmt.Printf("切片: %v\n", s)
	fmt.Printf("指针: %x, 长度: %d, 容量: %d\n", ptr, len_, cap_)

	s2 := s[1:3]
	ptr2, len2, cap2 := GetSliceHeader(s2)
	fmt.Printf("子切片: %v\n", s2)
	fmt.Printf("指针: %x, 长度: %d, 容量: %d\n", ptr2, len2, cap2)
	fmt.Println("注意: 指针偏移了一个int的大小")
}
