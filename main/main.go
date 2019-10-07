package main

import (
	"fmt"

	"github.com/greenmonn/tsp-go/graph"
	"github.com/greenmonn/tsp-go/operator"
	"github.com/greenmonn/tsp-go/solver"
)

func main() {
	filename := "rl11849"
	graph.SetGraphFromFile("problems/" + filename + ".tsp")

	GAFromRandomPopulation(filename)
}

func GAFromRandomPopulation(filename string) {
	tour := solver.SolveGA([]*graph.Tour{}, 50)

	fmt.Println("Distance: ", tour.Distance)
	fmt.Println("Path: ", graph.PathToIDs(tour.Path))

	tour.FromPath(tour.Path)
	operator.Optimize(tour)

	fmt.Println("Distance after Optimization: ", tour.Distance)

	n := tour.WritePathToFile(filename)

	fmt.Printf("%d Bytes Wrote\n", n)
}

func greedy(filename string) {
	tour := solver.SolveGreedy()

	fmt.Println("Distance: ", tour.Distance)

	operator.Optimize(tour)

	fmt.Println("Distance after Optimization: ", tour.Distance)

	n := tour.WritePathToFile(filename)

	fmt.Printf("%d Bytes Wrote\n", n)
}
