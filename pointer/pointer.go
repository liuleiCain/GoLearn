package pointer

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

func BasicPointer() {
	fmt.Println("=== 基础指针示例 ===")

	x := 42
	p := &x

	fmt.Printf("x的值: %d\n", x)
	fmt.Printf("x的地址: %p\n", &x)
	fmt.Printf("p的值(地址): %p\n", p)
	fmt.Printf("p指向的值: %d\n", *p)

	*p = 100
	fmt.Printf("通过指针修改后x的值: %d\n", x)
}

func PointerArithmetic() {
	fmt.Println("=== 指针运算示例 ===")

	arr := [5]int{10, 20, 30, 40, 50}

	p := &arr[0]
	fmt.Printf("arr[0]地址: %p, 值: %d\n", p, *p)

	fmt.Println("\nGo不支持指针算术运算，但可以使用unsafe包:")
	fmt.Println("注意: unsafe包可能影响可移植性和安全性")

	ptr := unsafe.Pointer(p)
	size := unsafe.Sizeof(arr[0])

	for i := 0; i < 5; i++ {
		newPtr := unsafe.Pointer(uintptr(ptr) + uintptr(i)*size)
		val := *(*int)(newPtr)
		fmt.Printf("arr[%d]: 地址=%p, 值=%d\n", i, newPtr, val)
	}
}

type Person struct {
	Name string
	Age  int
}

func (p *Person) SetNamePointer(name string) {
	p.Name = name
}

func (p Person) SetNameValue(name string) {
	p.Name = name
}

func ReceiverCompare() {
	fmt.Println("=== 值接收者vs指针接收者 ===")

	p := Person{Name: "张三", Age: 25}
	fmt.Printf("原始: %+v\n", p)

	p.SetNameValue("李四")
	fmt.Printf("值接收者修改后: %+v\n", p)

	p.SetNamePointer("王五")
	fmt.Printf("指针接收者修改后: %+v\n", p)

	fmt.Println("\n选择原则:")
	fmt.Println("- 需要修改接收者: 使用指针接收者")
	fmt.Println("- 结构体较大: 使用指针接收者(避免复制)")
	fmt.Println("- 一致性: 如果有指针接收者方法，所有方法都用指针接收者")
}

func NilPointer() {
	fmt.Println("=== nil指针示例 ===")

	var p *Person
	fmt.Printf("nil指针: %v, == nil: %v\n", p, p == nil)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("访问nil指针会panic:", r)
		}
	}()

	fmt.Println(p.Name)
}

func PointerToPointer() {
	fmt.Println("=== 指向指针的指针 ===")

	x := 42
	p1 := &x
	p2 := &p1

	fmt.Printf("x = %d\n", x)
	fmt.Printf("*p1 = %d\n", *p1)
	fmt.Printf("**p2 = %d\n", **p2)

	**p2 = 100
	fmt.Printf("通过**p2修改后, x = %d\n", x)
}

func newAndMake() {
	fmt.Println("=== new和make的区别 ===")

	p1 := new(int)
	fmt.Printf("new(int): 值=%d, 类型=%T\n", *p1, p1)

	p2 := new(Person)
	fmt.Printf("new(Person): %+v\n", *p2)

	s := make([]int, 3, 5)
	fmt.Printf("make slice: %v, len=%d, cap=%d\n", s, len(s), cap(s))

	m := make(map[string]int)
	m["a"] = 1
	fmt.Printf("make map: %v\n", m)

	ch := make(chan int, 2)
	ch <- 1
	fmt.Printf("make channel: %v\n", <-ch)

	fmt.Println("\n区别:")
	fmt.Println("- new: 分配内存，返回指针，值初始化为零值")
	fmt.Println("- make: 只用于slice/map/channel，返回初始化后的值(不是指针)")
}

