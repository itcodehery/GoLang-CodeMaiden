// First Golang code

package main

import "fmt"

func add_hello(name string) string {
	return "Hello " + name;
}

func main() {
	var kartik string = "Kartik"
	var a int = 10

	fmt.Println(add_hello(kartik))
	fmt.Println(a)
	fmt.Printf("%v %T",a,a)

	var (
		distro_name string = "Gentoo"
		distro_rating float32 = 4.9
		distro_cost int = 0
	)

	fmt.Printf("%v\n",distro_name)
	fmt.Printf("%v\n",distro_rating)
	fmt.Printf("%v\n",distro_cost)

	var b bool = true 
	fmt.Printf("%v\n",b)

	// Arrays
	// var marks[3] int
	// marks[0] = 10
	// marks[1] = 20
	// marks[2] = 30
	// or
	marks := [3]int {10,20,30}

	fmt.Println(marks)
	// Maps	
	ram_prices := map[string]float32 {
		"Corsair" : 7899.99,
		"Samsung" : 10999.99,
		"Crucial" : 8999.99,
	}

	fmt.Println(ram_prices)
}
