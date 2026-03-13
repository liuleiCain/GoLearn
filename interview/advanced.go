package interview

import (
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"unsafe"
)

// ============================================================
// 底层原理面试题
// ============================================================

// Q1: Go语言的内存模型是什么？
// A1: Go内存模型核心概念：
// 1. Happens-Before关系：定义操作之间的可见性顺序
// 2. 原子操作：sync/atomic包提供的原子操作
// 3. 同步原语：channel、mutex、WaitGroup等
// 4. 内存屏障：确保内存操作的顺序性
// 5. 数据竞争：多个goroutine同时访问同一变量且至少一个写入

// Q2: Go语言的垃圾回收机制？
// A2: Go的GC机制：
// 1. 三色标记-清除算法
// 2. 并发标记，与用户程序并行执行
// 3. 写屏障技术保证正确性
// 4. STW（Stop-The-World）时间很短
// 5. GOGC环境变量控制GC触发频率

func GCMechanism() {
	fmt.Println("=== GC机制 ===")

	// 获取GC统计信息
	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	fmt.Printf("GC次数: %d\n", stats.NumGC)
	fmt.Printf("暂停总时间: %v\n", stats.PauseTotal)

	// 内存分配统计
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("堆分配: %d bytes\n", m.HeapAlloc)
	fmt.Printf("系统分配: %d bytes\n", m.Sys)
	fmt.Printf("GC次数: %d\n", m.NumGC)

	// 手动触发GC
	runtime.GC()
	fmt.Println("手动触发GC完成")

	// 调整GOGC
	old := debug.SetGCPercent(100)
	fmt.Printf("GOGC: %d\n", old)
}

// Q3: Go语言的调度器原理？
// A3: Go调度器（GMP模型）：
// 1. G（Goroutine）：协程，包含栈、指令指针等信息
// 2. M（Machine）：系统线程，执行G
// 3. P（Processor）：逻辑处理器，包含运行队列
// 4. 调度策略：工作窃取、系统调用时切换
// 5. 抢占式调度：基于信号的抢占

func SchedulerPrinciple() {
	fmt.Println("\n=== 调度器原理 ===")

	// 获取goroutine数量
	fmt.Printf("当前goroutine数量: %d\n", runtime.NumGoroutine())

	// 获取CPU核心数
	fmt.Printf("CPU核心数: %d\n", runtime.NumCPU())

	// 获取GOMAXPROCS
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))

	// 让出CPU
	runtime.Gosched()
	fmt.Println("Gosched()让出CPU")
}

// Q4: Go语言的逃逸分析？
// A4: 逃逸分析决定变量分配在栈还是堆：
// 1. 栈分配：函数返回后不再使用的变量
// 2. 堆分配：变量在函数返回后仍被引用
// 3. 逃逸场景：返回局部变量指针、闭包捕获、接口转换
// 4. 使用 go build -gcflags="-m" 查看逃逸分析结果

func EscapeAnalysisDemo() {
	fmt.Println("\n=== 逃逸分析 ===")

	// 不逃逸：分配在栈上
	noEscape := func() int {
		x := 10
		return x // 返回值，不逃逸
	}
	fmt.Printf("不逃逸: %d\n", noEscape())

	// 逃逸：分配在堆上
	escape := func() *int {
		x := 10
		return &x // 返回指针，逃逸到堆
	}
	fmt.Printf("逃逸: %d\n", *escape())

	// 闭包捕获：逃逸
	closure := func() func() int {
		x := 10
		return func() int {
			return x // 闭包捕获，逃逸
		}
	}
	fn := closure()
	fmt.Printf("闭包捕获: %d\n", fn())

	// 接口转换：逃逸
	interfaceEscape := func() interface{} {
		x := 10
		return x // 接口转换，逃逸
	}
	fmt.Printf("接口转换: %v\n", interfaceEscape())
}

// Q5: Go语言的slice底层原理？
// A5: slice底层结构：
// 1. slice包含三个字段：指针、长度、容量
// 2. 指针指向底层数组
// 3. 长度表示元素个数
// 4. 容量表示底层数组的大小
// 5. 扩容时会分配新数组并复制数据

