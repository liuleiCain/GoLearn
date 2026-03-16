# Go 设计模式

本模块详细介绍 Go 语言中常用的设计模式，包含创建型、结构型和行为型三大类共 19 种设计模式。

## 目录

- [设计模式概述](#设计模式概述)
- [创建型模式](#创建型模式)
- [结构型模式](#结构型模式)
- [行为型模式](#行为型模式)
- [运行测试](#运行测试)

## 设计模式概述

设计模式是软件开发中常见问题的典型解决方案。Go 语言中常用的设计模式分为三大类：

| 类型 | 数量 | 说明 |
|------|------|------|
| 创建型模式 | 5 | 处理对象创建机制 |
| 结构型模式 | 6 | 处理类或对象的组合 |
| 行为型模式 | 8 | 处理对象间的通信 |

## 创建型模式

### 1. 单例模式（Singleton）

确保一个类只有一个实例，并提供一个全局访问点。

```go
// 使用场景：配置管理、日志记录、数据库连接池
func GetSingleton() *Singleton {
    if !singletonOnce.done {
        singletonInstance = &Singleton{name: "singleton"}
        singletonOnce.done = true
    }
    return singletonInstance
}
```

**优点**：保证唯一实例，减少内存开销
**缺点**：可能隐藏依赖关系，不利于测试

### 2. 工厂模式（Factory）

定义一个创建对象的接口，让子类决定实例化哪个类。

```go
// 使用场景：日志记录器、数据库驱动
func ShapeFactory(shapeType string) Shape {
    switch shapeType {
    case "circle":
        return &Circle{}
    case "rectangle":
        return &Rectangle{}
    }
    return nil
}
```

**优点**：解耦对象创建和使用
**缺点**：每增加一个产品需要修改工厂

### 3. 抽象工厂模式（Abstract Factory）

创建一系列相关或相互依赖对象的接口。

```go
// 使用场景：跨平台UI组件、数据库访问层
type GUIFactory interface {
    CreateButton() Button
    CreateCheckbox() Checkbox
}
```

**优点**：保证产品族的一致性
**缺点**：扩展新产品族困难

### 4. 建造者模式（Builder）

将复杂对象的构建与表示分离。

```go
// 使用场景：复杂配置对象、SQL查询构建器
computer := NewComputerBuilder().
    SetCPU("Intel i7").
    SetMemory("32GB").
    Build()
```

**优点**：分步创建复杂对象，支持链式调用
**缺点**：增加代码量

### 5. 原型模式（Prototype）

通过复制现有对象来创建新对象。

```go
// 使用场景：对象克隆、避免重复初始化
clone := original.Clone()
```

**优点**：无需知道对象创建细节
**缺点**：深拷贝可能复杂

## 结构型模式

### 1. 适配器模式（Adapter）

将一个类的接口转换成客户期望的另一个接口。

```go
// 使用场景：第三方库集成、旧系统兼容
type MediaAdapter struct {
    advancedPlayer AdvancedMediaPlayer
}
```

**优点**：提高复用性，解耦接口
**缺点**：增加系统复杂度

### 2. 装饰器模式（Decorator）

动态地给对象添加额外的职责。

```go
// 使用场景：日志增强、缓存装饰
coffee := NewMilkDecorator(NewSugarDecorator(&SimpleCoffee{}))
```

**优点**：比继承更灵活，可动态组合
**缺点**：可能产生很多小对象

### 3. 代理模式（Proxy）

为其他对象提供代理以控制对这个对象的访问。

```go
// 使用场景：远程代理、虚拟代理、保护代理
type ImageProxy struct {
    realImage *RealImage
}
```

**优点**：控制访问，延迟加载
**缺点**：增加间接层

### 4. 组合模式（Composite）

将对象组合成树形结构以表示"部分-整体"的层次结构。

```go
// 使用场景：文件系统、组织架构
type Directory struct {
    children []FileSystemNode
}
```

**优点**：统一处理简单和复杂元素
**缺点**：设计可能过于通用

### 5. 外观模式（Facade）

为子系统中的一组接口提供一个统一的入口。

```go
// 使用场景：简化复杂API、分层架构
func (f *ComputerFacade) Start() {
    f.cpu.Freeze()
    f.memory.Load(0, f.hardDrive.Read(0, 1024))
    f.cpu.Execute()
}
```

**优点**：简化接口，解耦客户端和子系统
**缺点**：可能成为"上帝对象"

### 6. 桥接模式（Bridge）

将抽象部分与实现部分分离。

```go
// 使用场景：跨平台UI、多种渲染方式
type Shape struct {
    renderer Renderer
}
```

**优点**：分离抽象和实现，扩展性好
**缺点**：增加系统复杂度

## 行为型模式

### 1. 策略模式（Strategy）

定义一系列算法，让它们可以互相替换。

```go
// 使用场景：支付方式、排序算法
type PaymentStrategy interface {
    Pay(amount float64)
}
```

**优点**：避免条件语句，易于扩展
**缺点**：客户端需要了解策略差异

### 2. 观察者模式（Observer）

定义对象间的一对多依赖关系。

```go
// 使用场景：事件系统、消息订阅
func (s *Subject) Notify(message string) {
    for _, observer := range s.observers {
        observer.Update(message)
    }
}
```

**优点**：解耦发布者和订阅者
**缺点**：可能导致内存泄漏

### 3. 命令模式（Command）

将请求封装为对象，支持撤销操作。

```go
// 使用场景：事务操作、宏命令
type Command interface {
    Execute()
    Undo()
}
```

**优点**：支持撤销，支持队列
**缺点**：增加系统复杂度

### 4. 责任链模式（Chain of Responsibility）

让多个对象都有机会处理请求。

```go
// 使用场景：审批流程、日志处理
func (h *ManagerHandler) Handle(request string) string {
    if request == "请假1天" {
        return "经理批准"
    }
    return h.HandleNext(request)
}
```

**优点**：解耦发送者和接收者
**缺点**：可能无人处理

### 5. 模板方法模式（Template Method）

定义算法骨架，将某些步骤延迟到子类。

```go
// 使用场景：框架设计、数据处理流程
func (t *DataMinerTemplate) Mine(miner DataMiner) {
    miner.OpenFile()
    miner.ExtractData()
    miner.ParseData()
}
```

**优点**：复用代码，控制扩展点
**缺点**：继承限制灵活性

### 6. 迭代器模式（Iterator）

提供一种方法顺序访问聚合对象中的元素。

```go
// 使用场景：集合遍历、自定义数据结构
type Iterator interface {
    HasNext() bool
    Next() interface{}
}
```

**优点**：统一遍历接口
**缺点**：增加额外类

### 7. 状态模式（State）

允许对象在内部状态改变时改变行为。

```go
// 使用场景：订单状态、游戏角色状态
func (o *Order) Process() {
    o.state.Handle()
}
```

**优点**：消除条件语句
**缺点**：增加类数量

### 8. 备忘录模式（Memento）

在不破坏封装的前提下捕获对象的内部状态。

```go
// 使用场景：撤销操作、事务回滚
func (o *Originator) Save() *Memento {
    return &Memento{state: o.state}
}
```

**优点**：保持封装性
**缺点**：可能消耗大量内存

## 运行测试

```bash
# 运行所有测试
go test -v ./designpattern

# 运行特定测试
go test -v ./designpattern -run TestSingleton
```

## 使用示例

```go
package main

import "github.com/yourusername/GoLearn/designpattern"

func main() {
    // 创建型模式
    designpattern.SingletonDemo()
    designpattern.FactoryDemo()
    designpattern.BuilderDemo()

    // 结构型模式
    designpattern.AdapterDemo()
    designpattern.DecoratorDemo()
    designpattern.ProxyDemo()

    // 行为型模式
    designpattern.StrategyDemo()
    designpattern.ObserverDemo()
    designpattern.CommandDemo()
}
```

## 设计原则

### SOLID原则

| 原则 | 说明 |
|------|------|
| S - 单一职责 | 一个类只负责一项职责 |
| O - 开闭原则 | 对扩展开放，对修改关闭 |
| L - 里氏替换 | 子类可以替换父类 |
| I - 接口隔离 | 接口要小而专一 |
| D - 依赖倒置 | 依赖抽象而非具体 |

### Go语言设计模式特点

1. **组合优于继承**：Go没有继承，使用组合实现代码复用
2. **接口隐式实现**：不需要显式声明实现接口
3. **函数是一等公民**：可以用函数替代简单的策略模式
4. **并发原语**：goroutine和channel提供了新的模式实现方式

## 相关资源

- [Go设计模式](https://github.com/tmrts/go-patterns)
- [设计模式：可复用面向对象软件的基础](https://book.douban.com/subject/1052241/)
- [Go语言设计与实现](https://draveness.me/golang/)
