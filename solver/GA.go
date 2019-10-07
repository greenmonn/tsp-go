package solver

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/greenmonn/tsp-go/graph"
	"github.com/greenmonn/tsp-go/operator"
)

const (
	fitnessEvaluation = 100
	tournamentSize    = 6
	mutationRate      = 0.01
	elitism           = true
)

func SolveGA(initialTours []*graph.Tour, populationNumber int) *graph.Tour {
	N := populationNumber
	var population *Population

	if len(initialTours) == 0 {
		population = NewRandomPopulation(N)
	} else {
		population = NewPopulation(N, initialTours)
	}

	for i := 0; i < fitnessEvaluation; i++ {
		fmt.Printf("\n%dth Generation\n", i+1)

		population = EvolvePopulation(population)
	}

	bestTour := population.bestTour()
	bestTour.FromPath(bestTour.Path)
	operator.Optimize(bestTour)

	return bestTour
}

func EvolvePopulation(p *Population) *Population {
	children := make([]*graph.Tour, p.N)

	offset := 0

	if elitism {
		elite := p.bestTour()

		fmt.Println("Current Best Distance: ", elite.Distance)

		children[offset] = elite
		offset++
	}

	var (
		parent1, parent2 *graph.Tour
		child1, child2   *graph.Tour
	)

	for offset < p.N {
		parent1, parent2 = selectParents(p)

		child1, child2 = crossover(parent1, parent2)

		mutate(child1)

		if child1 == nil {
			log.Fatalln("Invalid child")
		}

		children[offset] = child1
		offset++

		if offset == p.N {
			break
		}

		mutate(child2)

		if child2 == nil {
			log.Fatalln("Invalid child")
		}

		children[offset] = child2
		offset++
	}

	return NewPopulation(p.N, children)
}

func selectParents(p *Population) (parent1 *graph.Tour, parent2 *graph.Tour) {
	return selectTournament(p)
}

func selectTournament(p *Population) (parent1 *graph.Tour, parent2 *graph.Tour) {
	parents := make([]*graph.Tour, 2)

	var (
		tournament *Population
	)

	for i := 0; i < 2; i++ {
		candidates := make([]*graph.Tour, tournamentSize)

		var (
			randomIndex int
		)

		for j := 0; j < tournamentSize; j++ {
			randomIndex = rand.Intn(p.N)

			candidates[j] = p.Tours[randomIndex]
		}

		tournament = NewPopulation(tournamentSize, candidates)
		parents[i] = tournament.bestTour()
	}

	parent1 = parents[0]
	parent2 = parents[1]

	return
}

func crossover(parent1 *graph.Tour, parent2 *graph.Tour) (child1 *graph.Tour, child2 *graph.Tour) {
	return operator.OrderCrossover(parent1, parent2)
}

func mutate(t *graph.Tour) {
	operator.SwapPositionMutate(t, mutationRate)
}
