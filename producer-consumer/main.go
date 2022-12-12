package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {

	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func makePizza(pizzaNumber int) *PizzaOrder {

	// create a pizza order
	pizzaNumber++

	rnd := rand.Intn(12) + 1
	msg := ""
	success := false

	// check if we have made enough pizzas
	if pizzaNumber > NumberOfPizzas {
		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     "We have made enough pizzas! \n",
			success:     success,
		}
	}

	// make the pizza
	delay := rand.Intn(5) + 1
	fmt.Printf("Received order #%d\n", pizzaNumber)

	pizzasMade++

	// check if the pizza was burned or undercooked on two different cases
	switch {

	case rnd <= 2:
		msg = ("Pizza was undercooked! \n")

	case rnd > 2 && rnd < 5:
		msg = ("Pizza was burned!")

	}

	total++

	// check if the pizza was burned or undercooked
	if rnd < 5 {
		pizzasFailed++
		pizzasMade--
		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
	}

	fmt.Printf("Making pizza #%d, will take %d seconds... \n", pizzaNumber, delay)
	time.Sleep(time.Duration(delay) * time.Second)

	// return the pizza orders
	success = true
	msg = fmt.Sprintf("Pizza order #%d is ready!, pizzas successfully made: %d \n", pizzaNumber, pizzasMade)

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
		message:     msg,
		success:     success,
	}

}

func pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we are making
	var i = 0

	// run until we recevie a message on the quit channel

	// try to make pizzas

	for {
		// try to make a pizza
		currentPizza := makePizza(i)

		// decision
		if currentPizza == nil {
			break
		}

		// do not nest the select statement
		i = currentPizza.pizzaNumber
		select {

		case pizzaMaker.data <- *currentPizza:

		case quitChan := <-pizzaMaker.quit:
			// close the channel
			close(pizzaMaker.data)
			close(quitChan)
			return
		}

	}

}

func main() {

	//seed the random number generator
	rand.Seed(time.Now().UnixNano())

	//print out a message
	color.Cyan("The Pizzeria is open for business! Let's make some pizzas! \n")
	color.Cyan("---------------------------------------------------------- \n")

	//create a producer

	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	//run the producer in a goroutine
	go pizzeria(pizzaJob)

	//create a consumer
	for i := range pizzaJob.data {

		if !i.success {
			color.Red(i.message)
			color.Red("The customer of order: #%d will be really mad! \n", i.pizzaNumber)
		} else {
			color.Green(i.message)
			color.Green("Pizza order: #%d is out for delivery!\n", i.pizzaNumber)
		}

		if i.pizzaNumber > NumberOfPizzas {
			// print out the message
			color.Cyan("Done making %d pizzas \n", pizzasMade)

			err := pizzaJob.Close()
			if err != nil {
				color.Red("Error closing the pizzeria: %v \n", err)
			}

		}
	}

	// print out the ending message
	color.Cyan("The pizzeria is closed! \n")

	color.Cyan("We made %d pizzas, but failed to make %d, with %d attempts in total \n", pizzasMade, pizzasFailed, total)

	switch {
	case pizzasFailed > 9:
		color.Red("This pizzeria is a failure! \n")
	case pizzasFailed < 5 && pizzasFailed > 2:
		color.Yellow("This pizzeria have some room to improve! \n")
	default:
		color.Green("This pizzeria is a success! \n")
	}

}
