package operator

import "github.com/greenmonn/tsp-go/graph"

func Optimize(tour *graph.Tour) {
	// 2-opt pairwise exchange
	tour.UpdateDistance()

	N := graph.GetNodesCount()

	for i := 0; i < N; i++ {
		for j := i + 3; j < i+N-1; j++ {
			SwapTwoEdges(tour, i, j, true)
		}
	}
}

func SwapTwoEdges(tour *graph.Tour, i int, j int, onlyIfBetter bool) {
	D := graph.GetDistance

	a := tour.GetNode(i)
	b := tour.GetNode(i + 1)

	c := tour.GetNode(j)
	d := tour.GetNode(j + 1)

	if onlyIfBetter && D(a, b)+D(c, d) <= D(a, c)+D(b, d) {
		return
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
}
