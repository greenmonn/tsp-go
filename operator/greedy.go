package operator

import (
	"container/heap"
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/greenmonn/tsp-go/graph"
	"github.com/greenmonn/tsp-go/types"
)

func PartialRandomGreedy() (tour *graph.Tour) {
	fmt.Println("Start Solving Greedy")

	nodes := graph.CopyNodesFromGraph()

	sets, setsCount, connectedEdges := randomConnect(nodes)

	edges := &types.PriorityQueue{}
	heap.Init(edges)

	types.InitEdges(edges, nodes)

	fmt.Println("1: Make Edges Priority Queue")

	GreedyConnect(edges, nodes, sets, setsCount, connectedEdges)

	fmt.Println("2: Connect possible short edges")

	tour = graph.NewTour()
	tour.FromNodes(nodes)

	// For GX Crossover: maintain edges set / non-fixed edges set
	tour.Edges = connectedEdges
	tour.FlexEdges = make([]*graph.Edge, len(connectedEdges))
	index := 0

	for _, e := range connectedEdges {
		tour.FlexEdges[index] = e
		index++
	}

	fmt.Println("3: Construct tour from connected nodes")

	return

}

func GreedyConnect(edges *types.PriorityQueue, nodes []*graph.Node, sets map[int]*[]*graph.Node, setsCount int, connectedEdges map[string]*graph.Edge) {
	if sets == nil {
		sets = make(map[int]*[]*graph.Node)
		setsCount = graph.GetNodesCount()

		for i := 0; i < setsCount; i++ {
			node := nodes[i]
			sets[node.ID] = &[]*graph.Node{node}
		}
	}

	if connectedEdges == nil {
		connectedEdges = make(map[string]*graph.Edge)
	}

	for setsCount > 0 {
		e := heap.Pop(edges).(*graph.Edge)
		from := e.From
		to := e.To

		if from.Degree == 2 || to.Degree == 2 {
			continue
		}

		if setsCount > 1 && sets[from.ID] == sets[to.ID] {
			continue
		}

		connectedEdges[e.Hash()] = e

		e.From.Degree++
		e.To.Degree++

		e.From.Connected = append(e.From.Connected, to)
		e.To.Connected = append(e.To.Connected, from)

		mergedSet := append(*sets[from.ID], *sets[to.ID]...)

		for _, node := range mergedSet {
			sets[node.ID] = &mergedSet
		}

		setsCount--
	}

	return
}

func randomConnect(nodes []*graph.Node) (sets map[int]*[]*graph.Node, setsCount int, connectedEdges map[string]*graph.Edge) {
	sets = make(map[int]*[]*graph.Node)
	setsCount = graph.GetNodesCount()

	connectedEdges = make(map[string]*graph.Edge)

	for i := 0; i < setsCount; i++ {
		node := nodes[i]
		sets[node.ID] = &[]*graph.Node{node}
	}

	unvisitedNodes := make(map[int]*graph.Node)
	for _, node := range nodes {
		unvisitedNodes[node.ID] = node
	}
	var node *graph.Node

	// Randomly choose 1/4 of edges
	for setsCount > int(3*float64(graph.GetNodesCount())/4) {
		// iteration 1

		node = randomUnvisitedNode(unvisitedNodes)
		delete(unvisitedNodes, node.ID)

		nearestUnvisited, secondNearest := findNearUnvisited(node, unvisitedNodes)

		randomChoose := func(n1 *graph.Node, n2 *graph.Node) *graph.Node {
			r := rand.Float64()

			if r < 0.33 {
				return n2
			}
			return n1
		}

		neighbor := randomChoose(nearestUnvisited, secondNearest)
		delete(unvisitedNodes, neighbor.ID)

		// fmt.Printf("Randomly connected: (%d, %d)\n", node.ID, neighbor.ID)

		node.Degree++
		node.Connected = append(node.Connected, neighbor)

		neighbor.Degree++
		neighbor.Connected = append(neighbor.Connected, node)

		mergedSet := append(*sets[node.ID], *sets[neighbor.ID]...)

		for _, node := range mergedSet {
			sets[node.ID] = &mergedSet
		}

		setsCount--

		e := graph.NewEdge(node, neighbor)
		connectedEdges[e.Hash()] = e
	}

	return
}

func randomUnvisitedNode(unvisitedNodes map[int]*graph.Node) *graph.Node {
	i := rand.Intn(len(unvisitedNodes))

	for _, node := range unvisitedNodes {
		if i == 0 {
			return node
		}
		i--
	}

	log.Fatalln("randomUnvisitedNode")
	return nil
}

func findNearUnvisited(node *graph.Node, unvisitedNodes map[int]*graph.Node) (nearest *graph.Node, secondNearest *graph.Node) {
	D := graph.GetDistance

	minDistance := math.MaxFloat64
	secondMinDistance := minDistance

	for _, other := range unvisitedNodes {
		if D(node, other) < minDistance {
			secondNearest = nearest
			secondMinDistance = minDistance

			nearest = other
			minDistance = D(node, other)

		} else if D(node, other) < secondMinDistance {
			secondNearest = other
			secondMinDistance = D(node, other)
		}
	}

	return
}