func EscapeAnalysis() {
	fmt.Println("=== 逃逸分析示例 ===")

	fmt.Println("逃逸分析决定变量在栈还是堆上分配")
	fmt.Println("使用 'go build -gcflags=\"-m -m\"' 查看详细逃逸分析")

	fmt.Println("\n--- 示例1: 栈分配 ---")
	stackAlloc := func() int {
		x := 42
		return x
	}
	a := stackAlloc()
	fmt.Printf("栈分配结果: %d (值返回，不逃逸)\n", a)

	fmt.Println("\n--- 示例2: 堆分配(返回指针) ---")
	heapAlloc := func() *int {
		x := 42
		return &x
	}
	b := heapAlloc()
	fmt.Printf("堆分配结果: %d (返回指针，逃逸到堆)\n", *b)

	fmt.Println("\n--- 示例3: 闭包捕获 ---")
	closureEscape := func() func() int {
		x := 100
		return func() int {
			return x
		}
	}
	c := closureEscape()
	fmt.Printf("闭包捕获结果: %d (闭包变量逃逸)\n", c())

	fmt.Println("\n--- 示例4: 接口类型 ---")
	interfaceEscape := func() interface{} {
		x := 200
		return x
	}
	d := interfaceEscape()
	fmt.Printf("接口转换结果: %v (转换为interface{}逃逸)\n", d)

	fmt.Println("\n--- 示例5: 切片扩容 ---")
	sliceEscape := func() []int {
		s := make([]int, 0)
		for i := 0; i < 1000; i++ {
			s = append(s, i)
		}
		return s
	}
	e := sliceEscape()
	fmt.Printf("切片扩容结果: len=%d (大切片可能逃逸)\n", len(e))

	fmt.Println("\n--- 示例6: Channel发送 ---")
	ch := make(chan *int, 1)
	channelEscape := func(val int) {
		x := val
		ch <- &x
	}
	go func() {
		channelEscape(42)
	}()
	f := <-ch
	fmt.Printf("Channel发送结果: %d (发送指针到channel逃逸)\n", *f)

	fmt.Println("\n逃逸到堆的常见情况:")
	fmt.Println("1. 返回局部变量的指针")
	fmt.Println("2. 发送指针到channel")
	fmt.Println("3. 存储在全局变量")
	fmt.Println("4. 闭包捕获变量")
	fmt.Println("5. 转换为interface{}")
	fmt.Println("6. 大对象(超过32KB)")
	fmt.Println("7. 不确定大小的对象")
}

var globalSum int

