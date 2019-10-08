package operator

import (
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

func selectRandomParents(parents []*graph.Tour) *graph.Tour {
	unit := float64(1.0 / len(parents))
	probability := unit

	r := rand.Float64()

	for _, p := range parents {
		if r < probability {
			return p
		}

		probability += unit
	}

	return parents[len(parents)-1]
}

func EdgeRecombinationCrossover(parent1 *graph.Tour, parent2 *graph.Tour) (children []*graph.Tour) {

	N := graph.GetNodesCount()

	childPath := make([]*graph.Node, N)

	parent1.UpdateConnections()
	parent2.UpdateConnections()

	node := selectRandomParents([]*graph.Tour{parent1, parent2}).GetNode(0)

	index := 0
	connectionsUnion := make(map[int][]*graph.Node)

	usedNodes := make(map[int]bool)
	unUsedNodes := make(map[int]*graph.Node)

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

		if n.Connected[0].ID != connections[0].ID && n.Connected[1].ID != connections[0].ID {
			newConnections = append(newConnections, connections[0])
		}

		if n.Connected[0].ID != connections[1].ID && n.Connected[1].ID != connections[1].ID {
			newConnections = append(newConnections, connections[1])
		}

		newConnections = append(newConnections, n.Connected...)

		connectionsUnion[n.ID] = newConnections
	}

	for {
		childPath[index] = node
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
	child.FromPath(childPath)

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
	return []*graph.Tour{parent1, parent2}
}
