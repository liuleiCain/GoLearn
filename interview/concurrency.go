package interview

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================
// 并发编程面试题
// ============================================================

// Q1: Goroutine和线程的区别？
// A1: Goroutine和线程的主要区别：
// 1. 内存占用：Goroutine初始栈2KB，线程通常1MB+
// 2. 调度：Goroutine由Go运行时调度，线程由OS调度
// 3. 切换成本：Goroutine切换成本低（用户态），线程切换成本高（内核态）
// 4. 通信：Goroutine使用channel，线程使用共享内存+锁
// 5. 数量：可以轻松创建百万级goroutine，线程数量有限

func GoroutineVsThread() {
	fmt.Println("=== Goroutine vs 线程 ===")

	// 打印goroutine栈大小信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("初始goroutine数量: %d\n", runtime.NumGoroutine())

	// 创建大量goroutine
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
		}()
	}

	fmt.Printf("创建10000个goroutine后: %d\n", runtime.NumGoroutine())
	wg.Wait()

	runtime.ReadMemStats(&memStats)
	fmt.Printf("goroutine执行完毕后: %d\n", runtime.NumGoroutine())
}

// Q2: Channel的底层原理？
// A2: Channel的底层原理：
// 1. 底层是hchan结构体，包含循环缓冲区、等待队列等
// 2. 无缓冲channel：同步通信，发送和接收必须同时准备好
// 3. 有缓冲channel：异步通信，缓冲区满时发送阻塞，空时接收阻塞
// 4. channel的发送和接收都是原子操作
// 5. 关闭channel后，读取会返回零值和false

func ChannelPrinciple() {
	fmt.Println("\n=== Channel底层原理 ===")

	// 无缓冲channel - 同步
	ch1 := make(chan int)
	go func() {
		ch1 <- 1
		fmt.Println("无缓冲channel发送完成")
	}()
	<-ch1
	fmt.Println("无缓冲channel接收完成")

	// 有缓冲channel - 异步
	ch2 := make(chan int, 2)
	ch2 <- 1
	ch2 <- 2
	fmt.Printf("有缓冲channel: len=%d, cap=%d\n", len(ch2), cap(ch2))

	// 关闭channel
	close(ch2)
	for v := range ch2 {
		fmt.Printf("从已关闭channel读取: %d\n", v)
	}

	// 从关闭的channel读取
	v, ok := <-ch2
	fmt.Printf("关闭后读取: value=%d, ok=%v\n", v, ok)
}

// Q3: select的执行机制？
// A3: select的执行机制：
// 1. 如果只有一个case可以执行，执行该case
// 2. 如果多个case可以执行，随机选择一个执行
// 3. 如果没有case可以执行且没有default，阻塞等待
// 4. 如果没有case可以执行且有default，执行default
// 5. select不会对nil channel进行操作

func SelectMechanism() {
	fmt.Println("\n=== select执行机制 ===")

	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	ch1 <- 1
	ch2 <- 2

	// 多个case可执行时随机选择
	for i := 0; i < 5; i++ {
		select {
		case v := <-ch1:
			fmt.Printf("从ch1读取: %d\n", v)
			ch1 <- 1 // 放回去
		case v := <-ch2:
			fmt.Printf("从ch2读取: %d\n", v)
			ch2 <- 2 // 放回去
		}
	}

	// 非阻塞select
	select {
	case v := <-ch1:
		fmt.Printf("非阻塞读取ch1: %d\n", v)
	default:
		fmt.Println("没有数据可读")
	}
}

// Q4: Context的使用场景？
// A4: Context的使用场景：
// 1. 超时控制：WithTimeout
// 2. 取消信号：WithCancel
// 3. 截止时间：WithDeadline
// 4. 值传递：WithValue
// 5. 链式传递：子context继承父context的取消信号

func ContextUsage() {
	fmt.Println("\n=== Context使用场景 ===")

	// 超时控制
	ctx1, cancel1 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel1()

	select {
	case <-time.After(200 * time.Millisecond):
		fmt.Println("操作完成")
	case <-ctx1.Done():
		fmt.Printf("超时取消: %v\n", ctx1.Err())
	}

	// 取消信号
	ctx2, cancel2 := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel2()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		fmt.Println("操作完成")
	case <-ctx2.Done():
		fmt.Printf("手动取消: %v\n", ctx2.Err())
	}

	// 值传递
	ctx3 := context.WithValue(context.Background(), "key", "value")
	if v, ok := ctx3.Value("key").(string); ok {
		fmt.Printf("从context获取值: %s\n", v)
	}
}

// Q5: sync.Map和普通map的区别？
// A5: sync.Map和普通map的区别：
// 1. sync.Map是并发安全的，普通map不是
// 2. sync.Map适合读多写少的场景
// 3. sync.Map使用空间换时间，内部维护两份map
// 4. sync.Map的Range方法可以安全遍历
// 5. sync.Map没有len方法

