package memorymodel

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// DataRace 展示数据竞争问题
func DataRace() {
	fmt.Println("数据竞争示例:")

	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++
		}()
	}

	wg.Wait()
	fmt.Println("  期望结果: 1000, 实际结果:", counter)
}

// AtomicOperation 展示如何使用原子操作避免数据竞争
func AtomicOperation() {
	fmt.Println("原子操作示例:")

	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()
	fmt.Println("  期望结果: 1000, 实际结果:", counter)
}

// MutexLock 展示如何使用互斥锁避免数据竞争
func MutexLock() {
	fmt.Println("互斥锁示例:")

	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

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
	fmt.Println("  期望结果: 1000, 实际结果:", counter)
}

// HappensBefore 展示 happens-before 关系
func HappensBefore() {
	fmt.Println("Happens-Before 关系示例:")

	var a string
	var done bool

	go func() {
		a = "hello, world"
		done = true
	}()

	for !done {
		time.Sleep(1 * time.Millisecond)
	}

	fmt.Println("  读取到的值:", a)
}

// ChannelSynchronization 展示使用通道进行同步
func ChannelSynchronization() {
	fmt.Println("通道同步示例:")

	var a string
	c := make(chan bool, 1)

	go func() {
		a = "hello from goroutine"
		c <- true
	}()

	<-c
	fmt.Println("  读取到的值:", a)
}

// OnceInitialization 展示 sync.Once 的使用
func OnceInitialization() {
	fmt.Println("sync.Once 示例:")

	var (
		data string
		once sync.Once
	)

	initData := func() {
		fmt.Println("  初始化数据...")
		data = "initialized data"
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			once.Do(initData)
			fmt.Println("  Goroutine", id, "读取到:", data)
		}(i)
	}

	wg.Wait()
}

// WaitGroupSynchronization 展示 sync.WaitGroup 的使用
func WaitGroupSynchronization() {
	fmt.Println("sync.WaitGroup 示例:")

	var wg sync.WaitGroup
	results := make([]string, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			results[id] = fmt.Sprintf("result from goroutine %d", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("  所有 goroutine 完成，结果:", results)
}

// MemoryOrdering 展示内存重排序问题
func MemoryOrdering() {
	fmt.Println("内存重排序示例:")

	var x, y int
	var wg sync.WaitGroup

	const iterations = 10000
	reorderCount := 0

	for i := 0; i < iterations; i++ {
		x, y = 0, 0
		var r1, r2 int

		wg.Add(2)

		go func() {
			defer wg.Done()
			x = 1
			r1 = y
		}()

		go func() {
			defer wg.Done()
			y = 1
			r2 = x
		}()

		wg.Wait()

		if r1 == 0 && r2 == 0 {
			reorderCount++
		}
	}

	fmt.Printf("  在 %d 次迭代中，检测到 %d 次可能的重排序\n", iterations, reorderCount)
}

// RWMutexExample 展示读写锁的使用
func RWMutexExample() {
	fmt.Println("读写锁示例:")

	var (
		data  string
		rwMu  sync.RWMutex
		wg    sync.WaitGroup
	)

	write := func(id int) {
		defer wg.Done()
		rwMu.Lock()
		data = fmt.Sprintf("written by %d", id)
		fmt.Printf("  Goroutine %d 写入: %s\n", id, data)
		rwMu.Unlock()
	}

	read := func(id int) {
		defer wg.Done()
		rwMu.RLock()
		fmt.Printf("  Goroutine %d 读取: %s\n", id, data)
		rwMu.RUnlock()
	}

	wg.Add(3)
	go write(1)
	go write(2)
	go write(3)
	wg.Wait()

	wg.Add(5)
	for i := 1; i <= 5; i++ {
		go read(i)
	}
	wg.Wait()
}

// AtomicValue 展示 atomic.Value 的使用
func AtomicValue() {
	fmt.Println("atomic.Value 示例:")

	var config atomic.Value

	type Config struct {
		Timeout time.Duration
		MaxConn int
	}

	loadConfig := func() {
		cfg := Config{
			Timeout: 10 * time.Second,
			MaxConn: 100,
		}
		config.Store(cfg)
		fmt.Println("  配置已加载")
	}

	updateConfig := func() {
		cfg := Config{
			Timeout: 30 * time.Second,
			MaxConn: 200,
		}
		config.Store(cfg)
		fmt.Println("  配置已更新")
	}

	readConfig := func(id int) {
		cfg := config.Load().(Config)
		fmt.Printf("  Goroutine %d 读取配置: Timeout=%v, MaxConn=%d\n", id, cfg.Timeout, cfg.MaxConn)
	}

	var wg sync.WaitGroup

	loadConfig()

	wg.Add(3)
	for i := 1; i <= 3; i++ {
		go func(id int) {
			defer wg.Done()
			readConfig(id)
		}(i)
	}
	wg.Wait()

	updateConfig()

	wg.Add(3)
	for i := 4; i <= 6; i++ {
		go func(id int) {
			defer wg.Done()
			readConfig(id)
		}(i)
	}
	wg.Wait()
}
