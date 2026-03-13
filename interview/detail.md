# Go 语言面试题详解

本文档提供 Go 语言面试题的详细解答，涵盖基础语法、并发编程和底层原理。

---

## 一、基础语法面试题详解

### 1. Go语言的主要特点是什么？

**详细解答：**

Go语言是由Google开发的一种静态类型、编译型语言，具有以下核心特点：

#### 1.1 静态类型与编译型
```go
// 编译时类型检查
var x int = 10
x = "hello" // 编译错误：不能将string赋值给int

// 编译生成机器码，执行效率高
go build -o myapp main.go
```

#### 1.2 语法简洁
- 只有25个关键字
- 没有继承、重载、泛型（Go 1.18之前）
- 强制统一的代码风格（gofmt）

#### 1.3 内置并发支持
```go
// goroutine: 轻量级协程
go func() {
    fmt.Println("并发执行")
}()

// channel: 通信机制
ch := make(chan int)
go func() { ch <- 42 }()
value := <-ch
```

#### 1.4 垃圾回收
- 自动内存管理
- 三色标记-清除算法
- 低延迟STW

#### 1.5 快速编译
- 依赖分析优化
- 模块化编译
- 增量编译支持

---

### 2. make和new的区别？

**详细解答：**

#### 2.1 new 的实现原理
```go
// new(T) 的内部实现（伪代码）
func new(T) *T {
    var zero T           // 分配零值内存
    return &zero         // 返回指针
}

// 使用示例
p := new(int)       // *int, 值为0
s := new([]int)     // *[]int, 值为nil（切片的零值）
m := new(map[int]int) // *map[int]int, 值为nil
```

**关键点：**
- new 为任何类型分配内存
- 返回指向零值的指针
- 不会初始化引用类型（slice、map、channel）

#### 2.2 make 的实现原理
```go
// make 只能用于 slice、map、channel

// slice: make([]T, len, cap)
s := make([]int, 5, 10)
// 内部会：
// 1. 分配底层数组（大小为cap）
// 2. 创建slice结构体（ptr, len, cap）

// map: make(map[K]V, hint)
m := make(map[string]int, 100)
// 内部会：
// 1. 分配哈希表结构
// 2. 初始化bucket数组

// channel: make(chan T, size)
ch := make(chan int, 10)
// 内部会：
// 1. 分配hchan结构体
// 2. 初始化缓冲区
```

#### 2.3 为什么需要 make？
```go
// 如果用 new 创建 map
m := new(map[string]int)
*m = make(map[string]int) // 必须再初始化

// 直接用 make 更简洁
m := make(map[string]int)
```

---

### 3. 数组和切片的区别？

**详细解答：**

#### 3.1 内存布局

**数组：**
```
[3]int{1, 2, 3} 内存布局：
+-----+-----+-----+
|  1  |  2  |  3  |
+-----+-----+-----+
连续内存，值类型
```

**切片：**
```
slice := make([]int, 3, 5)

SliceHeader 结构：
+--------+--------+--------+
|  ptr   |  len   |  cap   |
+--------+--------+--------+
    |
    v
底层数组：
+-----+-----+-----+-----+-----+
|  0  |  0  |  0  |     |     |
+-----+-----+-----+-----+-----+
```

#### 3.2 值传递 vs 引用传递
```go
// 数组：复制整个数组
func modifyArray(arr [3]int) {
    arr[0] = 100  // 修改的是副本
}

// 切片：复制切片头（ptr, len, cap）
func modifySlice(s []int) {
    s[0] = 100    // 通过ptr修改底层数组
    s = append(s, 4) // 修改的是切片头副本，不影响原切片
}
```

#### 3.3 切片扩容详解
```go
// Go 1.18+ 扩容算法
func growslice(oldLen, oldCap, num int) int {
    newcap := oldCap
    doublecap := newcap + newcap
    if num > doublecap {
        newcap = num
    } else {
        const threshold = 256
        if oldCap < threshold {
            newcap = doublecap
        } else {
            // 平滑过渡
            for newcap < num {
                newcap += (newcap + 3*threshold) / 4
            }
        }
    }
    // 内存对齐...
    return newcap
}
```

---

### 4. defer的执行机制

**详细解答：**

#### 4.1 defer 的内部实现
```go
// defer 语句会被编译器转换为：
// 1. 在栈上分配一个 _defer 结构体
// 2. 将 defer 函数和参数保存到结构体
// 3. 将结构体链入 goroutine 的 defer 链表

type _defer struct {
    siz     int32   // 参数大小
    started bool
    openDefer bool
    sp      uintptr // 栈指针
    pc      uintptr // 程序计数器
    fn      *funcval // defer 函数
    _panic  *_panic
    link    *_defer // 链表指针
}
```

