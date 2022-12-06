package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

type Icome struct {
	Source string
	Amount int
}

func main() {

	// variable for bank balance
	var bankBalance int

	var balance sync.Mutex

	// print out statarting value
	fmt.Printf("Starting bank balance: $%d.00", bankBalance)
	fmt.Println()

	// define weekly revenue

	incomes := []Icome{
		{Source: "Paycheck", Amount: 1000},
		{Source: "Videos", Amount: 3000},
		{Source: "Sponsorship", Amount: 5000},
		{Source: "Tips", Amount: 1000},
		{Source: "Invesments", Amount: 100},
	}

	wg.Add(len(incomes))

	// loop trough 52 weeks and print out how much is made, keep a running total
	for i, income := range incomes {
		go func(i int, income Icome) {
			defer wg.Done()
			for week := 1; week <= 52; week++ {

				balance.Lock()
				temp := bankBalance
				temp += income.Amount
				bankBalance = temp
				balance.Unlock()

				fmt.Printf("On week %d, you earned $%d.00 from %s\n", week, income.Amount, income.Source)

			}

		}(i, income)

	}

	wg.Wait()
	// print out final balance
	fmt.Printf("Final bank balance: $%d.00\n", bankBalance)
}
