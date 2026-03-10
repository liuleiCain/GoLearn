package learn_map

import (
	"fmt"
	"sort"
	"sync"
)

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

// MapOperations 演示map的基本操作
func MapOperations() {
	fmt.Println("=== Map基本操作演示 ===")

	// 创建map
	m := make(map[string]int)

	// 插入
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3
	fmt.Printf("插入后: %v\n", m)

	// 获取
	v := m["a"]
	fmt.Printf("获取m[\"a\"]: %d\n", v)

	// 获取不存在的key
	v = m["notexist"]
	fmt.Printf("获取不存在的key: %d (返回零值)\n", v)

	// 检查key是否存在
	value, ok := m["a"]
	fmt.Printf("m[\"a\"]存在: value=%d, ok=%v\n", value, ok)

	value, ok = m["notexist"]
	fmt.Printf("m[\"notexist\"]存在: value=%d, ok=%v\n", value, ok)

	// 修改
	m["a"] = 100
	fmt.Printf("修改后: %v\n", m)

	// 删除
	delete(m, "b")
	fmt.Printf("删除b后: %v\n", m)

	// 删除不存在的key - 不会报错
	delete(m, "notexist")
	fmt.Printf("删除不存在的key后: %v\n", m)

	// 获取长度
	fmt.Printf("map长度: %d\n", len(m))
}

// MapLiteral 演示map字面量初始化
func MapLiteral() {
	fmt.Println("=== Map字面量初始化 ===")

	// 字面量初始化
	m1 := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	fmt.Printf("字面量初始化: %v\n", m1)

	// 空map
	m2 := map[string]int{}
	fmt.Printf("空map: %v, len=%d\n", m2, len(m2))

	// nil map
	var m3 map[string]int
	fmt.Printf("nil map: %v, len=%d\n", m3, len(m3))

	// 使用make指定初始容量
	m4 := make(map[string]int, 100)
	fmt.Printf("指定容量的map: len=%d\n", len(m4))
}

// MapKeyType 演示map的key类型限制
func MapKeyType() {
	fmt.Println("=== Map Key类型限制 ===")

	// 可比较类型可以作为key
	m := make(map[interface{}]int)

	// int作为key
	m[1] = 100
	m[2] = 200

	// string作为key
	m["hello"] = 300

	// bool作为key
	m[true] = 400

	// float作为key
	m[3.14] = 500

	// 数组作为key
	m[[3]int{1, 2, 3}] = 600

	// 指针作为key
	x := 10
	m[&x] = 700

	// 结构体作为key（所有字段都可比较）
	type Point struct{ X, Y int }
	m[Point{1, 2}] = 800

	fmt.Printf("各种类型作为key: %v\n", m)

	// 不能作为key的类型：
	// - 切片
	// - map
	// - 函数
	// 这些类型不可比较
}

// MapNested 演示嵌套map
func MapNested() {
	fmt.Println("=== 嵌套Map演示 ===")

	// map的value是map
	nestedMap := make(map[string]map[string]int)

	// 需要初始化内层map
	nestedMap["group1"] = make(map[string]int)
	nestedMap["group1"]["a"] = 1
	nestedMap["group1"]["b"] = 2

	nestedMap["group2"] = make(map[string]int)
	nestedMap["group2"]["c"] = 3

	fmt.Printf("嵌套map: %v\n", nestedMap)

	// map的value是切片
	mapSlice := make(map[string][]int)
	mapSlice["nums"] = []int{1, 2, 3}
	mapSlice["evens"] = []int{2, 4, 6}

	fmt.Printf("value是切片: %v\n", mapSlice)

	// 切片的元素是map
	sliceMap := []map[string]int{
		{"a": 1, "b": 2},
		{"c": 3, "d": 4},
	}
	fmt.Printf("切片元素是map: %v\n", sliceMap)
}

