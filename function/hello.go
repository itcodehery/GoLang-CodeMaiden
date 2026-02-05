package main

import "fmt"

func test(a int, b int) (r1 int, r2 int) {
	r1 = a + b
	r2 = a - b
	return
}

func foo(x int) int {
	defer fmt.Println("Checking order")
	defer fmt.Println("Checking order ")
	defer fmt.Println("Checking order 1st time")
	fmt.Println("Checking order again")
	return x * x
}

// Anonymous Functions
func main() {
	// func(n int) {
	// 	fmt.Println("Hello World!!!", n)
	// }(23)
	x := func(n int) {
		fmt.Println("Hello World!!!")
	}
	fmt.Printf("Type of x: %T\n", x)

	foo(2)
}
