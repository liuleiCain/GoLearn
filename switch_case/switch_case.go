package switch_case

import (
	"fmt"
	"runtime"
	"time"
)

func BasicSwitch() {
	fmt.Println("=== 基础Switch示例 ===")

	for i := 0; i < 4; i++ {
		switch i {
		case 0:
			fmt.Printf("i=%d: case 0\n", i)
		case 1:
			fmt.Printf("i=%d: case 1\n", i)
		case 2:
			fmt.Printf("i=%d: case 2\n", i)
		default:
			fmt.Printf("i=%d: default case\n", i)
		}
	}
}

func FallthroughDemo() {
	fmt.Println("=== Fallthrough穿透示例 ===")

	fmt.Println("fallthrough会穿透到下一个case，不判断条件:")
	for i := 0; i < 4; i++ {
		switch i {
		case 0:
			fmt.Printf("i=%d: case 0 -> fallthrough\n", i)
			fallthrough
		case 1:
			fmt.Printf("i=%d: case 1 (从case 0穿透过来)\n", i)
		case 2:
			fmt.Printf("i=%d: case 2\n", i)
			fallthrough
		case 3:
			fmt.Printf("i=%d: case 3 (从case 2穿透过来)\n", i)
			fallthrough
		default:
			fmt.Printf("i=%d: default (从case 3穿透过来)\n", i)
		}
		fmt.Println("---")
	}
}

func MultiCaseDemo() {
	fmt.Println("=== 多条件Case示例 ===")

	for i := 0; i < 6; i++ {
		switch i {
		case 0, 1:
			fmt.Printf("i=%d: case 0 或 1\n", i)
		case 2, 3:
			fmt.Printf("i=%d: case 2 或 3\n", i)
		case 4, 5:
			fmt.Printf("i=%d: case 4 或 5\n", i)
		default:
			fmt.Printf("i=%d: default\n", i)
		}
	}
}

func NoExpressionSwitch() {
	fmt.Println("=== 无表达式Switch示例 ===")

	num := 42
	switch {
	case num < 0:
		fmt.Printf("%d 是负数\n", num)
	case num == 0:
		fmt.Printf("%d 是零\n", num)
	case num > 0 && num < 100:
		fmt.Printf("%d 是正数且小于100\n", num)
	default:
		fmt.Printf("%d 大于等于100\n", num)
	}

	score := 85
	switch {
	case score >= 90:
		fmt.Printf("分数 %d: 优秀\n", score)
	case score >= 80:
		fmt.Printf("分数 %d: 良好\n", score)
	case score >= 70:
		fmt.Printf("分数 %d: 中等\n", score)
	case score >= 60:
		fmt.Printf("分数 %d: 及格\n", score)
	default:
		fmt.Printf("分数 %d: 不及格\n", score)
	}
}

func TypeSwitch() {
	fmt.Println("=== 类型Switch示例 ===")

	checkType := func(x interface{}) {
		switch v := x.(type) {
		case nil:
			fmt.Printf("类型: nil\n")
		case int:
			fmt.Printf("类型: int, 值: %d\n", v)
		case string:
			fmt.Printf("类型: string, 值: %s\n", v)
		case bool:
			fmt.Printf("类型: bool, 值: %t\n", v)
		case []int:
			fmt.Printf("类型: []int, 长度: %d\n", len(v))
		case map[string]int:
			fmt.Printf("类型: map[string]int, 长度: %d\n", len(v))
		default:
			fmt.Printf("未知类型: %T\n", v)
		}
	}

	checkType(nil)
	checkType(42)
	checkType("hello")
	checkType(true)
	checkType([]int{1, 2, 3})
	checkType(map[string]int{"a": 1})
	checkType(3.14)
}

func SwitchBreak() {
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
	}

	fmt.Println("\n使用标签break跳出外层循环:")
outerLoop:
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			fmt.Printf("i=%d: 继续外层循环\n", i)
		case 1:
			fmt.Printf("i=%d: 跳出外层循环\n", i)
			break outerLoop
		}
	}
}

func SwitchInFor() {
	fmt.Println("=== for循环中的Switch示例 ===")

	fmt.Println("使用标签continue跳过当前迭代:")
outerLoop:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch j {
			case 1:
				fmt.Printf("i=%d, j=%d: continue外层循环\n", i, j)
				continue outerLoop
			default:
				fmt.Printf("i=%d, j=%d: 正常处理\n", i, j)
			}
		}
	}
}

func SwitchInitialization() {
	fmt.Println("=== Switch初始化语句示例 ===")

	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Printf("运行在 macOS 上 (%s)\n", os)
	case "linux":
		fmt.Printf("运行在 Linux 上 (%s)\n", os)
	case "windows":
		fmt.Printf("运行在 Windows 上 (%s)\n", os)
	default:
		fmt.Printf("运行在其他系统上 (%s)\n", os)
	}

	switch now := time.Now(); now.Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Printf("今天是周末 (%s)\n", now.Weekday())
	default:
		fmt.Printf("今天是工作日 (%s)\n", now.Weekday())
	}
}

func SwitchPitfalls() {
	fmt.Println("=== Switch常见陷阱 ===")

	fmt.Println("陷阱1: case条件必须是可比较的类型")
	fmt.Println("陷阱2: fallthrough不会判断下一个case条件")
	fmt.Println("陷阱3: 无表达式switch的case必须是布尔表达式")

	fmt.Println("\n正确示例 - 无表达式switch:")
	x := 10
	switch {
	case x > 5:
		fmt.Printf("x=%d > 5\n", x)
		fallthrough
	case x > 20:
		fmt.Printf("注意: fallthrough不判断条件，x=%d 不大于20但仍执行\n", x)
	}

	fmt.Println("\n陷阱: fallthrough后不能是最后一个case")
	fmt.Println("错误: case x: fallthrough 后面没有case会编译错误")
}
