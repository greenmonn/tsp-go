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

	greedy()

	duration := time.Now().Sub(startTime)

	fmt.Println("Duration: ", duration)
}

func GAFromRandomPopulation() {
	tour := solver.SolveGA([]*graph.Tour{}, populationNumber, generations)

	fmt.Println("Distance: ", tour.Distance)

	for i := 0; i < optimizationCount; i++ {
		operator.LocalSearchOptimize(tour)
	}

	fmt.Println("Distance after Optimization: ", tour.Distance)

	n := tour.WritePathToFile(filename)

	fmt.Printf("%d Bytes Wrote\n", n)
}

// func GAWithGreedyPopulation() {
// 	greedyTour := solver.SolveGreedy()

// }

func greedy() {
	tour := solver.SolveGreedy()

	fmt.Println("Distance: ", tour.Distance)

	tour = operator.LocalSearchOptimize(tour)

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
