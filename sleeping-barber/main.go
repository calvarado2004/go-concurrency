package main

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// variables

var seatingCapacity = 10
var arrivalRate = 100 * time.Millisecond
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

	// start the barbershop

	// add customers

	// block until barbershop is closed
	time.Sleep(5 * time.Second)

	// print goodbye message
}
