package main

import "fmt"

func arbiter(kwargs ...int) {
	sum := 0
	for _, val := range kwargs {
		sum = sum + val
	}
	fmt.Println(kwargs, sum)
}

func main() {
	arbiter(1, 2)
	arbiter(1, 2, 3, 4, 5, 6, 7)
}