func SliceInternal() {
	fmt.Println("\n=== Slice底层原理 ===")

	// slice结构
	s := make([]int, 3, 5)
	fmt.Printf("slice: len=%d, cap=%d\n", len(s), cap(s))

	// 获取slice的底层数组指针
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Printf("底层数组指针: %x\n", hdr.Data)
	fmt.Printf("长度: %d\n", hdr.Len)
	fmt.Printf("容量: %d\n", hdr.Cap)

	// 多个slice共享底层数组
	s1 := s[0:2]
	s2 := s[1:3]
	s1[1] = 100
	fmt.Printf("s1=%v, s2=%v (共享底层数组)\n", s1, s2)

	// append可能导致底层数组变化
	oldPtr := hdr.Data
	s = append(s, 1, 2, 3) // 超过容量，重新分配
	hdr = (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Printf("扩容前指针: %x, 扩容后指针: %x\n", oldPtr, hdr.Data)
}

// Q6: Go语言的interface底层原理？
// A6: interface底层实现：
// 1. 空接口（interface{}）：包含类型指针和数据指针
// 2. 非空接口：包含itab（接口表）和数据指针
// 3. itab包含接口的类型信息和方法表
// 4. 类型断言通过itab实现
// 5. nil接口和nil值接口的区别

func InterfaceInternal() {
	fmt.Println("\n=== Interface底层原理 ===")

	// 空接口结构
	type eface struct {
		_type unsafe.Pointer
		data  unsafe.Pointer
	}

	// 非空接口结构
	type iface struct {
		tab  unsafe.Pointer
		data unsafe.Pointer
	}

	// 空接口
	var i interface{} = 42
	fmt.Printf("空接口: type=%T, value=%v\n", i, i)

	// 非空接口
	var r fmt.Stringer = myString("hello")
	fmt.Printf("非空接口: %s\n", r.String())

	// nil接口 vs nil值接口
	var p *int
	var e interface{} = p
	fmt.Printf("p==nil: %v, e==nil: %v\n", p == nil, e == nil)
}

type myString string

func (s myString) String() string {
	return string(s)
}

// Q7: Go语言的map底层原理？
// A7: map底层实现：
// 1. 哈希表实现，使用链地址法解决冲突
// 2. bucket（桶）存储键值对
// 3. 扩容时渐进式迁移
// 4. 不支持并发读写
// 5. 遍历顺序随机

func MapInternal() {
	fmt.Println("\n=== Map底层原理 ===")

	// map结构
	m := make(map[int]int, 8)
	for i := 0; i < 10; i++ {
		m[i] = i * i
	}

	// 获取map类型信息
	t := reflect.TypeOf(m)
	fmt.Printf("map类型: %v\n", t)
	fmt.Printf("key类型: %v\n", t.Key())
	fmt.Printf("value类型: %v\n", t.Elem())

	// 遍历顺序随机
	fmt.Print("第一次遍历: ")
	for k := range m {
		fmt.Printf("%d ", k)
	}
	fmt.Println()

	fmt.Print("第二次遍历: ")
	for k := range m {
		fmt.Printf("%d ", k)
	}
	fmt.Println()
}

// Q8: Go语言的channel底层原理？
// A8: channel底层结构（hchan）：
// 1. 循环缓冲区：存储数据
// 2. 发送等待队列：recvq
// 3. 接收等待队列：sendq
// 4. 互斥锁：保护共享数据
// 5. 关闭标志：closed

func ChannelInternal() {
	fmt.Println("\n=== Channel底层原理 ===")

	// 无缓冲channel
	ch1 := make(chan int)
	fmt.Printf("无缓冲channel: len=%d, cap=%d\n", len(ch1), cap(ch1))

	// 有缓冲channel
	ch2 := make(chan int, 3)
	ch2 <- 1
	ch2 <- 2
	fmt.Printf("有缓冲channel: len=%d, cap=%d\n", len(ch2), cap(ch2))

	// channel的内部状态
	t := reflect.TypeOf(ch2)
	fmt.Printf("channel类型: %v\n", t)
	fmt.Printf("元素类型: %v\n", t.Elem())
}

// Q9: unsafe包的使用场景？
// A9: unsafe包的使用场景：
// 1. 指针类型转换
// 2. 获取结构体字段偏移量
// 3. 修改私有字段
// 4. 零拷贝类型转换
// 5. 与C语言交互

func UnsafeUsage() {
	fmt.Println("\n=== unsafe包使用 ===")

	// Sizeof: 获取大小
	fmt.Printf("int大小: %d\n", unsafe.Sizeof(int(0)))
	fmt.Printf("int64大小: %d\n", unsafe.Sizeof(int64(0)))
	fmt.Printf("string大小: %d\n", unsafe.Sizeof("hello"))

	// Offsetof: 获取字段偏移
	type Person struct {
		Name string
		Age  int
	}
	p := Person{Name: "Alice", Age: 30}
	nameOffset := unsafe.Offsetof(p.Name)
	ageOffset := unsafe.Offsetof(p.Age)
	fmt.Printf("Name偏移: %d, Age偏移: %d\n", nameOffset, ageOffset)

	// Alignof: 获取对齐值
	fmt.Printf("int对齐: %d\n", unsafe.Alignof(int(0)))
	fmt.Printf("string对齐: %d\n", unsafe.Alignof("hello"))

	// 指针转换
	s := "hello"
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	fmt.Printf("string指针: %x, 长度: %d\n", hdr.Data, hdr.Len)
}

// Q10: Go语言的反射机制？
// A10: 反射机制：
// 1. reflect.TypeOf: 获取类型信息
// 2. reflect.ValueOf: 获取值信息
// 3. 通过反射可以修改变量值
// 4. 通过反射可以调用方法
// 5. 反射性能较低，谨慎使用

func ReflectionDemo() {
	fmt.Println("\n=== 反射机制 ===")

	// TypeOf
	x := 42
	t := reflect.TypeOf(x)
	fmt.Printf("类型: %v, 名称: %v\n", t, t.Name())

	// ValueOf
	v := reflect.ValueOf(x)
	fmt.Printf("值: %v, 类型: %v\n", v, v.Type())

	// 修改值（需要传指针）
	ptr := reflect.ValueOf(&x)
	ptr.Elem().SetInt(100)
	fmt.Printf("修改后: %d\n", x)

	// 结构体反射
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	user := User{Name: "Alice", Age: 30}
	t = reflect.TypeOf(user)
	v = reflect.ValueOf(user)

	// 遍历字段
	fmt.Println("结构体字段:")
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fmt.Printf("  %s: %v (tag: %s)\n", field.Name, value, field.Tag.Get("json"))
	}
}

