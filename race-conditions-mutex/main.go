package main

import (
	"fmt"
	"sync"
)

var msg string

var wg sync.WaitGroup

func updateMessage(s string) {
	defer wg.Done()
	msg = s
}

func main() {

	msg = "Hello, World!"

	wg.Add(2)
	go updateMessage("Hello, Go!")
	go updateMessage("Hello, Gophers!")
	wg.Wait()

	fmt.Println(msg)
}
