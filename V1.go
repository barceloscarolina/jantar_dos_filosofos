
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
	timeout         = 5 * time.Minute
)

type Metrics struct {
	mu        sync.Mutex
	totalWait time.Duration
	maxWait   time.Duration
	waitCount int
	eatingNow int
	maxEating int
	cycles    [numPhilosophers]int
	states    [numPhilosophers]string
}

func (m *Metrics) updateState(id int, state string) {
	m.mu.Lock()
	m.states[id] = state
	m.mu.Unlock()
}

func philosopher(id int, forks []sync.Mutex, wg *sync.WaitGroup, m *Metrics) {
	defer wg.Done()

	left := id
	right := (id + 1) % numPhilosophers

	for iteration := 1; iteration <= rounds; iteration++ {
		m.updateState(id, "pensando")
		fmt.Printf("P%d: pensando, ciclo %d\n", id+1, iteration)
		time.Sleep(timeToThink)

		startWait := time.Now()

		m.updateState(id, fmt.Sprintf("aguardando G%d", left+1))
		forks[left].Lock()

		m.updateState(id, fmt.Sprintf(
			"possui G%d, aguardando G%d", left+1, right+1,
		))
		fmt.Printf("P%d: pegou G%d\n", id+1, left+1)

		forks[right].Lock()

		wait := time.Since(startWait)

		m.mu.Lock()
		m.totalWait += wait
		m.waitCount++
		if wait > m.maxWait {
			m.maxWait = wait
		}
		m.states[id] = "comendo"
		m.eatingNow++
		if m.eatingNow > m.maxEating {
			m.maxEating = m.eatingNow
		}
		m.mu.Unlock()

		fmt.Printf("P%d: comendo, ciclo %d\n", id+1, iteration)
		time.Sleep(timeToEat)

		m.mu.Lock()
		m.eatingNow--
		m.cycles[id]++
		m.states[id] = "liberando garfos"
		m.mu.Unlock()

		forks[right].Unlock()
		forks[left].Unlock()
		fmt.Printf("P%d: terminou ciclo %d\n", id+1, iteration)
	}

	m.updateState(id, "concluido")
}

func printMetrics(m *Metrics, elapsed time.Duration, completed bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Println("\n========== RESULTADOS V1 ==========")

	if completed {
		fmt.Printf("Tempo total: %.4f s\n", elapsed.Seconds())
	} else {
		fmt.Printf("Tempo ate o limite: %.4f s\n", elapsed.Seconds())
		fmt.Println("Tempo total de conclusao: N/A")
	}

	if m.waitCount > 0 {
		average := m.totalWait.Seconds() / float64(m.waitCount)
		fmt.Printf("Tempo medio de espera: %.4f s\n", average)
	} else {
		fmt.Println("Tempo medio de espera: N/A")
	}

	fmt.Printf("Maior espera concluida: %.4f s\n", m.maxWait.Seconds())
	fmt.Printf("Maximo comendo simultaneamente: %d\n", m.maxEating)

	fmt.Print("Ciclos concluidos (P1 a P5): ")
	for i, count := range m.cycles {
		if i > 0 {
			fmt.Print("/")
		}
		fmt.Print(count)
	}
	fmt.Println()

	fmt.Println("\nEstado dos filosofos:")
	for i, state := range m.states {
		fmt.Printf("P%d: %s\n", i+1, state)
	}
}

func main() {
	var forks [numPhilosophers]sync.Mutex
	var wg sync.WaitGroup
	var metrics Metrics

	start := time.Now()
	done := make(chan struct{})

	fmt.Println("V1 - Filosofos Jantando")
	fmt.Println("Versao concorrente sujeita a deadlock")

	for id := 0; id < numPhilosophers; id++ {
		wg.Add(1)
		go philosopher(id, forks[:], &wg, &metrics)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("\nExecucao concluida normalmente.")
		printMetrics(&metrics, time.Since(start), true)

	case <-time.After(timeout):
		fmt.Println("\nLimite de observacao atingido.")
		fmt.Println("Verificar estados para identificar possivel deadlock.")
		printMetrics(&metrics, time.Since(start), false)
	}
}