func SyncMapVsMap() {
	fmt.Println("\n=== sync.Map vs map ===")

	var sm sync.Map

	// 存储
	sm.Store("a", 1)
	sm.Store("b", 2)
	sm.Store("c", 3)

	// 读取
	if v, ok := sm.Load("a"); ok {
		fmt.Printf("读取: a=%v\n", v)
	}

	// 删除
	sm.Delete("b")

	// 遍历
	fmt.Print("遍历sync.Map: ")
	sm.Range(func(key, value interface{}) bool {
		fmt.Printf("%v=%v ", key, value)
		return true
	})
	fmt.Println()

	// LoadOrStore
	actual, loaded := sm.LoadOrStore("a", 100)
	fmt.Printf("LoadOrStore: actual=%v, loaded=%v\n", actual, loaded)
}

// Q6: WaitGroup的使用场景？
// A6: WaitGroup的使用场景：
// 1. 等待一组goroutine完成
// 2. Add()增加计数器
// 3. Done()减少计数器
// 4. Wait()阻塞等待计数器归零
// 5. 计数器不能为负数

func WaitGroupUsage() {
	fmt.Println("\n=== WaitGroup使用 ===")

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d 执行\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("所有goroutine执行完毕")
}

// Q7: Mutex和RWMutex的区别？
// A7: Mutex和RWMutex的区别：
// 1. Mutex是互斥锁，同一时间只能一个goroutine持有
// 2. RWMutex是读写锁，允许多个读者同时读取
// 3. RWMutex写操作是互斥的
// 4. RWMutex适合读多写少的场景
// 5. RWMutex比Mutex开销大

func MutexVsRWMutex() {
	fmt.Println("\n=== Mutex vs RWMutex ===")

	var (
		mu    sync.Mutex
		rwMu  sync.RWMutex
		data  int
	)

	// Mutex
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			data++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("Mutex写入1000次: data=%d, 耗时=%v\n", data, time.Since(start))

	// RWMutex
	data = 0
	start = time.Now()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rwMu.Lock()
			data++
			rwMu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("RWMutex写入1000次: data=%d, 耗时=%v\n", data, time.Since(start))
}

// Q8: 原子操作有哪些？
// A8: sync/atomic包提供的原子操作：
// 1. Add: 原子加法
// 2. CompareAndSwap (CAS): 比较并交换
// 3. Swap: 交换
// 4. Load: 原子读取
// 5. Store: 原子存储
// 6. atomic.Value: 存储任意类型的值

func AtomicOperations() {
	fmt.Println("\n=== 原子操作 ===")

	var counter int64

	// Add
	atomic.AddInt64(&counter, 1)
	fmt.Printf("Add: counter=%d\n", counter)

	// CompareAndSwap
	swapped := atomic.CompareAndSwapInt64(&counter, 1, 10)
	fmt.Printf("CAS: swapped=%v, counter=%d\n", swapped, counter)

	// Swap
	old := atomic.SwapInt64(&counter, 100)
	fmt.Printf("Swap: old=%d, counter=%d\n", old, counter)

	// Load
	value := atomic.LoadInt64(&counter)
	fmt.Printf("Load: value=%d\n", value)

	// Store
	atomic.StoreInt64(&counter, 200)
	fmt.Printf("Store: counter=%d\n", counter)

	// atomic.Value
	var config atomic.Value
	type Config struct {
		Name string
	}
	config.Store(Config{Name: "config1"})
	cfg := config.Load().(Config)
	fmt.Printf("atomic.Value: config.Name=%s\n", cfg.Name)
}

// Q9: 如何实现一个并发安全的单例模式？
// A9: 实现并发安全单例的几种方式：
// 1. sync.Once（推荐）
// 2. 双重检查锁定
// 3. sync.Mutex
// 4. init函数

type Singleton struct {
	name string
}

var (
	instance     *Singleton
	once         sync.Once
	instanceMu   sync.Mutex
)

func GetSingletonOnce() *Singleton {
	once.Do(func() {
		instance = &Singleton{name: "singleton"}
	})
	return instance
}

func GetSingletonDoubleCheck() *Singleton {
	if instance == nil {
		instanceMu.Lock()
		defer instanceMu.Unlock()
		if instance == nil {
			instance = &Singleton{name: "singleton"}
		}
	}
	return instance
}

func SingletonPattern() {
	fmt.Println("\n=== 单例模式 ===")

	// sync.Once方式
	s1 := GetSingletonOnce()
	s2 := GetSingletonOnce()
	fmt.Printf("sync.Once: s1==s2: %v\n", s1 == s2)

	// 双重检查锁定方式
	instance = nil // 重置
	s3 := GetSingletonDoubleCheck()
	s4 := GetSingletonDoubleCheck()
	fmt.Printf("双重检查: s3==s4: %v\n", s3 == s4)
}

// Q10: 如何优雅地关闭goroutine？
// A10: 优雅关闭goroutine的方式：
// 1. 使用context的取消信号
// 2. 使用channel发送停止信号
// 3. 使用close关闭channel
// 4. 使用sync.WaitGroup等待goroutine退出

func GracefulShutdown() {
	fmt.Println("\n=== 优雅关闭goroutine ===")

	// 方式1: context
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("goroutine收到取消信号，退出")
				return
			default:
				// 模拟工作
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	time.Sleep(200 * time.Millisecond)
	cancel()
	wg.Wait()
	fmt.Println("goroutine已退出")

	// 方式2: channel
	stop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				fmt.Println("收到stop信号，退出")
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()
	fmt.Println("goroutine已退出")
}
