package operator

import (
	"math/rand"

	"github.com/greenmonn/tsp-go/graph"
)

func SwapPositionMutate(t *graph.Tour, mutationRate float64) {
	N := graph.GetNodesCount()

	for i := 0; i < N; i++ {
		if rand.Float64() >= mutationRate {
			continue
		}

		swapIndex := rand.Intn(N)

		temp := t.GetNode(i)
		t.SetNode(i, t.GetNode(swapIndex))
		t.SetNode(swapIndex, temp)
	}

	t.UpdateDistance()
}

func EdgeExchangeMutate(t *graph.Tour, mutationRate float64) {
	N := graph.GetNodesCount()

	t.UpdateConnections()

	for i := 0; i < N; i++ {
		if rand.Float64() >= mutationRate {
			continue
		}

		gap := rand.Intn(N - 3)

		SwapTwoEdges(t, i, i+2+gap, false)
	}

	t.UpdateDistance()
}

func EdgeExchangeMutateForGX(t *graph.Tour, mutationRate float64) {
	N := len(t.FlexEdges)

	t.UpdateConnections()

	for i, e1 := range t.FlexEdges {
		if rand.Float64() >= mutationRate {
			continue
		}

		gap := rand.Intn(N - 1)
		j := i + 1 + gap

		for j >= N {
			j -= N
		}

		e2 := t.FlexEdges[j]

		if e1.From.ID == e2.From.ID || e1.From.ID == e2.To.ID || e1.To.ID == e2.From.ID || e1.To.ID == e2.To.ID {
			continue
		}

		a, b := e1.From, e1.To
		c, d := findOrderedVertices(a, b, e2)

		replace(a, b, c)
		replace(b, a, d)
		replace(c, d, a)
		replace(d, c, b)

		t.FlexEdges[i] = graph.NewEdge(a, c)
		t.FlexEdges[j] = graph.NewEdge(b, d)
	}

	t.FromNodes(t.Path)
	t.UpdateDistance()
	//t.UpdateEdges()
}
