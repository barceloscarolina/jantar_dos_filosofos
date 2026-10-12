
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numPhilosophers = 5
	rounds          = 20
	timeToThink     = 2 * time.Second
	timeToEat       = 3 * time.Second
)

func philosopher(id int, forks []sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	left := id
	right := (id + 1) % numPhilosophers

	for iteration := 1; iteration <= rounds; iteration++ {
		// Pensar
		fmt.Printf("Filosofo %d pensando (rodada %d)\n", id+1, iteration)
		time.Sleep(timeToThink)

		// Pegar o garfo esquerdo
		forks[left].Lock()
		fmt.Printf("Filosofo %d pegou o garfo esquerdo\n", id+1)

		// Tentar pegar o garfo direito
		forks[right].Lock()
		fmt.Printf("Filosofo %d pegou o garfo direito\n", id+1)

		// Comer
		fmt.Printf("Filosofo %d comendo (rodada %d)\n", id+1, iteration)
		time.Sleep(timeToEat)

		// Liberar os garfos
		forks[right].Unlock()
		forks[left].Unlock()

		fmt.Printf("Filosofo %d terminou de comer\n", id+1)
	}

	fmt.Printf("Filosofo %d terminou todas as rodadas\n", id+1)
}

func main() {
	var forks [numPhilosophers]sync.Mutex
	var wg sync.WaitGroup

	start := time.Now()

	fmt.Println("V1 - Filosofos Jantando (possibilidade de deadlock)")
	fmt.Println("-----------------------------------------------")

	for id := 0; id < numPhilosophers; id++ {
		wg.Add(1)
		go philosopher(id, forks[:], &wg)
	}

	wg.Wait()

	fmt.Printf("\nTempo total: %v\n", time.Since(start))
	fmt.Println("Todos os filosofos terminaram.")
}
