package internals

import (
	"fmt"
	"unsafe"
)

// ============================================================
// Channel 底层实现原理
// ============================================================

// hchan 是 channel 的核心数据结构
type hchan struct {
	qcount   uint           // 队列中元素总数
	dataqsiz uint           // 循环队列大小（缓冲区容量）
	buf      unsafe.Pointer // 指向循环队列的指针
	elemsize uint16         // 元素大小
	closed   uint32         // channel是否已关闭
	elemtype unsafe.Pointer // 元素类型
	sendx    uint           // 发送索引
	recvx    uint           // 接收索引
	recvq    waitq          // 接收等待队列
	sendq    waitq          // 发送等待队列
	lock     mutex          // 互斥锁
}

// waitq 等待队列
type waitq struct {
	first *sudog
	last  *sudog
}

// sudog 表示等待的goroutine
type sudog struct {
	g       unsafe.Pointer // goroutine指针
	elem    unsafe.Pointer // 数据元素指针
	next    *sudog
	prev    *sudog
}

// mutex 简化的互斥锁
type mutex struct{}

// DemoChannelStructure 演示channel的内存结构
func DemoChannelStructure() {
	fmt.Println("=== Channel内存结构 ===")

	// 无缓冲channel
	ch1 := make(chan int)
	fmt.Printf("无缓冲channel: len=%d, cap=%d\n", len(ch1), cap(ch1))

	// 有缓冲channel
	ch2 := make(chan int, 3)
	fmt.Printf("有缓冲channel: len=%d, cap=%d\n", len(ch2), cap(ch2))

	fmt.Println("\nhchan结构:")
	fmt.Println("+--------------------+")
	fmt.Println("| qcount   (uint)    | 队列元素数")
	fmt.Println("| dataqsiz (uint)    | 缓冲区大小")
	fmt.Println("| buf      (unsafe.P)| 循环队列指针")
	fmt.Println("| elemsize (uint16)  | 元素大小")
	fmt.Println("| closed   (uint32)  | 关闭标志")
	fmt.Println("| elemtype (unsafe.P)| 元素类型")
	fmt.Println("| sendx    (uint)    | 发送索引")
	fmt.Println("| recvx    (uint)    | 接收索引")
	fmt.Println("| recvq    (waitq)   | 接收等待队列")
	fmt.Println("| sendq    (waitq)   | 发送等待队列")
	fmt.Println("| lock     (mutex)   | 互斥锁")
	fmt.Println("+--------------------+")

	fmt.Println("\n循环队列示意图 (cap=3):")
	fmt.Println("        sendx=1")
	fmt.Println("           |")
	fmt.Println("           v")
	fmt.Println("+---+---+---+")
	fmt.Println("|   | 2 | 3 |  <- buf")
	fmt.Println("+---+---+---+")
	fmt.Println("  ^")
	fmt.Println("  |")
	fmt.Println("recvx=0")
}

// DemoChannelSend 演示channel发送过程
func DemoChannelSend() {
	fmt.Println("\n=== Channel发送过程 ===")

	fmt.Println("发送流程 (ch <- x):")
	fmt.Println("1. 获取channel锁")
	fmt.Println("2. 检查channel是否已关闭")
	fmt.Println("   - 已关闭: panic(\"send on closed channel\")")
	fmt.Println("3. 检查是否有等待的接收者 (recvq)")
	fmt.Println("   - 有: 直接发送给接收者，唤醒接收者")
	fmt.Println("4. 检查缓冲区是否有空间")
	fmt.Println("   - 有: 放入缓冲区，更新sendx")
	fmt.Println("5. 否则，阻塞当前goroutine")
	fmt.Println("   - 创建sudog，加入sendq")
	fmt.Println("   - 释放锁，挂起goroutine")
	fmt.Println("6. 释放锁")

	fmt.Println("\n发送流程图:")
	fmt.Println("ch <- x")
	fmt.Println("   |")
	fmt.Println("   v")
	fmt.Println("+--------+     是      +----------------+")
	fmt.Println("|已关闭? |----------->| panic          |")
	fmt.Println("+--------+            +----------------+")
	fmt.Println("   |否")
	fmt.Println("   v")
	fmt.Println("+------------+  是    +----------------+")
	fmt.Println("|recvq有等待?|------->|直接发送+唤醒   |")
	fmt.Println("+------------+        +----------------+")
	fmt.Println("   |否")
	fmt.Println("   v")
	fmt.Println("+------------+  是    +----------------+")
	fmt.Println("|缓冲区有空间?|------->|放入缓冲区      |")
	fmt.Println("+------------+        +----------------+")
	fmt.Println("   |否")
	fmt.Println("   v")
	fmt.Println("+----------------+")
	fmt.Println("|阻塞，加入sendq |")
	fmt.Println("+----------------+")
}

