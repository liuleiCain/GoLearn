package goroutine

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func BasicGoroutine() {
	fmt.Println("=== 基础Goroutine示例 ===")

	fmt.Println("主函数开始")

	go func() {
		fmt.Println("这是第一个goroutine")
	}()

	go func(msg string) {
		fmt.Println("这是带参数的goroutine:", msg)
	}("Hello Goroutine")

	time.Sleep(100 * time.Millisecond)
	fmt.Println("主函数结束")
}

func WaitGroupDemo() {
	fmt.Println("=== WaitGroup同步示例 ===")

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d 开始工作\n", id)
			time.Sleep(time.Duration(id*100) * time.Millisecond)
			fmt.Printf("Worker %d 完成工作\n", id)
		}(i)
	}

	fmt.Println("等待所有worker完成...")
	wg.Wait()
	fmt.Println("所有worker已完成")
}

func GoroutineLeakDemo() {
	fmt.Println("=== Goroutine泄漏示例 ===")

	leak := func() <-chan int {
		ch := make(chan int)
		go func() {
			ch <- 42
		}()
		return ch
	}

	ch := leak()
	select {
	case val := <-ch:
		fmt.Println("接收到值:", val)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("超时，可能导致goroutine泄漏")
	}
}

func GoroutineNumberDemo() {
	fmt.Println("=== Goroutine数量监控 ===")

	fmt.Printf("初始goroutine数量: %d\n", runtime.NumGoroutine())

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("Goroutine %d 运行中, 当前总数: %d\n", id, runtime.NumGoroutine())
		}(i)
	}

	wg.Wait()
	fmt.Printf("最终goroutine数量: %d\n", runtime.NumGoroutine())
}

func GoschedDemo() {
	fmt.Println("=== Gosched让出CPU示例 ===")

	fmt.Println("不使用Gosched:")
	for i := 0; i < 3; i++ {
		go func(id int) {
			fmt.Printf("Goroutine %d\n", id)
		}(i)
	}
	time.Sleep(10 * time.Millisecond)

	fmt.Println("\n使用Gosched:")
	for i := 0; i < 3; i++ {
		go func(id int) {
			runtime.Gosched()
			fmt.Printf("Goroutine %d (Gosched)\n", id)
		}(i)
	}
	time.Sleep(10 * time.Millisecond)
}

func GOMAXPROCSDemo() {
	fmt.Println("=== GOMAXPROCS示例 ===")

	oldProcs := runtime.GOMAXPROCS(0)
	fmt.Printf("当前CPU核心数: %d\n", oldProcs)
	fmt.Printf("逻辑CPU数: %d\n", runtime.NumCPU())

	runtime.GOMAXPROCS(2)
	fmt.Printf("设置GOMAXPROCS为2, 当前值: %d\n", runtime.GOMAXPROCS(0))

	runtime.GOMAXPROCS(oldProcs)
}

func OnceDemo() {
	fmt.Println("=== sync.Once示例 ===")

	var once sync.Once
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			once.Do(func() {
				fmt.Printf("这个函数只会执行一次, 调用者: %d\n", id)
			})
			fmt.Printf("Goroutine %d 执行完成\n", id)
		}(i)
	}

	wg.Wait()
}

// MutexDemo 演示互斥锁保护共享资源
func MutexDemo() {
	fmt.Println("=== 互斥锁示例 ===")

	var mu sync.Mutex
	counter := 0
	var wg sync.WaitGroup

	// 不使用锁的情况
	fmt.Println("不使用锁:")
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++
		}()
	}
	wg.Wait()
	fmt.Printf("期望值: 1000, 实际值: %d\n", counter)

	// 使用锁的情况
	counter = 0
	fmt.Println("\n使用锁:")
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("期望值: 1000, 实际值: %d\n", counter)
}

