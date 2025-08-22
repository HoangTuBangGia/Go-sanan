package main

import "fmt"

func main() {
	// var i, j string = "Hello", "World!"
	// var i, j = 10, 20
	// var i string = "Hello"
	// var j int = 15

	// fmt.Print(i, "\n")
	// fmt.Print(j, "\n")

	// fmt.Print(i, "\n", j)

	// fmt.Print(i, j, "\n")
	// fmt.Println(i, j)
	// fmt.Printf("i has value: %v and type: %T\n", i, i)
	// fmt.Printf("j has value: %v and type: %T\n", j, j)

	var i = 15.5
	var txt = "Hello World!"

	fmt.Printf("%v\n", i)
	fmt.Printf("%#v\n", i)
	fmt.Printf("%v%%\n", i)
	fmt.Printf("%T\n", i)

	fmt.Printf("%v\n", txt)
	fmt.Printf("%#v\n", txt)
	fmt.Printf("%T\n", txt)
}
