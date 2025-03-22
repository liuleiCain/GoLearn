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
			delete(m, "c")
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
