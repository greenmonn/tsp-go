package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
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

	printLog bool
)

func main() {
	parseArguments()

	rand.Seed(time.Now().UnixNano())

	if !printLog {
		fpLog, err := os.OpenFile("log-"+filename+".txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer fpLog.Close()

		log.SetOutput(fpLog)
	}

	graph.SetGraphFromFile("problems/" + filename + ".tsp")

	startTime := time.Now()

	tour := LocalSearchFromPartialGreedyTour() // Change here

	fmt.Println("Distance: ", tour.Distance)

	n := tour.WritePathToFile(filename)

	log.Printf("%d Bytes Wrote\n", n)

	duration := time.Now().Sub(startTime)

	log.Println("Duration: ", duration)
}

// Random Population Based

func GAFromRandomPopulation() *graph.Tour {
	/* recommended: -p=50, -f=100000 (for fl1400.tsp)
	Best Performance: Edge Recombination Crossover + Edge Exchange Mutation */
	tour := solver.SolveGA([]*graph.Tour{}, populationNumber, generations)

	log.Println("Distance: ", tour.Distance)

	for i := 0; i < optimizationCount; i++ {
		operator.Optimize(tour)
	}

	log.Println("Distance after Optimization: ", tour.Distance)

	return tour
}

func LocalSearchFromRandomTour() *graph.Tour {
	tour := graph.NewRandomTour()

	operator.LocalSearchOptimize(tour, -1)

	log.Println("Distance: ", tour.Distance)

	return tour
}

func GAOptimizeFinalPopulation() *graph.Tour {
	/* recommended: -p=50, -f=100000 (for fl1400.tsp) */
	population := solver.GAOptimize([]*graph.Tour{}, populationNumber, generations)

	log.Println("Best Distance: ", population.BestTour().Distance)

	for _, tour := range population.Tours {
		operator.Optimize(tour)
	}

	best := population.BestTour()

	log.Println("Best Distance after Optimization: ", best.Distance)

	return best
}

// Partially Greedy Population Based

func GAFromGreedyPopulation() *graph.Tour {
	/* recommended: -p=10, -f=10000 (for fl1400.tsp) */
	tours := make([]*graph.Tour, populationNumber)
	graph.SetNearestNeighbors(5)

	for i := 0; i < populationNumber; i++ {
		tour := operator.PartialRandomGreedy()

		tours[i] = tour
	}

	tour := solver.SolveGA(tours, populationNumber, generations)

	log.Println("Distance: ", tour.Distance)

	for i := 0; i < optimizationCount; i++ {
		operator.Optimize(tour)
	}

	log.Println("Distance after Optimization: ", tour.Distance)

	return tour
}

func MAFromGreedyPopulation() *graph.Tour {
	/* recommended: -p=10, -f=10 */

	optimizeGap := 1

	tours := make([]*graph.Tour, populationNumber)

	for i := 0; i < populationNumber; i++ {
		tour := operator.PartialRandomGreedy()

		tours[i] = tour
	}

	optTour := solver.SolveMA(tours, populationNumber, generations, optimizeGap)

	log.Println("Final Best Distance: ", optTour.Distance)

	return optTour
}

func LocalSearchFromPartialGreedyTour() *graph.Tour {
	/* recommended: -p=10 */

	tours := make([]*graph.Tour, populationNumber)
	graph.SetNearestNeighbors(5)

	for i := 0; i < populationNumber; i++ {
		tour := operator.PartialRandomGreedy()

		tours[i] = tour
	}

	for _, tour := range tours {
		for i := 0; i < optimizationCount; i++ {
			operator.Optimize(tour)
		}
	}

	for _, tour := range tours {
		operator.LocalSearchOptimize(tour, -1)
	}

	population := solver.NewPopulation(populationNumber, tours)

	best := population.BestTour()
	log.Println("Distance: ", best.Distance)

	// operator.LocalSearchOptimize(best, -1)
	// log.Println("Distance after final optimization: ", best.Distance)

	return best

}

// Deterministic Optimization

func IterativeOptimization() *graph.Tour {
	tour := graph.NewRandomTour()

	for i := 0; i < optimizationCount; i++ {
		operator.Optimize(tour)
	}

	log.Println("Distance: ", tour.Distance)

	return tour
}

func Greedy() *graph.Tour {
	tour := solver.SolveGreedy()

	log.Println("Distance: ", tour.Distance)

	operator.LocalSearchOptimize(tour, -1)

	log.Println("Distance after Optimization: ", tour.Distance)

	return tour
}

func parseArguments() {
	filenamePtr := flag.String("filename", "rl11849", "filename of a TSP instance from TSPLIB symmetric euc2d problems set")

	populationNumberPtr := flag.Int("p", 50, "population number")
	generationsPtr := flag.Int("f", 100, "generations for fitness evaluations")
	optimizationCountPtr := flag.Int("o", 2, "count of iterative optimizations for final solution")

	printLogPtr := flag.Bool("v", true, "print log to stdout")

	flag.Parse()
	filename = *filenamePtr
	populationNumber = *populationNumberPtr
	generations = *generationsPtr
	optimizationCount = *optimizationCountPtr
	printLog = *printLogPtr
}