#### 4.2 参数预计算
```go
func deferArgs() {
    i := 1
    defer fmt.Println("defer:", i)  // 参数i在此时已计算为1
    i = 2
    // 输出: defer: 1
}

func deferClosure() {
    i := 1
    defer func() {
        fmt.Println("defer:", i)  // 闭包捕获变量i
    }()
    i = 2
    // 输出: defer: 2
}
```

#### 4.3 defer 与 return
```go
func deferReturn() (result int) {
    defer func() {
        result++  // 可以修改命名返回值
    }()
    return 0  // 1. result = 0, 2. defer执行 result++, 3. 返回
}
// 返回值: 1

// 执行顺序：
// 1. 返回值赋值
// 2. 执行 defer
// 3. 执行 RET 指令
```

---

### 5. Go语言只有值传递

**详细解答：**

#### 5.1 值传递的本质
```go
// 所有参数传递都是复制

// 基本类型
func passInt(x int) {
    x = 100  // 修改的是副本
}

// 指针类型
func passPointer(x *int) {
    *x = 100  // 通过指针修改原值
    x = nil   // 修改的是指针副本，不影响原指针
}

// 切片类型
func passSlice(s []int) {
    s[0] = 100     // 通过切片头中的ptr修改底层数组
    s = append(s, 4) // 修改的是切片头副本
}
```

#### 5.2 切片传递详解
```go
// 切片头结构
type SliceHeader struct {
    Data uintptr  // 底层数组指针
    Len  int      // 长度
    Cap  int      // 容量
}

// 传递切片时，复制的是这个24字节的结构体
// Data 指针指向的底层数组是共享的
```

---

### 6. slice扩容机制详解

**详细解答：**

#### 6.1 扩容算法（Go 1.18+）
```go
// 源码简化版
func nextGrow(oldCap, num int) int {
    newcap := oldCap
    doublecap := newcap + newcap
    
    if num > doublecap {
        newcap = num
    } else {
        if oldCap < 256 {
            newcap = doublecap  // 2倍增长
        } else {
            // 1.25倍平滑增长
            newcap += (newcap + 3*256) / 4
        }
    }
    
    // 内存对齐调整
    // ...
    return newcap
}
```

#### 6.2 扩容示例
```
容量变化序列：
0 -> 1 -> 2 -> 4 -> 8 -> 16 -> 32 -> 64 -> 128 -> 256
-> 512 -> 848 -> 1280 -> ...

注意：256之后不再是严格的2倍增长
```

---

### 7. map底层实现详解

**详细解答：**

#### 7.1 map的数据结构
```go
// hmap: map的头部结构
type hmap struct {
    count     int    // 元素个数
    flags     uint8
    B         uint8  // bucket数量 = 2^B
    noverflow uint16 // 溢出bucket数量
    hash0     uint32 // hash种子
    
    buckets    unsafe.Pointer // bucket数组
    oldbuckets unsafe.Pointer // 扩容时的旧bucket
    nevacuate  uintptr        // 扩容进度
    
    extra *mapextra // 溢出bucket链表
}

// bmap: bucket结构
type bmap struct {
    tophash [8]uint8  // 高8位hash值
    // 后面跟着8个key和8个value
    // 最后是overflow指针
}
```

#### 7.2 查找过程
```
1. 计算key的hash值
2. 取低B位，定位bucket
3. 取高8位，在bucket中查找tophash
4. 匹配则返回对应value
5. 不匹配则继续查找overflow bucket
```

#### 7.3 扩容机制
```
触发条件：
1. 负载因子 > 6.5（等量扩容）
2. overflow bucket过多（增量扩容）

扩容过程：
1. 创建新bucket数组
2. 渐进式迁移（每次操作迁移少量bucket）
3. 迁移完成后释放旧bucket
```

---

### 8. rune和byte详解

**详细解答：**

#### 8.1 字符串底层
```go
// string的底层结构
type StringHeader struct {
    Data uintptr  // 字节数组指针
    Len  int      // 字节长度
}

// 字符串是不可变的字节序列
s := "Hello世界"
// 内存：H e l l o 世 界
// 字节：72 101 108 108 111 228 184 150 231 149 140
// 长度：11字节（不是7）
```

