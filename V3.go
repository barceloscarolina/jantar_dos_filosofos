// V3.go - Jantar dos Filósofos sem deadlock (ordenação dos garfos)
// Grupo B

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

// Métricas
var counterMutex sync.Mutex // mutex que protege o acesso a eatingNow e maxEating
var eatingNow int           // quantos estão comendo agora
var maxEating int           // maior valor de eatingNow observado

// Cada filósofo escreve apenas na sua própria posição
var waitTotal [philosophers]time.Duration
var waitMax [philosophers]time.Duration
var cycles [philosophers]int

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
		fmt.Printf("P%d tentando pegar G%d e G%d\n", id+1, first+1, second+1)
		startWait := time.Now()

		forks[first].Lock()
		fmt.Printf("P%d pegou G%d\n", id+1, first+1)

		fmt.Printf("P%d possui G%d e espera G%d\n", id+1, first+1, second+1)
		forks[second].Lock()

		wait := time.Since(startWait)
		waitTotal[id] += wait
		if wait > waitMax[id] {
			waitMax[id] = wait
		}

		// comendo
		counterMutex.Lock()
		eatingNow++
		if eatingNow > maxEating {
			maxEating = eatingNow
		}
		counterMutex.Unlock()

		fmt.Printf("P%d: E%d\n", id+1, iteration)
		time.Sleep(time.Duration(timeToEat) * time.Second)

		counterMutex.Lock()
		eatingNow--
		counterMutex.Unlock()

		// soltando garfos
		forks[second].Unlock()
		forks[first].Unlock()

		cycles[id]++
		fmt.Printf("P%d liberou G%d e G%d\n", id+1, first+1, second+1)
	}
}

func main() {
	var wg sync.WaitGroup

	fmt.Println("\n[P1] [P2] [P3] [P4] [P5]")
	fmt.Println("V3 - Jantar dos Filósofos sem deadlock (ordenação dos garfos)")

	start := time.Now()

	wg.Add(philosophers)
	for id := 0; id < philosophers; id++ {
		go Philosopher(id, &wg)
	}
	wg.Wait()

	elapsed := time.Since(start)

	// cálculo das métricas
	var somaEspera time.Duration
	var maiorEspera time.Duration
	for i := 0; i < philosophers; i++ {
		somaEspera += waitTotal[i]
		if waitMax[i] > maiorEspera {
			maiorEspera = waitMax[i]
		}
	}
	esperaMedia := somaEspera / time.Duration(philosophers*rounds)

	fmt.Println("\n===== RESUMO V3 =====")
	fmt.Printf("Tempo total: %s\n", elapsed)
	fmt.Printf("Tempo médio de espera pelos garfos: %s\n", esperaMedia)
	fmt.Printf("Maior tempo de espera: %s\n", maiorEspera)
	fmt.Printf("Máximo de filósofos comendo simultaneamente: %d\n", maxEating)
	for i := 0; i < philosophers; i++ {
		fmt.Printf("P%d concluiu %d ciclos\n", i+1, cycles[i])
	}
}
