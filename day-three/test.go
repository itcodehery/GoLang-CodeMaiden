package main

import "fmt"

func main() {
	a := 1
	b := 2.2
	// ./test.go:8:12: invalid operation: a + b (mismatched types int and float64)
	fmt.Printf("Brother the answer is: %v\n", a+int(b))
	v := "string"
	Define()

	fmt.Printf("The length of the string is: %v", len(v))
}
