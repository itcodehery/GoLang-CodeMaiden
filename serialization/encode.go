package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   uint8  `json:"age"`
}

func main() {

	user := User{Name: "Arden", Email: "arden@diago.com", Age: 30}
	data, err := json.Marshal(user)

	fmt.Println(user)
	fmt.Println(string(data))
	fmt.Println(err)

}
