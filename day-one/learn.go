package main

import "fmt" 
import "math/rand"

func add(a , b int) int {
	return a + b
}

func main() {
	var res int = rand.Intn(100)

	fmt.Printf("%v is the result right now\n",res)
	res = add(2,3)

	fmt.Printf("%v is the result after the addition\n",res)
}
