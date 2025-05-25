package for_range

import (
	"time"
)

func BreakForSelect() {
	exit := make(chan interface{})

	go func() {
		i := 1
		for {
			select {
			case <-time.After(time.Millisecond * 100):
				i++
				println("ping....", i)
			case <-exit:
				println("receive exit", i)
				break
			}
			// 防止死循环
			if i > 20 {
				break
			}
		}
		println("退出for循环...", i)
	}()

	time.Sleep(time.Second * 1)
	exit <- struct{}{}

	// 延迟1.5秒，确保i可以达到最大值退出
	time.Sleep(time.Millisecond * 1500)
}

func BreakForSelectLabel() {
	exit := make(chan interface{})

	go func() {
		i := 1
	loopLabel:
		for {
			select {
			case <-time.After(time.Millisecond * 100):
				i++
				println("ping....", i)
			case <-exit:
				println("receive exit", i)
				break loopLabel
			}
			// 防止死循环
			if i > 20 {
				break
			}
		}
		println("退出for循环...", i)
	}()

	time.Sleep(time.Second * 1)
	exit <- struct{}{}

	// 延迟1.5秒，确保i可以达到最大值退出
	time.Sleep(time.Millisecond * 1500)
}

func BreakForLabel() {
	for k := 0; k < 1; k++ {
		println("第", k, "次循环")
	forLoop:
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				println("i=", i, "j=", j)
				if i == 1 && j == 1 {
					break forLoop
				}
			}
		}
	}
}