func PointerPerformanceTime() {
	fmt.Println("=== 指针时间性能测试 ===")

	type SmallStruct struct {
		a, b int
	}

	type LargeStruct struct {
		data [10000]int
	}

	fmt.Println("\n--- 测试1: 小结构体时间开销 ---")
	fmt.Println("小结构体(16字节): 2个int字段")
	small := SmallStruct{a: 1, b: 2}

	smallValuePass := func(s SmallStruct) int {
		return s.a + s.b
	}
	smallPointerPass := func(s *SmallStruct) int {
		return s.a + s.b
	}

	const iterations = 1000000
	start := time.Now()
	for i := 0; i < iterations; i++ {
		globalSum += smallValuePass(small)
	}
	valueDuration := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		globalSum += smallPointerPass(&small)
	}
	pointerDuration := time.Since(start)

	fmt.Printf("  值传递:   %v (%.2f ns/op)\n", valueDuration, float64(valueDuration.Nanoseconds())/float64(iterations))
	fmt.Printf("  指针传递: %v (%.2f ns/op)\n", pointerDuration, float64(pointerDuration.Nanoseconds())/float64(iterations))
	fmt.Println()
	if valueDuration <= pointerDuration {
		fmt.Println("结论: 小结构体值传递更快或相当")
		fmt.Println("原因: 复制16字节开销很小，指针解引用有额外开销")
	} else {
		fmt.Println("结论: 小结构体指针传递更快")
	}

	fmt.Println("\n--- 测试2: 大结构体时间开销 ---")
	fmt.Println("大结构体(80KB): 10000个int字段")
	large := LargeStruct{}
	for i := range large.data {
		large.data[i] = i
	}

	largeValuePass := func(s LargeStruct) int {
		result := 0
		for j := 0; j < len(s.data); j++ {
			result += s.data[j]
		}
		return result
	}
	largePointerPass := func(s *LargeStruct) int {
		result := 0
		for j := 0; j < len(s.data); j++ {
			result += s.data[j]
		}
		return result
	}

	const largeIterations = 10000
	start = time.Now()
	globalSum = 0
	for i := 0; i < largeIterations; i++ {
		globalSum += largeValuePass(large)
	}
	valueDuration = time.Since(start)

	start = time.Now()
	globalSum = 0
	for i := 0; i < largeIterations; i++ {
		globalSum += largePointerPass(&large)
	}
	pointerDuration = time.Since(start)

	fmt.Printf("  值传递:   %v (%.2f μs/op)\n", valueDuration, float64(valueDuration.Microseconds())/float64(largeIterations))
	fmt.Printf("  指针传递: %v (%.2f μs/op)\n", pointerDuration, float64(pointerDuration.Microseconds())/float64(largeIterations))
	fmt.Println()
	if pointerDuration < valueDuration {
		speedup := float64(valueDuration) / float64(pointerDuration)
		fmt.Printf("结论: 大结构体指针传递快 %.1fx\n", speedup)
		fmt.Println("原因: 避免每次调用复制80KB数据")
	} else {
		fmt.Println("结论: 大结构体值传递更快")
	}

	fmt.Println("\n--- 时间性能总结 ---")
	fmt.Println("┌─────────────┬──────────┬──────────┐")
	fmt.Println("│ 结构体大小  │ 推荐方式 │ 原因     │")
	fmt.Println("├─────────────┼──────────┼──────────┤")
	fmt.Println("│ <= 32字节   │ 值传递   │ 复制快   │")
	fmt.Println("│ > 32字节    │ 指针传递 │ 避免复制 │")
	fmt.Println("│ > 1KB       │ 指针传递 │ 大幅提升 │")
	fmt.Println("└─────────────┴──────────┴──────────┘")
}

