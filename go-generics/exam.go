package main

import (
	"cmp"
	"fmt"
)

type Number interface {
	int | int64 | float64
}

type Stack[T any] struct {
	items []T
}

func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}

	return b
}

func Sum[T Number](arr []T) T {
	var total T

	for _, v := range arr {
		total += v
	}

	return total
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() T {
	n := len(s.items)
	v := s.items[n-1]
	s.items = s.items[:n-1]

	return v
}

func main() {
	fmt.Println(Min(3, 5))
	fmt.Println(Min(1.999, 2.001))
	fmt.Println(Min("a", "b"))

	fmt.Println(Sum([]int{1, 2, 3, 4}))
	fmt.Println(Sum([]float64{3.4, 6.4}))

	intStack := Stack[int]{}
	intStack.Push(10)
	intStack.Push(20)
	fmt.Println(intStack.Pop())

	stringStack := Stack[string]{}
	stringStack.Push("Hưng")
	stringStack.Push("Thảo")
	fmt.Println(stringStack.Pop())
}
