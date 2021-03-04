package main

import (
	"fmt"
)

func main() {
	a := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	i := 3

	a = append(a[:i], a[i+1:]...)
	fmt.Println(a)
	for _, b := range a {
		println(b)
	}
}
