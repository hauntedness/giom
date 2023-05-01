package main

import "fmt"

type Num[T any] int

func main() {
	var n1 Num[string] = 1
	var n2 Num[string] = 2
	fmt.Println(n1 + n2)
}
