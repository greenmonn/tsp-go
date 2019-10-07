package solver

import (
	"math"

	"github.com/greenmonn/tsp-go/graph"
)

type Population struct {
	Tours []*graph.Tour
	N     int
}

func NewPopulation(N int, tours []*graph.Tour) *Population {
	return &Population{Tours: tours, N: N}
}

func NewRandomPopulation(N int) *Population {
	tours := make([]*graph.Tour, N)

	for i := 0; i < N; i++ {
		tours[i] = graph.NewRandomTour()
	}

	return &Population{Tours: tours, N: N}
}

func (p *Population) bestTour() *graph.Tour {
	minDistance := math.MaxFloat64
	var currentBest *graph.Tour

	for _, t := range p.Tours {
		if t.Distance < minDistance {
			minDistance = t.Distance
			currentBest = t
		}
	}

	return currentBest
}
