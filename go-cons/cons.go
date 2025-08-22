package main

import "fmt"

// const PI = 3.14
// const A int = 1
// const A = 1

const (
	A int = 1
	B     = 3.14
	C     = "Hi!"
)

func main() {
	// fmt.Println(PI)
	// const A = 1
	// A := 2 cannot assign to A
	fmt.Println(A)
	fmt.Println(B)
	fmt.Println(C)
}
