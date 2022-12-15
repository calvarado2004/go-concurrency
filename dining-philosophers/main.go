package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// The Dining Philosophers problem is well known in computer science circles.
// Five philosophers, numbered from 0 through 4, live in a house where the
// table is laid for them; each philosopher has their own place at the table.
// Their only difficulty – besides those of philosophy – is that the dish
// served is a very difficult kind of spaghetti which has to be eaten with
// two forks. There are two forks next to each plate, so that presents no
// difficulty. As a consequence, however, this means that no two neighbours
// may be eating simultaneously, since there are five philosophers and five forks.
//
// This is a simple implementation of Dijkstra's solution to the "Dining
// Philosophers" dilemma.

// Philosopher is a struct that represents a philosopher.
type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// List of philosophers.
var philosophers = []Philosopher{
	{name: "Aristotle", rightFork: 4, leftFork: 0},
	{name: "Kant", rightFork: 0, leftFork: 1},
	{name: "Marx", rightFork: 1, leftFork: 2},
	{name: "Russell", rightFork: 2, leftFork: 3},
	{name: "Sartre", rightFork: 3, leftFork: 4},
}

// how many times each philosopher will eat.
var eatCount = 3
var eatTime = 1 * time.Second
var thinkTime = 3 * time.Second
var sleepTime = 2 * time.Second

var orderMutex = &sync.Mutex{}
var orderFinished []string

func main() {

	fmt.Println("Dining Philosophers")
	fmt.Println("-------------------")
	fmt.Println("The table is empty.")

	time.Sleep(sleepTime)

	// start the meal
	dine()

	// finish the meal
	fmt.Println("The table is empty.")

	time.Sleep(sleepTime)
	fmt.Printf("Order finished: %s.\n", strings.Join(orderFinished, ", "))

}

func dine() {

	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// forks is a map of all forks.
	forks := make(map[int]*sync.Mutex)

	// create a mutex for each fork.
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// start the meal

	for i := 0; i < len(philosophers); i++ {
		go diningProblem(philosophers[i], wg, forks, seated)
	}

	wg.Wait()

}

func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	defer wg.Done()

	// seat the philosopher at the table.
	fmt.Printf("%s sits down at the table.\n", philosopher.name)
	seated.Done()

	// eat three times.
	for i := eatCount; i > 0; i-- {

		// get a lock on both forks.

		if philosopher.leftFork > philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s picks up the fork on their right.\n", philosopher.name)
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s picks up the fork on their left.\n", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s picks up the fork on their left.\n", philosopher.name)
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s picks up the fork on their right.\n", philosopher.name)
		}

		fmt.Printf("\t%s is eating.\n", philosopher.name)
		time.Sleep(eatTime)

		fmt.Printf("\t%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)

		// release the forks.
		forks[philosopher.rightFork].Unlock()
		fmt.Printf("\t%s puts down the fork on their right.\n", philosopher.name)
		forks[philosopher.leftFork].Unlock()
		fmt.Printf("\t%s puts down the fork on their left.\n", philosopher.name)

	}

	fmt.Println(philosopher.name, "is done eating.")
	fmt.Println(philosopher.name, "left the table.")

	orderMutex.Lock()
	orderFinished = append(orderFinished, philosopher.name)
	orderMutex.Unlock()

}