// RWMutexDemo 演示读写锁
func RWMutexDemo() {
	fmt.Println("=== 读写锁示例 ===")

	var rwmu sync.RWMutex
	data := make(map[string]int)
	var wg sync.WaitGroup

	// 写操作
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			rwmu.Lock()
			data[fmt.Sprintf("key%d", id)] = id
			fmt.Printf("写入: key%d = %d\n", id, id)
			rwmu.Unlock()
		}(i)
	}

	time.Sleep(10 * time.Millisecond)

	// 读操作
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			rwmu.RLock()
			if v, ok := data[fmt.Sprintf("key%d", id)]; ok {
				fmt.Printf("读取: key%d = %d\n", id, v)
			}
			rwmu.RUnlock()
		}(i)
	}

	wg.Wait()
}

// AtomicDemo 演示原子操作
func AtomicDemo() {
	fmt.Println("=== 原子操作示例 ===")

	var counter int64
	var wg sync.WaitGroup

	// 使用原子操作
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()
	fmt.Printf("原子计数器: %d\n", counter)

	// 原子比较交换
	var value int64 = 100
	swapped := atomic.CompareAndSwapInt64(&value, 100, 200)
	fmt.Printf("CAS操作: swapped=%v, value=%d\n", swapped, value)

	// 原子加载和存储
	atomic.StoreInt64(&value, 300)
	loaded := atomic.LoadInt64(&value)
	fmt.Printf("原子加载: %d\n", loaded)
}

// ContextCancel 演示使用context取消goroutine
func ContextCancel() {
	fmt.Println("=== Context取消示例 ===")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// 启动多个worker
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Worker %d: 收到取消信号，退出\n", id)
					return
				default:
					fmt.Printf("Worker %d: 工作中...\n", id)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}(i)
	}

	// 主goroutine等待一段时间后取消
	time.Sleep(300 * time.Millisecond)
	fmt.Println("主函数: 发送取消信号")
	cancel()

	wg.Wait()
	fmt.Println("所有worker已退出")
}

// ContextTimeout 演示使用context超时控制
func ContextTimeout() {
	fmt.Println("=== Context超时示例 ===")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	result := make(chan string, 1)

	// 模拟耗时操作
	go func() {
		time.Sleep(300 * time.Millisecond) // 模拟耗时操作
		result <- "操作完成"
	}()

	select {
	case res := <-result:
		fmt.Println("结果:", res)
	case <-ctx.Done():
		fmt.Println("操作超时:", ctx.Err())
	}
}

// WorkerPool 演示工作池模式
func WorkerPool() {
	fmt.Println("=== 工作池模式 ===")

	jobs := make(chan int, 20)
	results := make(chan int, 20)

	// 启动3个worker
	worker := func(id int, jobs <-chan int, results chan<- int) {
		for j := range jobs {
			fmt.Printf("Worker %d: 处理任务 %d\n", id, j)
			time.Sleep(50 * time.Millisecond)
			results <- j * 2
		}
	}

	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(w)
	}

	// 发送任务
	go func() {
		for j := 1; j <= 10; j++ {
			jobs <- j
		}
		close(jobs)
	}()

	// 等待所有worker完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	for result := range results {
		fmt.Printf("结果: %d\n", result)
	}
}

// Pipeline 演示管道模式
func Pipeline() {
	fmt.Println("=== 管道模式 ===")

	// 第一阶段: 生成数据
	generator := func(done <-chan struct{}, nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				select {
				case out <- n:
				case <-done:
					return
				}
			}
		}()
		return out
	}

	// 第二阶段: 处理数据
	square := func(done <-chan struct{}, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				select {
				case out <- n * n:
				case <-done:
					return
				}
			}
		}()
		return out
	}

	// 第三阶段: 消费数据
	consumer := func(done <-chan struct{}, in <-chan int) {
		for n := range in {
			select {
			case <-done:
				return
			default:
				fmt.Printf("结果: %d\n", n)
			}
		}
	}

	done := make(chan struct{})
	defer close(done)

	nums := generator(done, 1, 2, 3, 4, 5)
	squared := square(done, nums)
	consumer(done, squared)
}