// MapSortedIteration 演示有序遍历map
func MapSortedIteration() {
	fmt.Println("=== 有序遍历Map ===")

	m := map[string]int{
		"banana": 2,
		"apple":  1,
		"cherry": 3,
		"date":   4,
	}

	// 按key排序遍历
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("按key排序遍历:")
	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, m[k])
	}

	// 按value排序遍历
	type kv struct {
		Key   string
		Value int
	}
	var pairs []kv
	for k, v := range m {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value < pairs[j].Value
	})

	fmt.Println("\n按value排序遍历:")
	for _, p := range pairs {
		fmt.Printf("  %s: %d\n", p.Key, p.Value)
	}
}

// MapConcurrent 演示map并发问题与解决方案
func MapConcurrent() {
	fmt.Println("=== Map并发安全 ===")

	// 普通map并发写会panic
	fmt.Println("普通map并发写（演示错误，已注释）:")
	// m := make(map[int]int)
	// var wg sync.WaitGroup
	// for i := 0; i < 10; i++ {
	//     wg.Add(1)
	//     go func(n int) {
	//         defer wg.Done()
	//         m[n] = n // 并发写会panic
	//     }(i)
	// }
	// wg.Wait()

	// 解决方案1: 使用互斥锁
	fmt.Println("\n方案1: 使用sync.Mutex")
	m1 := make(map[int]int)
	var mu sync.Mutex
	var wg1 sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg1.Add(1)
		go func(n int) {
			defer wg1.Done()
			mu.Lock()
			m1[n] = n * 10
			mu.Unlock()
		}(i)
	}
	wg1.Wait()
	fmt.Printf("使用Mutex: %v\n", m1)

	// 解决方案2: 使用读写锁
	fmt.Println("\n方案2: 使用sync.RWMutex")
	m2 := make(map[int]int)
	var rwmu sync.RWMutex
	var wg2 sync.WaitGroup

	// 写操作
	for i := 0; i < 5; i++ {
		wg2.Add(1)
		go func(n int) {
			defer wg2.Done()
			rwmu.Lock()
			m2[n] = n * 10
			rwmu.Unlock()
		}(i)
	}
	wg2.Wait()

	// 读操作
	rwmu.RLock()
	fmt.Printf("使用RWMutex: %v\n", m2)
	rwmu.RUnlock()

	// 解决方案3: 使用sync.Map
	fmt.Println("\n方案3: 使用sync.Map")
	var syncMap sync.Map
	var wg3 sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg3.Add(1)
		go func(n int) {
			defer wg3.Done()
			syncMap.Store(n, n*10)
		}(i)
	}
	wg3.Wait()

	fmt.Print("使用sync.Map: ")
	syncMap.Range(func(key, value interface{}) bool {
		fmt.Printf("%v:%v ", key, value)
		return true
	})
	fmt.Println()
}

// SyncMapDemo 详细演示sync.Map
func SyncMapDemo() {
	fmt.Println("=== sync.Map详细演示 ===")

	var m sync.Map

	// Store: 存储键值对
	m.Store("name", "张三")
	m.Store("age", 25)
	m.Store("active", true)

	// Load: 读取值
	if v, ok := m.Load("name"); ok {
		fmt.Printf("Load name: %v\n", v)
	}

	// LoadOrStore: 读取或存储
	actual, loaded := m.LoadOrStore("name", "李四")
	fmt.Printf("LoadOrStore name: actual=%v, loaded=%v\n", actual, loaded)

	actual, loaded = m.LoadOrStore("city", "北京")
	fmt.Printf("LoadOrStore city: actual=%v, loaded=%v\n", actual, loaded)

	// Delete: 删除
	m.Delete("age")
	if _, ok := m.Load("age"); !ok {
		fmt.Println("Delete age: 已删除")
	}

	// Range: 遍历
	fmt.Println("Range遍历:")
	m.Range(func(key, value interface{}) bool {
		fmt.Printf("  %v: %v\n", key, value)
		return true
	})

	// LoadAndDelete: 读取并删除
	v, loaded := m.LoadAndDelete("active")
	fmt.Printf("LoadAndDelete active: value=%v, loaded=%v\n", v, loaded)
}

