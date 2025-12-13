package main

import "fmt"

func Digitize(n int) []int {
	var list []int
	for n > 0 {
		n = n % 10
		list = append(list, n)
	}
	for i := 0; i < len(list); i++ {
		if list[i] < list[i-1] {
			list = append(list, list[i])
		}
	}
	return list
}

func main() {
	fmt.Println(Digitize(32541))
}
