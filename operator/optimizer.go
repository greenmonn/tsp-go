package operator

import (
	"fmt"

	"github.com/greenmonn/tsp-go/graph"
)

func Optimize(tour *graph.Tour) {
	// 2-opt pairwise exchange
	tour.UpdateConnections()

	N := graph.GetNodesCount()

	for i := 0; i < N; i++ {
		for j := i + 2; j < i+N-1; j++ {
			changed := SwapTwoEdges(tour, i, j, true)

			// Update path from connections
			if changed {
				tour.FromNodes(tour.Path)
			}
		}
	}
}

func LocalSearchOptimize(tour *graph.Tour) {
	// Search neighbors until no better neighbors exist
	tour.UpdateConnections()

	iteration := 0

	for {
		found := find2OptBetterMove(tour)

		iteration++
		if iteration%10 == 0 {
			fmt.Println("\nIteration count: ", iteration)
			fmt.Println(tour.Distance)
		}

		if !found {
			fmt.Println("\nIteration count: ", iteration)
			return
		}
	}
}

func find2OptBetterMove(tour *graph.Tour) (found bool) {
	N := graph.GetNodesCount()

	found = false

	for i := 0; i < N; i++ {
		for j := i + 2; j < i+N-1; j++ {
			found = SwapTwoEdges(tour, i, j, true)

			if found {
				tour.FromNodes(tour.Path)
				return
			}
		}
	}

	return
}

func findLKOptFirstBetterNeighbor(tour *graph.Tour) (neighbor *graph.Tour, found bool) {
	return
}

func SwapTwoEdges(tour *graph.Tour, i int, j int, onlyIfBetter bool) bool {
	D := graph.GetDistance

	a := tour.GetNode(i)
	b := tour.GetNode(i + 1)

	c := tour.GetNode(j)
	d := tour.GetNode(j + 1)

	if onlyIfBetter && D(a, b)+D(c, d) <= D(a, c)+D(b, d) {
		return false
	}

	replace := func(node *graph.Node, from *graph.Node, to *graph.Node) {
		for i, n := range node.Connected {
			if n.ID != from.ID {
				continue
			}
			node.Connected[i] = to
		}
	}

	replace(a, b, c)
	replace(b, a, d)
	replace(c, d, a)
	replace(d, c, b)

	return true
}
