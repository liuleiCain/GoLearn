package generics

import (
	"cmp"
	"fmt"
	"strings"
)

func Print[T any](value T) {
	fmt.Printf("值: %v, 类型: %T\n", value, value)
}

func BasicGeneric() {
	fmt.Println("=== 基础泛型函数示例 ===")

	Print(42)
	Print("hello")
	Print(3.14)
	Print([]int{1, 2, 3})
}

func MapSlice[T, U any](s []T, f func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

func Filter[T any](s []T, f func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func GenericSlice() {
	fmt.Println("=== 泛型切片函数示例 ===")

	numbers := []int{1, 2, 3, 4, 5}
	squares := MapSlice(numbers, func(n int) int {
		return n * n
	})
	fmt.Printf("原切片: %v\n", numbers)
	fmt.Printf("平方后: %v\n", squares)

	evens := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("偶数: %v\n", evens)
}

func PrintAny[T any](v T) {
	fmt.Printf("  %v (%T)\n", v, v)
}

func IndexOf[T comparable](s []T, x T) int {
	for i, v := range s {
		if v == x {
			return i
		}
	}
	return -1
}

func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func TypeConstraint() {
	fmt.Println("=== 类型约束示例 ===")

	fmt.Println("any约束 - 接受任何类型:")
	PrintAny(42)
	PrintAny("hello")

	fmt.Println("\ncomparable约束 - 只接受可比较类型:")
	fmt.Printf("  indexOf([1,2,3], 2) = %d\n", IndexOf([]int{1, 2, 3}, 2))
	fmt.Printf("  indexOf(['a','b','c'], 'd') = %d\n", IndexOf([]string{"a", "b", "c"}, "d"))

	fmt.Println("\nOrdered约束 (cmp.Ordered) - 只接受有序类型:")
	fmt.Printf("  min(3, 5) = %v\n", Min(3, 5))
	fmt.Printf("  min(3.14, 2.71) = %v\n", Min(3.14, 2.71))
	fmt.Printf("  min('a', 'b') = %v\n", Min('a', 'b'))
}

type Number interface {
	int | int64 | float64 | float32
}

func Sum[T Number](s []T) T {
	var total T
	for _, v := range s {
		total += v
	}
	return total
}

func GenericSum() {
	fmt.Println("=== 泛型求和示例 ===")

	ints := []int{1, 2, 3, 4, 5}
	floats := []float64{1.1, 2.2, 3.3}

	fmt.Printf("int求和: %v\n", Sum(ints))
	fmt.Printf("float求和: %v\n", Sum(floats))
}

type Stack[T any] struct {
	elements []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{elements: make([]T, 0)}
}

func (s *Stack[T]) Push(v T) {
	s.elements = append(s.elements, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.elements) == 0 {
		return zero, false
	}
	index := len(s.elements) - 1
	element := s.elements[index]
	s.elements = s.elements[:index]
	return element, true
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.elements) == 0 {
		return zero, false
	}
	return s.elements[len(s.elements)-1], true
}

func (s *Stack[T]) Size() int {
	return len(s.elements)
}

func GenericType() {
	fmt.Println("=== 泛型类型示例 ===")

	intStack := NewStack[int]()
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)

	fmt.Printf("栈大小: %d\n", intStack.Size())
	if v, ok := intStack.Peek(); ok {
		fmt.Printf("栈顶元素: %v\n", v)
	}
	if v, ok := intStack.Pop(); ok {
		fmt.Printf("弹出元素: %v\n", v)
	}

	strStack := NewStack[string]()
	strStack.Push("hello")
	strStack.Push("world")
	if v, ok := strStack.Pop(); ok {
		fmt.Printf("字符串栈弹出: %v\n", v)
	}
}

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func (p Pair[K, V]) String() string {
	return fmt.Sprintf("(%v: %v)", p.Key, p.Value)
}

func MultiTypeParam() {
	fmt.Println("=== 多类型参数示例 ===")

	p1 := Pair[string, int]{Key: "age", Value: 25}
	p2 := Pair[int, string]{Key: 1, Value: "first"}

	fmt.Printf("p1: %v\n", p1)
	fmt.Printf("p2: %v\n", p2)
}

type MyInt int

func (m MyInt) String() string {
	return fmt.Sprintf("MyInt(%d)", m)
}

type Stringer interface {
	String() string
}

func PrintStringer[T Stringer](v T) {
	fmt.Printf("Stringer: %s\n", v.String())
}

func GenericMethod() {
	fmt.Println("=== 泛型方法示例 ===")

	var m MyInt = 42
	PrintStringer(m)
}

type Processor[T any] interface {
	Process(T) T
}

type UpperProcessor struct{}

func (p UpperProcessor) Process(s string) string {
	return strings.ToUpper(s)
}

type DoubleProcessor struct{}

func (p DoubleProcessor) Process(n int) int {
	return n * 2
}

func GenericInterface() {
	fmt.Println("=== 泛型接口示例 ===")

	var strProcessor Processor[string] = UpperProcessor{}
	fmt.Printf("处理字符串: %v\n", strProcessor.Process("hello"))

	var intProcessor Processor[int] = DoubleProcessor{}
	fmt.Printf("处理整数: %v\n", intProcessor.Process(21))
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Float interface {
	~float32 | ~float64
}

type Integer interface {
	Signed | Unsigned
}

type Numeric interface {
	Integer | Float
}

func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Clamp[T cmp.Ordered](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ConstraintUnion() {
	fmt.Println("=== 类型联合约束示例 ===")

	fmt.Printf("Max(3, 5) = %v\n", Max(3, 5))
	fmt.Printf("Max(3.14, 2.71) = %v\n", Max(3.14, 2.71))
	fmt.Printf("Max('a', 'z') = %v\n", Max('a', 'z'))

	fmt.Printf("Clamp(50, 0, 100) = %v\n", Clamp(50, 0, 100))
	fmt.Printf("Clamp(-10, 0, 100) = %v\n", Clamp(-10, 0, 100))
	fmt.Printf("Clamp(150, 0, 100) = %v\n", Clamp(150, 0, 100))
}
