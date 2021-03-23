package main

import "strings"

func main() {
	var iistr []string
	iistr = strings.Split("1 11 12 13 21 23 29 31 32 40 44 49 55 58 59 63 64 75 77 86 88 92 94 98 99", " ")

	println(len(iistr))
}
