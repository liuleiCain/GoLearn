package for_range

import (
	"fmt"
	"sync"
	"time"
)

func AnalysisBasic() {
	fmt.Println("Start test for range")
	// go 1.22后支持 遍历数字变量
	n := 5
	//for i := range n {
	//	n = 2
	//	fmt.Println(i)
	//}
	fmt.Println("n = ", n)

	// 遍历字符串变量
	s := "hello, 你好"
	for i, k := range s {
		s = "world，世界"
		fmt.Println(i, string(k), s[i])
	}
	fmt.Println("s = ", s)
}

func AnalysisArr() {
	fmt.Println("Start test for range")
	arr := [5]int{1, 2, 3, 4, 5}
	for i, k := range arr {
		arr[i] = arr[i] + 100
		fmt.Println("index=", i, "value=", k)
	}

	fmt.Println("arr = ", arr)
}

func AnalysisSlice() {
	fmt.Println("Start test for range")
	arr := []int{1, 2, 3}
	fmt.Println("arr = ", cap(arr), len(arr))
	// 改变切片值
	for i, k := range arr {
		arr[i] = arr[i] + 100
		fmt.Println("index=", i, "value=", k)
	}
	fmt.Println("arr = ", arr)

	// 改变切片大小
	fmt.Printf("初始切片长度：%d \n", len(arr))
	for i, k := range arr {
		arr = append(arr, k+100)
		fmt.Println("index=", i, "value=", k)
	}
	fmt.Printf("切片长度变化：%d \n", len(arr))
	fmt.Println("arr = ", arr)
}

func AnalysisMap() {
	fmt.Println("Start test for range")
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	// 改变map的值
	fmt.Println("m = ", m)
	for k, v := range m {
		m["a"] = 100
		m["b"] = 101
		m["c"] = 102
		fmt.Println("key=", k, "value=", v)
	}
	fmt.Println("m = ", m)

	// 改变map的大小，多次运行结果不一样
	for i := 0; i < 10; i++ {
		m = map[string]int{"a": 1, "b": 2, "c": 3}
		fmt.Printf("第%d次执行 \n", i)
		for k, v := range m {
			delete(m, "a")
			fmt.Println("key=", k, "value=", v)
		}
		fmt.Println("m = ", m)
		fmt.Println("--------------")
	}
}

// Go 1.22 之前，所有 goroutine 打印最后一个元素；现在会正确打印每个元素。
func AnalysisRangeSliceIndex() {
	fmt.Println("Start test for range")
	arr := []int{100, 200, 300}
	var wg sync.WaitGroup
	wg.Add(len(arr) * 2)
	for i, k := range arr {
		go func() {
			defer wg.Done()
			fmt.Printf("并发---- 第%d次goroutine执行: %d \n", i, k)
		}()
		fmt.Printf("并发---- 第%d次打印: %d  \n", i, k)
	}

	for i, k := range arr {
		go func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			fmt.Println("goroutine-i-k", time.Now(), i, k)
			fmt.Printf("延时==== 第%d次goroutine执行: %d  \n", i, k)
		}()
		fmt.Println("i-k", time.Now(), i, k)
		fmt.Printf("延时==== 第%d次打印: %d  \n", i, k)
	}

	wg.Wait()
}

type Person struct {
	Name string
	Age  int
}

func AnalysisRangeStruct() {
	list := []Person{
		{Name: "张三", Age: 18},
		{Name: "李四", Age: 20},
		{Name: "王五", Age: 22},
	}
	var wg sync.WaitGroup
	wg.Add(len(list))
	// 改变遍历结构体元素变量的值
	for i, k := range list {
		go func() {
			defer wg.Done()
			k.Age = k.Age + 100
			fmt.Println("goroutine index=", i, "value=", k)
		}()
		fmt.Println("index=", i, "value=", k)
	}
	fmt.Println("list = ", list)
	wg.Wait()
	fmt.Println("执行完成后list = ", list)

	pointerList := []*Person{
		{Name: "张三", Age: 18},
		{Name: "李四", Age: 20},
		{Name: "王五", Age: 22},
	}
	wg.Add(len(pointerList))
	// 改变遍历结构体元素变量的值
	for i, k := range pointerList {
		go func() {
			defer wg.Done()
			k.Age = k.Age + 100
			fmt.Println("goroutine index=", i, "value=", *k)
		}()
		fmt.Println("index=", i, "value=", *k)
	}
	fmt.Printf("pointerList = %v \n", printList(pointerList))
	wg.Wait()
	fmt.Printf("执行完成后pointerList = %v \n", printList(pointerList))
}

