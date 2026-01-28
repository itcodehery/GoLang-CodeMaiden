package main

import (
	"fmt"
)

func main() {
	// Creating a Slice
	// slc := make([]int, 5)
	// for i := range slc {
	// 	slc[i] = i * 2
	// }
	// fmt.Println("Slice contents: ", slc)

	// Creating a Map
	// mp := make(map[string]int)
	// mp["one"] = 1
	// mp["one"] = 11
	// mp["two"] = 2
	// mp["three"] = 3
	// mp["four"] = 4
	// fmt.Println(mp)

	// Creating a Struct
	type Person struct {
		Name  string
		Age   int
		Email string
	}

	p1 := Person{"Praneeth", 21, "praneeth.m@mca.christuniversity.in"}

	fmt.Printf("%+v", p1)
}
