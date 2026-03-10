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

### 1. 使用json.RawMessage延迟解析

`json.RawMessage` 是一个字节切片，可以延迟JSON的解析，适用于以下场景：

#### 场景一：根据类型动态解析
当JSON数据的某个字段类型不确定时，可以先解析类型字段，再根据类型解析数据：

```go
type Request struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"` // 延迟解析
}

type LoginData struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type RegisterData struct {
    Email    string `json:"email"`
    Username string `json:"username"`
    Password string `json:"password"`
}

func handleRequest(jsonStr string) {
    var req Request
    json.Unmarshal([]byte(jsonStr), &req)
    
    switch req.Type {
    case "login":
        var login LoginData
        json.Unmarshal(req.Data, &login)
        fmt.Printf("登录: %s\n", login.Username)
    case "register":
        var reg RegisterData
        json.Unmarshal(req.Data, &reg)
        fmt.Printf("注册: %s\n", reg.Email)
    }
}
```

#### 场景二：部分解析大JSON
只解析需要的字段，其他字段保持原样：

```go
type PartialRequest struct {
    ID   int             `json:"id"`
    Meta json.RawMessage `json:"meta"` // 不解析meta字段
}

// 后续需要时再解析
var meta map[string]interface{}
json.Unmarshal(req.Meta, &meta)
```

### 2. 使用json.Decoder/Encoder处理流式JSON

适用于处理来自网络连接、文件等数据流的JSON数据。

#### 从Reader解码JSON
```go
func decodeFromReader(r io.Reader) {
    decoder := json.NewDecoder(r)
    
    // 单个对象
    var user User
    if err := decoder.Decode(&user); err != nil {
        log.Fatal(err)
    }
    
    // 多个连续JSON对象
    for {
        var u User
        if err := decoder.Decode(&u); err == io.EOF {
            break
        } else if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("用户: %s\n", u.Name)
    }
}
```

#### 向Writer编码JSON
```go
func encodeToWriter(w io.Writer, user User) {
    encoder := json.NewEncoder(w)
    
    // 设置缩进格式
    encoder.SetIndent("", "  ")
    
    if err := encoder.Encode(user); err != nil {
        log.Fatal(err)
    }
}
```

#### 处理HTTP请求/响应
```go
func handleHTTP(w http.ResponseWriter, r *http.Request) {
    // 解码请求
    var req Request
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // 编码响应
    w.Header().Set("Content-Type", "application/json")
    encoder := json.NewEncoder(w)
    encoder.Encode(map[string]string{"status": "ok"})
}
```

### 3. 自定义序列化/反序列化

实现 `json.Marshaler` 和 `json.Unmarshaler` 接口：

```go
type Person struct {
    Name string
    Age  int
}

// 自定义序列化
func (p Person) MarshalJSON() ([]byte, error) {
    type Alias Person
    return json.Marshal(&struct {
        DisplayName string `json:"display_name"`
        *Alias
    }{
        DisplayName: p.Name,
        Alias:       (*Alias)(&p),
    })
}

// 自定义反序列化
func (p *Person) UnmarshalJSON(data []byte) error {
    type Alias Person
    aux := &struct {
        DisplayName string `json:"display_name"`
        *Alias
    }{
        Alias: (*Alias)(p),
    }
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    p.Name = aux.DisplayName
    return nil
}
```

### 4. 处理未知结构的JSON

使用 `map[string]interface{}` 或 `interface{}`：

```go
// 使用 map[string]interface{}
var result map[string]interface{}
json.Unmarshal(jsonData, &result)

// 类型断言访问
if name, ok := result["name"].(string); ok {
    fmt.Println(name)
}

// 使用 interface{} + 类型判断
var data interface{}
json.Unmarshal(jsonData, &data)

switch v := data.(type) {
case map[string]interface{}:
    fmt.Println("对象:", v)
case []interface{}:
    fmt.Println("数组:", v)
case string:
    fmt.Println("字符串:", v)
case float64:
    fmt.Println("数字:", v)
}
```

### 5. 使用json.Number精确处理数字

避免大数字精度丢失：

```go
decoder := json.NewDecoder(r)
decoder.UseNumber() // 启用json.Number

var result map[string]interface{}
decoder.Decode(&result)

// 获取原始数字字符串
if num, ok := result["id"].(json.Number); ok {
    fmt.Println("原始数字:", num.String())
    
    // 转换为int64
    n, _ := num.Int64()
    fmt.Println("int64:", n)
}
```

### 6. 处理JSON流中的多个对象

从包含多个JSON对象的流中读取：

```go
func decodeMultipleJSON(r io.Reader) {
    decoder := json.NewDecoder(r)
    
    for decoder.More() { // 检查是否还有更多数据
        var obj map[string]interface{}
        if err := decoder.Decode(&obj); err != nil {
            if err == io.EOF {
                break
            }
            log.Fatal(err)
        }
        fmt.Printf("解析到: %v\n", obj)
    }
}
```
