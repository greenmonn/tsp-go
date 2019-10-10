package solver

import (
	"log"

	"github.com/greenmonn/tsp-go/graph"
	"github.com/greenmonn/tsp-go/operator"
)

func SolveMA(initialTours []*graph.Tour, populationNumber int, generations int, optimizeGap int) *graph.Tour {
	// Memetic Algorithm (GA + Local Search)

	N := populationNumber
	var population *Population

	if len(initialTours) == 0 {
		population = NewRandomPopulation(N)
	} else {
		population = NewPopulation(N, initialTours)
	}

	for i := 0; i < generations; i++ {
		log.Printf("\n%dth Generation\n", i+1)

		population = EvolvePopulation(population)

		// Optimize whole population: individuals would be 'near' local optimum
		if i%optimizeGap != 0 {
			continue
		}

		for _, tour := range population.Tours {
			operator.FastLocalSearchOptimize(tour)

		}
	}

	operator.LocalSearchOptimize(population.BestTour(), -1)

	return population.BestTour()
}
