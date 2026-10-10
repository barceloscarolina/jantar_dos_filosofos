// V2.go - Jantar dos Filósofos -  sem deadlock - sincronização excessiva
// FPPD - Trabalho 1
//
// V2: 1 filósofo por vez pode pegar 2 garfos, evitando deadlock.
// Mantém alguma possibilidade de concorrência, mas com sincronização excessiva.
// 
//


package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	philosophers = 5
	rounds        = 20
	timeToThink   = 2
	timeToEat     = 3
)

// Cada Mutex representa um garfo.
// Um garfo pode ser utilizado por apenas um filosofo de cada vez.
var forks [philosophers]sync.Mutex
//1 filósofo por vez pode pegar 2 garfos. Escolher 2 garfos é um mutex
var pickForksMutex sync.Mutex

func Philosopher(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	//O filósofo utiliza os garfos ao lado dele
	leftFork := id
	rightFork := (id + 1) % philosophers

	for iteration := 1; iteration <= rounds; iteration++ {

		// Pensando
		fmt.Printf("P%d: T%d\n", id+1, iteration)
		time.Sleep(time.Duration(timeToThink) * time.Second)

		//Adquirindo garfos
		fmt.Printf("P%d tenta pegar G%d e G%d\n", id+1, leftFork+1, rightFork+1)

		//1 filósofo bloqueia os demais, enquanto está pegando os DOIS garfos
		pickForksMutex.Lock()

		forks[leftFork].Lock()
		forks[rightFork].Lock()

		pickForksMutex.Unlock()

		//Conseguiu os garfos
		fmt.Printf("P%d pegou G%d e G%d\n", id+1, leftFork+1, rightFork+1)

		//Comendo
		fmt.Printf("P%d: E%d\n", id+1, iteration)
		time.Sleep(time.Duration(timeToEat) * time.Second)

		//Libera os garfos após comer
		forks[leftFork].Unlock()
		forks[rightFork].Unlock()

	}
}

func main() {
	var wg sync.WaitGroup

	fmt.Println("\n[P1] [P2] [P3] [P4] [P5]")
	fmt.Println("V2.go - Jantar dos Filósofos - sem deadlock - sincronização excessiva\n")

	start := time.Now()

	wg.Add(philosophers)

	// Run philosophers concurrently
	for id := 0; id < philosophers; id++ {
		go Philosopher(id, &wg)
	}

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Printf("\nV2 Dinner took %s\n\n", elapsed)
}