package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// UserInfo 用户信息结构体
type UserInfo struct {
	Id   int
	Name string
}

// UnmarshalArr JSON数组反序列化到切片
func UnmarshalArr() {
	list := make([]UserInfo, 0)
	list = append(list, UserInfo{Id: 1, Name: "tom1"})
	list = append(list, UserInfo{Id: 2, Name: "tom2"})
	list = append(list, UserInfo{Id: 3, Name: "tom3"})
	_ = json.Unmarshal([]byte(`[{"id":1,"name":"tom"},{"id":2,"name":"jerry"}]`), &list)
}

// BasicMarshal 基本序列化示例
func BasicMarshal() {
	fmt.Println("=== 基本序列化示例 ===")

	// 序列化结构体
	user := UserInfo{Id: 1, Name: "Tom"}
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("序列化错误: %v\n", err)
		return
	}
	fmt.Printf("序列化结果: %s\n", data)

	// 序列化切片
	users := []UserInfo{{1, "Tom"}, {2, "Jerry"}}
	data, _ = json.Marshal(users)
	fmt.Printf("序列化切片: %s\n", data)

	// 序列化map
	m := map[string]int{"a": 1, "b": 2}
	data, _ = json.Marshal(m)
	fmt.Printf("序列化map: %s\n", data)
}

// BasicUnmarshal 基本反序列化示例
func BasicUnmarshal() {
	fmt.Println("=== 基本反序列化示例 ===")

	// 反序列化到结构体
	jsonStr := `{"id":1,"name":"Tom"}`
	var user UserInfo
	err := json.Unmarshal([]byte(jsonStr), &user)
	if err != nil {
		fmt.Printf("反序列化错误: %v\n", err)
		return
	}
	fmt.Printf("反序列化结果: %+v\n", user)

	// 反序列化到切片
	jsonArr := `[{"id":1,"name":"Tom"},{"id":2,"name":"Jerry"}]`
	var users []UserInfo
	json.Unmarshal([]byte(jsonArr), &users)
	fmt.Printf("反序列化切片: %+v\n", users)

	// 反序列化到map
	jsonMap := `{"a":1,"b":2}`
	var m map[string]int
	json.Unmarshal([]byte(jsonMap), &m)
	fmt.Printf("反序列化map: %+v\n", m)
}

// StructTags 演示结构体标签
func StructTags() {
	fmt.Println("=== 结构体标签示例 ===")

	type User struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Age      int    `json:"age,omitempty"`
		Password string `json:"-"`
		Email    string `json:"email,omitempty"`
	}

	// omitempty: 空值不输出
	user := User{ID: 1, Name: "Tom", Password: "secret"}
	data, _ := json.Marshal(user)
	fmt.Printf("omitempty效果: %s\n", data)

	// -: 字段不序列化
	user2 := User{ID: 2, Name: "Jerry", Password: "hidden", Email: "jerry@example.com"}
	data, _ = json.Marshal(user2)
	fmt.Printf("忽略字段效果: %s\n", data)
}

// MarshalIndent 格式化输出
func MarshalIndent() {
	fmt.Println("=== 格式化输出示例 ===")

	user := UserInfo{Id: 1, Name: "Tom"}

	// 紧凑输出
	compact, _ := json.Marshal(user)
	fmt.Printf("紧凑输出: %s\n", compact)

	// 格式化输出
	indented, _ := json.MarshalIndent(user, "", "  ")
	fmt.Printf("格式化输出:\n%s\n", indented)
}

// Request 请求结构体，用于演示RawMessage
type Request struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// LoginData 登录数据
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterData 注册数据
type RegisterData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// RawMessageDemo 演示json.RawMessage延迟解析
func RawMessageDemo() {
	fmt.Println("=== RawMessage延迟解析示例 ===")

	// 登录请求
	loginJSON := `{"type":"login","data":{"username":"tom","password":"123456"}}`
	handleRequest(loginJSON)

	// 注册请求
	registerJSON := `{"type":"register","data":{"email":"tom@example.com","username":"tom","password":"123456"}}`
	handleRequest(registerJSON)
}