func printList(list []*Person) string {
	res := ""
	for i, v := range list {
		res += fmt.Sprintf("index=%d, value=%v ", i, v)
	}
	res += fmt.Sprintf("\n")
	return res
}

// RangeChannel 演示遍历channel
func RangeChannel() {
	fmt.Println("=== 遍历Channel演示 ===")

	// 创建带缓冲的channel
	ch := make(chan int, 5)

	// 发送数据
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		close(ch) // 必须关闭channel，否则range会死锁
	}()

	// 使用range遍历channel
	for v := range ch {
		fmt.Printf("接收到: %d\n", v)
	}

	fmt.Println("Channel已关闭，遍历结束")
}

// RangeStringDetail 演示字符串遍历的细节
func RangeStringDetail() {
	fmt.Println("=== 字符串遍历细节 ===")

	s := "Hello世界"

	// 使用range遍历 - 按rune遍历
	fmt.Println("使用range遍历(按rune):")
	for i, r := range s {
		fmt.Printf("  位置: %d, rune: %c, Unicode: %U\n", i, r, r)
	}

	// 使用普通for循环 - 按字节遍历
	fmt.Println("使用普通for遍历(按字节):")
	for i := 0; i < len(s); i++ {
		fmt.Printf("  位置: %d, 字节: %x, 字符: %c\n", i, s[i], s[i])
	}

	// 转换为rune切片
	fmt.Println("转换为rune切片后遍历:")
	runes := []rune(s)
	for i, r := range runes {
		fmt.Printf("  索引: %d, rune: %c\n", i, r)
	}
}

// RangeWithBreak 演示range中使用break和continue
func RangeWithBreak() {
	fmt.Println("=== Range中的break和continue ===")

	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 使用break提前退出
	fmt.Println("找到第一个偶数后退出:")
	for _, v := range nums {
		if v%2 == 0 {
			fmt.Printf("找到偶数: %d\n", v)
			break
		}
		fmt.Printf("检查: %d\n", v)
	}

	// 使用continue跳过
	fmt.Println("\n跳过偶数:")
	for _, v := range nums {
		if v%2 == 0 {
			continue
		}
		fmt.Printf("奇数: %d\n", v)
	}

	// 使用标签break退出外层循环
	fmt.Println("\n嵌套循环中使用标签break:")
outer:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			if i*j > 4 {
				fmt.Printf("i=%d, j=%d 时退出\n", i, j)
				break outer
			}
			fmt.Printf("i=%d, j=%d\n", i, j)
		}
	}
}

// RangeNilCollection 演示遍历nil集合
func RangeNilCollection() {
	fmt.Println("=== 遍历nil集合 ===")

	// nil切片 - 安全，不执行循环体
	var nilSlice []int
	fmt.Printf("nil切片长度: %d\n", len(nilSlice))
	for i, v := range nilSlice {
		fmt.Printf("索引: %d, 值: %d\n", i, v) // 不会执行
	}
	fmt.Println("nil切片遍历完成，无输出")

	// nil map - 安全，不执行循环体
	var nilMap map[string]int
	fmt.Printf("nil map长度: %d\n", len(nilMap))
	for k, v := range nilMap {
		fmt.Printf("键: %s, 值: %d\n", k, v) // 不会执行
	}
	fmt.Println("nil map遍历完成，无输出")

	// nil channel - 会永久阻塞！
	fmt.Println("nil channel遍历会永久阻塞，不演示")
}

