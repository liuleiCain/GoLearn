package goroutine

import (
	"fmt"
	"runtime"
	"sync"
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
