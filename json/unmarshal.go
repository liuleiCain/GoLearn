package json

import "encoding/json"

func UnmarshalArr() {
	list := make([]UserInfo, 0)
	list = append(list, UserInfo{Id: 1, Name: "tom1"})
	list = append(list, UserInfo{Id: 2, Name: "tom2"})
	list = append(list, UserInfo{Id: 3, Name: "tom3"})
	_ = json.Unmarshal([]byte(`[{"id":1,"name":"tom"},{"id":2,"name":"jerry"}]`), &list)
}

type UserInfo struct {
	Id   int
	Name string
}
