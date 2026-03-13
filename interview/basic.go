package interview

import "fmt"

// ============================================================
// 基础语法面试题
// ============================================================

// Q1: Go语言的主要特点是什么？
// A1: Go语言的主要特点包括：
// 1. 静态类型、编译型语言，执行效率高
// 2. 语法简洁，学习曲线平缓
// 3. 内置并发支持（goroutine和channel）
// 4. 垃圾回收机制，自动内存管理
// 5. 编译速度快，跨平台编译
// 6. 丰富的标准库
// 7. 接口设计灵活，隐式实现

// Q2: Go语言中的数据类型有哪些？
// A2: Go语言的数据类型分为：
// 1. 基本类型：bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, byte, rune, float32, float64, complex64, complex128
// 2. 复合类型：array, struct, slice, map, channel
// 3. 引用类型：slice, map, channel, interface, function
// 4. 指针类型：*T

// Q3: Go语言中make和new的区别？
// A3: make和new的主要区别：
// 1. new(T) 为类型T分配零值内存，返回*T（指针）
// 2. make(T, args) 只用于slice、map、channel的初始化，返回T（不是指针）
// 3. new返回指针，make返回引用类型本身
// 4. new分配的内存会被初始化为零值

func MakeVsNew() {
	fmt.Println("=== make vs new ===")

	// new 返回指针
	p := new(int)
	fmt.Printf("new(int): type=%T, value=%v, *value=%d\n", p, p, *p)

	// make 返回引用类型本身
	s := make([]int, 5)
	fmt.Printf("make([]int, 5): type=%T, value=%v, len=%d\n", s, s, len(s))

	m := make(map[string]int)
	fmt.Printf("make(map[string]int): type=%T, value=%v\n", m, m)

	ch := make(chan int, 10)
	fmt.Printf("make(chan int, 10): type=%T, value=%v\n", ch, ch)
}

// Q4: Go语言中数组和切片的区别？
// A4: 数组和切片的主要区别：
// 1. 数组是值类型，切片是引用类型
// 2. 数组长度固定，切片长度可变
// 3. 数组作为参数传递会复制整个数组，切片只复制切片头
// 4. 切片有容量概念，可以动态扩容

func ArrayVsSlice() {
	fmt.Println("\n=== 数组 vs 切片 ===")

	// 数组：长度固定，值类型
	arr := [3]int{1, 2, 3}
	arrCopy := arr
	arrCopy[0] = 100
	fmt.Printf("数组: arr=%v, arrCopy=%v (修改arrCopy不影响arr)\n", arr, arrCopy)

	// 切片：长度可变，引用类型
	slice := []int{1, 2, 3}
	sliceRef := slice
	sliceRef[0] = 100
	fmt.Printf("切片: slice=%v, sliceRef=%v (修改sliceRef会影响slice)\n", slice, sliceRef)

	// 切片的扩容
	s := make([]int, 0, 3)
	fmt.Printf("初始: len=%d, cap=%d\n", len(s), cap(s))
	for i := 0; i < 10; i++ {
		s = append(s, i)
		fmt.Printf("追加%d: len=%d, cap=%d\n", i, len(s), cap(s))
	}
}

// Q5: Go语言中defer的执行顺序？
// A5: defer的执行顺序：
// 1. LIFO（后进先出）顺序执行
// 2. defer语句在函数返回前执行
// 3. defer语句的参数在声明时就确定了
// 4. defer可以修改命名返回值

func DeferOrder() {
	fmt.Println("\n=== defer执行顺序 ===")

	defer func() {
		fmt.Println("defer 1")
	}()
	defer func() {
		fmt.Println("defer 2")
	}()
	defer func() {
		fmt.Println("defer 3")
	}()

	fmt.Println("函数体")
	// 输出顺序：函数体 -> defer 3 -> defer 2 -> defer 1
}

// Q6: Go语言中值传递和引用传递？
// A6: Go语言只有值传递：
// 1. 所有参数传递都是值传递
// 2. 传递指针时，复制的是指针的值（地址）
// 3. 传递切片、map、channel时，复制的是引用的值
// 4. 传递结构体时，复制整个结构体

