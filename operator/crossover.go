package operator

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/greenmonn/tsp-go/graph"
)

func OrderCrossover(parent1 *graph.Tour, parent2 *graph.Tour) []*graph.Tour {
	N := graph.GetNodesCount()

	child1 := graph.NewTour()
	child2 := graph.NewTour()

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

	return []*graph.Tour{child1, child2}
}

func EdgeRecombinationCrossover(parent1 *graph.Tour, parent2 *graph.Tour) (children []*graph.Tour) {
	parent1.UpdateConnections()
	parent2.UpdateConnections()

	N := graph.GetNodesCount()

	K := make([]*graph.Node, N)

	node := parent1.GetNode(0)

	if rand.Float64() < 0.5 {
		node = parent2.GetNode(0)
	}

	index := 0
	usedNodes := make(map[int]bool)
	unUsedNodes := make(map[int]*graph.Node)
	connectionsUnion := make(map[int][]*graph.Node)

	for i := 0; i < N; i++ {
		n := parent1.GetNode(i)

		unUsedNodes[n.ID] = n
		connectionsUnion[n.ID] = make([]*graph.Node, 2)
		copy(connectionsUnion[n.ID], n.Connected)
	}

	for i := 0; i < N; i++ {
		n := parent2.GetNode(i)
		connections := connectionsUnion[n.ID]
		newConnections := make([]*graph.Node, 0, 4)

		for _, n2 := range n.Connected {
			for _, n1 := range connections {
				if n2.ID == n1.ID {
					n2.ID = -n2.ID
				} else {
					newConnections = append(newConnections, n1)
				}
				newConnections = append(newConnections, n2)
			}
		}
		connectionsUnion[n.ID] = newConnections
	}

	for {
		K[index] = node
		index++

		if index == N {
			break
		}

		usedNodes[node.ID] = true
		delete(unUsedNodes, node.ID)

		neighbor := false
		sort.Slice(connectionsUnion[node.ID], func(i, j int) bool {
			return connectionsUnion[node.ID][i].ID < connectionsUnion[node.ID][j].ID
		})
		for _, next := range connectionsUnion[node.ID] {
			if usedNodes[next.ID] == true {
				continue
			}

			if next.ID < 0 {
				next.ID = -next.ID
			}
			neighbor = true
			node = next
			break
		}

		if neighbor == false {
			_, node = chooseRandomID(unUsedNodes)
		}
	}

	child := graph.NewTour()
	child.FromPath(K)

	return []*graph.Tour{child}
}

func chooseRandomID(unUsedNodes map[int]*graph.Node) (int, *graph.Node) {
	i := rand.Intn(len(unUsedNodes))
	for key, node := range unUsedNodes {
		if i == 0 {
			return key, node
		}
		i--
	}
	return -1, nil
}

func GXCrossover(parent1 *graph.Tour, parent2 *graph.Tour) (child1 *graph.Tour, child2 *graph.Tour) {
	return
}

func NoCrossover(parent1 *graph.Tour, parent2 *graph.Tour) []*graph.Tour {
	fmt.Println(graph.PathToIDs(parent1.Path))
	fmt.Println(graph.PathToIDs(parent2.Path))
	return []*graph.Tour{parent1, parent2}
}
