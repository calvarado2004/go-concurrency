package main

import (
	"fmt"
	"sync"
)

func printSomething(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(s)
}

func main() {

	var wg sync.WaitGroup

	words := []string{"alpha", "beta", "delta", "gamma", "epsilon"}

	wg.Add(len(words))

	for i, word := range words {
		go printSomething(fmt.Sprintf("%d: %s", i, word), &wg)
	}

	wg.Wait()

	go printSomething("This is a string!", &wg)

	//awful solution to wait for the goroutine to finish
	//time.Sleep(1 * time.Second)

	wg.Add(1)

	printSomething("This is another string!", &wg)
}
