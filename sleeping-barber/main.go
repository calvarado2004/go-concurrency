package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// variables

var seatingCapacity = 10
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second

func main() {
	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	// print welcome message
	color.Yellow("The Sleeping Barber Problem")
	color.Yellow("----------------------------")

	// create channels
	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// create barbershop
	barbershop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		BarbersDoneChan: doneChan,
		ClientsChan:     clientChan,
		Open:            true,
	}

	color.Green("Barbershop is open!")

	// add barbers
	barbershop.addBarber("Frank")
	barbershop.addBarber("Bob")
	barbershop.addBarber("Joe")
	barbershop.addBarber("Sally")
	barbershop.addBarber("Sue")

	// start the barbershop as a go routine

	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		barbershop.closeShopForTheDay()
		closed <- true

	}()

	// add customers
	i := 1
	go func() {
		for {
			// get a random number with average arrival rate
			randomMilliseconds := rand.Int() % (2 * arrivalRate)

			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMilliseconds)):
				barbershop.addCustomer(fmt.Sprintf("Customer %d", i))
				i++
			}
		}
	}()

	// block until barbershop is closed
	<-closed
	// print goodbye message
	color.Yellow("Barbershop is closed!")
}
