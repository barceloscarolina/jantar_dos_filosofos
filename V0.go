// Example code filo_seq.go - Sequential Dinning Philosophers
// Since is sequential, elapsed timne could be simply calculated by (roubnds*philosophers*(time_to_think+time_to_eat))
// By Cesar De Rose
// acrescentadas as métricas do Trabalho 1 (FPPD)

package main

import (
	"fmt"
	"time"
)

var eatingNow int // quantos estão comendo agora
var maxEating int // maior valor de eatingNow observado
var cycles [5]int // ciclos concluídos por filósofo

func Philosopher(id, iteration int) {
	time_to_think := 2
	time_to_eat := 3

	// think
	for s := 1; s <= id*6; s++ {
		fmt.Printf(" ")
	}
	fmt.Printf(" T%d\n", iteration) // Thinking
	time.Sleep(time.Duration(time_to_think) * time.Second)

	// eat
	eatingNow++
	if eatingNow > maxEating {
		maxEating = eatingNow
	}

	for s := 1; s <= id*6; s++ {
		fmt.Printf(" ")
	}
	fmt.Printf(" E%d\n", iteration) // Eating
	time.Sleep(time.Duration(time_to_eat) * time.Second)

	eatingNow--
	cycles[id]++
}

func main() {
	philosophers := 5
	rounds := 20

	fmt.Println("\n[P1] [P2] [P3] [P4] [P5]\n")

	// run philosophers

	start := time.Now()

	for i := 1; i <= rounds; i++ {
		for j := 0; j < philosophers; j++ {
			Philosopher(j, i)
		}
	}
	elapsed := time.Since(start)

	var esperaMedia time.Duration = 0
	var maiorEspera time.Duration = 0

	fmt.Println("\n===== RESUMO V0 =====")
	fmt.Printf("Tempo total: %s\n", elapsed)
	fmt.Printf("Tempo médio de espera pelos garfos: %s\n", esperaMedia)
	fmt.Printf("Maior tempo de espera: %s\n", maiorEspera)
	fmt.Printf("Máximo de filósofos comendo simultaneamente: %d\n", maxEating)
	for j := 0; j < philosophers; j++ {
		fmt.Printf("P%d concluiu %d ciclos\n", j+1, cycles[j])
	}
}
