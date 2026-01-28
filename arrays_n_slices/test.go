package main

import "fmt"

func arraysnslices() {
	// in Go, all variables
	// are initialized while declared with zero values

	// Array
	var nums [5]int = [5]int{1, 2, 3, 4, 5}

	// Slice (Own individual entity)
	// Created using an array
	s := nums[1:4]
	fmt.Println(s)
	// Len tells us the length of the slice
	fmt.Println(len(s))
	// Cap tells us the total capacity of the slice
	fmt.Println(cap(s))

	// Reassignment
	s = nums[1:]
	fmt.Println(s)
	fmt.Println(len(s))
	fmt.Println(cap(s))

	// Creating a slice without an array
	var sumnums []int = []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println(sumnums)
	fmt.Println(len(sumnums))
	fmt.Println(cap(sumnums))
	sumnums.append(8)
}
