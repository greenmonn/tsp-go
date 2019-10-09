package main

import (
	"flag"
	"fmt"
	"math/rand"
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

	rand.Seed(time.Now().UnixNano())

	graph.SetGraphFromFile("problems/" + filename + ".tsp")

	startTime := time.Now()

	GAFromGreedyPopulation()

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

func GAFromGreedyPopulation() {
	tours := make([]*graph.Tour, populationNumber)

	for i := 0; i < populationNumber; i++ {
		tour := solver.PartialRandomGreedy()
		fmt.Println("Random Greedy: ", tour.Distance)

		tours[i] = tour
	}

	tour := solver.SolveGA(tours, populationNumber, generations)

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

	operator.LocalSearchOptimize(tour, -1)

	fmt.Println("Distance: ", tour.Distance)

}

func IterativeOptimization() {
	tour := graph.NewRandomTour()

	for i := 0; i < optimizationCount; i++ {
		operator.Optimize(tour)
	}

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
	optimizeGap := 1

	tours := make([]*graph.Tour, populationNumber)

	for i := 0; i < populationNumber; i++ {
		tour := solver.PartialRandomGreedy()
		fmt.Println("Random Greedy: ", tour.Distance)

		tours[i] = tour
	}

	optTour := solver.SolveMA(tours, populationNumber, generations, optimizeGap)

	fmt.Println("Final Best Distance: ", optTour.Distance)

	n := optTour.WritePathToFile(filename)

	fmt.Printf("%d Bytes Wrote\n", n)
}

func MAWithRandomPopulation() {

	optTour := solver.SolveMA([]*graph.Tour{}, populationNumber, generations, 10)

	fmt.Println("Final Best Distance: ", optTour.Distance)

	n := optTour.WritePathToFile(filename)

	fmt.Printf("%d Bytes Wrote\n", n)
}

func greedy() {
	tour := solver.SolveGreedy()

	fmt.Println("Distance: ", tour.Distance)

	operator.LocalSearchOptimize(tour, -1)

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
