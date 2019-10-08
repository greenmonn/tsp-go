package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/greenmonn/tsp-go/graph"
	"github.com/greenmonn/tsp-go/operator"
	"github.com/greenmonn/tsp-go/solver"
)

var (
	populationNumber  int
	generations       int
	optimizationCount int

	filename string
)

func main() {
	parseArguments()

	graph.SetGraphFromFile("problems/" + filename + ".tsp")

	startTime := time.Now()

	MAWithGreedyPopulation()

	duration := time.Now().Sub(startTime)

	fmt.Println("Duration: ", duration)
}

func GAFromRandomPopulation() {
	tour := solver.SolveGA([]*graph.Tour{}, populationNumber, generations)

	fmt.Println("Distance: ", tour.Distance)

	for i := 0; i < optimizationCount; i++ {
		operator.Optimize(tour)
	}

	fmt.Println("Distance after Optimization: ", tour.Distance)

	n := tour.WritePathToFile(filename)

	fmt.Printf("%d Bytes Wrote\n", n)
}

func LocalSearchFromRandomTour() {
	tour := graph.NewRandomTour()

	operator.LocalSearchOptimize(tour)

	fmt.Println("Distance: ", tour.Distance)

}

func GAOptimizeFinalPopulation() {
	population := solver.GAOptimize([]*graph.Tour{}, populationNumber, generations)

	fmt.Println("Best Distance: ", population.BestTour().Distance)

	for _, tour := range population.Tours {
		operator.Optimize(tour)
	}

	best := population.BestTour()

	fmt.Println("Best Distance after Optimization: ", best.Distance)

	n := best.WritePathToFile(filename)

	fmt.Printf("%d Bytes Wrote\n", n)
}

func MAWithGreedyPopulation() {
	greedyTour := solver.SolveGreedy()
	N := graph.GetNodesCount()

	tours := make([]*graph.Tour, N)

	for i := 0; i < N; i++ {
		tour := graph.NewTour()
		tour.FromPath(greedyTour.Path)

		operator.EdgeExchangeMutate(tour, 0.5)
		tours[i] = tour
	}

	optTour := solver.SolveMA(tours, populationNumber, generations)

	fmt.Println("Final Best Distance: ", optTour.Distance)

	n := optTour.WritePathToFile(filename)

	fmt.Printf("%d Bytes Wrote\n", n)
}

func greedy() {
	tour := solver.SolveGreedy()

	fmt.Println("Distance: ", tour.Distance)

	operator.LocalSearchOptimize(tour)

	fmt.Println("Distance after Optimization: ", tour.Distance)

	n := tour.WritePathToFile(filename)

	fmt.Printf("%d Bytes Wrote\n", n)
}

func parseArguments() {
	filenamePtr := flag.String("filename", "rl11849", "filename of a TSP instance from TSPLIB symmetric euc2d problems set")

	populationNumberPtr := flag.Int("p", 50, "population number")
	generationsPtr := flag.Int("f", 100, "generations for fitness evaluations")
	optimizationCountPtr := flag.Int("o", 2, "count of iterative optimizations for final solution")

	flag.Parse()
	filename = *filenamePtr
	populationNumber = *populationNumberPtr
	generations = *generationsPtr
	optimizationCount = *optimizationCountPtr
}
