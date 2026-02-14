package main

import (
	"errors"
	"fmt"
)

// Fix: return not used
func add(a int, b int) int {
	sum := a + b
	return sum
}

// Fix: wrong return type
func calculate(a, b int) int {
	return a + b
}

// TODO: Handle Error
// Fix: handled error
func divide(a, b int) (int, error) {
	return a / b
}

func rectangle(l, w float64) (area float64) {
	area := l * w
	return
}

func sumAll(nums ...int) int {
	total := 0
	for i := 0; i < len(nums); i++ {

		total += i
	}
	return total
}

func increment(n int) {
	n = n + 1
}

func counter() func() int {
	count := 0
	return func() {
		count++
	}
}

func factorial(n int) int {
	return n * factorial(n-1)
}
func apply(f int) int {
	return f(5)
}

// func main() {
// 	fmt.Println("Add:", add(5, 3))
// 	sum, diff := calculate(10, 5)
// 	fmt.Println("Sum:", sum, "Diff:", diff)
// 	result, err := divide(10, 0)
// 	fmt.Println(result, err)
// 	fmt.Println("Rectangle:", rectangle(5, 4))
// 	fmt.Println("SumAll:", sumAll(1, 2, 3, 4))
// 	x := 10
// 	increment(x)
// 	fmt.Println("Increment:", x)
//
// 	c := counter()
// 	fmt.Println(c())
// 	fmt.Println(c())
// 	fmt.Println("Factorial:", factorial(5))
// 	var multiply func(int, int) int
// 	fmt.Println(multiply(2, 3))
// 	fmt.Println(apply(add))
// }