// DemoChannelReceive 演示channel接收过程
func DemoChannelReceive() {
	fmt.Println("\n=== Channel接收过程 ===")

	fmt.Println("接收流程 (x := <-ch):")
	fmt.Println("1. 获取channel锁")
	fmt.Println("2. 检查是否有等待的发送者 (sendq)")
	fmt.Println("   - 有且缓冲区为空: 直接从发送者接收")
	fmt.Println("   - 有且缓冲区不为空: 从缓冲区接收，发送者数据放入缓冲区")
	fmt.Println("3. 检查缓冲区是否有数据")
	fmt.Println("   - 有: 从缓冲区接收，更新recvx")
	fmt.Println("4. 检查channel是否已关闭")
	fmt.Println("   - 已关闭: 返回零值和false")
	fmt.Println("5. 否则，阻塞当前goroutine")
	fmt.Println("   - 创建sudog，加入recvq")
	fmt.Println("   - 释放锁，挂起goroutine")
	fmt.Println("6. 释放锁")

	fmt.Println("\n接收流程图:")
	fmt.Println("<-ch")
	fmt.Println("   |")
	fmt.Println("   v")
	fmt.Println("+------------+  是    +----------------+")
	fmt.Println("|sendq有等待?|------->|接收+唤醒发送者 |")
	fmt.Println("+------------+        +----------------+")
	fmt.Println("   |否")
	fmt.Println("   v")
	fmt.Println("+------------+  是    +----------------+")
	fmt.Println("|缓冲区有数据?|------->|从缓冲区接收    |")
	fmt.Println("+------------+        +----------------+")
	fmt.Println("   |否")
	fmt.Println("   v")
	fmt.Println("+--------+     是      +----------------+")
	fmt.Println("|已关闭? |----------->|返回零值        |")
	fmt.Println("+--------+            +----------------+")
	fmt.Println("   |否")
	fmt.Println("   v")
	fmt.Println("+----------------+")
	fmt.Println("|阻塞，加入recvq |")
	fmt.Println("+----------------+")
}

// DemoChannelClose 演示channel关闭过程
func DemoChannelClose() {
	fmt.Println("\n=== Channel关闭过程 ===")

	fmt.Println("关闭流程 (close(ch)):")
	fmt.Println("1. 获取channel锁")
	fmt.Println("2. 检查channel是否已关闭")
	fmt.Println("   - 已关闭: panic(\"close of closed channel\")")
	fmt.Println("3. 检查channel是否为nil")
	fmt.Println("   - nil: panic(\"close of nil channel\")")
	fmt.Println("4. 设置closed标志")
	fmt.Println("5. 释放所有等待的接收者")
	fmt.Println("   - 唤醒recvq中的所有goroutine")
	fmt.Println("   - 它们会收到零值")
	fmt.Println("6. 释放所有等待的发送者")
	fmt.Println("   - 唤醒sendq中的所有goroutine")
	fmt.Println("   - 它们会panic")
	fmt.Println("7. 释放锁")

	fmt.Println("\n关闭后的行为:")
	fmt.Println("  发送: panic(\"send on closed channel\")")
	fmt.Println("  接收: 返回零值和false")
	fmt.Println("  关闭: panic(\"close of closed channel\")")

	fmt.Println("\n正确用法:")
	fmt.Println("  1. 由发送者关闭channel")
	fmt.Println("  2. 不要关闭已关闭的channel")
	fmt.Println("  3. 不要向已关闭的channel发送")
	fmt.Println("  4. 使用 ok 判断是否关闭:")
	fmt.Println("     v, ok := <-ch")
	fmt.Println("     if !ok { /* channel已关闭 */ }")
}