// FanOutFanIn 演示扇出扇入模式
func FanOutFanIn() {
	fmt.Println("=== 扇出扇入模式 ===")

	// 生产者
	producer := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				out <- n
			}
		}()
		return out
	}

	// worker
	worker := func(name string, in <-chan int) <-chan string {
		out := make(chan string)
		go func() {
			defer close(out)
			for n := range in {
				time.Sleep(30 * time.Millisecond)
				out <- fmt.Sprintf("%s: %d -> %d", name, n, n*n)
			}
		}()
		return out
	}

	// 合并器
	merger := func(channels ...<-chan string) <-chan string {
		var wg sync.WaitGroup
		out := make(chan string)

		output := func(c <-chan string) {
			defer wg.Done()
			for s := range c {
				out <- s
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go output(c)
		}

		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}

	// 扇出: 多个worker读取同一输入
	input := producer(1, 2, 3, 4, 5)
	c1 := worker("Worker1", input)
	c2 := worker("Worker2", input)
	c3 := worker("Worker3", input)

	// 扇入: 合并多个输出
	for result := range merger(c1, c2, c3) {
		fmt.Println(result)
	}
}

// TimeoutPattern 演示超时模式
func TimeoutPattern() {
	fmt.Println("=== 超时模式 ===")

	// 模拟慢操作
	slowOperation := func() <-chan string {
		out := make(chan string)
		go func() {
			defer close(out)
			time.Sleep(200 * time.Millisecond)
			out <- "操作完成"
		}()
		return out
	}

	// 使用select实现超时
	select {
	case result := <-slowOperation():
		fmt.Println("结果:", result)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("操作超时!")
	}
}

// GracefulShutdown 演示优雅关闭
func GracefulShutdown() {
	fmt.Println("=== 优雅关闭示例 ===")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup

	// 启动多个服务
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("服务 %d: 正在关闭...\n", id)
					time.Sleep(50 * time.Millisecond) // 模拟清理工作
					fmt.Printf("服务 %d: 已关闭\n", id)
					return
				default:
					fmt.Printf("服务 %d: 处理请求\n", id)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("所有服务已优雅关闭")
}

// GoroutineLocal 演示goroutine本地存储（使用map模拟）
func GoroutineLocal() {
	fmt.Println("=== Goroutine本地存储模拟 ===")

	var mu sync.Mutex
	storage := make(map[int]interface{})

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 设置本地值
			mu.Lock()
			storage[id] = fmt.Sprintf("data-for-goroutine-%d", id)
			mu.Unlock()

			time.Sleep(10 * time.Millisecond)

			// 获取本地值
			mu.Lock()
			if v, ok := storage[id]; ok {
				fmt.Printf("Goroutine %d: 本地值 = %v\n", id, v)
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
}

// SelectDemo 演示select多路复用
func SelectDemo() {
	fmt.Println("=== Select多路复用示例 ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- "来自ch1"
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch2 <- "来自ch2"
	}()

	// 使用select等待多个channel
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("收到:", msg1)
		case msg2 := <-ch2:
			fmt.Println("收到:", msg2)
		}
	}
}

// TimerDemo 演示定时器
func TimerDemo() {
	fmt.Println("=== 定时器示例 ===")

	// 创建定时器
	timer := time.NewTimer(200 * time.Millisecond)
	fmt.Println("定时器已启动")

	<-timer.C
	fmt.Println("定时器触发")

	// 停止定时器
	timer2 := time.NewTimer(500 * time.Millisecond)
	go func() {
		<-timer2.C
		fmt.Println("这个不会执行")
	}()

	stopped := timer2.Stop()
	fmt.Printf("定时器已停止: %v\n", stopped)

	// 重置定时器
	timer.Reset(100 * time.Millisecond)
	<-timer.C
	fmt.Println("重置后的定时器触发")
}

// TickerDemo 演示定时器周期执行
func TickerDemo() {
	fmt.Println("=== Ticker周期执行示例 ===")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan bool)

	go func() {
		time.Sleep(350 * time.Millisecond)
		done <- true
	}()

	count := 0
	for {
		select {
		case <-done:
			fmt.Printf("Ticker执行了 %d 次\n", count)
			return
		case <-ticker.C:
			count++
			fmt.Printf("Ticker触发第 %d 次\n", count)
		}
	}
}
