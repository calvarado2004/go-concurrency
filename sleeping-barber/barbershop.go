package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for clients.", barber)

		for {
			if len(shop.ClientsChan) == 0 {

				color.Yellow("%s is sleeping.", barber)
				isSleeping = true

			}

			client, shopOpen := <-shop.ClientsChan

			if shopOpen {
				if isSleeping {
					color.Yellow("%s is waking up.", barber)
					isSleeping = false
				}

				shop.cutHair(barber, client)

			} else {

				shop.sendBarberHome(barber)
				return

			}

		}
	}()
}

func (shop *BarberShop) cutHair(barber, client string) {
	color.Green("%s is cutting %s's hair.", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is done cutting %s's hair.", barber, client)

}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home.", barber)
	shop.BarbersDoneChan <- true
}

func (shop *BarberShop) closeShopForTheDay() {
	color.Cyan("Barbershop is closing for the day.")
	shop.Open = false

	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarbersDoneChan
	}

	close(shop.BarbersDoneChan)

	color.Green("Barbershop is now closed for the day.")
}

func (shop *BarberShop) addCustomer(customer string) {

	// print a message
	color.Green("%s has arrived at the barbershop.", customer)

	if shop.Open {
		select {
		case shop.ClientsChan <- customer:
			color.Blue("%s is waiting for a barber.", customer)
		default:
			color.Red("%s has left the barbershop because it is full.", customer)
		}
	} else {
		color.Red("%s has left the barbershop because it is closed.", customer)
		return
	}
}