// DemoChannelBuffered 演示有缓冲channel的工作原理
func DemoChannelBuffered() {
	fmt.Println("\n=== 有缓冲Channel工作原理 ===")

	fmt.Println("循环队列实现:")
	fmt.Println("- buf: 指向预分配的循环队列")
	fmt.Println("- sendx: 下一个发送位置")
	fmt.Println("- recvx: 下一个接收位置")
	fmt.Println("- qcount: 当前元素数")
	fmt.Println("- dataqsiz: 队列容量")

	fmt.Println("\n操作示例 (cap=3):")

	// 初始状态
	fmt.Println("\n初始状态:")
	fmt.Println("sendx=0, recvx=0, qcount=0")
	fmt.Println("+---+---+---+")
	fmt.Println("|   |   |   |")
	fmt.Println("+---+---+---+")

	// 发送3个元素
	fmt.Println("\n发送3个元素 (1,2,3):")
	fmt.Println("sendx=0->1->2->3, recvx=0, qcount=3")
	fmt.Println("+---+---+---+")
	fmt.Println("| 1 | 2 | 3 |")
	fmt.Println("+---+---+---+")

	// 接收1个元素
	fmt.Println("\n接收1个元素:")
	fmt.Println("sendx=3, recvx=0->1, qcount=2")
	fmt.Println("+---+---+---+")
	fmt.Println("|   | 2 | 3 |")
	fmt.Println("+---+---+---+")
	fmt.Println("  ^")
	fmt.Println("  recvx=1")

	// 再发送1个元素（循环）
	fmt.Println("\n再发送1个元素 (4):")
	fmt.Println("sendx=3->0, recvx=1, qcount=3")
	fmt.Println("+---+---+---+")
	fmt.Println("| 4 | 2 | 3 |")
	fmt.Println("+---+---+---+")
	fmt.Println("  ^")
	fmt.Println("  sendx=0 (循环回来)")
}

// DemoChannelUnbuffered 演示无缓冲channel的工作原理
func DemoChannelUnbuffered() {
	fmt.Println("\n=== 无缓冲Channel工作原理 ===")

	fmt.Println("特点:")
	fmt.Println("  - buf = nil (没有缓冲区)")
	fmt.Println("  - dataqsiz = 0")
	fmt.Println("  - 发送和接收必须同步进行")

	fmt.Println("\n工作流程:")
	fmt.Println("1. 发送者到达，发现没有接收者")
	fmt.Println("   -> 创建sudog，加入sendq")
	fmt.Println("   -> 阻塞等待")
	fmt.Println("")
	fmt.Println("2. 接收者到达，发现有发送者等待")
	fmt.Println("   -> 直接从发送者复制数据")
	fmt.Println("   -> 唤醒发送者")
	fmt.Println("")
	fmt.Println("3. 接收者先到达，发送者后到达")
	fmt.Println("   -> 同理，接收者等待，发送者唤醒它")

	fmt.Println("\n同步语义:")
	fmt.Println("  发送 happens-before 接收完成")
	fmt.Println("  接收 happens-before 发送完成")
	fmt.Println("  (对于无缓冲channel)")
}

