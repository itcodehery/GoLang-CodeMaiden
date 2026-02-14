package main

import "fmt"

func test2(y func(int) int) int {
	return y(10)
}

func returnFunc(x string) func() {
	return func() {
		fmt.Println(x)
	}
}

func main() {
	test := func(x int) int {
		return x * 2
	}

	var func_from_api func(int) int = fetchFuncFromAPI()
	func_from_api()

	fmt.Println(test2(test))
	returnFunc("Hello")()
}
