package main

import "fmt"

// Warm up Activity
func operate(x ...int) int {
	sum := 0
	for _, ele := range x {
		sum += ele
	}
	return sum
}

// Activity Task 1:
func getOperation(op string) func(...int) float64 {
	// switch op {
	// case "add":
	// 	return func(x ...int) float64 {
	// 		sum := 0.0
	// 		for _, j := range x {
	// 			sum += float64(j)
	// 		}
	// 		return sum
	// 	}
	// case "mult":
	// 	return func(x ...int) float64 {
	// 		product := 0.0
	// 		for _, j := range x {
	// 			product *= float64(j)
	// 		}
	// 		return product
	// 	}
	// case "avg":
	// 	return func(x ...int) float64 {
	// 		return getOperation("add")(x...) / float64(len(x))
	// 	}
	// default:
	// 	return func(...int) float64 { return 0 }
	// }

	operations := make(map[string]func(...int) float64)
	operations["add"] = func(x ...int) float64 {
		sum := 0.0
		for _, j := range x {
			sum += float64(j)
		}
		return sum
	}

	operations["mult"] = func(x ...int) float64 {
		product := 0.0
		for _, j := range x {
			product *= float64(j)
		}
		return product
	}

	operations["avg"] = func(x ...int) float64 {
		return getOperation("add")(x...) / float64(len(x))
	}

	return operations[op]

}

// Activity Task 2
type Processor interface {
	Process(nums ...int) float64
}

type SumProcessor struct{}
type AvgProcessor struct{}
type MaxProcessor struct{}

func (s SumProcessor) Process(nums ...int) float64 {
	sum := 0.0
	for _, j := range nums {
		sum += float64(j)
	}
	return sum
}

func (s AvgProcessor) Process(nums ...int) float64 {
	sum := 0.0
	for _, j := range nums {
		sum += float64(j)
	}
	return sum / float64(len(nums))
}

func (s MaxProcessor) Process(nums ...int) float64 {
	max := nums[0]
	for i, _ := range nums {
		if nums[i] > max {
			max = nums[i]
		}
	}
	return float64(max)
}

func getProcessor(op string) Processor {
	switch op {
	case "add":
		return SumProcessor{}
	case "max":
		return MaxProcessor{}
	case "avg":
		return AvgProcessor{}
	default:
		return SumProcessor{}
	}
}

// Final Challenge
func runPipeline(pipeline []Processor, nums ...int) {
	for _, j := range pipeline {
		fmt.Println(j.Process(nums...))
	}
}

func main() {
	// result := getProcessor("avg")
	// fmt.Println(result.Process(2, 3, 4, 5))
	// sum := getProcessor("add")
	// fmt.Println(sum.Process(2, 3, 4, 5))

	pipeline := []Processor{
		SumProcessor{},
		AvgProcessor{},
		MaxProcessor{},
	}

	runPipeline(pipeline, 1, 2, 3)
}
