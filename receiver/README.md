# Receiver 方法接收者模块

## 学习目标
理解Go语言方法接收者的概念，掌握值接收者和指针接收者的区别与选择。

## 核心概念

### 1. 方法定义
Go语言的方法是带接收者的函数：
```go
func (receiver Type) MethodName(params) returns {
    // 方法体
}
```

### 2. 接收者类型
- 值接收者: `func (p Person) Method()`
- 指针接收者: `func (p *Person) Method()`

### 3. 方法集规则
| 接收者类型 | T值可调用 | *T值可调用 |
|-----------|----------|-----------|
| 值接收者 | ✓ | ✓ |
| 指针接收者 | ✗ | ✓ |

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| ReceiverTypeCompare | 值接收者vs指针接收者对比 | 初级 |
| ReceiverListTypeCompare | 切片遍历中的接收者行为 | 中级 |

## 值接收者 vs 指针接收者

### 值接收者
```go
func (p Person) SetName(name string) {
    p.Name = name  // 修改的是副本
}

p := Person{Name: "张三"}
p.SetName("李四")
fmt.Println(p.Name)  // 输出: 张三 (未改变)
```

### 指针接收者
```go
func (p *Person) SetName(name string) {
    p.Name = name  // 修改的是原值
}

p := Person{Name: "张三"}
p.SetName("李四")  // 自动取地址
fmt.Println(p.Name)  // 输出: 李四 (已改变)
```

## 选择原则

### 使用指针接收者的场景
1. **需要修改接收者**
2. **结构体较大** - 避免复制开销
3. **一致性** - 其他方法使用指针接收者
4. **实现接口** - 方法需要修改状态

### 使用值接收者的场景
1. **不需要修改接收者**
2. **小型结构体**
3. **值语义** - 希望方法不改变原值
4. **不可变对象**

## 方法集与接口

### 接口实现规则
```go
type Setter interface {
    SetName(string)
}

// 值接收者: T和*T都实现接口
func (p Person) SetName(string) {}

// 指针接收者: 只有*T实现接口
func (p *Person) SetName(string) {}
```

### 示例
```go
type Setter interface {
    SetName(string)
}

func main() {
    var s Setter
    
    // 值接收者时，两者都可以
    s = Person{}      // ✓
    s = &Person{}     // ✓
    
    // 指针接收者时，只有指针可以
    s = Person{}      // ✗ 编译错误
    s = &Person{}     // ✓
}
```

## 切片遍历中的陷阱

### 问题示例
```go
data := []Person{{"张三"}, {"李四"}, {"王五"}}
for _, v := range data {
    go v.printName()  // v是副本
}

data := []*Person{{"张三"}, {"李四"}, {"王五"}}
for _, v := range data {
    go v.printName()  // v是指针，指向原数据
}
```

## 常见陷阱

### 1. 混用接收者类型
```go
// 不推荐: 混用值和指针接收者
func (p Person) GetName() string { return p.Name }
func (p *Person) SetName(name string) { p.Name = name }

// 推荐: 统一使用指针接收者
func (p *Person) GetName() string { return p.Name }
func (p *Person) SetName(name string) { p.Name = name }
```

### 2. 值接收者修改无效
```go
func (p Person) Modify() {
    p.Name = "modified"  // 无效，修改的是副本
}
```

## 最佳实践

1. **一致性**: 同一类型的所有方法使用同一种接收者
2. **默认指针**: 不确定时使用指针接收者
3. **小型只读**: 小型、只读结构体可用值接收者
4. **避免混用**: 不要混用值和指针接收者
5. **文档说明**: 说明接收者选择的理由
