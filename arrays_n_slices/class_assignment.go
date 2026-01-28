package main

import "fmt"

func main() {

	var arr [6]int = [6]int{1, 2, 3, 4, 4, 5}

	fmt.Println("Array initialized with 6 elements: ", arr)
	for i, element1 := range arr {
		for j, element2 := range arr {
			if element1 == element2 && i != j {
				fmt.Println("Found duplicate:", element1, " at indices ", i, " and ", j)
			}
		}
	}
}
