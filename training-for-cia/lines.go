package main

import (
	"fmt"
)

func lines() {
	var students [5]int = [5]int{10, 11, 12, 13, 14}
	var marks []int = []int{87, 86, 84, 89, 92}

	studentmap := make(map[int]int)

	for index, val := range students {
		studentmap[val] = marks[index]
	}

	fmt.Println(studentmap)

	var sum int
	for _, val := range marks {
		sum = sum + val
	}

	avg := sum / len(marks)

	fmt.Println("Average of all marks: ", avg)

	var max int = marks[0]
	for _, i := range marks {
		if i > max {
			max = i
		}
	}

	fmt.Println("Max of all marks: ", max)

	var markSlice []int
	for _, i := range marks {
		if i > avg {
			markSlice = append(markSlice, i)
		}
	}

	fmt.Println("Mark Slice: ", markSlice)
}
