package main

import "fmt"

type Speaker interface {
	Speak() string
}

type Person struct {
	Name string
}

func (p Person) Speak() string {
	return "Xin chào, tôi là " + p.Name
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "Gâu gâu, giống chó " + d.Name
}

func saySomeThing(s Speaker) {
	fmt.Println(s.Speak())
}

func printAnyThing(x interface{}) {
	fmt.Println(x)
}

func main() {
	var s Speaker

	s = Person{"Hưng"}
	fmt.Println(s.Speak())

	s = Dog{"Border Corli"}
	fmt.Println(s.Speak())

	saySomeThing(Person{"Hoàng Phương Thảo"})
	saySomeThing(Dog{"Golden"})

	printAnyThing(123)
	printAnyThing("Hoàng Phương Thảo")

	var i interface{} = "Xin chào"
	str, ok := i.(string)
	if ok {
		fmt.Println("Chuỗi:", str)
	}

}
