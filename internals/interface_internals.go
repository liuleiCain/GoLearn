package internals

import (
	"fmt"
	"reflect"
	"unsafe"
)

// ============================================================
// Interface 底层实现原理
// ============================================================

// eface 空接口 (interface{}) 的内部表示
type eface struct {
	_type *_type    // 类型信息
	data  unsafe.Pointer // 数据指针
}

// iface 非空接口的内部表示
type iface struct {
	tab  *itab      // 接口表
	data unsafe.Pointer // 数据指针
}

// _type 类型信息
type _type struct {
	size       uintptr // 类型大小
	ptrdata    uintptr // 包含指针的内存大小
	hash       uint32  // 类型hash
	tflag      uint8   // 类型标志
	align      uint8   // 对齐
	fieldAlign uint8   // 字段对齐
	kind       uint8   // 类型种类
	equal      func(unsafe.Pointer, unsafe.Pointer) bool // 比较函数
	gcdata     *byte   // GC数据
	str        int32   // 类型名称字符串偏移
	ptrToThis  int32   // 指向此类型的指针
}

// itab 接口表
type itab struct {
	inter *interfacetype // 接口类型
	_type *_type         // 具体类型
	hash  uint32         // 类型hash (copy of _type.hash)
	_     [4]byte        // 填充
	fun   [1]uintptr     // 方法表
}

// interfacetype 接口类型信息
type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod // 接口方法
}

type name struct{}

type imethod struct{}

// DemoInterfaceStructure 演示interface的内存结构
func DemoInterfaceStructure() {
	fmt.Println("=== Interface内存结构 ===")

	// 空接口
	var i interface{} = 42
	fmt.Printf("空接口: %v (type: %T)\n", i, i)

	// 非空接口
	var s fmt.Stringer = myString("hello")
	fmt.Printf("非空接口: %v\n", s)

	fmt.Println("\n空接口 (eface) 结构:")
	fmt.Println("+----------------+")
	fmt.Println("| _type (*_type) | 类型信息指针")
	fmt.Println("| data (unsafe.P)| 数据指针")
	fmt.Println("+----------------+")

	fmt.Println("\n非空接口 (iface) 结构:")
	fmt.Println("+----------------+")
	fmt.Println("| tab (*itab)    | 接口表指针")
	fmt.Println("| data (unsafe.P)| 数据指针")
	fmt.Println("+----------------+")

	fmt.Println("\nitab 结构:")
	fmt.Println("+---------------------+")
	fmt.Println("| inter (*interfacetype)| 接口类型")
	fmt.Println("| _type (*_type)      | 具体类型")
	fmt.Println("| hash (uint32)       | 类型hash")
	fmt.Println("| fun [1]uintptr      | 方法表")
	fmt.Println("+---------------------+")
}

type myString string

func (s myString) String() string {
	return string(s)
}

// DemoInterfaceType 演示接口的类型信息
func DemoInterfaceType() {
	fmt.Println("\n=== 接口类型信息 ===")

	// 空接口存储不同类型
	var i interface{}

	i = 42
	fmt.Printf("存储int: %v, 类型: %T\n", i, i)

	i = "hello"
	fmt.Printf("存储string: %v, 类型: %T\n", i, i)

	i = []int{1, 2, 3}
	fmt.Printf("存储slice: %v, 类型: %T\n", i, i)

	// _type结构包含的信息
	fmt.Println("\n_type结构包含:")
	fmt.Println("  size: 类型大小")
	fmt.Println("  hash: 类型hash值")
	fmt.Println("  kind: 类型种类 (int, string, struct...)")
	fmt.Println("  equal: 比较函数")
	fmt.Println("  gcdata: GC扫描信息")
}