func PointerPerformanceMemory() {
	fmt.Println("=== 指针内存性能测试 ===")

	type LargeStruct struct {
		data [10000]int
	}

	fmt.Println("\n--- 测试1: 切片存储内存对比 ---")
	fmt.Println("场景: 存储1000个元素")

	fmt.Println("\n方式A: 值类型切片 []LargeStruct")
	fmt.Println("  每个元素80KB，需要完整存储")
	valueSlice := make([]LargeStruct, 1000)
	for i := range valueSlice {
		valueSlice[i].data[0] = i
	}

	fmt.Println("\n方式B: 指针类型切片 []*LargeStruct")
	fmt.Println("  每个指针8字节，指向共享数据")
	pointerSlice := make([]*LargeStruct, 1000)
	sharedData := &LargeStruct{}
	for i := range pointerSlice {
		pointerSlice[i] = sharedData
	}

	fmt.Println("\n内存占用计算:")
	valueMemory := 80 * 1000
	pointerMemory := 8 * 1000
	fmt.Printf("  值类型切片:   80KB × 1000 = %d KB = %.1f MB\n", valueMemory, float64(valueMemory)/1024)
	fmt.Printf("  指针类型切片: 8字节 × 1000 = %d KB\n", pointerMemory/1024)
	fmt.Printf("  内存节省:     %.2f%%\n", float64(valueMemory-pointerMemory)/float64(valueMemory)*100)

	fmt.Println("\n--- 测试2: 函数调用内存分配 ---")
	fmt.Println("场景: 函数内部创建并返回结构体")

	fmt.Println("\n方式A: 返回值类型")
	fmt.Println("  每次调用在栈上分配，不产生堆内存")
	returnValue := func() LargeStruct {
		var s LargeStruct
		s.data[0] = 1
		return s
	}

	fmt.Println("\n方式B: 返回指针类型")
	fmt.Println("  每次调用在堆上分配，产生堆内存")
	returnPointer := func() *LargeStruct {
		s := &LargeStruct{}
		s.data[0] = 1
		return s
	}

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	valueStartAlloc := m1.TotalAlloc

	for i := 0; i < 100; i++ {
		_ = returnValue()
	}

	runtime.ReadMemStats(&m2)
	valueTotalAlloc := m2.TotalAlloc - valueStartAlloc

	runtime.GC()
	runtime.ReadMemStats(&m1)
	pointerStartAlloc := m1.TotalAlloc

	for i := 0; i < 100; i++ {
		_ = returnPointer()
	}

	runtime.ReadMemStats(&m2)
	pointerTotalAlloc := m2.TotalAlloc - pointerStartAlloc

	fmt.Printf("\n堆内存分配对比 (100次调用):\n")
	fmt.Printf("  返回值类型:   %d bytes\n", valueTotalAlloc)
	fmt.Printf("  返回指针类型: %d bytes\n", pointerTotalAlloc)
	if pointerTotalAlloc > valueTotalAlloc {
		fmt.Printf("  指针方式多分配: %d bytes (%.1f KB)\n", pointerTotalAlloc-valueTotalAlloc, float64(pointerTotalAlloc-valueTotalAlloc)/1024)
	}

	fmt.Println("\n--- 测试3: 结构体大小与内存对齐 ---")
	fmt.Println("Go编译器会进行内存对齐优化")

	type Aligned1 struct {
		a bool
		b int64
		c bool
	}

	type Aligned2 struct {
		b int64
		a bool
		c bool
	}

	fmt.Printf("  Aligned1 (bool, int64, bool): %d 字节\n", unsafe.Sizeof(Aligned1{}))
	fmt.Printf("  Aligned2 (int64, bool, bool): %d 字节\n", unsafe.Sizeof(Aligned2{}))
	fmt.Println("  相同字段不同顺序，大小不同！")

	fmt.Println("\n--- 内存性能总结 ---")
	fmt.Println("┌─────────────────┬────────────┬────────────────┐")
	fmt.Println("│ 场景            │ 推荐方式   │ 原因           │")
	fmt.Println("├─────────────────┼────────────┼────────────────┤")
	fmt.Println("│ 存储大量元素    │ 指针切片   │ 节省内存       │")
	fmt.Println("│ 函数内创建返回  │ 值返回     │ 避免堆分配     │")
	fmt.Println("│ 共享数据        │ 指针       │ 避免重复存储   │")
	fmt.Println("│ 只读遍历        │ 值类型     │ 缓存友好       │")
	fmt.Println("└─────────────────┴────────────┴────────────────┘")

	fmt.Println("\n内存优化建议:")
	fmt.Println("1. 大数据集合使用指针切片")
	fmt.Println("2. 频繁创建的小对象用值返回")
	fmt.Println("3. 注意结构体字段顺序对齐")
	fmt.Println("4. 使用 go build -gcflags='-m' 查看逃逸分析")
}

func PointerSafety() {
	fmt.Println("=== 指针安全性 ===")

	fmt.Println("Go指针的安全特性:")
	fmt.Println("1. 不支持指针算术(除非使用unsafe)")
	fmt.Println("2. 自动垃圾回收，无悬空指针")
	fmt.Println("3. nil指针检查")

	fmt.Println("\nunsafe包的使用场景:")
	fmt.Println("- 与C代码交互")
	fmt.Println("- 性能优化")
	fmt.Println("- 底层系统编程")

	fmt.Println("\nunsafe包的风险:")
	fmt.Println("- 可能导致程序崩溃")
	fmt.Println("- 不可移植")
	fmt.Println("- 不受Go 1兼容性保证")
}
