package _interface

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func PolymorphismDemo() {
	fmt.Println("=== 多态示例 ===")

	shapes := []Shape{
		Rectangle{Width: 3, Height: 4},
		Circle{Radius: 5},
	}

	for i, shape := range shapes {
		fmt.Printf("图形 %d: 面积=%.2f, 周长=%.2f\n",
			i+1, shape.Area(), shape.Perimeter())
	}
}

func EmptyInterface() {
	fmt.Println("=== 空接口示例 ===")

	var any interface{}

	any = 42
	fmt.Printf("int: %v, 类型: %T\n", any, any)

	any = "hello"
	fmt.Printf("string: %v, 类型: %T\n", any, any)

	any = []int{1, 2, 3}
	fmt.Printf("slice: %v, 类型: %T\n", any, any)

	any = struct {
		Name string
		Age  int
	}{"张三", 25}
	fmt.Printf("struct: %v, 类型: %T\n", any, any)
}

func TypeAssertion() {
	fmt.Println("=== 类型断言示例 ===")

	var i interface{} = "hello"

	s, ok := i.(string)
	fmt.Printf("字符串断言: value=%s, ok=%v\n", s, ok)

	n, ok := i.(int)
	fmt.Printf("整数断言: value=%d, ok=%v\n", n, ok)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("类型断言失败(不安全方式):", r)
		}
	}()

	var x interface{} = 42
	s2 := x.(string)
	fmt.Println("这行不会执行:", s2)
}

func TypeSwitch() {
	fmt.Println("=== 类型switch示例 ===")

	checkType := func(v interface{}) {
		switch val := v.(type) {
		case int:
			fmt.Printf("int类型: %d\n", val)
		case string:
			fmt.Printf("string类型: %s\n", val)
		case bool:
			fmt.Printf("bool类型: %t\n", val)
		case []int:
			fmt.Printf("[]int类型: %v\n", val)
		case map[string]int:
			fmt.Printf("map[string]int类型: %v\n", val)
		default:
			fmt.Printf("未知类型: %T, 值: %v\n", val, val)
		}
	}

	checkType(42)
	checkType("hello")
	checkType(true)
	checkType([]int{1, 2, 3})
	checkType(map[string]int{"a": 1})
	checkType(3.14)
}

type Reader interface {
	Read() string
}

type Writer interface {
	Write(data string)
}

type ReadWriter interface {
	Reader
	Writer
}

type File struct {
	content string
}

func (f *File) Read() string {
	return f.content
}

func (f *File) Write(data string) {
	f.content = data
}

func InterfaceComposition() {
	fmt.Println("=== 接口组合示例 ===")

	file := &File{}
	var rw ReadWriter = file

	rw.Write("Hello, Interface!")
	fmt.Println("读取内容:", rw.Read())

	var r Reader = file
	fmt.Println("作为Reader读取:", r.Read())
}

func InterfaceNil() {
	fmt.Println("=== 接口nil判断示例 ===")

	var s Shape
	fmt.Printf("nil接口: %v, == nil: %v\n", s, s == nil)

	var r *Rectangle
	fmt.Printf("nil指针: %v, == nil: %v\n", r, r == nil)

	var s2 Shape = r
	fmt.Printf("包含nil指针的接口: %v, == nil: %v\n", s2, s2 == nil)

	fmt.Println("\n正确判断接口是否为nil:")
	if s2 == nil {
		fmt.Println("接口为nil")
	} else {
		fmt.Println("接口不为nil，但内部值可能为nil")
	}
}

type Stringer interface {
	String() string
}

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s (%d岁)", p.Name, p.Age)
}

func InterfaceImplementation() {
	fmt.Println("=== 隐式接口实现示例 ===")

	var s Stringer = Person{Name: "张三", Age: 25}
	fmt.Println("Stringer接口:", s.String())

	fmt.Println("\nGo语言接口是隐式实现的:")
	fmt.Println("- 不需要显式声明 implements")
	fmt.Println("- 只要类型实现了接口的所有方法即可")
}

func SmallInterface() {
	fmt.Println("=== 小接口原则示例 ===")

	fmt.Println("Go标准库中的小接口示例:")
	fmt.Println("- io.Reader: 只有Read方法")
	fmt.Println("- io.Writer: 只有Write方法")
	fmt.Println("- fmt.Stringer: 只有String方法")
	fmt.Println("- sort.Interface: Len, Less, Swap三个方法")

	fmt.Println("\n小接口的优点:")
	fmt.Println("1. 更容易实现")
	fmt.Println("2. 更灵活的组合")
	fmt.Println("3. 更好的可测试性")
}
