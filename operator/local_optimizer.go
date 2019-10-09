package operator

import (
	"fmt"
	"log"
	"math"

	"github.com/greenmonn/tsp-go/graph"
)

func Optimize(tour *graph.Tour) {
	// 2-opt pairwise exchange (iterate over a path)
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

func FastLocalSearchOptimize(tour *graph.Tour) {
	// Search for only non-fixed edges
	// Set by GX Crossover (Only compatible with the tour returned from GX Crossover)
	iteration := 0
	for {
		found := find2OptBetterMoveFromEdges(tour)

		iteration++
		if iteration%10 == 0 {
			fmt.Println("\nIteration count: ", iteration)
			tour.FromNodes(tour.Path)
			fmt.Println("Distance: ", tour.Distance)
		}

		if !found {
			fmt.Println("\nFINISH - Iteration count: ", iteration)
			// make path from connections
			tour.FromNodes(tour.Path)

			// make edges
			tour.UpdateEdges()
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

func find2OptBetterMoveFromEdges(tour *graph.Tour) (found bool) {
	if tour.FlexEdges == nil {
		log.Fatalln("Uncompatible tour for edge based local search")
	}

	N := len(tour.FlexEdges)

	D := graph.GetDistance

	found = false

	for i := 0; i < N; i++ {
		for j := i + 1; j < N; j++ {
			e1 := tour.FlexEdges[i]
			e2 := tour.FlexEdges[j]

			if e1.From.ID == e2.From.ID || e1.From.ID == e2.To.ID || e1.To.ID == e2.From.ID || e1.To.ID == e2.To.ID {
				continue
			}

			a, b := e1.From, e1.To

			c, d := findOrderedVertices(a, b, e2)

			if D(a, b)+D(c, d) <= D(a, c)+D(b, d) {
				continue
			}

			replace(a, b, c)
			replace(b, a, d)
			replace(c, d, a)
			replace(d, c, b)

			// tour.FromNodes(tour.Path)
			// idPath = graph.PathToIDs(tour.Path)
			// sort.Ints(idPath)
			// if !areEqualIntSlices(idPath, makeRange(1, 1400)) {
			// 	log.Fatalln("Invalid path: after replacement")
			// }

			tour.FlexEdges[i] = graph.NewEdge(a, c)
			tour.FlexEdges[j] = graph.NewEdge(b, d)

			return true
		}
	}

	return
}

func findOrderedVertices(prev *graph.Node, node *graph.Node, e *graph.Edge) (c *graph.Node, d *graph.Node) {
	for {
		for _, next := range node.Connected {
			if next.ID == prev.ID {
				continue
			}

			if next.ID == e.From.ID {
				c = e.From
				d = e.To
				return
			}

			if next.ID == e.To.ID {
				c = e.To
				d = e.From
				return
			}
			prev = node
			node = next
			break
		}
	}
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

	replace(a, b, c)
	replace(b, a, d)
	replace(c, d, a)
	replace(d, c, b)

	tour.FromNodes(tour.Path)

	return true
}

func replace(node *graph.Node, from *graph.Node, to *graph.Node) {
	for i, n := range node.Connected {
		if n.ID != from.ID {
			continue
		}
		node.Connected[i] = to
	}
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func areEqualIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