func handleRequest(jsonStr string) {
	var req Request
	json.Unmarshal([]byte(jsonStr), &req)

	switch req.Type {
	case "login":
		var login LoginData
		json.Unmarshal(req.Data, &login)
		fmt.Printf("登录请求: 用户名=%s\n", login.Username)
	case "register":
		var reg RegisterData
		json.Unmarshal(req.Data, &reg)
		fmt.Printf("注册请求: 邮箱=%s, 用户名=%s\n", reg.Email, reg.Username)
	}
}

// PartialParse 部分解析JSON
func PartialParse() {
	fmt.Println("=== 部分解析JSON示例 ===")

	type PartialRequest struct {
		ID   int             `json:"id"`
		Meta json.RawMessage `json:"meta"`
	}

	jsonStr := `{"id":123,"meta":{"version":"1.0","timestamp":1234567890,"extra":"data"}}`

	var req PartialRequest
	json.Unmarshal([]byte(jsonStr), &req)
	fmt.Printf("ID: %d\n", req.ID)
	fmt.Printf("Meta原始: %s\n", req.Meta)

	// 后续需要时再解析meta
	var meta map[string]interface{}
	json.Unmarshal(req.Meta, &meta)
	fmt.Printf("Meta解析后: %+v\n", meta)
}

// StreamDecode 演示流式解码
func StreamDecode() {
	fmt.Println("=== 流式解码示例 ===")

	// 模拟JSON流
	jsonStream := `{"name":"Alice","age":30}
{"name":"Bob","age":25}
{"name":"Charlie","age":35}`

	decoder := json.NewDecoder(strings.NewReader(jsonStream))

	for {
		var user map[string]interface{}
		if err := decoder.Decode(&user); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("解码错误: %v\n", err)
			break
		}
		fmt.Printf("解析到用户: %v\n", user)
	}
}

// StreamEncode 演示流式编码
func StreamEncode() {
	fmt.Println("=== 流式编码示例 ===")

	users := []UserInfo{
		{1, "Alice"},
		{2, "Bob"},
		{3, "Charlie"},
	}

	var buf strings.Builder
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")

	for _, user := range users {
		encoder.Encode(user)
	}

	fmt.Printf("编码结果:\n%s", buf.String())
}

// Person 用于演示自定义序列化
type Person struct {
	Name string
	Age  int
}

// MarshalJSON 自定义序列化
func (p Person) MarshalJSON() ([]byte, error) {
	type Alias Person
	return json.Marshal(&struct {
		DisplayName string `json:"display_name"`
		Age         int    `json:"age"`
	}{
		DisplayName: p.Name,
		Age:         p.Age,
	})
}

