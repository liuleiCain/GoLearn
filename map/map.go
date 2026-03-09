package learn_map

import "fmt"

func InitMap() {
	// 增加recover
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic err=", err)
		}
	}()
	// 初始化的map
	b := make(map[string]int)
	b["test"] = 1
	fmt.Println("设置b成功")
	// 只定义没有初始化的map
	var a map[string]int
	a["aaa"] = 1
	fmt.Println("设置a成功")
}

// RangeMap 演示map遍历的顺序
func RangeMap() {
	// 创建一个map
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}

	// 多次遍历map，观察顺序是否一致
	fmt.Println("第一次遍历:")
	for k, v := range m {
		fmt.Printf("key: %s, value: %d\n", k, v)
	}

	fmt.Println("\n第二次遍历:")
	for k, v := range m {
		fmt.Printf("key: %s, value: %d\n", k, v)
	}

	fmt.Println("\n第三次遍历:")
	for k, v := range m {
		fmt.Printf("key: %s, value: %d\n", k, v)
	}

	// 注意：Go语言中map的遍历顺序是随机的，每次运行可能会得到不同的顺序
	// 这是Go语言的设计决定，目的是防止程序员依赖map的遍历顺序
}
