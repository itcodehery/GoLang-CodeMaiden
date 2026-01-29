package main

import (
	"fmt"
	"sort"
	"strings"
)

func isExpressCargo(s int) bool {
	if s%2 == 0 && s%10 != 0 {
		return true
	}
	return false
}

func isSecure(s string) bool {
	chars := strings.Split(s, "")
	if chars[0] == chars[len(s)-1] {
		return true
	}
	return false
}

func displayPriorityReport(expCargo []int, secConts map[string]int, bayMap map[string][]int, hvyCargo []int) {
	fmt.Println("---------------")
	fmt.Println("Priority Report")
	fmt.Println("---------------")
	fmt.Printf("Express Cargo:\n%v\n", expCargo)
	fmt.Printf("Secure Containers:\n%v\n", secConts)
	fmt.Printf("Containers in Bays:\n%v\n", bayMap)
	fmt.Printf("Heavy Cargo:\n%v\n", hvyCargo[:2])
	fmt.Println("---------------")
}

func main() {
	var cargoIds [9]int = [9]int{50, 42, 100, 12, 18, 55, 60, 24, 5}
	var contTags [6]string = [6]string{"level", "gamma", "area", "trust", "radar", "sense"}

	// Express Cargo slice
	var expressCargo []int = []int{}

	// Secured Containers map I think
	secContainers := make(map[string]int)

	// Bay Maps (sounds like baymax lol)
	bayMap := make(map[string][]int)

	for _, ele := range cargoIds {
		if isExpressCargo(ele) {
			expressCargo = append(expressCargo, ele)
		}
	}

	for _, element := range contTags {
		if isSecure(element) {
			secContainers[element] = len(element)
		}
	}

	for _, elem := range expressCargo {
		if elem > 20 {
			var a = bayMap["Bay-A"]
			a = append(a, elem)
			bayMap["Bay-A"] = a
			continue
		}
		var b = bayMap["Bay-B"]
		b = append(b, elem)
		bayMap["Bay-B"] = b
	}

	// Heavy Cargo slice
	var hvyCargo []int = []int{}

	for _, car_element := range cargoIds {
		if car_element > 40 {
			hvyCargo = append(hvyCargo, car_element)
		}
	}

	sort.Slice(hvyCargo, func(i, j int) bool {
		if i > j {
			return true
		}
		return false
	})

	displayPriorityReport(expressCargo, secContainers, bayMap, hvyCargo)
}
