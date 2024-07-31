package main

import (
	"encoding/json"
	"fmt"
)

type MyData struct {
	ID       int64  `json:"id,string"`
	Username string `json:"username"`
}

func main() {

	//结构体转成json
	/*
		data := MyData{
			ID:       math.MaxInt64,
			Username: "cisco001",
		}
		jsonStr, err := json.Marshal(data)
		if err != nil {
			return
		}
		fmt.Println(string(jsonStr))

	*/

	//json转成结构体
	jsonStr := `{"id":"9223372036854775807","username":"cisco001"}`
	var d2 MyData
	json.Unmarshal([]byte(jsonStr), &d2)
	fmt.Printf("%#v\n", d2)

}