// UnmarshalJSON 自定义反序列化
func (p *Person) UnmarshalJSON(data []byte) error {
	type Alias Person
	aux := &struct {
		DisplayName string `json:"display_name"`
		Age         int    `json:"age"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.Name = aux.DisplayName
	p.Age = aux.Age
	return nil
}

// CustomMarshal 演示自定义序列化
func CustomMarshal() {
	fmt.Println("=== 自定义序列化示例 ===")

	p := Person{Name: "张三", Age: 25}

	// 使用自定义序列化
	data, _ := json.Marshal(p)
	fmt.Printf("自定义序列化: %s\n", data)

	// 使用自定义反序列化
	jsonStr := `{"display_name":"李四","age":30}`
	var p2 Person
	json.Unmarshal([]byte(jsonStr), &p2)
	fmt.Printf("自定义反序列化: %+v\n", p2)
}

// UnknownJSON 演示处理未知结构的JSON
func UnknownJSON() {
	fmt.Println("=== 处理未知结构JSON示例 ===")

	// 使用 map[string]interface{}
	jsonStr := `{"name":"Tom","age":30,"active":true,"score":95.5}`
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &result)

	fmt.Println("使用 map[string]interface{}:")
	if name, ok := result["name"].(string); ok {
		fmt.Printf("  name: %s\n", name)
	}
	if age, ok := result["age"].(float64); ok {
		fmt.Printf("  age: %.0f\n", age)
	}
	if active, ok := result["active"].(bool); ok {
		fmt.Printf("  active: %v\n", active)
	}
}

// JSONTypeCheck 演示JSON类型判断
func JSONTypeCheck() {
	fmt.Println("=== JSON类型判断示例 ===")

	jsonData := []byte(`{"name":"Tom","age":30,"scores":[90,85,95],"meta":null}`)

	var data interface{}
	json.Unmarshal(jsonData, &data)

	checkType(data, "root")
}

func checkType(v interface{}, name string) {
	switch val := v.(type) {
	case map[string]interface{}:
		fmt.Printf("%s 是对象，包含 %d 个字段\n", name, len(val))
		for k, v := range val {
			checkType(v, k)
		}
	case []interface{}:
		fmt.Printf("%s 是数组，包含 %d 个元素\n", name, len(val))
		for i, v := range val {
			checkType(v, fmt.Sprintf("%s[%d]", name, i))
		}
	case string:
		fmt.Printf("%s 是字符串: %s\n", name, val)
	case float64:
		fmt.Printf("%s 是数字: %.0f\n", name, val)
	case bool:
		fmt.Printf("%s 是布尔值: %v\n", name, val)
	case nil:
		fmt.Printf("%s 是null\n", name)
	}
}

// JSONNumber 演示json.Number精确处理数字
func JSONNumber() {
	fmt.Println("=== json.Number精确处理数字示例 ===")

	// 大数字可能精度丢失
	jsonStr := `{"id":12345678901234567890,"name":"test"}`

	// 普通解析（精度丢失）
	var normal map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &normal)
	fmt.Printf("普通解析: id = %v (类型: %T)\n", normal["id"], normal["id"])

	// 使用json.Number
	decoder := json.NewDecoder(strings.NewReader(jsonStr))
	decoder.UseNumber()

	var withNumber map[string]interface{}
	decoder.Decode(&withNumber)

	if num, ok := withNumber["id"].(json.Number); ok {
		fmt.Printf("json.Number: id = %s\n", num.String())
		n, _ := num.Int64()
		fmt.Printf("转为int64: %d\n", n)
	}
}

// NestedJSON 演示嵌套JSON处理
func NestedJSON() {
	fmt.Println("=== 嵌套JSON处理示例 ===")

	type Address struct {
		City    string `json:"city"`
		Country string `json:"country"`
	}

	type User struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Address Address  `json:"address"`
		Tags    []string `json:"tags"`
	}

	jsonStr := `{
		"name": "Tom",
		"age": 30,
		"address": {
			"city": "Beijing",
			"country": "China"
		},
		"tags": ["developer", "golang"]
	}`

	var user User
	json.Unmarshal([]byte(jsonStr), &user)
	fmt.Printf("用户: %+v\n", user)
	fmt.Printf("地址: %+v\n", user.Address)
	fmt.Printf("标签: %v\n", user.Tags)
}

// JSONValidation 演示JSON验证
func JSONValidation() {
	fmt.Println("=== JSON验证示例 ===")

	// 验证是否为有效JSON
	validJSON := `{"name":"Tom","age":30}`
	invalidJSON := `{"name":"Tom","age":30,}` // 多余逗号

	fmt.Printf("有效JSON: %v\n", json.Valid([]byte(validJSON)))
	fmt.Printf("无效JSON: %v\n", json.Valid([]byte(invalidJSON)))
}

// CompactJSON 演示JSON压缩
func CompactJSON() {
	fmt.Println("=== JSON压缩示例 ===")

	indented := `{
  "name": "Tom",
  "age": 30
}`

	var buf bytes.Buffer
	err := json.Compact(&buf, []byte(indented))
	if err != nil {
		fmt.Printf("压缩错误: %v\n", err)
		return
	}
	fmt.Printf("压缩后: %s\n", buf.String())
}