// Q11: Go语言的内存对齐？
// A11: 内存对齐规则：
// 1. 成员对齐：每个成员的偏移量是其大小的整数倍
// 2. 整体对齐：结构体大小是最大成员对齐值的整数倍
// 3. 空结构体大小为0
// 4. 合理排列字段可以减少内存占用

func MemoryAlignment() {
	fmt.Println("\n=== 内存对齐 ===")

	// 不合理的字段顺序
	type Bad struct {
		a bool   // 1 byte + 7 padding
		b int64  // 8 bytes
		c bool   // 1 byte + 7 padding
	}

	// 合理的字段顺序
	type Good struct {
		b int64  // 8 bytes
		a bool   // 1 byte
		c bool   // 1 byte + 6 padding
	}

	fmt.Printf("Bad结构体大小: %d\n", unsafe.Sizeof(Bad{}))
	fmt.Printf("Good结构体大小: %d\n", unsafe.Sizeof(Good{}))

	// 空结构体
	type Empty struct{}
	fmt.Printf("空结构体大小: %d\n", unsafe.Sizeof(Empty{}))

	// 空结构体作为字段
	type WithEmpty struct {
		x int
		_ Empty
		y int
	}
	fmt.Printf("含空结构体的结构体大小: %d\n", unsafe.Sizeof(WithEmpty{}))
}

// Q12: Go语言的闭包原理？
// A12: 闭包原理：
// 1. 闭包捕获外部变量，形成引用
// 2. 捕获的是变量本身，不是值
// 3. 循环中的闭包陷阱
// 4. 闭包会导致变量逃逸到堆

func ClosurePrinciple() {
	fmt.Println("\n=== 闭包原理 ===")

	// 闭包捕获变量
	x := 10
	fn := func() int {
		x++ // 捕获外部变量
		return x
	}
	fmt.Printf("第一次调用: %d\n", fn())
	fmt.Printf("第二次调用: %d\n", fn())
	fmt.Printf("外部x: %d\n", x) // x已被修改

	// 循环中的闭包陷阱
	// 注意：Go 1.22+ 已修复此问题，循环变量每次迭代都会创建新副本
	// Go 1.22之前版本会打印: i=3, i=3, i=3
	// Go 1.22+版本会打印: i=0, i=1, i=2
	fmt.Println("\n闭包陷阱（Go 1.22+已修复）:")
	var funcs []func()
	for i := 0; i < 3; i++ {
		funcs = append(funcs, func() {
			fmt.Printf("  i=%d\n", i) // Go 1.22+每次迭代创建新的i
		})
	}
	for _, f := range funcs {
		f()
	}

	// 正确做法：传参
	fmt.Println("\n正确做法:")
	funcs = nil
	for i := 0; i < 3; i++ {
		funcs = append(funcs, func(n int) func() {
			return func() {
				fmt.Printf("  n=%d\n", n)
			}
		}(i))
	}
	for _, f := range funcs {
		f()
	}
}