// DemoInterfaceNil 演示nil接口与nil值接口的区别
func DemoInterfaceNil() {
	fmt.Println("\n=== nil接口 vs nil值接口 ===")

	// nil接口
	var nilInterface interface{}
	fmt.Printf("nil接口: %v, nil=%v\n", nilInterface, nilInterface == nil)

	// nil值接口 (接口包含类型信息)
	var p *int
	var nilValueInterface interface{} = p
	fmt.Printf("nil值接口: %v, nil=%v\n", nilValueInterface, nilValueInterface == nil)

	// 查看内部结构
	fmt.Println("\n内部结构对比:")
	fmt.Println("nil接口:")
	fmt.Println("  eface{_type: nil, data: nil}")
	fmt.Println("")
	fmt.Println("nil值接口:")
	fmt.Println("  eface{_type: *int, data: nil}")
	fmt.Println("  类型指针不为nil，所以接口不等于nil")

	// 正确的nil判断
	fmt.Println("\n正确的nil判断:")
	fmt.Println("  if p == nil && nilValueInterface == nil {")
	fmt.Println("      // 两个条件都满足才是真正的nil")
	fmt.Println("  }")

	// 使用反射判断
	fmt.Println("\n使用反射判断:")
	v := reflect.ValueOf(nilValueInterface)
	fmt.Printf("  reflect.Value.IsNil(): %v\n", v.IsNil())
	fmt.Printf("  reflect.Value.IsValid(): %v\n", v.IsValid())
}

// DemoInterfaceAssert 演示类型断言的实现
func DemoInterfaceAssert() {
	fmt.Println("\n=== 类型断言实现 ===")

	var i interface{} = 42

	// 类型断言
	v, ok := i.(int)
	fmt.Printf("i.(int): v=%d, ok=%v\n", v, ok)

	v2, ok2 := i.(string)
	fmt.Printf("i.(string): v=%v, ok=%v\n", v2, ok2)

	fmt.Println("\n类型断言实现原理:")
	fmt.Println("1. 检查eface._type是否等于目标类型")
	fmt.Println("2. 如果相等，返回data指针")
	fmt.Println("3. 如果不等，返回零值和false")

	fmt.Println("\n实现伪代码:")
	fmt.Println("func assertI2I(t *_type, i interface{}) (r interface{}) {")
	fmt.Println("    e := (*eface)(unsafe.Pointer(&i))")
	fmt.Println("    if e._type == t {")
	fmt.Println("        r = e.data")
	fmt.Println("    }")
	fmt.Println("    return")
	fmt.Println("}")
}

// DemoInterfaceSwitch 演示类型switch的实现
func DemoInterfaceSwitch() {
	fmt.Println("\n=== 类型switch实现 ===")

	var i interface{} = "hello"

	switch v := i.(type) {
	case int:
		fmt.Printf("int: %d\n", v)
	case string:
		fmt.Printf("string: %s\n", v)
	case []int:
		fmt.Printf("[]int: %v\n", v)
	default:
		fmt.Printf("unknown: %T\n", v)
	}

	fmt.Println("\n类型switch实现原理:")
	fmt.Println("1. 编译器生成类型比较代码")
	fmt.Println("2. 按顺序比较_type.hash")
	fmt.Println("3. 匹配后执行对应分支")

	fmt.Println("\n编译后伪代码:")
	fmt.Println("switch e._type.kind {")
	fmt.Println("case kindInt:")
	fmt.Println("    v := *(*int)(e.data)")
	fmt.Println("case kindString:")
	fmt.Println("    v := *(*string)(e.data)")
	fmt.Println("}")
}

// DemoInterfaceMethod 演示接口方法调用
func DemoInterfaceMethod() {
	fmt.Println("\n=== 接口方法调用 ===")

	var s fmt.Stringer = myString("hello")
	result := s.String()
	fmt.Printf("调用String(): %s\n", result)

	fmt.Println("\n方法调用实现原理:")
	fmt.Println("1. 从itab.fun获取方法地址")
	fmt.Println("2. 将data作为接收者参数")
	fmt.Println("3. 调用方法")

	fmt.Println("\nitab.fun数组:")
	fmt.Println("  fun[0]: 第一个方法的地址")
	fmt.Println("  fun[1]: 第二个方法的地址")
	fmt.Println("  ...")

	fmt.Println("\n调用伪代码:")
	fmt.Println("func callMethod(i iface, methodIdx int, args ...) {")
	fmt.Println("    fn := i.tab.fun[methodIdx]")
	fmt.Println("    receiver := i.data")
	fmt.Println("    fn(receiver, args...)")
	fmt.Println("}")
}

