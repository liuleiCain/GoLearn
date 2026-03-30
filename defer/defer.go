package learn_defer

import (
	"fmt"
	"sync"
	"time"
)

func DeferFor() {
	for i := 1; i < 10; i++ {
		// 直接使用defer容易内存泄露，建议使用闭包
		defer fmt.Println("第1轮打印", i)
	}
	for i := 1; i < 10; i++ {
		func(i2 int) {
			defer fmt.Println("第2轮打印", i2)
		}(i)
	}
	for i := 1; i < 10; i++ {
		func() {
			defer fmt.Println("第3轮打印", i)
		}()
	}
}

func DeferForParams() {
	x := 1
	defer fmt.Println("defer x=", x)
	defer func(x2 int) {
		fmt.Println("defer2 x=", x2)
	}(x)
	defer func() {
		fmt.Println("defer3 x=", x)
	}()
	x++
}

func DeferForReturn() int {
	x := 1
	defer func() {
		x += 1
		fmt.Println("defer x=", x)
	}()
	defer func() {
		x += 2
		fmt.Println("defer2 x=", x)
	}()
	defer func() {
		x += 3
		fmt.Println("defer3 x=", x)
	}()
	fmt.Println("返回结果为x初始值")
	return x
}
func DeferForReturn2() (x int) {
	x = 1
	defer func() {
		x++
		fmt.Println("defer x=", x)
	}()
	defer func() {
		x++
		fmt.Println("defer2 x=", x)
	}()
	defer func() {
		x++
		fmt.Println("defer3 x=", x)
	}()
	fmt.Println("返回结果为defer计算完的x值")
	return x
}

func DeferForPanic() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic err=", err)
		}
	}()
	var m map[string]bool
	defer fmt.Println("panic前的defer1")
	defer func() {
		fmt.Println("panic前的闭包defer")
	}()
	defer fmt.Println("panic前的defer2")
	m["a"] = true
	defer fmt.Println("panic后的defer")
}

func DeferForOldParams() {
	x := 1
	oldX := x
	fmt.Println("初始x=", x)
	defer func() {
		fmt.Println("defer x=", x, "oldX=", oldX)
		x = oldX
		fmt.Println("defer2 x=", x, "oldX=", oldX)
	}()
	x++
	fmt.Println("改变后x=", x)
}

// DeferTiming 演示defer的执行时机
// defer在外层函数返回前执行，包括return语句赋值之后
func DeferTiming() {
	fmt.Println("=== defer执行时机演示 ===")

	// 示例1: defer在return之后执行
	result := deferTimingHelper()
	fmt.Printf("函数返回值: %d\n", result)
}

func deferTimingHelper() int {
	fmt.Println("函数开始执行")
	defer fmt.Println("defer执行: 在return之后")
	fmt.Println("函数即将返回")
	return 42
}

// DeferStackTrace 演示defer在panic时打印堆栈信息
func DeferStackTrace() {
	fmt.Println("=== defer打印堆栈信息 ===")
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("捕获到panic: %v\n", err)
		}
	}()
	panic("主动触发panic")
}

// DeferMutexUnlock 演示defer用于互斥锁解锁
// 这是defer最常见的使用场景之一
func DeferMutexUnlock() {
	fmt.Println("=== defer解锁互斥锁 ===")
	var mu sync.Mutex
	var counter int

	// 使用defer确保锁一定会被释放
	increment := func() {
		mu.Lock()
		defer mu.Unlock() // 无论是否发生panic，都会解锁
		counter++
		fmt.Printf("counter增加到: %d\n", counter)
	}

	for i := 0; i < 3; i++ {
		increment()
	}
	fmt.Printf("最终counter值: %d\n", counter)
}

// DeferFileOperation 模拟文件操作的defer使用
func DeferFileOperation() {
	fmt.Println("=== defer文件操作模拟 ===")

	// 模拟文件打开和关闭
	openFile := func(name string) *mockFile {
		fmt.Printf("打开文件: %s\n", name)
		return &mockFile{name: name}
	}

	// 使用defer确保资源释放
	processFile := func(name string) {
		f := openFile(name)
		defer f.close() // 确保文件关闭

		fmt.Printf("处理文件内容: %s\n", name)
		// 模拟处理过程
	}

	processFile("test.txt")
	processFile("data.json")
}

type mockFile struct {
	name string
}

func (f *mockFile) close() {
	fmt.Printf("关闭文件: %s\n", f.name)
}

// DeferMethod 演示defer调用方法
func DeferMethod() {
	fmt.Println("=== defer调用方法 ===")

	p := &person{name: "张三", age: 25}
	defer p.sayHello() // defer可以调用方法
	defer p.grow()     // 多个defer按LIFO顺序执行

	fmt.Println("函数主体执行中...")
}

type person struct {
	name string
	age  int
}

func (p *person) sayHello() {
	fmt.Printf("你好，我是%s，今年%d岁\n", p.name, p.age)
}

func (p *person) grow() {
	p.age++
	fmt.Printf("%s长大了一岁，现在%d岁\n", p.name, p.age)
}

// DeferPerformance 演示defer的性能影响
func DeferPerformance() {
	fmt.Println("=== defer性能对比 ===")

	// 不使用defer
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		withoutDefer()
	}
	withoutDeferTime := time.Since(start)

	// 使用defer
	start = time.Now()
	for i := 0; i < 1000000; i++ {
		withDefer()
	}
	withDeferTime := time.Since(start)

	fmt.Printf("不使用defer: %v\n", withoutDeferTime)
	fmt.Printf("使用defer: %v\n", withDeferTime)
	fmt.Printf("性能差异: %.2f%%\n", float64(withDeferTime-withoutDeferTime)/float64(withoutDeferTime)*100)
}

func withoutDefer() {
	x := 1
	x++
	_ = x
}

func withDefer() {
	x := 1
	defer func() { _ = x }()
	x++
}

// DeferNilFunction 演示defer nil函数会panic
func DeferNilFunction() {
	fmt.Println("=== defer nil函数演示 ===")
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("捕获到panic: %v\n", err)
		}
	}()

	var fn func()
	defer fn() // fn是nil，panic在这行触发
	fmt.Println("这行会执行，因为panic还没发生")
}

// DeferArguments 演示defer参数的计算时机
func DeferArguments() {
	fmt.Println("=== defer参数计算时机 ===")

	// 参数在defer声明时就已经计算
	x := 1
	defer fmt.Printf("defer1: x=%d\n", x) // 输出1

	x = 2
	defer fmt.Printf("defer2: x=%d\n", x) // 输出2

	// 闭包捕获的是引用
	x = 3
	defer func() {
		fmt.Printf("defer3(闭包): x=%d\n", x) // 输出4
	}()

	x = 4
	fmt.Printf("函数末尾: x=%d\n", x)
}

// DeferNamedReturnMultiple 演示多个命名返回值与defer
func DeferNamedReturnMultiple() (a int, b string, c bool) {
	fmt.Println("=== 多命名返回值与defer ===")

	a = 10
	b = "hello"
	c = true

	defer func() {
		fmt.Printf("defer修改前: a=%d, b=%s, c=%v\n", a, b, c)
		a *= 2
		b += " world"
		c = false
		fmt.Printf("defer修改后: a=%d, b=%s, c=%v\n", a, b, c)
	}()

	fmt.Printf("返回前: a=%d, b=%s, c=%v\n", a, b, c)
	return
}
