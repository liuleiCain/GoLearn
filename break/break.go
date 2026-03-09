package _break

import (
	"fmt"
	"time"
)

func BasicBreak() {
	fmt.Println("=== 基础Break示例 ===")

	fmt.Println("普通for循环中的break:")
	for i := 0; i < 10; i++ {
		if i == 5 {
			fmt.Printf("i=%d时break跳出循环\n", i)
			break
		}
		fmt.Printf("i=%d\n", i)
	}
	fmt.Println("循环结束")

	fmt.Println("\n遍历切片时的break:")
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, n := range nums {
		if n > 5 {
			fmt.Printf("找到第一个大于5的数: %d，停止遍历\n", n)
			break
		}
		fmt.Printf("检查: %d\n", n)
	}
}

func BreakInSwitch() {
	fmt.Println("=== Switch中的Break示例 ===")

	fmt.Println("switch中的break跳出switch块:")
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			fmt.Printf("i=%d: case 0\n", i)
			break
		case 1:
			fmt.Printf("i=%d: case 1\n", i)
		case 2:
			fmt.Printf("i=%d: case 2\n", i)
		}
		fmt.Printf("i=%d: switch后的代码\n", i)
	}

	fmt.Println("\n注意: Go的switch默认不需要break")
	fmt.Println("显式break只在需要提前退出case时有用")
}

func BreakVsContinue() {
	fmt.Println("=== Break vs Continue对比 ===")

	fmt.Println("使用break - 遇到条件完全退出循环:")
	for i := 0; i < 5; i++ {
		if i == 3 {
			fmt.Printf("i=%d时break，退出循环\n", i)
			break
		}
		fmt.Printf("处理 i=%d\n", i)
	}

	fmt.Println("\n使用continue - 跳过当前迭代继续循环:")
	for i := 0; i < 5; i++ {
		if i == 3 {
			fmt.Printf("i=%d时continue，跳过本次迭代\n", i)
			continue
		}
		fmt.Printf("处理 i=%d\n", i)
	}

	fmt.Println("\n对比总结:")
	fmt.Println("  break: 完全退出循环")
	fmt.Println("  continue: 跳过当前迭代，继续下一次")
}

func BreakForSelect() {
	fmt.Println("=== For-Select中的Break陷阱 ===")

	fmt.Println("问题: select中的break只跳出select，不跳出for循环")
	fmt.Println("演示:")

	i := 0
	for {
		select {
		case <-time.After(30 * time.Millisecond):
			i++
			fmt.Printf("  ping %d\n", i)
			if i >= 5 {
				fmt.Println("  达到条件，跳出了select块")
				break
			}
		case <-time.After(100 * time.Millisecond):
			fmt.Println("  超时退出演示")
			break
		}
		if i >= 6 {
			fmt.Println("  达到条件，真正退出for循环")
			break
		}
	}
	fmt.Println("结论: select中的break只跳出了select块，需要额外的break退出for")
}

func BreakForSelectLabel() {
	fmt.Println("=== 使用标签跳出For-Select ===")

	fmt.Println("解决方案: 使用标签break跳出外层for循环")
	fmt.Println("演示:")

	i := 0
loopLabel:
	for {
		select {
		case <-time.After(30 * time.Millisecond):
			i++
			fmt.Printf("  ping %d\n", i)
			if i >= 2 {
				fmt.Println("  达到条件，break loopLabel跳出for循环")
				break loopLabel
			}
		case <-time.After(150 * time.Millisecond):
			fmt.Println("  超时退出")
			break loopLabel
		}
	}
	fmt.Println("  已退出for循环")
}

func BreakNestedLoop() {
	fmt.Println("=== 多重嵌套循环的Break ===")

	fmt.Println("问题: 普通break只能跳出最内层循环")
	fmt.Println("演示:")
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				fmt.Printf("  i=%d, j=%d时break，只跳出内层循环\n", i, j)
				break
			}
			fmt.Printf("  内层: i=%d, j=%d\n", i, j)
		}
		fmt.Printf("  外层循环继续: i=%d\n", i)
	}

	fmt.Println("\n解决方案: 使用标签跳出指定层级的循环")
	fmt.Println("演示:")
outerLoop:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				fmt.Printf("  i=%d, j=%d时break outerLoop，跳出外层循环\n", i, j)
				break outerLoop
			}
			fmt.Printf("  内层: i=%d, j=%d\n", i, j)
		}
	}
	fmt.Println("  已跳出所有循环")
}

func BreakWithLabelNames() {
	fmt.Println("=== 标签命名最佳实践 ===")

	fmt.Println("推荐使用有意义的标签名:")
	fmt.Println("  - outerLoop: 外层循环")
	fmt.Println("  - readLoop: 读取循环")
	fmt.Println("  - processLoop: 处理循环")
	fmt.Println("  - mainLoop: 主循环")

	fmt.Println("\n示例 - 搜索二维数组:")
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	target := 5

searchLoop:
	for i, row := range matrix {
		for j, val := range row {
			if val == target {
				fmt.Printf("  找到目标值%d，位置: [%d][%d]\n", target, i, j)
				break searchLoop
			}
		}
	}
}

func BreakVsReturn() {
	fmt.Println("=== Break vs Return对比 ===")

	processWithBreak := func(items []int, max int) int {
		sum := 0
		for _, item := range items {
			if item > max {
				fmt.Printf("  遇到%d > %d，break退出循环\n", item, max)
				break
			}
			sum += item
		}
		fmt.Printf("  break后继续执行，sum=%d\n", sum)
		return sum
	}

	processWithReturn := func(items []int, max int) int {
		sum := 0
		for _, item := range items {
			if item > max {
				fmt.Printf("  遇到%d > %d，return退出函数\n", item, max)
				return sum
			}
			sum += item
		}
		fmt.Println("  这行不会执行")
		return sum
	}

	items := []int{1, 2, 3, 10, 4, 5}
	fmt.Println("使用break:")
	result1 := processWithBreak(items, 5)
	fmt.Printf("  最终结果: %d\n", result1)

	fmt.Println("\n使用return:")
	result2 := processWithReturn(items, 5)
	fmt.Printf("  最终结果: %d\n", result2)

	fmt.Println("\n对比总结:")
	fmt.Println("  break: 退出循环，继续执行函数剩余代码")
	fmt.Println("  return: 退出整个函数")
}

func InfiniteLoopBreak() {
	fmt.Println("=== 无限循环中的Break ===")

	fmt.Println("无限循环 for{} 配合break实现条件退出:")

	count := 0
	for {
		count++
		fmt.Printf("  迭代 %d\n", count)

		if count >= 3 {
			fmt.Println("  达到条件，break退出无限循环")
			break
		}

		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println("  无限循环已退出")
}

func BreakPitfalls() {
	fmt.Println("=== Break常见陷阱 ===")

	fmt.Println("陷阱1: 误以为break能跳出所有嵌套循环")
	fmt.Println("  错误理解: break会跳出所有循环")
	fmt.Println("  正确行为: break只跳出最内层循环")

	fmt.Println("\n陷阱2: select中的break陷阱")
	fmt.Println("  错误理解: break能跳出for-select")
	fmt.Println("  正确行为: break只跳出select，需要标签才能跳出for")

	fmt.Println("\n陷阱3: 标签作用域")
	fmt.Println("  标签必须定义在break所在的循环外")
	fmt.Println("  不能跨函数使用标签")

	fmt.Println("\n陷阱4: break在switch中的冗余")
	fmt.Println("  Go的switch默认break，不需要显式break")
	fmt.Println("  除非需要提前退出case块")
}
