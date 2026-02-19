package main

import "fmt"

type Hit interface {
	Utility() int
}

type Supes struct {
	name  string
	power int
	saves int
}

type Villains struct {
	name  string
	power int
	kills int
}

// A Struct Method!
func (s Supes) Utility() int {
	return s.power * s.saves
}

func (v Villains) Utility() int {
	return v.power * (1 / v.kills)
}

// Cannot define methods for non-local types
// func (int) belowHundred() int {
// 	return 10
// }

// Pass by reference instead of copying
// Sets the Power field of the Supe to an int value
func (s *Supes) setPower(power int) {
	s.power = power
}

func main() {
	Batman := Supes{"Batman", 3, 500}
	Superman := Supes{"Superman", 9, 800}

	Joker := Villains{"Joker", 2, 1000}
	Darkseid := Villains{"Darkseid", 10, 300}

	JusticeLeague := []Supes{Batman, Superman}
	InjusticeLeague := []Villains{Joker, Darkseid}

	HitSquad := []Hit{Batman, Superman, Joker, Darkseid}

	fmt.Println(Batman.Utility())
	fmt.Println(Superman.Utility())
	Batman.setPower(1)
	fmt.Println(Batman.Utility())
	fmt.Println(JusticeLeague)

	fmt.Println(Joker.Utility())
	fmt.Println(InjusticeLeague)
	fmt.Println(HitSquad)
}
