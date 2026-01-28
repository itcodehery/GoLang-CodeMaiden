package main

import "fmt"

func define() {
	type Empl struct {
		Name   string
		Age    int
		Remote bool
	}

	var job Empl = Empl{Name: "Praneeth", Age: 21, Remote: true}
	fmt.Println("%t", job)
}
