// V3.go - Jantar dos Filósofos sem deadlock (ordenação de recursos)

package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	philosophers = 5
	rounds       = 20
	timeToThink  = 2
	timeToEat    = 3
)

var forks [philosophers]sync.Mutex

func Philosopher(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	leftFork := id
	rightFork := (id + 1) % philosophers

	// ordenacao
	// sempre pega primeiro o garfo de menor índice.
	first, second := leftFork, rightFork
	if first > second {
		first, second = second, first
	}

	for iteration := 1; iteration <= rounds; iteration++ {

		// pensa
		fmt.Printf("P%d: T%d\n", id+1, iteration)
		time.Sleep(time.Duration(timeToThink) * time.Second)

		// pega os garfos (em ordem crescente)
		fmt.Printf("P%d tentando pegar G%d\n", id+1, first+1)
		forks[first].Lock()
		fmt.Printf("P%d pegou G%d\n", id+1, first+1)

		fmt.Printf("P%d possui G%d e espera G%d\n", id+1, first+1, second+1)
		forks[second].Lock()

		// Eat
		fmt.Printf("P%d: E%d\n", id+1, iteration)
		time.Sleep(time.Duration(timeToEat) * time.Second)

		// Release forks
		forks[second].Unlock()
		forks[first].Unlock()

		fmt.Printf("P%d liberou G%d e G%d\n", id+1, first+1, second+1)
	}
}

func main() {
	var wg sync.WaitGroup

	fmt.Println("\n[P1] [P2] [P3] [P4] [P5]")
	fmt.Println("V3 - Dining Philosophers com ordenação de recursos")

	start := time.Now()

	wg.Add(philosophers)
	for id := 0; id < philosophers; id++ {
		go Philosopher(id, &wg)
	}
	wg.Wait()

	fmt.Printf("\nConcurrent Dinner took %s\n\n", time.Since(start))
}
