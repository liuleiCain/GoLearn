package learn_defer

import "fmt"

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
