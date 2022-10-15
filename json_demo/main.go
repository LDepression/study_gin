package main

import (
	"encoding/json"
	"fmt"
	"math"
)

type Test struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}

func main() {
	t1 := &Test{
		ID:   math.MaxInt64,
		Name: "yyy",
	}
	//序列化:将go语言类型转换为json类型
	b, err := json.Marshal(t1)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(string(b))
	s := `{"id":"9223372036854775807","name":"yyy"}`
	var t2 Test
	err = json.Unmarshal([]byte(s), &t2)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(t2)
}
