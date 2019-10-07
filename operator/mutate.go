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
