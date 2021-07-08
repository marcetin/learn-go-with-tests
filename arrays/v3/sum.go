package main

// Sum израчунава укупан износ из кришки бројева.
func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
