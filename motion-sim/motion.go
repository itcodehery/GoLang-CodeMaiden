package main

import (
	"fmt"
)

func motion_sim() {
	// Variables
	var init_pos float64
	var init_vel float64
	var acceleration float64
	var mass float64
	var num_steps int

	// User input
	fmt.Println("Enter initial position: ")
	fmt.Scanln(&init_pos)
	fmt.Println("Enter initial velocity: ")
	fmt.Scanln(&init_vel)
	fmt.Println("Enter acceleration: ")
	fmt.Scanln(&acceleration)
	fmt.Println("Enter mass: ")
	fmt.Scanln(&mass)
	fmt.Println("Enter number of time steps: ")
	fmt.Scanln(&num_steps)

	// Computation
	ke := 0.0
	for t := 0; t <= num_steps; t++ {
		init_vel = init_vel + acceleration
		init_pos = init_pos + init_vel
		ke = 0.5 * mass * init_vel * init_vel
		fmt.Println("t = ", t+1, " | x = ", init_pos, " | v = ", init_vel, " | ke = ", ke)
	}

	fmt.Println()

	// Result
	if init_vel < 0 {
		fmt.Println("Object reversed direction")
	}

	if ke < 0.1 {
		fmt.Println("Motion nearly stopped")
	}

	if init_pos > 20.0 {
		fmt.Println("Object left simulation boundary")
	}
}

func main() {
	fmt.Println("Motion Simulator")
	fmt.Println("----------------")
	motion_sim()
}
