package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

// Person 示例结构体
type Person struct {
	Name    string `json:"name" db:"name_column"`
	Age     int    `json:"age" db:"age_column"`
	private bool
}

func (p Person) SayHello(greeting string) string {
	return fmt.Sprintf("%s, I'm %s, %d years old", greeting, p.Name, p.Age)
}

func (p *Person) SetName(name string) {
	p.Name = name
}

// BasicTypeCheck 基础类型检查
func BasicTypeCheck() {
	fmt.Println("=== 基础类型检查 ===")

	var x int = 42
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	fmt.Printf("类型: %s\n", t.Name())
	fmt.Printf("种类: %s\n", t.Kind())
	fmt.Printf("值: %v\n", v.Interface())
	fmt.Printf("是否可设置: %v\n", v.CanSet())
	fmt.Println()

	// 指针类型
	var p *int = &x
	tPtr := reflect.TypeOf(p)
	vPtr := reflect.ValueOf(p)

	fmt.Printf("指针类型: %s\n", tPtr)
	fmt.Printf("指针种类: %s\n", tPtr.Kind())
	fmt.Printf("解指针后的类型: %s\n", tPtr.Elem().Name())
	fmt.Printf("解指针后的种类: %s\n", tPtr.Elem().Kind())
	fmt.Printf("解指针后的值: %v\n", vPtr.Elem().Interface())
	fmt.Printf("指针解指针后是否可设置: %v\n", vPtr.Elem().CanSet())
}

// StructInspect 结构体检查
func StructInspect() {
	fmt.Println("=== 结构体检查 ===")

	p := Person{Name: "Alice", Age: 25, private: true}
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)

	fmt.Printf("类型: %s\n", t.Name())
	fmt.Printf("字段数量: %d\n", t.NumField())
	fmt.Println()

	fmt.Println("字段详情:")
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		fmt.Printf("  字段名: %s\n", field.Name)
		fmt.Printf("  类型: %s\n", field.Type)
		fmt.Printf("  可见性: %s\n", func() string {
			if strings.ToUpper(field.Name[:1]) == field.Name[:1] {
				return "导出"
			}
			return "未导出"
		}())
		fmt.Printf("  标签: %s\n", field.Tag)
		fmt.Printf("  json标签: %s\n", field.Tag.Get("json"))
		fmt.Printf("  db标签: %s\n", field.Tag.Get("db"))
		fmt.Printf("  值: %v\n", value.Interface())
		fmt.Println()
	}
}

// MethodCall 方法调用
func MethodCall() {
	fmt.Println("=== 方法调用 ===")

	p := Person{Name: "Bob", Age: 30}
	v := reflect.ValueOf(p)
	t := reflect.TypeOf(p)

	fmt.Printf("方法数量: %d\n", v.NumMethod())
	fmt.Println()

	fmt.Println("方法列表:")
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Printf("  方法名: %s\n", method.Name)
		fmt.Printf("  方法类型: %s\n", method.Type)
		fmt.Printf("  参数数量: %d\n", method.Type.NumIn()-1)
		fmt.Printf("  返回值数量: %d\n", method.Type.NumOut())
		fmt.Println()
	}

	fmt.Println("调用SayHello方法:")
	method := v.MethodByName("SayHello")
	args := []reflect.Value{reflect.ValueOf("Hi")}
	results := method.Call(args)
	fmt.Printf("结果: %s\n", results[0].String())
}

// SetValue 通过反射设置值
func SetValue() {
	fmt.Println("=== 通过反射设置值 ===")

	p := &Person{Name: "Charlie", Age: 35}
	v := reflect.ValueOf(p).Elem()

	fmt.Printf("修改前: %+v\n", p)

	nameField := v.FieldByName("Name")
	if nameField.CanSet() && nameField.Kind() == reflect.String {
		nameField.SetString("David")
	}

	ageField := v.FieldByName("Age")
	if ageField.CanSet() && ageField.Kind() == reflect.Int {
		ageField.SetInt(40)
	}

	fmt.Printf("修改后: %+v\n", p)

	// 尝试设置未导出字段
	fmt.Println()
	fmt.Println("尝试设置未导出字段private:")
	privateField := v.FieldByName("private")
	fmt.Printf("  是否可设置: %v\n", privateField.CanSet())
	if privateField.CanSet() {
		privateField.SetBool(true)
	} else {
		fmt.Println("  未导出字段无法设置")
	}
}

// CallPointerMethod 调用指针接收者方法
func CallPointerMethod() {
	fmt.Println("=== 调用指针接收者方法 ===")

	p := &Person{Name: "Eve", Age: 45}
	v := reflect.ValueOf(p)

	fmt.Printf("调用前: %+v\n", p)

	method := v.MethodByName("SetName")
	args := []reflect.Value{reflect.ValueOf("Frank")}
	method.Call(args)

	fmt.Printf("调用后: %+v\n", p)
}

