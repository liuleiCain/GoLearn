# Map 映射模块

## 学习目标
深入掌握Go语言map的初始化、使用方法、遍历特性、并发安全和高级应用。

## 核心概念

### 1. map基础
map是Go语言的内置哈希表实现，存储键值对。

### 2. map声明与初始化
```go
// 声明但未初始化 (nil map)
var m map[string]int

// 使用make初始化
m := make(map[string]int)

// 字面量初始化
m := map[string]int{"a": 1, "b": 2}
```

### 3. map操作
- 插入: `m[key] = value`
- 获取: `value := m[key]`
- 删除: `delete(m, key)`
- 检查存在: `value, ok := m[key]`

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| InitMap | map初始化与nil map陷阱 | 初级 |
| RangeMap | map遍历顺序随机性 | 初级 |
| MapOperations | map基本操作 | 初级 |
| MapLiteral | map字面量初始化 | 初级 |
| MapKeyType | map key类型限制 | 中级 |
| MapNested | 嵌套map | 中级 |
| MapSortedIteration | 有序遍历map | 中级 |
| MapConcurrent | map并发安全 | 高级 |
| SyncMapDemo | sync.Map详细演示 | 高级 |
| MapCopy | map拷贝 | 中级 |
| MapSet | 使用map实现Set | 中级 |
| MapCounter | 使用map实现计数器 | 中级 |
| MapClear | 清空map的方法 | 初级 |

## nil map vs 空map

### 关键区别
```go
// nil map - 不能写入！
var m map[string]int
m["key"] = 1  // panic: assignment to entry in nil map

// 空map - 可以正常使用
m := make(map[string]int)
m["key"] = 1  // 正常
```

### 对比表
| 操作 | nil map | 空map |
|------|---------|-------|
| 读取 | 返回零值 | 返回零值 |
| 写入 | panic | 正常 |
| 删除 | 无操作 | 正常 |
| len() | 0 | 0 |
| range | 不执行 | 不执行 |

## 新增示例详解

### MapOperations - 基本操作
```go
m := make(map[string]int)

// 插入
m["a"] = 1

// 获取
v := m["a"]

// 检查存在
value, ok := m["a"]

// 删除
delete(m, "a")

// 长度
len(m)
```

### MapKeyType - Key类型限制
可以作为key的类型：
- 基本类型：int, string, bool, float
- 数组
- 指针
- 结构体（所有字段可比较）

不能作为key的类型：
- 切片
- map
- 函数

### MapNested - 嵌套Map
```go
// map的value是map
nestedMap := make(map[string]map[string]int)
nestedMap["group1"] = make(map[string]int)
nestedMap["group1"]["a"] = 1

// map的value是切片
mapSlice := make(map[string][]int)
mapSlice["nums"] = []int{1, 2, 3}
```

### MapSortedIteration - 有序遍历
```go
m := map[string]int{"banana": 2, "apple": 1, "cherry": 3}

// 按key排序
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, m[k])
}
```

### MapConcurrent - 并发安全
```go
// 方案1: sync.Mutex
var mu sync.Mutex
m := make(map[int]int)
mu.Lock()
m[1] = 10
mu.Unlock()

// 方案2: sync.RWMutex（读多写少）
var rwmu sync.RWMutex
rwmu.RLock()  // 读锁
v := m[1]
rwmu.RUnlock()

// 方案3: sync.Map
var syncMap sync.Map
syncMap.Store("key", "value")
v, ok := syncMap.Load("key")
```

### SyncMapDemo - sync.Map详解
```go
var m sync.Map

// 存储
m.Store("name", "张三")

// 读取
v, ok := m.Load("name")

// 读取或存储
actual, loaded := m.LoadOrStore("name", "李四")

// 删除
m.Delete("name")

// 遍历
m.Range(func(key, value interface{}) bool {
    fmt.Println(key, value)
    return true
})

// 读取并删除
v, loaded := m.LoadAndDelete("name")
```

### MapCopy - Map拷贝
```go
// 浅拷贝
original := map[string]int{"a": 1, "b": 2}
shallowCopy := make(map[string]int)
for k, v := range original {
    shallowCopy[k] = v
}

// 引用拷贝（指向同一个map）
referenceCopy := original

// 深拷贝（value是指针时）
deepCopy := make(map[string]*Person)
for k, v := range originalPtr {
    deepCopy[k] = &Person{Name: v.Name, Age: v.Age}
}
```

### MapSet - 实现集合
```go
// 使用map[T]struct{}实现集合（struct{}不占内存）
set := make(map[string]struct{})

// 添加
set["apple"] = struct{}{}

// 检查存在
_, exists := set["apple"]

// 删除
delete(set, "apple")
```

### MapCounter - 实现计数器
```go
counter := make(map[string]int)
for _, word := range words {
    counter[word]++
}
```

### MapClear - 清空Map
```go
// 方法1: 遍历删除
for k := range m {
    delete(m, k)
}

// 方法2: 重新赋值
m = make(map[string]int)

// 方法3: Go 1.21+ clear函数
clear(m)
```

## map遍历顺序

### 随机性设计
Go语言故意将map遍历顺序设计为随机的，原因：
1. 防止程序依赖特定顺序
2. 避免哈希碰撞攻击
3. 不保证插入顺序

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}

// 每次遍历顺序可能不同
for k, v := range m {
    fmt.Println(k, v)
}
```

### 有序遍历
```go
// 如需有序遍历，先排序key
m := map[string]int{"c": 3, "a": 1, "b": 2}

keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)

for _, k := range keys {
    fmt.Println(k, m[k])
}
```

## 常见陷阱

### 1. nil map写入
```go
var m map[string]int
m["key"] = 1  // panic!

// 解决方案
m := make(map[string]int)
```

### 2. 并发读写
```go
// 错误: map不是并发安全的
go func() {
    m["key"] = 1  // 并发写
}()
m["key"] = 2  // 并发写

// 解决方案1: 使用sync.RWMutex
var mu sync.RWMutex
mu.Lock()
m["key"] = 1
mu.Unlock()

// 解决方案2: 使用sync.Map
var m sync.Map
m.Store("key", 1)
v, ok := m.Load("key")
```

### 3. 遍历时修改
```go
m := map[string]int{"a": 1, "b": 2}
for k := range m {
    m["c"] = 3  // 可能导致不确定行为
}
```

### 4. 嵌套map未初始化
```go
// 错误
m := make(map[string]map[string]int)
m["a"]["b"] = 1  // panic: nil map

// 正确
m := make(map[string]map[string]int)
m["a"] = make(map[string]int)
m["a"]["b"] = 1
```

## 最佳实践

1. **初始化**: 使用`make`初始化map
2. **检查存在**: 使用`value, ok := m[key]`
3. **并发安全**: 使用`sync.Map`或加锁
4. **有序遍历**: 先获取key并排序
5. **预估容量**: `make(map[K]V, hint)`减少扩容
6. **集合实现**: 使用`map[T]struct{}`节省内存
7. **嵌套map**: 记得初始化内层map

## sync.Map

适用于以下场景：
- 读多写少
- key相对稳定
- 需要并发安全

```go
var m sync.Map

// 存储
m.Store("key", "value")

// 读取
v, ok := m.Load("key")

// 删除
m.Delete("key")

// 遍历
m.Range(func(k, v interface{}) bool {
    fmt.Println(k, v)
    return true
})
```

## 运行测试

```bash
go test -v ./map/...
```

## 参考资料

- [Go语言Map详解](https://go.dev/blog/maps)
- [Go并发编程之sync.Map](https://pkg.go.dev/sync#Map)