// DemoChannelSelect 演示select的实现原理
func DemoChannelSelect() {
	fmt.Println("\n=== Select实现原理 ===")

	fmt.Println("select语句编译后:")
	fmt.Println("1. 编译器将select转换为runtime.selectgo调用")
	fmt.Println("2. 随机打乱case顺序")
	fmt.Println("3. 检查所有case:")
	fmt.Println("   - 如果有case可以立即执行，随机选择一个")
	fmt.Println("   - 如果都没有准备好，阻塞等待")
	fmt.Println("   - 如果有default，执行default")

	fmt.Println("\n实现伪代码:")
	fmt.Println("func selectgo(cases []scase) (int, bool) {")
	fmt.Println("    // 随机打乱顺序")
	fmt.Println("    shuffle(cases)")
	fmt.Println("    ")
	fmt.Println("    // 检查是否有case可以立即执行")
	fmt.Println("    for _, c := range cases {")
	fmt.Println("        if canExecute(c) {")
	fmt.Println("            execute(c)")
	fmt.Println("            return c.index, true")
	fmt.Println("        }")
	fmt.Println("    }")
	fmt.Println("    ")
	fmt.Println("    // 没有case准备好，阻塞等待")
	fmt.Println("    return block(cases)")
	fmt.Println("}")

	fmt.Println("\n性能注意:")
	fmt.Println("  - select的case数量影响性能")
	fmt.Println("  - 编译器对小数量case有优化")
	fmt.Println("  - 避免在热路径使用大量case的select")
}

// DemoChannelNil 演示nil channel的行为
func DemoChannelNil() {
	fmt.Println("\n=== Nil Channel行为 ===")

	var nilCh chan int

	fmt.Println("nil channel的特点:")
	fmt.Printf("  nilCh == nil: %v\n", nilCh == nil)
	fmt.Printf("  len(nilCh) = %d\n", len(nilCh))
	fmt.Printf("  cap(nilCh) = %d\n", cap(nilCh))

	fmt.Println("\n操作行为:")
	fmt.Println("  发送: 永久阻塞")
	fmt.Println("  接收: 永久阻塞")
	fmt.Println("  关闭: panic(\"close of nil channel\")")

	fmt.Println("\n在select中使用nil channel:")
	fmt.Println("  - nil channel的case永远不会被选中")
	fmt.Println("  - 可以用来禁用某个case")
	fmt.Println("")
	fmt.Println("  示例:")
	fmt.Println("  var ch chan int  // nil")
	fmt.Println("  select {")
	fmt.Println("  case v := <-ch:  // 永远不会执行")
	fmt.Println("  case <-done:")
	fmt.Println("  }")
}

// DemoChannelCommonPatterns 演示channel常见模式
func DemoChannelCommonPatterns() {
	fmt.Println("\n=== Channel常见模式 ===")

	fmt.Println("1. 生产者-消费者模式")
	fmt.Println("   producer -> ch -> consumer")

	fmt.Println("\n2. 工作池模式")
	fmt.Println("   jobs -> [worker1, worker2, worker3] -> results")

	fmt.Println("\n3. 扇出模式")
	fmt.Println("   source -> [ch1, ch2, ch3] -> 多个消费者")

	fmt.Println("\n4. 扇入模式")
	fmt.Println("   多个生产者 -> [ch1, ch2, ch3] -> 单个消费者")

	fmt.Println("\n5. 超时控制")
	fmt.Println("   select {")
	fmt.Println("   case v := <-ch:")
	fmt.Println("   case <-time.After(timeout):")
	fmt.Println("   }")

	fmt.Println("\n6. 取消信号")
	fmt.Println("   done := make(chan struct{})")
	fmt.Println("   go func() {")
	fmt.Println("       select {")
	fmt.Println("       case <-done: return")
	fmt.Println("       }")
	fmt.Println("   }()")
	fmt.Println("   close(done) // 取消")
}
