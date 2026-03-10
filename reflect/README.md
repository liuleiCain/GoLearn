# Reflect 反射模块

## 学习目标
深入理解Go语言反射机制，掌握动态类型检查、方法调用、值修改等高级特性。

## 核心概念

### 1. 反射基础
反射允许程序在运行时检查类型和值，动态操作对象。

```go
reflect.TypeOf(v)    // 获取类型
reflect.ValueOf(v)   // 获取值
```

### 2. Type和Value
- Type: 表示Go类型的信息
- Value: 表示Go值的运行时表示

### 3. Kind
类型的基本分类：
```go
reflect.Int
reflect.String
reflect.Struct
reflect.Ptr
reflect.Slice
reflect.Map
```

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| BasicTypeCheck | 基础类型检查 | 初级 |
| StructInspect | 结构体检查 | 初级 |
| MethodCall | 方法调用 | 中级 |
| SetValue | 通过反射设置值 | 中级 |
| CallPointerMethod | 调用指针接收者方法 | 中级 |
| DynamicStruct | 动态创建结构体 | 高级 |
| MakeSlice | 通过反射创建切片 | 中级 |
| MakeMap | 通过反射创建map | 中级 |
| MakeFunc | 通过反射创建函数 | 高级 |
| InstanceOf | 判断类型 | 初级 |
| DeepEqual | 深度相等比较 | 初级 |
| ZeroValue | 零值 | 初级 |

## 新增示例详解

### BasicTypeCheck 基础类型检查
```go
t := reflect.TypeOf(x)
v := reflect.ValueOf(x)

fmt.Println("类型:", t.Name())
fmt.Println("种类:", t.Kind())
fmt.Println("值:", v.Interface())
```

### StructInspect 结构体检查
```go
type Person struct {
    Name string `json:"name"`
    Age  int
}

for i := 0; i < t.NumField(); i++ {
    field := t.Field(i)
    fmt.Printf("字段名: %s\n", field.Name)
    fmt.Printf("标签: %s\n", field.Tag)
    fmt.Printf("json标签: %s\n", field.Tag.Get("json"))
}
```

### MethodCall 方法调用
```go
v := reflect.ValueOf(p)
method := v.MethodByName("SayHello")
args := []reflect.Value{reflect.ValueOf("Hi")}
results := method.Call(args)
```

### SetValue 通过反射设置值
```go
p := &Person{Name: "Charlie"}
v := reflect.ValueOf(p).Elem()

nameField := v.FieldByName("Name")
if nameField.CanSet() {
    nameField.SetString("David")
}
```

### DynamicStruct 动态创建结构体
```go
fields := []reflect.StructField{
    {Name: "ID", Type: reflect.TypeOf(0), Tag: `json:"id"`},
    {Name: "Title", Type: reflect.TypeOf(""), Tag: `json:"title"`},
}

typ := reflect.StructOf(fields)
inst := reflect.New(typ).Elem()
```

### MakeSlice 通过反射创建切片
```go
intType := reflect.TypeOf(0)
sliceType := reflect.SliceOf(intType)

slice := reflect.MakeSlice(sliceType, 3, 5)
for i := 0; i < slice.Len(); i++ {
    slice.Index(i).SetInt(int64(i * 10))
}
```

### MakeFunc 通过反射创建函数
```go
funcType := reflect.FuncOf(
    []reflect.Type{intType, intType},
    []reflect.Type{intType},
    false,
)

addFunc := reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {
    a := args[0].Int()
    b := args[1].Int()
    return []reflect.Value{reflect.ValueOf(int(a + b))}
})
```

## 常见陷阱

### 1. 可设置性
```go
// 错误: 不能设置非指针或不可导出字段
v := reflect.ValueOf(x)
v.SetString("hello")  // panic

// 正确: 需要指针且是可导出字段
v := reflect.ValueOf(&x).Elem()
v.SetString("hello")  // 正确
```

### 2. 类型安全
```go
// 错误: 设置错误类型
v.SetString(123)  // panic

// 正确: 类型检查
if v.Kind() == reflect.String {
    v.SetString("hello")
}
```

### 3. 性能
反射操作比直接操作慢很多，避免在性能关键路径使用。

## 最佳实践

1. **避免过度使用**: 反射降低可读性和性能
2. **类型检查**: 始终检查类型和可设置性
3. **错误处理**: 处理可能的panic
4. **考虑替代方案**: 使用接口或代码生成代替反射
5. **限制范围**: 将反射限制在小范围内

## 运行测试

```bash
go test -v ./reflect/...
```

## 参考资料

- [Go语言反射](https://go.dev/blog/laws-of-reflection)
- [reflect包文档](https://pkg.go.dev/reflect)