// DynamicStruct 动态创建结构体
func DynamicStruct() {
	fmt.Println("=== 动态创建结构体 ===")

	fields := []reflect.StructField{
		{
			Name: "ID",
			Type: reflect.TypeOf(0),
			Tag:  reflect.StructTag(`json:"id"`),
		},
		{
			Name: "Title",
			Type: reflect.TypeOf(""),
			Tag:  reflect.StructTag(`json:"title"`),
		},
	}

	typ := reflect.StructOf(fields)
	fmt.Printf("动态创建的结构体类型: %s\n", typ)

	// 创建结构体实例
	inst := reflect.New(typ).Elem()

	idField := inst.FieldByName("ID")
	idField.SetInt(123)

	titleField := inst.FieldByName("Title")
	titleField.SetString("Hello Reflect")

	fmt.Printf("结构体实例: %+v\n", inst.Interface())

	// 序列化为JSON演示（需手动实现）
	fmt.Println()
	fmt.Println("结构体字段值:")
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := inst.Field(i)
		fmt.Printf("  %s: %v\n", field.Name, value.Interface())
	}
}

// MakeSlice 通过反射创建切片
func MakeSlice() {
	fmt.Println("=== 通过反射创建切片 ===")

	// 创建int类型切片
	intType := reflect.TypeOf(0)
	sliceType := reflect.SliceOf(intType)

	// 创建长度3、容量5的切片
	slice := reflect.MakeSlice(sliceType, 3, 5)

	fmt.Printf("切片类型: %s\n", slice.Type())
	fmt.Printf("长度: %d, 容量: %d\n", slice.Len(), slice.Cap())

	// 设置元素
	for i := 0; i < slice.Len(); i++ {
		slice.Index(i).SetInt(int64(i * 10))
	}

	fmt.Printf("设置后: %v\n", slice.Interface())

	// 追加元素
	slice = reflect.Append(slice, reflect.ValueOf(40))
	slice = reflect.Append(slice, reflect.ValueOf(50))
	fmt.Printf("追加后: %v\n", slice.Interface())
}

// MakeMap 通过反射创建map
func MakeMap() {
	fmt.Println("=== 通过反射创建map ===")

	keyType := reflect.TypeOf("")
	valueType := reflect.TypeOf(0)
	mapType := reflect.MapOf(keyType, valueType)

	m := reflect.MakeMap(mapType)

	m.SetMapIndex(reflect.ValueOf("a"), reflect.ValueOf(1))
	m.SetMapIndex(reflect.ValueOf("b"), reflect.ValueOf(2))
	m.SetMapIndex(reflect.ValueOf("c"), reflect.ValueOf(3))

	fmt.Printf("map: %v\n", m.Interface())

	// 遍历map
	fmt.Println("遍历map:")
	iter := m.MapRange()
	for iter.Next() {
		fmt.Printf("  %s: %v\n", iter.Key().Interface(), iter.Value().Interface())
	}
}

// MakeFunc 通过反射创建函数
func MakeFunc() {
	fmt.Println("=== 通过反射创建函数 ===")

	intType := reflect.TypeOf(0)
	funcType := reflect.FuncOf([]reflect.Type{intType, intType}, []reflect.Type{intType}, false)

	addFunc := reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {
		a := args[0].Int()
		b := args[1].Int()
		return []reflect.Value{reflect.ValueOf(int(a + b))}
	})

	result := addFunc.Call([]reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)})
	fmt.Printf("10 + 20 = %d\n", result[0].Int())
}

// InstanceOf 判断类型
func InstanceOf() {
	fmt.Println("=== 判断类型 ===")

	var x interface{} = Person{Name: "Grace", Age: 50}

	t := reflect.TypeOf(x)
	fmt.Printf("类型: %s\n", t.Name())

	// 类型断言
	if p, ok := x.(Person); ok {
		fmt.Printf("类型断言成功: %+v\n", p)
	}

	// 判断是否实现接口
	var iface interface{} = (*fmt.Stringer)(nil)
	ifaceType := reflect.TypeOf(iface).Elem()
	fmt.Printf("\n是否实现fmt.Stringer: %v\n", t.Implements(ifaceType))
}

// DeepEqual 深度相等比较
func DeepEqual() {
	fmt.Println("=== 深度相等比较 ===")

	p1 := Person{Name: "Henry", Age: 60}
	p2 := Person{Name: "Henry", Age: 60}
	p3 := Person{Name: "Ivy", Age: 70}

	fmt.Printf("p1 & p2 深度相等: %v\n", reflect.DeepEqual(p1, p2))
	fmt.Printf("p1 & p3 深度相等: %v\n", reflect.DeepEqual(p1, p3))

	// 切片比较
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s3 := []int{3, 2, 1}
	fmt.Printf("\ns1 & s2 深度相等: %v\n", reflect.DeepEqual(s1, s2))
	fmt.Printf("s1 & s3 深度相等: %v\n", reflect.DeepEqual(s1, s3))
}

// ZeroValue 零值
func ZeroValue() {
	fmt.Println("=== 零值 ===")

	intType := reflect.TypeOf(0)
	stringType := reflect.TypeOf("")
	boolType := reflect.TypeOf(false)
	personType := reflect.TypeOf(Person{})

	fmt.Printf("int零值: %v\n", reflect.Zero(intType).Interface())
	fmt.Printf("string零值: %v\n", reflect.Zero(stringType).Interface())
	fmt.Printf("bool零值: %v\n", reflect.Zero(boolType).Interface())
	fmt.Printf("Person零值: %+v\n", reflect.Zero(personType).Interface())
}