// MapCopy 演示map的拷贝
func MapCopy() {
	fmt.Println("=== Map拷贝演示 ===")

	// 浅拷贝
	original := map[string]int{"a": 1, "b": 2}
	shallowCopy := make(map[string]int)
	for k, v := range original {
		shallowCopy[k] = v
	}
	shallowCopy["a"] = 100
	fmt.Printf("原始map: %v\n", original)
	fmt.Printf("浅拷贝: %v\n", shallowCopy)

	// 引用拷贝（指向同一个map）
	referenceCopy := original
	referenceCopy["b"] = 200
	fmt.Printf("引用拷贝修改后，原始map: %v\n", original)

	// 深拷贝（value是指针或引用类型时需要）
	type Person struct {
		Name string
		Age  int
	}
	originalPtr := map[string]*Person{
		"p1": {Name: "张三", Age: 20},
	}
	deepCopy := make(map[string]*Person)
	for k, v := range originalPtr {
		deepCopy[k] = &Person{Name: v.Name, Age: v.Age}
	}
	deepCopy["p1"].Age = 30
	fmt.Printf("深拷贝后，原始: %v, 深拷贝: %v\n", originalPtr["p1"].Age, deepCopy["p1"].Age)
}

// MapSet 使用map实现集合(Set)
func MapSet() {
	fmt.Println("=== 使用Map实现Set ===")

	// 使用map[T]struct{}实现集合（struct{}不占内存）
	set := make(map[string]struct{})

	// 添加元素
	set["apple"] = struct{}{}
	set["banana"] = struct{}{}
	set["cherry"] = struct{}{}

	// 检查元素是否存在
	if _, exists := set["apple"]; exists {
		fmt.Println("apple存在于集合中")
	}

	if _, exists := set["grape"]; !exists {
		fmt.Println("grape不存在于集合中")
	}

	// 删除元素
	delete(set, "banana")

	// 集合大小
	fmt.Printf("集合大小: %d\n", len(set))

	// 遍历集合
	fmt.Print("集合元素: ")
	for item := range set {
		fmt.Printf("%s ", item)
	}
	fmt.Println()
}

// MapCounter 使用map实现计数器
func MapCounter() {
	fmt.Println("=== 使用Map实现计数器 ===")

	// 统计单词出现次数
	counter := make(map[string]int)
	for _, word := range []string{"hello", "world", "hello", "go", "world", "go", "go"} {
		counter[word]++
	}

	fmt.Println("单词计数:")
	for word, count := range counter {
		fmt.Printf("  %s: %d\n", word, count)
	}

	// 使用sync.Map实现并发安全计数器
	fmt.Println("\n并发安全计数器:")
	var syncCounter sync.Map

	var wg sync.WaitGroup
	words := []string{"a", "b", "a", "c", "b", "a", "d"}

	for _, word := range words {
		wg.Add(1)
		go func(w string) {
			defer wg.Done()
			for {
				if v, ok := syncCounter.Load(w); ok {
					if syncCounter.CompareAndSwap(w, v, v.(int)+1) {
						break
					}
				} else {
					if _, loaded := syncCounter.LoadOrStore(w, 1); !loaded {
						break
					}
				}
			}
		}(word)
	}
	wg.Wait()

	syncCounter.Range(func(key, value interface{}) bool {
		fmt.Printf("  %v: %v\n", key, value)
		return true
	})
}

// MapClear 演示清空map的方法
func MapClear() {
	fmt.Println("=== 清空Map的方法 ===")

	// 方法1: 遍历删除
	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Printf("原始map: %v\n", m1)
	for k := range m1 {
		delete(m1, k)
	}
	fmt.Printf("遍历删除后: %v\n", m1)

	// 方法2: 重新赋值（原map会被GC回收）
	m2 := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Printf("\n原始map: %v\n", m2)
	m2 = make(map[string]int)
	fmt.Printf("重新赋值后: %v\n", m2)

	// 方法3: Go 1.21+ 使用clear函数
	m3 := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Printf("\n原始map: %v\n", m3)
	// clear(m3) // Go 1.21+
	// fmt.Printf("clear函数后: %v\n", m3)
	fmt.Println("Go 1.21+ 可以使用 clear(m3)")
}