#### 8.2 rune处理
```go
// rune是int32的别名，表示Unicode码点
s := "Hello世界"

// 错误方式：按字节遍历
for i := 0; i < len(s); i++ {
    fmt.Printf("%c ", s[i]) // 会输出乱码
}

// 正确方式：按rune遍历
for _, r := range s {
    fmt.Printf("%c ", r) // H e l l o 世 界
}

// 获取字符数
charCount := utf8.RuneCountInString(s) // 7
```

---

## 二、并发编程面试题详解

### 1. Goroutine与线程的区别

**详细解答：**

| 特性 | Goroutine | 系统线程 |
|------|-----------|----------|
| 栈大小 | 2KB起，动态增长 | 固定1MB+ |
| 调度 | Go运行时调度 | OS调度 |
| 切换成本 | ~200ns | ~1μs |
| 创建成本 | 低 | 高 |
| 通信 | Channel | 共享内存+锁 |

#### 1.1 Goroutine栈
```go
// 栈增长过程
// 初始: 2KB
// 增长: 2KB -> 4KB -> 8KB -> 16KB -> ...
// 最大: 1GB (64位系统)

// 栈分裂检查
func growStack() {
    // 编译器会在函数入口插入栈检查
    // 如果栈不够，调用morestack增长
}
```

#### 1.2 调度开销对比
```
Goroutine切换：
1. 保存3个寄存器（PC/SP/DX）
2. 无需进入内核态
3. 约200-300ns

线程切换：
1. 保存所有寄存器（16+个）
2. 进入内核态
3. 约1-2μs
```

---

### 2. Channel底层原理

**详细解答：**

#### 2.1 hchan结构
```go
type hchan struct {
    qcount   uint           // 队列中元素数量
    dataqsiz uint           // 循环队列大小
    buf      unsafe.Pointer // 循环队列指针
    elemsize uint16         // 元素大小
    closed   uint32         // 关闭标志
    elemtype *_type         // 元素类型
    
    sendx    uint   // 发送索引
    recvx    uint   // 接收索引
    recvq    waitq  // 接收等待队列
    sendq    waitq  // 发送等待队列
    
    lock mutex  // 互斥锁
}
```

#### 2.2 发送过程
```
1. 获取锁
2. 如果有等待的接收者：
   - 直接将数据复制给接收者
   - 唤醒接收者
3. 如果有缓冲区空间：
   - 将数据复制到缓冲区
4. 否则：
   - 将当前goroutine加入sendq
   - 阻塞等待
5. 释放锁
```

#### 2.3 接收过程
```
1. 获取锁
2. 如果有等待的发送者：
   - 从发送者获取数据
   - 唤醒发送者
3. 如果缓冲区有数据：
   - 从缓冲区读取数据
4. 如果channel已关闭：
   - 返回零值
5. 否则：
   - 将当前goroutine加入recvq
   - 阻塞等待
6. 释放锁
```

---

### 3. Context使用详解

**详细解答：**

#### 3.1 Context接口
```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

#### 3.2 使用模式
```go
// 超时控制
func doWork(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case result := <-doSomething():
        return nil
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := doWork(ctx); err != nil {
        log.Fatal(err)
    }
}
```

#### 3.3 最佳实践
1. 不要将Context存储在结构体中
2. Context应该作为函数的第一个参数
3. 不要传递nil Context
4. Context.Value只用于请求范围的数据

---

## 三、底层原理面试题详解

### 1. GMP调度模型

**详见**：[GMP调度器详解](./gmp_scheduler.md)

### 2. 垃圾回收机制

**详见**：[GC垃圾回收详解](./gc_detail.md)

### 3. 内存模型

**详见**：[内存模型详解](./memory_model_detail.md)

---

## 四、常见陷阱与最佳实践

### 1. 循环变量捕获
```go
// Go 1.22之前版本的问题（Go 1.22+已修复）
for i := 0; i < 3; i++ {
    defer fmt.Println(i)  // Go 1.22之前: 3, 3, 3; Go 1.22+: 2, 1, 0
}

// 兼容所有版本的写法
for i := 0; i < 3; i++ {
    defer func(n int) { fmt.Println(n) }(i)  // 输出: 2, 1, 0
}
```

### 2. 接口nil判断
```go
// 错误
var p *int
var i interface{} = p
fmt.Println(i == nil) // false!

// 原因：接口包含(type, value)，p的类型是*int
// 正确判断
if p == nil && i == nil {
    // ...
}
```

### 3. 切片内存泄漏
```go
// 内存泄漏
func leak() []int {
    large := make([]int, 1000000)
    return large[:10]  // 仍引用整个大数组
}

// 正确做法
func noLeak() []int {
    large := make([]int, 1000000)
    result := make([]int, 10)
    copy(result, large[:10])
    return result
}
```
