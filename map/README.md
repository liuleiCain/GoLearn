# Map 映射模块

## 学习目标
掌握Go语言map的初始化、使用方法和遍历特性。

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

## 最佳实践

1. **初始化**: 使用`make`初始化map
2. **检查存在**: 使用`value, ok := m[key]`
3. **并发安全**: 使用`sync.Map`或加锁
4. **有序遍历**: 先获取key并排序
5. **预估容量**: `make(map[K]V, hint)`减少扩容

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
