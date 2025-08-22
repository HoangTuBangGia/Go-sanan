package main

import "fmt"

func main() {
	s := "Hello"
	fmt.Println(len(s))

	fmt.Println()

	s = "Xin chào"
	fmt.Println(len(s))
	fmt.Println(len([]rune(s)))

	for i, r := range s {
		fmt.Printf("%d: %c\n", i, r)
	}

	fmt.Println()

	fmt.Println(s[0])
	fmt.Println(string(s[0]))

	fmt.Println()
	runes := []rune(s)
	fmt.Println(runes)
	fmt.Println(string(runes))
}
