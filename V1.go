// V1.go - Concurrent Dining Philosophers subject to deadlock
// FPPD - Trabalho 1
//
// V1: cada filosofo executa concorrentemente.
// Cada filosofo pega primeiro o garfo da esquerda e depois
// tenta pegar o garfo da direita.
//
// No primeiro ciclo, uma barreira garante que todos os filosofos
// tenham pego o garfo esquerdo antes de tentarem pegar o direito.
// Assim, a situacao de deadlock pode ser demonstrada claramente.

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

// Canais utilizados como barreira para tornar
// o deadlock reproduzivel no primeiro ciclo.
var leftForkAcquired = make(chan struct{}, philosophers)
var startRightFork = make(chan struct{})

func Philosopher(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	leftFork := id
	rightFork := (id + 1) % philosophers

	for iteration := 1; iteration <= rounds; iteration++ {

		// Think
		fmt.Printf("P%d: T%d\n", id+1, iteration)
		time.Sleep(time.Duration(timeToThink) * time.Second)

		// Acquire left fork
		fmt.Printf("P%d tentando pegar G%d\n", id+1, leftFork+1)

		forks[leftFork].Lock()

		fmt.Printf("P%d pegou G%d\n", id+1, leftFork+1)

		// No primeiro ciclo, todos devem pegar o garfo esquerdo
		// antes de tentar pegar o garfo direito.
		if iteration == 1 {
			leftForkAcquired <- struct{}{}
			<-startRightFork
		}

		// Acquire right fork
		fmt.Printf(
			"P%d possui G%d e espera G%d\n",
			id+1,
			leftFork+1,
			rightFork+1,
		)

		forks[rightFork].Lock()

		// Eat
		fmt.Printf("P%d: E%d\n", id+1, iteration)
		time.Sleep(time.Duration(timeToEat) * time.Second)

		// Release forks
		forks[rightFork].Unlock()
		forks[leftFork].Unlock()

		fmt.Printf(
			"P%d liberou G%d e G%d\n",
			id+1,
			leftFork+1,
			rightFork+1,
		)
	}
}

func main() {
	var wg sync.WaitGroup

	fmt.Println("\n[P1] [P2] [P3] [P4] [P5]")
	fmt.Println("V1 - Concurrent Dining Philosophers subject to deadlock\n")

	start := time.Now()

	wg.Add(philosophers)

	// Run philosophers concurrently
	for id := 0; id < philosophers; id++ {
		go Philosopher(id, &wg)
	}

	// Espera os cinco filosofos pegarem seus garfos esquerdos.
	for i := 0; i < philosophers; i++ {
		<-leftForkAcquired
	}

	fmt.Println("\nTodos os filosofos pegaram o garfo esquerdo.")
	fmt.Println("Agora todos tentarao pegar o garfo direito.\n")

	close(startRightFork)

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Printf("\nConcurrent Dinner took %s\n\n", elapsed)
}