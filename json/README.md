# JSON 序列化模块

## 学习目标
掌握Go语言JSON序列化与反序列化的基本用法和常见陷阱。

## 核心概念

### 1. 序列化 (Marshal)
将Go数据结构转换为JSON字节切片：
```go
data, err := json.Marshal(v)
```

### 2. 反序列化 (Unmarshal)
将JSON字节切片转换为Go数据结构：
```go
err := json.Unmarshal(data, &v)
```

### 3. 结构体标签
使用`json`标签控制字段映射：
```go
type User struct {
    Name string `json:"name"`      // 字段名映射
    Age  int    `json:"age,omitempty"` // 空值忽略
    Pass string `json:"-"`         // 忽略字段
}
```

## 示例说明

| 函数 | 说明 | 难度 |
|------|------|------|
| UnmarshalArr | JSON数组反序列化到切片 | 初级 |

## 常见用法

### 基本序列化
```go
user := UserInfo{Id: 1, Name: "Tom"}
data, _ := json.Marshal(user)
fmt.Println(string(data))  // {"Id":1,"Name":"Tom"}
```

### 基本反序列化
```go
var user UserInfo
json.Unmarshal([]byte(`{"Id":1,"Name":"Tom"}`), &user)
```

### 结构体标签
```go
type User struct {
    ID   int    `json:"id"`           // 映射为id
    Name string `json:"name"`         // 映射为name
    Age  int    `json:"age,omitempty"` // 为空时不输出
    Pass string `json:"-"`            // 不序列化
}
```

## 常见陷阱

### 1. 大小写问题
```go
type User struct {
    Name string  // 默认映射为 "Name" (首字母大写)
}

// 解决方案: 使用json标签
type User struct {
    Name string `json:"name"`  // 映射为 "name"
}
```

### 2. 未导出字段
```go
type User struct {
    name string  // 小写开头，无法序列化
    Name string  // 大写开头，可以序列化
}
```

### 3. 空值处理
```go
type User struct {
    Name string `json:"name,omitempty"`  // 空字符串时不输出
    Age  int    `json:"age"`             // 0也会输出
}
```

### 4. 反序列化到非指针
```go
var user UserInfo
json.Unmarshal(data, user)   // 错误: 必须传指针
json.Unmarshal(data, &user)  // 正确
```

## 最佳实践

1. **使用指针**: 反序列化时必须传递指针
2. **添加标签**: 为字段添加`json`标签控制映射
3. **处理错误**: 始终检查`Marshal`和`Unmarshal`的错误
4. **使用struct tag**: 控制字段名和空值行为
5. **数字类型**: JSON数字默认解析为float64

## 进阶用法

### 使用json.RawMessage
延迟解析JSON的某部分：
```go
type Request struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}
```

### 使用json.Decoder/Encoder
处理流式JSON：
```go
decoder := json.NewDecoder(reader)
encoder := json.NewEncoder(writer)
```
