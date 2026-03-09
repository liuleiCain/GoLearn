# Interface 接口模块

## 学习目标
理解Go语言接口的核心概念、多态实现、类型断言和接口设计原则。

## 核心概念

### 1. 接口定义
接口是一组方法的集合，定义了对象的行为：
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

### 2. 隐式实现
Go语言中接口是隐式实现的，不需要显式声明`implements`：
```go
type File struct{}

// File隐式实现了Reader接口
func (f *File) Read(p []byte) (n int, err error) {
    // 实现
}
```

### 3. 空接口
`interface{}`或`any`(Go 1.18+)可以接受任意类型的值。

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| PolymorphismDemo | 多态示例 - 不同类型实现同一接口 | 初级 |
| EmptyInterface | 空接口接受任意类型 | 初级 |
| TypeAssertion | 类型断言 - 安全和不安全方式 | 初级 |
| TypeSwitch | 类型switch - 处理多种类型 | 初级 |
| InterfaceComposition | 接口组合 - 小接口组合成大接口 | 中级 |
| InterfaceNil | 接口nil判断陷阱 | 中级 |
| InterfaceImplementation | 隐式接口实现 | 初级 |
| SmallInterface | 小接口设计原则 | 中级 |

## 接口内部结构

接口值由两部分组成：
- 类型信息 (type)
- 值信息 (value)

```
interface{} = (type, value)
```

### nil接口 vs 包含nil指针的接口
```go
var s Shape          // s == nil (type=nil, value=nil)
var r *Rectangle     // r == nil
var s2 Shape = r     // s2 != nil (type=*Rectangle, value=nil)
```

## 类型断言

### 安全方式 (推荐)
```go
value, ok := i.(Type)
if ok {
    // 类型断言成功
}
```

### 不安全方式 (可能panic)
```go
value := i.(Type)  // 如果断言失败会panic
```

### 类型switch
```go
switch v := i.(type) {
case int:
    fmt.Println("int:", v)
case string:
    fmt.Println("string:", v)
default:
    fmt.Println("unknown:", v)
}
```

## 常见陷阱

### 1. 接口nil判断
```go
// 错误示例
var r *Rectangle  // nil
var s Shape = r   // s不是nil!

// 正确判断
if s == nil {
    // 只有当type和value都为nil时才为true
}
```

### 2. 值接收者 vs 指针接收者
```go
type Setter interface {
    Set(string)
}

type Foo struct {
    value string
}

// 值接收者 - 修改的是副本
func (f Foo) Set(s string) {
    f.value = s  // 不会影响原值
}

// 指针接收者 - 修改的是原值
func (f *Foo) Set(s string) {
    f.value = s  // 会影响原值
}
```

## 最佳实践

1. **小接口原则**: 接口应该尽可能小，只包含必要的方法
2. **接受接口，返回结构体**: 函数参数使用接口，返回值使用具体类型
3. **接口在使用方定义**: 谁使用谁定义，而不是在实现方定义
4. **避免空接口**: 尽量使用具体类型或泛型，而不是`interface{}`
5. **注意nil判断**: 理解接口的nil语义

## 标准库常用接口

| 接口 | 方法 | 用途 |
|------|------|------|
| io.Reader | Read | 读取数据 |
| io.Writer | Write | 写入数据 |
| fmt.Stringer | String | 字符串表示 |
| sort.Interface | Len/Less/Swap | 排序 |
| error | Error | 错误处理 |
| http.Handler | ServeHTTP | HTTP处理 |
