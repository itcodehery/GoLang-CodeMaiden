package main

import (
	"fmt"
	"time"
)

func main() {
	counter := 0

	for i := 0 ; i < 10 ; i++ {
		go func() {
			for j := 0 ; j < 10 ; j++ {
				counter++
			}
		}()
	}

	time.Sleep(time.Second)
	fmt.Printf("%v",counter)
}
