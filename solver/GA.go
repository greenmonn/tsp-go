package solver

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/greenmonn/tsp-go/graph"
	"github.com/greenmonn/tsp-go/operator"
)

const (
	tournamentSize = 6
	mutationRate   = 0.01
	elitism        = true
)

func SolveGA(initialTours []*graph.Tour, populationNumber int, generations int) *graph.Tour {
	population := GAOptimize(initialTours, populationNumber, generations)

	return population.BestTour()
}

func GAOptimize(initialTours []*graph.Tour, populationNumber int, generations int) *Population {
	N := populationNumber
	var population *Population

	if len(initialTours) == 0 {
		population = NewRandomPopulation(N)
	} else {
		population = NewPopulation(N, initialTours)
	}

	for i := 0; i < generations; i++ {
		fmt.Printf("\n%dth Generation\n", i+1)

		population = EvolvePopulation(population)
	}

	return population
}

func SolveMA(initialTours []*graph.Tour, populationNumber int, generations int) *graph.Tour {
	// Memetic Algorithm (GA + Local Search)

	N := populationNumber
	var population *Population

	if len(initialTours) == 0 {
		population = NewRandomPopulation(N)
	} else {
		population = NewPopulation(N, initialTours)
	}

	for i := 0; i < generations; i++ {
		fmt.Printf("\n%dth Generation\n", i+1)

		population = EvolvePopulation(population)

		// Optimize whole population: individuals would be 'near' local optimum
		for _, tour := range population.Tours {
			operator.Optimize(tour)
		}
	}

	// Local Search
	operator.LocalSearchOptimize(population.BestTour())

	return population.BestTour()
}

func EvolvePopulation(p *Population) *Population {
	tours := make([]*graph.Tour, p.N)

	offset := 0

	if elitism {
		elite := p.BestTour()

		fmt.Println("Current Best Distance: ", elite.Distance)

		tours[offset] = elite
		offset++
	}

	var (
		parent1, parent2 *graph.Tour
		children         []*graph.Tour
	)

	for offset < p.N {
		parent1, parent2 = selectParents(p)

		children = crossover(parent1, parent2)

		for _, child := range children {
			mutate(child)

			if child == nil {
				log.Fatalln("Invalid child")
			}

			tours[offset] = child
			offset++

			if offset == p.N {
				return NewPopulation(p.N, tours)
			}
		}
	}

	return NewPopulation(p.N, tours)
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
		parents[i] = tournament.BestTour()
	}

	parent1 = parents[0]
	parent2 = parents[1]

	return
}

func crossover(parent1 *graph.Tour, parent2 *graph.Tour) (children []*graph.Tour) {
	return operator.EdgeRecombinationCrossover(parent1, parent2)
}

func mutate(t *graph.Tour) {
	operator.EdgeExchangeMutate(t, mutationRate)
}
