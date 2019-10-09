package operator

import (
	"container/heap"
	"log"
	"math/rand"
	"sort"

	"github.com/greenmonn/tsp-go/container"

	"github.com/greenmonn/tsp-go/graph"
)

func GXCrossover(parent1 *graph.Tour, parent2 *graph.Tour, cRate float64, nRate float64, iRate float64) (children []*graph.Tour) {
	N := graph.GetNodesCount()

	nodes := graph.CopyNodesFromGraph()

	candidateEdges := graph.NewEdges(nodes)

	childEdges := make(map[string]*graph.Edge)
	flexEdges := make([]*graph.Edge, 0, N)

	sets := make(map[int]*[]*graph.Node)
	setsCount := graph.GetNodesCount()

	for i := 0; i < N; i++ {
		node := nodes[i]
		sets[node.ID] = &[]*graph.Node{node}
	}

	// Copy common edges
	for id, e := range parent1.Edges {
		_, exist := parent2.Edges[id]

		if exist && rand.Float64() < cRate {
			childEdges[id] = graph.NewEdge(nodes[e.From.ID-1], nodes[e.To.ID-1])

			// if childEdges[id].From.ID != e.From.ID {
			// 	log.Fatalln("ID mismatch")
			// }

			// if childEdges[id].To.ID != e.To.ID {
			// 	log.Fatalln("ID mismatch")
			// }

			childEdges[id].UpdateNodes()
			delete(candidateEdges, id)

			mergedSet := append(*sets[e.From.ID], *sets[e.To.ID]...)

			for _, node := range mergedSet {
				sets[node.ID] = &mergedSet
				delete(candidateEdges, graph.EdgeID(e.From.ID, node.ID))
				delete(candidateEdges, graph.EdgeID(e.To.ID, node.ID))
			}

			setsCount--
		}
	}

	// Insert new edges
	for k := 0; k < int(float64(setsCount)*nRate); k++ {
		randIndex := int(float64(N) * rand.Float64())
		i := nodes[randIndex]

		nearests := graph.NearestNeighbors(i.ID)
		id := nearests[rand.Intn(len(nearests))].ID
		// TODO: find nodes by ID
		// Currently guaranteed: ID = index + 1
		j := nodes[id-1]

		// (i, j) feasible?
		if i.Degree == 2 || j.Degree == 2 {
			continue
		}

		if setsCount > 1 && sets[i.ID] == sets[j.ID] {
			continue
		}

		e := graph.NewEdge(i, j)
		edgeID := e.Hash()
		// (i, j) not in a
		_, exist := parent1.Edges[edgeID]
		if exist {
			continue
		}
		// (i, j) not in b
		_, exist = parent2.Edges[edgeID]
		if exist {
			continue
		}

		childEdges[edgeID] = e
		flexEdges = append(flexEdges, e)
		e.UpdateNodes()
		delete(candidateEdges, edgeID)

		mergedSet := append(*sets[i.ID], *sets[j.ID]...)

		for _, node := range mergedSet {
			sets[node.ID] = &mergedSet
			delete(candidateEdges, graph.EdgeID(i.ID, node.ID))
			delete(candidateEdges, graph.EdgeID(j.ID, node.ID))
		}

		setsCount--

	}

	// Inherit edges from parents
	for k := 0; k < int(float64(setsCount)*iRate); k++ {
		parent := selectRandomTour([]*graph.Tour{parent1, parent2})

		hasCandidateEdges := false
		inherits := make([]*graph.Edge, 0)
		for id := range parent.Edges {
			if candidateEdges[id] != nil {
				inherits = append(inherits, candidateEdges[id])
				hasCandidateEdges = true
			}
		}

		if !hasCandidateEdges {
			continue
		}

		var e *graph.Edge

		if len(inherits) == 1 {
			e = graph.NewEdge(nodes[inherits[0].From.ID-1], nodes[inherits[0].To.ID-1])
		} else {

			sort.Slice(inherits, func(i, j int) bool {
				return inherits[i].Distance < inherits[j].Distance
			})

			e1 := inherits[0]
			e2 := inherits[1]

			e = e1
			if rand.Float64() < 0.5 {
				e = e2
			}
		}

		if e.From.Degree == 2 || e.To.Degree == 2 {
			continue
		}

		edgeID := e.Hash()

		childEdges[edgeID] = e
		flexEdges = append(flexEdges, e)
		e.UpdateNodes()
		delete(candidateEdges, edgeID)

		mergedSet := append(*sets[e.From.ID], *sets[e.To.ID]...)

		for _, node := range mergedSet {
			sets[node.ID] = &mergedSet
			delete(candidateEdges, graph.EdgeID(e.From.ID, node.ID))
			delete(candidateEdges, graph.EdgeID(e.To.ID, node.ID))
		}

		setsCount--

	}

	candidates := &container.PriorityQueue{}
	heap.Init(candidates)

	for _, edge := range candidateEdges {
		heap.Push(candidates, edge)
	}

	// greedy completion
	for setsCount > 0 {
		var e *graph.Edge

		e1 := heap.Pop(candidates)
		if candidates.Len() < 1 {
			e = e1.(*graph.Edge)
		} else {
			e2 := heap.Pop(candidates)
			if rand.Float64() < 0.5 {
				e = e2.(*graph.Edge)
				heap.Push(candidates, e1)
			} else {
				e = e1.(*graph.Edge)
				heap.Push(candidates, e2)
			}
		}

		if e.From.Degree == 2 || e.To.Degree == 2 {
			continue
		}

		if setsCount > 1 && sets[e.From.ID] == sets[e.To.ID] {
			continue
		}

		edgeID := e.Hash()

		childEdges[edgeID] = e
		flexEdges = append(flexEdges, e)
		e.UpdateNodes()
		delete(candidateEdges, edgeID)

		mergedSet := append(*sets[e.From.ID], *sets[e.To.ID]...)

		for _, node := range mergedSet {
			sets[node.ID] = &mergedSet
			delete(candidateEdges, graph.EdgeID(e.From.ID, node.ID))
			delete(candidateEdges, graph.EdgeID(e.To.ID, node.ID))
		}

		setsCount--
	}

	if len(childEdges) != N {
		log.Fatalln("Insufficient child edges")
	}

	child := graph.NewTour()
	child.FromNodes(nodes)

	child.Edges = childEdges

	child.FlexEdges = flexEdges

	return []*graph.Tour{child}
}

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

	N := graph.GetNodesCount()

	childPath := make([]*graph.Node, N)

	parent1.UpdateConnections()
	parent2.UpdateConnections()

	node := selectRandomTour([]*graph.Tour{parent1, parent2}).GetNode(0)

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

func NoCrossover(parent1 *graph.Tour, parent2 *graph.Tour) []*graph.Tour {
	return []*graph.Tour{parent1, parent2}
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

func selectFromTwoShortest(candidates []*graph.Edge) (shortest *graph.Edge) {
	shortest = candidates[0]
	secondShortest := candidates[0]

	for _, e := range candidates[1:] {
		if e.Distance < shortest.Distance {
			secondShortest = shortest
			shortest = e
			break
		}
		if e.Distance < secondShortest.Distance {
			secondShortest = e
		}
	}

	if 0.5 < rand.Float64() {
		return secondShortest
	}
	return shortest
}

func selectRandomTour(parents []*graph.Tour) *graph.Tour {
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

func selectRandomNode(nodes []*graph.Node) *graph.Node {
	unit := float64(1.0 / len(nodes))
	probability := unit

	r := rand.Float64()

	for _, p := range nodes {
		if r < probability {
			return p
		}

		probability += unit
	}

	return nodes[len(nodes)-1]
}