func ValueVsReference() {
	fmt.Println("\n=== 值传递 vs 引用传递 ===")

	// 值传递
	x := 10
	modifyValue := func(n int) {
		n = 100
	}
	modifyValue(x)
	fmt.Printf("值传递: x=%d (未改变)\n", x)

	// 指针传递（本质也是值传递，传递的是指针的值）
	modifyPointer := func(n *int) {
		*n = 100
	}
	modifyPointer(&x)
	fmt.Printf("指针传递: x=%d (已改变)\n", x)

	// 切片传递（传递的是切片头的副本）
	slice := []int{1, 2, 3}
	modifySlice := func(s []int) {
		s[0] = 100
	}
	modifySlice(slice)
	fmt.Printf("切片传递: slice=%v (已改变)\n", slice)

	// append操作不会影响原切片
	appendSlice := func(s []int) {
		s = append(s, 4)
	}
	appendSlice(slice)
	fmt.Printf("append后: slice=%v (未改变)\n", slice)
}

// Q7: Go语言中init函数的特点？
// A7: init函数的特点：
// 1. 每个包可以有多个init函数
// 2. init函数在包被导入时自动执行
// 3. init函数没有参数和返回值
// 4. 执行顺序：常量->变量->init函数
// 5. 包的init函数执行顺序：依赖包->当前包

// Q8: Go语言中slice的扩容机制？
// A8: slice扩容机制（Go 1.18+）：
// 1. 如果新容量 > 2倍旧容量，直接使用新容量
// 2. 如果旧容量 < 256，新容量 = 2倍旧容量
// 3. 如果旧容量 >= 256，新容量 = 旧容量 + (旧容量+3*256)/4
// 4. 最后根据元素大小进行内存对齐

func SliceGrowth() {
	fmt.Println("\n=== 切片扩容机制 ===")

	s := make([]int, 0)
	oldCap := 0

	for i := 0; i < 1000; i++ {
		s = append(s, i)
		if cap(s) != oldCap {
			fmt.Printf("len=%d, oldCap=%d -> newCap=%d\n", len(s), oldCap, cap(s))
			oldCap = cap(s)
		}
	}
}

// Q9: Go语言中map的特点？
// A9: map的特点：
// 1. map是无序的，每次遍历顺序可能不同
// 2. map不是并发安全的，并发读写需要加锁
// 3. map的key必须是可以比较的类型
// 4. map的零值是nil，需要make初始化
// 5. 从map中取不存在的key，返回value类型的零值

func MapFeatures() {
	fmt.Println("\n=== map特点 ===")

	// map的零值是nil
	var m map[string]int
	fmt.Printf("nil map: m=%v, m==nil: %v\n", m, m == nil)

	// 需要make初始化
	m = make(map[string]int)
	m["a"] = 1
	fmt.Printf("初始化后: m=%v\n", m)

	// 取不存在的key返回零值
	value := m["b"]
	fmt.Printf("不存在的key: m[\"b\"]=%d\n", value)

	// 判断key是否存在
	if v, ok := m["b"]; ok {
		fmt.Printf("key存在: value=%d\n", v)
	} else {
		fmt.Println("key不存在")
	}

	// map遍历顺序不确定
	m["c"] = 3
	m["d"] = 4
	fmt.Print("遍历map: ")
	for k, v := range m {
		fmt.Printf("%s=%d ", k, v)
	}
	fmt.Println()
}

// Q10: Go语言中rune和byte的区别？
// A10: rune和byte的区别：
// 1. byte是uint8的别名，表示一个字节
// 2. rune是int32的别名，表示一个Unicode码点
// 3. 处理ASCII字符用byte，处理中文等Unicode字符用rune
// 4. len(string)返回字节数，不是字符数

func RuneVsByte() {
	fmt.Println("\n=== rune vs byte ===")

	str := "Hello世界"

	// byte遍历（按字节）
	fmt.Print("byte遍历: ")
	for i := 0; i < len(str); i++ {
		fmt.Printf("%d ", str[i])
	}
	fmt.Printf("\n字节数: len=%d\n", len(str))

	// rune遍历（按字符）
	fmt.Print("rune遍历: ")
	for _, r := range str {
		fmt.Printf("%c ", r)
	}
	fmt.Printf("\n字符数: len([]rune)=%d\n", len([]rune(str)))
}
