# Go 语言面试题集

本模块整理了 Go 语言常见的面试题，包含问题和详细解答，涵盖基础语法、并发编程、底层原理等方面。

## 目录

- [基础语法面试题](#基础语法面试题)
- [并发编程面试题](#并发编程面试题)
- [底层原理面试题](#底层原理面试题)
- [运行测试](#运行测试)

## 基础语法面试题

### 1. Go语言的主要特点是什么？

**答案：**
- 静态类型、编译型语言，执行效率高
- 语法简洁，学习曲线平缓
- 内置并发支持（goroutine和channel）
- 垃圾回收机制，自动内存管理
- 编译速度快，跨平台编译
- 丰富的标准库
- 接口设计灵活，隐式实现

### 2. Go语言中make和new的区别？

**答案：**
| 特性 | new(T) | make(T, args) |
|------|--------|---------------|
| 用途 | 分配内存 | 初始化slice、map、channel |
| 返回值 | *T（指针） | T（引用类型本身） |
| 适用类型 | 所有类型 | 仅slice、map、channel |
| 初始化 | 零值初始化 | 完整初始化 |

### 3. Go语言中数组和切片的区别？

**答案：**
| 特性 | 数组 | 切片 |
|------|------|------|
| 类型 | 值类型 | 引用类型 |
| 长度 | 固定 | 可变 |
| 传递 | 复制整个数组 | 复制切片头 |
| 容量 | 无概念 | 有容量概念 |
| 声明 | [3]int | []int |

### 4. Go语言中defer的执行顺序？

**答案：**
- LIFO（后进先出）顺序执行
- defer语句在函数返回前执行
- defer语句的参数在声明时就确定了
- defer可以修改命名返回值

```go
func example() {
    defer fmt.Println("1")  // 最后执行
    defer fmt.Println("2")  // 第二执行
    defer fmt.Println("3")  // 最先执行
}
// 输出: 3 -> 2 -> 1
```

### 5. Go语言中slice的扩容机制？

**答案：**（Go 1.18+）
1. 如果新容量 > 2倍旧容量，直接使用新容量
2. 如果旧容量 < 256，新容量 = 2倍旧容量
3. 如果旧容量 >= 256，新容量 = 旧容量 + (旧容量+3*256)/4
4. 最后根据元素大小进行内存对齐

### 6. Go语言中map的特点？

**答案：**
- map是无序的，每次遍历顺序可能不同
- map不是并发安全的，并发读写需要加锁
- map的key必须是可以比较的类型
- map的零值是nil，需要make初始化
- 从map中取不存在的key，返回value类型的零值

### 7. Go语言中rune和byte的区别？

**答案：**
| 类型 | 别名 | 大小 | 用途 |
|------|------|------|------|
| byte | uint8 | 1字节 | 处理ASCII字符 |
| rune | int32 | 4字节 | 处理Unicode字符 |

## 并发编程面试题

### 1. Goroutine和线程的区别？

**答案：**
| 特性 | Goroutine | 线程 |
|------|-----------|------|
| 内存占用 | 2KB初始栈 | 1MB+栈空间 |
| 调度方式 | Go运行时调度 | OS调度 |
| 切换成本 | 低（用户态） | 高（内核态） |
| 通信方式 | Channel | 共享内存+锁 |
| 数量限制 | 百万级 | 有限 |

### 2. Channel的底层原理？

**答案：**
- 底层是hchan结构体，包含循环缓冲区、等待队列等
- 无缓冲channel：同步通信，发送和接收必须同时准备好
- 有缓冲channel：异步通信，缓冲区满时发送阻塞，空时接收阻塞
- channel的发送和接收都是原子操作
- 关闭channel后，读取会返回零值和false

### 3. select的执行机制？

**答案：**
1. 如果只有一个case可以执行，执行该case
2. 如果多个case可以执行，随机选择一个执行
3. 如果没有case可以执行且没有default，阻塞等待
4. 如果没有case可以执行且有default，执行default
5. select不会对nil channel进行操作

### 4. Context的使用场景？

**答案：**
| 函数 | 用途 |
|------|------|
| WithTimeout | 超时控制 |
| WithCancel | 取消信号 |
| WithDeadline | 截止时间 |
| WithValue | 值传递 |

### 5. sync.Map和普通map的区别？

**答案：**
| 特性 | sync.Map | 普通map |
|------|----------|---------|
| 并发安全 | 是 | 否 |
| 适用场景 | 读多写少 | 单goroutine |
| 内部实现 | 双map | 单map |
| Range方法 | 有 | 无 |
| len方法 | 无 | 有 |

### 6. Mutex和RWMutex的区别？

**答案：**
| 特性 | Mutex | RWMutex |
|------|-------|---------|
| 类型 | 互斥锁 | 读写锁 |
| 读操作 | 互斥 | 共享 |
| 写操作 | 互斥 | 互斥 |
| 适用场景 | 读写均衡 | 读多写少 |
| 开销 | 低 | 较高 |

### 7. 原子操作有哪些？

**答案：**
- Add: 原子加法
- CompareAndSwap (CAS): 比较并交换
- Swap: 交换
- Load: 原子读取
- Store: 原子存储
- atomic.Value: 存储任意类型的值

### 8. 如何实现并发安全的单例模式？

**答案：**
```go
// 推荐方式：sync.Once
var (
    instance *Singleton
    once     sync.Once
)

func GetSingleton() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}
```

### 9. 如何优雅地关闭goroutine？

**答案：**
1. 使用context的取消信号
2. 使用channel发送停止信号
3. 使用close关闭channel
4. 使用sync.WaitGroup等待goroutine退出

## 底层原理面试题

### 1. Go语言的垃圾回收机制？

**答案：**
- 三色标记-清除算法
- 并发标记，与用户程序并行执行
- 写屏障技术保证正确性
- STW（Stop-The-World）时间很短
- GOGC环境变量控制GC触发频率

### 2. Go语言的调度器原理？

**答案：** GMP模型
- G（Goroutine）：协程，包含栈、指令指针等信息
- M（Machine）：系统线程，执行G
- P（Processor）：逻辑处理器，包含运行队列
- 调度策略：工作窃取、系统调用时切换
- 抢占式调度：基于信号的抢占

### 3. Go语言的逃逸分析？

**答案：**
- 栈分配：函数返回后不再使用的变量
- 堆分配：变量在函数返回后仍被引用
- 逃逸场景：
  - 返回局部变量指针
  - 闭包捕获
  - 接口转换
  - 大对象
- 使用 `go build -gcflags="-m"` 查看逃逸分析结果

### 4. Go语言的slice底层原理？

**答案：**
- slice包含三个字段：指针、长度、容量
- 指针指向底层数组
- 长度表示元素个数
- 容量表示底层数组的大小
- 扩容时会分配新数组并复制数据

### 5. Go语言的interface底层原理？

**答案：**
- 空接口（interface{}）：包含类型指针和数据指针
- 非空接口：包含itab（接口表）和数据指针
- itab包含接口的类型信息和方法表
- 类型断言通过itab实现
- nil接口和nil值接口的区别

### 6. Go语言的map底层原理？

**答案：**
- 哈希表实现，使用链地址法解决冲突
- bucket（桶）存储键值对
- 扩容时渐进式迁移
- 不支持并发读写
- 遍历顺序随机

### 7. Go语言的channel底层原理？

**答案：** hchan结构体包含：
- 循环缓冲区：存储数据
- 发送等待队列：recvq
- 接收等待队列：sendq
- 互斥锁：保护共享数据
- 关闭标志：closed

### 8. unsafe包的使用场景？

**答案：**
- 指针类型转换
- 获取结构体字段偏移量
- 修改私有字段
- 零拷贝类型转换
- 与C语言交互

### 9. Go语言的反射机制？

**答案：**
- reflect.TypeOf: 获取类型信息
- reflect.ValueOf: 获取值信息
- 通过反射可以修改变量值
- 通过反射可以调用方法
- 反射性能较低，谨慎使用

### 10. Go语言的内存对齐？

**答案：**
- 成员对齐：每个成员的偏移量是其大小的整数倍
- 整体对齐：结构体大小是最大成员对齐值的整数倍
- 空结构体大小为0
- 合理排列字段可以减少内存占用

### 11. Go语言的闭包原理？

**答案：**
- 闭包捕获外部变量，形成引用
- 捕获的是变量本身，不是值
- 循环中的闭包陷阱
- 闭包会导致变量逃逸到堆

## 运行测试

```bash
# 运行所有测试
go test -v ./interview

# 运行特定测试
go test -v ./interview -run TestMakeVsNew
```

## 使用示例

```go
package main

import "github.com/yourusername/GoLearn/interview"

func main() {
    // 基础语法
    interview.MakeVsNew()
    interview.ArrayVsSlice()
    interview.DeferOrder()
    
    // 并发编程
    interview.ChannelPrinciple()
    interview.ContextUsage()
    interview.AtomicOperations()
    
    // 底层原理
    interview.GCMechanism()
    interview.EscapeAnalysisDemo()
    interview.ReflectionDemo()
}
```

## 注意事项

1. **理解原理**：不仅要会使用，更要理解底层原理
2. **实践验证**：通过代码验证理论知识
3. **性能考量**：了解各种操作的底层开销
4. **并发安全**：始终考虑并发场景下的安全性
5. **最佳实践**：遵循Go语言的惯用法

## 相关资源

- [Go官方文档](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go语言圣经](https://gopl-zh.github.io/)
- [Go语言设计与实现](https://draveness.me/golang/)
