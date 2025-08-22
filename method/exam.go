package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

// Method gắn vào type Person
func (p Person) Greet() {
	fmt.Printf("Xin chào, tôi là %s, %d tuổi.\n", p.Name, p.Age)
}

func (p *Person) GrowUp() {
	p.Age += 1
	fmt.Println(p.Age)
}

func main() {
	p := Person{"Hưng", 21}

	p.Greet() // Gọi Method
	p.GrowUp()

	fmt.Println(p.Age)
}
