package closure

import (
	"fmt"
	"sync"
	"time"
)

func BasicClosure() {
	fmt.Println("=== 基础闭包示例 ===")

	x := 10

	increment := func() int {
		x++
		return x
	}

	fmt.Println("第一次调用:", increment())
	fmt.Println("第二次调用:", increment())
	fmt.Println("第三次调用:", increment())
	fmt.Println("外部x值:", x)
}

func ClosureFactory() {
	fmt.Println("=== 闭包工厂示例 ===")

	adder := func(base int) func(int) int {
		return func(x int) int {
			return base + x
		}
	}

	add5 := adder(5)
	add10 := adder(10)

	fmt.Println("add5(3):", add5(3))
	fmt.Println("add5(7):", add5(7))
	fmt.Println("add10(3):", add10(3))

	multiplier := func(factor int) func(int) int {
		return func(x int) int {
			return x * factor
		}
	}

	double := multiplier(2)
	triple := multiplier(3)

	fmt.Println("double(5):", double(5))
	fmt.Println("triple(5):", triple(5))
}

func LoopTrap() {
	fmt.Println("=== 循环变量捕获陷阱 ===")

	fmt.Println("问题示例 - 捕获变量引用:")
	var funcs []func()
	for i := 0; i < 3; i++ {
		funcs = append(funcs, func() {
			fmt.Printf("i = %d\n", i)
		})
	}
	for _, f := range funcs {
		f()
	}

	fmt.Println("\n解决方案1 - 传递参数:")
	var funcs2 []func()
	for i := 0; i < 3; i++ {
		funcs2 = append(funcs2, func(n int) func() {
			return func() {
				fmt.Printf("n = %d\n", n)
			}
		}(i))
	}
	for _, f := range funcs2 {
		f()
	}

	fmt.Println("\n解决方案2 - 创建局部变量:")
	var funcs3 []func()
	for i := 0; i < 3; i++ {
		i := i
		funcs3 = append(funcs3, func() {
			fmt.Printf("i(local) = %d\n", i)
		})
	}
	for _, f := range funcs3 {
		f()
	}
}

func GoroutineTrap() {
	fmt.Println("=== Goroutine中的闭包陷阱 ===")

	fmt.Println("问题示例:")
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("goroutine i = %d\n", i)
		}()
	}
	wg.Wait()

	fmt.Println("\n解决方案 - 传递参数:")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Printf("goroutine n = %d\n", n)
		}(i)
	}
	wg.Wait()

	fmt.Println("\nGo 1.22+ 自动修复此问题")
}

func ClosureState() {
	fmt.Println("=== 闭包状态封装示例 ===")

	counter := func() func() int {
		count := 0
		return func() int {
			count++
			return count
		}
	}()

	fmt.Println("计数器:", counter())
	fmt.Println("计数器:", counter())
	fmt.Println("计数器:", counter())

	bank := func() (func(int), func() int) {
		balance := 0
		deposit := func(amount int) {
			balance += amount
		}
		getBalance := func() int {
			return balance
		}
		return deposit, getBalance
	}

	deposit, getBalance := bank()
	deposit(100)
	fmt.Println("余额:", getBalance())
	deposit(50)
	fmt.Println("余额:", getBalance())
}

func ClosureRecursion() {
	fmt.Println("=== 闭包递归示例 ===")

	var fib func(int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}

	for i := 0; i <= 10; i++ {
		fmt.Printf("fib(%d) = %d\n", i, fib(i))
	}
}

func ClosurePerformance() {
	fmt.Println("=== 闭包性能考量 ===")

	directCall := func(x int) int {
		return x * 2
	}

	closureCall := func() func(int) int {
		factor := 2
		return func(x int) int {
			return x * factor
		}
	}()

	fmt.Println("直接调用:", directCall(5))
	fmt.Println("闭包调用:", closureCall(5))

	fmt.Println("\n性能说明:")
	fmt.Println("- 闭包有轻微的性能开销")
	fmt.Println("- 每次创建闭包会分配堆内存")
	fmt.Println("- 在性能敏感场景考虑直接传参")
}

func ClosureDefer() {
	fmt.Println("=== defer中的闭包 ===")

	x := 1
	defer func() {
		fmt.Println("defer闭包捕获x:", x)
	}()

	defer func(v int) {
		fmt.Println("defer参数传递x:", v)
	}(x)

	x = 100
	fmt.Println("修改x为:", x)
}

func ClosureTimer() {
	fmt.Println("=== 闭包实现计时器 ===")

	timer := func(name string) func() {
		start := time.Now()
		return func() {
			fmt.Printf("%s 耗时: %v\n", name, time.Since(start))
		}
	}

	defer timer("主函数")()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("执行一些操作...")
	time.Sleep(50 * time.Millisecond)
}
