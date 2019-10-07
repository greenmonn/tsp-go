package operator

import (
	"math/rand"

	"github.com/greenmonn/tsp-go/graph"
)

func OrderCrossover(parent1 *graph.Tour, parent2 *graph.Tour) (child1 *graph.Tour, child2 *graph.Tour) {
	N := graph.GetNodesCount()

	child1 = graph.NewTour()
	child2 = graph.NewTour()

	startIndex := rand.Intn(N)
	endIndex := rand.Intn(N)

	if startIndex > endIndex {
		startIndex, endIndex = endIndex, startIndex
	}

	UsedNodeIds1 := make(map[int]bool)
	UsedNodeIds2 := make(map[int]bool)

	for i := startIndex; i <= endIndex; i++ {
		node1 := parent1.GetNode(i)
		UsedNodeIds1[node1.ID] = true

		node2 := parent2.GetNode(i)
		UsedNodeIds2[node2.ID] = true

		child1.SetNode(i, node1)
		child2.SetNode(i, node2)
	}

	parentIndex := endIndex + 1

	offset1 := endIndex + 1
	offset2 := endIndex + 1
	for {
		node1 := parent2.GetNode(parentIndex)
		node2 := parent1.GetNode(parentIndex)

		if UsedNodeIds1[node1.ID] != true {
			child1.SetNode(offset1, node1)
			offset1++
		}

		if UsedNodeIds2[node2.ID] != true {
			child2.SetNode(offset2, node2)
			offset2++
		}

		if offset1 == (N+startIndex) && offset2 == (N+startIndex) {
			break
		}

		parentIndex++
	}

	child1.UpdateDistance()
	child2.UpdateDistance()

	return
}

func NoCrossover(parent1 *graph.Tour, parent2 *graph.Tour) (child1 *graph.Tour, child2 *graph.Tour) {
	child1 = parent1
	child2 = parent2
	return
}
