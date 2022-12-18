package main

import (
	"fmt"
	"strings"
)

// Channels allows to communicate between goroutines

func shout(ping <-chan string, pong chan<- string) {

	for {
		s := <-ping

		pong <- fmt.Sprintf("%s!!!", strings.ToUpper(s))
	}
}

func main() {

	//create two channels
	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)

	fmt.Println("Type something and press enter (enter q to quit):")

	for {

		fmt.Print("->")
		var userInput string
		_, _ = fmt.Scanln(&userInput)

		if userInput == strings.ToLower("q") {
			break
		}

		ping <- userInput
		response := <-pong
		fmt.Println("Response: ", response)
	}

	fmt.Println("Bye! closing channels...")

	close(ping)
	close(pong)

}