// DemoInterfaceConversion 演示接口转换
func DemoInterfaceConversion() {
	fmt.Println("\n=== 接口转换 ===")

	fmt.Println("接口转换规则:")
	fmt.Println("  大接口 -> 小接口: 可以")
	fmt.Println("  小接口 -> 大接口: 需要类型断言")

	fmt.Println("\n转换实现:")
	fmt.Println("1. 检查原接口是否实现了目标接口")
	fmt.Println("2. 查找或创建新的itab")
	fmt.Println("3. 复制data指针")
}

// DemoInterfaceBoxing 演示接口装箱
func DemoInterfaceBoxing() {
	fmt.Println("\n=== 接口装箱 ===")

	// 值类型装箱
	x := 42
	var i interface{} = x
	fmt.Printf("值类型装箱: %v\n", i)

	// 指针类型装箱
	p := &x
	var j interface{} = p
	fmt.Printf("指针类型装箱: %v\n", j)

	fmt.Println("\n装箱过程:")
	fmt.Println("1. 分配内存存储值")
	fmt.Println("2. 设置_type指向类型信息")
	fmt.Println("3. 设置data指向值")

	fmt.Println("\n值类型 vs 指针类型:")
	fmt.Println("  值类型: 复制值到新内存")
	fmt.Println("  指针类型: 只复制指针")

	// 修改原值
	*p = 100
	fmt.Printf("修改后: i=%v, j=%v\n", i, j)
}

// DemoInterfaceCompare 演示接口比较
func DemoInterfaceCompare() {
	fmt.Println("\n=== 接口比较 ===")

	// 可比较的类型
	var i1 interface{} = 42
	var i2 interface{} = 42
	fmt.Printf("42 == 42: %v\n", i1 == i2)

	// 不可比较的类型
	// var s1 interface{} = []int{1, 2, 3}
	// var s2 interface{} = []int{1, 2, 3}
	// fmt.Printf("slice == slice: %v\n", s1 == s2) // panic!

	fmt.Println("\n比较规则:")
	fmt.Println("1. 类型必须相同")
	fmt.Println("2. 类型必须可比较")
	fmt.Println("3. 调用_type.equal函数比较data")

	fmt.Println("\n可比较类型:")
	fmt.Println("  ✓ 基本类型: int, float, string, bool")
	fmt.Println("  ✓ 指针")
	fmt.Println("  ✓ 数组 (元素可比较)")
	fmt.Println("  ✓ 结构体 (字段可比较)")

	fmt.Println("\n不可比较类型:")
	fmt.Println("  ✗ 切片")
	fmt.Println("  ✗ map")
	fmt.Println("  ✗ 函数")

	fmt.Println("\n比较实现伪代码:")
	fmt.Println("func equal(e1, e2 eface) bool {")
	fmt.Println("    if e1._type != e2._type {")
	fmt.Println("        return false")
	fmt.Println("    }")
	fmt.Println("    return e1._type.equal(e1.data, e2.data)")
	fmt.Println("}")
}

// DemoInterfaceEscape 演示接口导致的逃逸
func DemoInterfaceEscape() {
	fmt.Println("\n=== 接口与逃逸分析 ===")

	fmt.Println("接口装箱会导致逃逸:")
	fmt.Println("  值类型 -> 接口: 逃逸到堆")
	fmt.Println("  原因: 接口需要存储值的指针")

	fmt.Println("\n示例:")
	fmt.Println("func escape() {")
	fmt.Println("    x := 42        // 栈上")
	fmt.Println("    var i interface{} = x  // x逃逸到堆")
	fmt.Println("}")

	fmt.Println("\n避免逃逸:")
	fmt.Println("  1. 使用指针: var i interface{} = &x")
	fmt.Println("  2. 避免不必要的接口转换")
	fmt.Println("  3. 使用泛型替代接口 (Go 1.18+)")

	fmt.Println("\n逃逸分析命令:")
	fmt.Println("  go build -gcflags='-m'")
}
