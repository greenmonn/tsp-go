package operator

import (
	"fmt"
	"math"

	"github.com/greenmonn/tsp-go/graph"
)

func Optimize(tour *graph.Tour) {
	// 2-opt pairwise exchange (iterate over one path)
	tour.UpdateConnections()

	N := graph.GetNodesCount()

	for i := 0; i < N; i++ {
		for j := i + 2; j < i+N-1; j++ {
			SwapTwoEdges(tour, i, j, true)
		}
	}
}

func FastOptimize(tour *graph.Tour) {
	// still O(N^3) but half time (less efficient)
	tour.UpdateConnections()

	N := graph.GetNodesCount()

	for i := 0; i < N; i++ {
		for j := i + 2; j < N; j++ {
			SwapTwoEdges(tour, i, j, true)
		}
	}
}

func LocalSearchOptimize(tour *graph.Tour, iterationLimit int) {
	// Search neighbors until no better neighbors exist
	if iterationLimit < 0 {
		iterationLimit = math.MaxInt64
	}

	tour.UpdateConnections()

	iteration := 0

	for iteration < iterationLimit {
		found := find2OptBetterMoveFromConnections(tour)

		iteration++
		if iteration%100 == 0 {
			fmt.Println("\nIteration count: ", iteration)
			tour.FromNodes(tour.Path)
			fmt.Println("Distance: ", tour.Distance)
		}

		if !found {
			fmt.Println("\nFINISH - Iteration count: ", iteration)
			tour.FromNodes(tour.Path)
			return
		}
	}
}

func find2OptBetterMove(tour *graph.Tour) (found bool) {
	N := graph.GetNodesCount()

	found = false

	for i := 0; i < N; i++ {
		for j := i + 2; j < N; j++ {
			found = SwapTwoEdges(tour, i, j, true)

			if found {
				return
			}
		}
	}

	return
}

func find2OptBetterMoveFromConnections(tour *graph.Tour) (found bool) {
	// Does not need path restoration
	N := graph.GetNodesCount()

	found = false

	for i := 0; i < N; i++ {
		edge1To := tour.GetNode(i)
		edge1From := edge1To.Connected[0]

		prev := edge1To
		node := edge1To.Connected[1]

		for {
			edge2From := node
			var edge2To *graph.Node
			for _, next := range edge2From.Connected {
				if next.ID == prev.ID {
					continue
				}

				edge2To = next
				break
			}

			if edge2To.ID == edge1From.ID {
				break
			}

			found = SwapTwoEdgesByNodes(tour, edge1From, edge1To, edge2From, edge2To, true)

			if found {
				return
			}

			prev = edge2From
			node = edge2To
		}
	}

	return
}

func findLKOptBetterMove(tour *graph.Tour) (found bool) {
	// TODO
	return
}

func SwapTwoEdgesByNodes(tour *graph.Tour, a *graph.Node, b *graph.Node, c *graph.Node, d *graph.Node, onlyIfBetter bool) bool {
	D := graph.GetDistance

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

	tour.FromNodes(tour.Path)

	return true
}
