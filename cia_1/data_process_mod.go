package main

import (
	"fmt"
	"sort"
	"strings"
)

func categorize_year(slice []int) (z1 []int, z2 []int) {
	for _, val := range slice {
		if val%2 == 0 {
			z1 = append(z1, val)
		} else {
			z2 = append(z2, val)
		}
	}
	return z1, z2
}

func identifier_len(str_slice []string) map[int][]string {
	var len_group map[int][]string = make(map[int][]string)
	for _, val := range str_slice {

		x := len(val)
		var curr_slice []string

		for _, val := range str_slice {
			if len(val) == x {
				curr_slice = append(curr_slice, val)
			}
		}

		len_group[x] = curr_slice
	}
	return len_group
}

func year_digit(rec_years []int) map[int]int {
	var len_group map[int]int = make(map[int]int)
	for _, val := range rec_years {

		last_digit_of_year := val % 10

		var curr_count int = 0

		for _, val := range rec_years {
			if val%10 == last_digit_of_year {
				curr_count++
			}
		}

		len_group[int(last_digit_of_year)] = curr_count
	}
	return len_group
}

func recent_year_select(rec_years []int) []int {
	rec_copy := rec_years
	sort.Slice(rec_copy, func(i, j int) bool {
		if i > j {
			return true
		}
		return false
	})
	return rec_copy[:3]
}

func contains(char_array []string, item string) bool {
	for _, val := range char_array {
		if val == item {
			return true
		}
	}
	return false
}

func alpha_filter(asset_ids []string) (ret []string) {

	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for _, val := range asset_ids {
		var str string = val
		// if len gt than 4
		if len(val) > 4 {
			// filter by alphabetical only
			x := strings.Split(val, "")
			// For each char
			for _, num := range numbers {
				if !contains(x, num) {
					str = val
					break
				}
			}
			ret = append(ret, str)
		}
	}
	return
}

func main() {
	RecordYears := []int{1900, 1947, 2000, 2012, 2023, 2024, 2100, 1985, 1975}

	AssetIDs := []string{"naman", "12321", "sweets", "malayalam", "mela", "101", "isro"}

	fmt.Println("\nYear Categorization")
	fmt.Println(categorize_year(RecordYears))
	fmt.Println("\nIdentifier Length Grouping")
	fmt.Println(identifier_len(AssetIDs))
	fmt.Println("\nYear Ending Frequency")
	fmt.Println(year_digit(RecordYears))
	fmt.Println("\nRecent Year Select ")
	fmt.Println(recent_year_select(RecordYears))
	fmt.Println("\nRecent Year Selection")
	fmt.Println(alpha_filter(AssetIDs))
}