// RangeCopyValue 演示range值是副本
func RangeCopyValue() {
	fmt.Println("=== Range值是副本演示 ===")

	type Item struct {
		Value int
	}

	items := []Item{{1}, {2}, {3}}

	// 修改副本不影响原切片
	fmt.Println("修改副本:")
	for _, item := range items {
		item.Value *= 10
	}
	fmt.Printf("原切片: %v\n", items)

	// 通过索引修改
	fmt.Println("\n通过索引修改:")
	for i := range items {
		items[i].Value *= 10
	}
	fmt.Printf("修改后: %v\n", items)

	// 使用指针切片
	fmt.Println("\n使用指针切片:")
	ptrItems := []*Item{{1}, {2}, {3}}
	for _, item := range ptrItems {
		item.Value *= 10
	}
	fmt.Printf("指针切片修改后: ")
	for _, item := range ptrItems {
		fmt.Printf("%d ", item.Value)
	}
	fmt.Println()
}

// RangeTwoValues 演示range的不同返回值用法
func RangeTwoValues() {
	fmt.Println("=== Range返回值用法 ===")

	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// 只获取key
	fmt.Println("只获取key:")
	for k := range m {
		fmt.Printf("  key: %s\n", k)
	}

	// 只获取value
	fmt.Println("\n只获取value:")
	for _, v := range m {
		fmt.Printf("  value: %d\n", v)
	}

	// 获取key和value
	fmt.Println("\n获取key和value:")
	for k, v := range m {
		fmt.Printf("  key: %s, value: %d\n", k, v)
	}

	// 切片同理
	s := []int{10, 20, 30}
	fmt.Println("\n只获取索引:")
	for i := range s {
		fmt.Printf("  索引: %d, 值: %d\n", i, s[i])
	}
}

// RangeModifyDuringIteration 演示遍历时修改集合
func RangeModifyDuringIteration() {
	fmt.Println("=== 遍历时修改集合 ===")

	// 切片：追加元素不会影响当前遍历
	fmt.Println("切片遍历时追加:")
	s := []int{1, 2, 3}
	for i, v := range s {
		s = append(s, v+10)
		fmt.Printf("索引: %d, 值: %d, 切片长度: %d\n", i, v, len(s))
	}
	fmt.Printf("最终切片: %v\n", s)

	// Map：修改已有key的值会影响后续读取
	fmt.Println("\nMap遍历时修改值:")
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range m {
		if k == "a" {
			m["b"] = 100 // 修改已有key
		}
		fmt.Printf("key: %s, value: %d\n", k, v)
	}
	fmt.Printf("最终map: %v\n", m)

	// Map：删除元素
	fmt.Println("\nMap遍历时删除:")
	m2 := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	for k := range m2 {
		if k == "a" {
			delete(m2, "b") // 删除其他key
		}
		fmt.Printf("key: %s\n", k)
	}
	fmt.Printf("最终map: %v\n", m2)
}

// RangePerformance 演示range性能考量
func RangePerformance() {
	fmt.Println("=== Range性能考量 ===")

	// 大切片遍历
	largeSlice := make([]int, 1000000)
	for i := range largeSlice {
		largeSlice[i] = i
	}

	// 方式1: range遍历
	start := time.Now()
	sum1 := 0
	for _, v := range largeSlice {
		sum1 += v
	}
	duration1 := time.Since(start)

	// 方式2: 索引遍历
	start = time.Now()
	sum2 := 0
	for i := 0; i < len(largeSlice); i++ {
		sum2 += largeSlice[i]
	}
	duration2 := time.Since(start)

	fmt.Printf("Range遍历: %v, sum=%d\n", duration1, sum1)
	fmt.Printf("索引遍历: %v, sum=%d\n", duration2, sum2)

	// 大结构体切片
	type BigStruct struct {
		data [100]int
	}
	bigSlice := make([]BigStruct, 10000)

	// 遍历大结构体 - 值拷贝开销
	start = time.Now()
	for _, v := range bigSlice {
		_ = v.data[0]
	}
	duration3 := time.Since(start)

	// 遍历大结构体 - 索引访问
	start = time.Now()
	for i := range bigSlice {
		_ = bigSlice[i].data[0]
	}
	duration4 := time.Since(start)

	fmt.Printf("大结构体Range遍历: %v\n", duration3)
	fmt.Printf("大结构体索引遍历: %v\n", duration4)
}
