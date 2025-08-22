package main

import "fmt"

func main() {
	// ch := make(chan int)

	// go func() {
	// 	ch <- 42
	// }()

	// value := <-ch

	// fmt.Println(value)

	// ch := make(chan string)

	// go func() {
	// 	ch <- "ping"
	// }()

	// msg := <-ch
	// fmt.Println(msg)

	// ch := make(chan int, 2)

	// ch <- 10
	// ch <- 20

	// close(ch)

	// // fmt.Println(<-ch)
	// // fmt.Println(<-ch)

	// for v := range ch {
	// 	fmt.Println(v)
	// }

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() { ch1 <- "from ch1" }()
	go func() { ch2 <- "from ch2" }()

	select {
	case msg1 := <-ch1:
		fmt.Println("Received:", msg1)
	case msg2 := <-ch2:
		fmt.Println("Received:", msg2)
	}

	ch := make(chan string)

	go func() { ch <- "A" }()
	go func() { ch <- "B" }()

	fmt.Println(<-ch)
	fmt.Println(<-ch)

}
