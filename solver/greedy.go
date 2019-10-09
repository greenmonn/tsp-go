package solver

import (
	"container/heap"
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/greenmonn/tsp-go/graph"
)

func SolveGreedy() (tour *graph.Tour) {
	fmt.Println("Start Solving Greedy")

	nodes := copyNodesFromGraph()

	edges := &priorityQueue{}
	heap.Init(edges)

	initEdges(edges, nodes)

	fmt.Println("1: Make Edges Priority Queue")

	connect(edges, nodes, nil, -1)

	fmt.Println("2: Connect possible short edges")

	tour = graph.NewTour()
	tour.FromNodes(nodes)

	fmt.Println("3: Construct tour from connected nodes")

	return

}

func PartialRandomGreedy() (tour *graph.Tour) {
	fmt.Println("Start Solving Greedy")

	nodes := copyNodesFromGraph()

	sets, setsCount := randomConnect(nodes)

	edges := &priorityQueue{}
	heap.Init(edges)

	initEdges(edges, nodes)

	fmt.Println("1: Make Edges Priority Queue")

	connect(edges, nodes, sets, setsCount)

	fmt.Println("2: Connect possible short edges")

	tour = graph.NewTour()
	tour.FromNodes(nodes)

	fmt.Println("3: Construct tour from connected nodes")

	return

}

func connect(edges *priorityQueue, nodes []*graph.Node, sets map[int]*[]*graph.Node, setsCount int) {
	if sets == nil {
		sets = make(map[int]*[]*graph.Node)
		setsCount = graph.GetNodesCount()

		for i := 0; i < setsCount; i++ {
			node := nodes[i]
			sets[node.ID] = &[]*graph.Node{node}
		}
	}

	for {
		e := heap.Pop(edges).(*graph.Edge)
		from := e.From
		to := e.To

		if from.Degree == 2 || to.Degree == 2 {
			continue
		}

		if setsCount > 1 && sets[from.ID] == sets[to.ID] {
			continue
		}

		e.From.Degree++
		e.To.Degree++

		e.From.Connected = append(e.From.Connected, to)
		e.To.Connected = append(e.To.Connected, from)

		if setsCount == 1 {
			break
		}

		mergedSet := append(*sets[from.ID], *sets[to.ID]...)

		for _, node := range mergedSet {
			sets[node.ID] = &mergedSet
		}

		setsCount--
	}

	return
}

func randomConnect(nodes []*graph.Node) (sets map[int]*[]*graph.Node, setsCount int) {
	sets = make(map[int]*[]*graph.Node)
	setsCount = graph.GetNodesCount()

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
	}

	return
}

func copyNodesFromGraph() []*graph.Node {
	N := graph.GetNodesCount()

	nodes := make([]*graph.Node, N)
	for i := 0; i < N; i++ {
		id := graph.GetNode(i).ID
		x := graph.GetNode(i).X
		y := graph.GetNode(i).Y
		nodes[i] = graph.NewNode(id, x, y)
	}

	return nodes
}

func initEdges(edges *priorityQueue, nodes []*graph.Node) {
	N := graph.GetNodesCount()
	D := graph.GetDistanceByIndex

	for i := 0; i < N; i++ {
		skipID := -1
		if nodes[i].Degree > 0 {
			skipID = nodes[i].Connected[0].ID

		}

		for j := 0; j < i; j++ {
			if skipID == nodes[j].ID {
				continue
			}
			edge := &graph.Edge{From: nodes[i], To: nodes[j], Distance: D(i, j)}
			heap.Push(edges, edge)
		}
	}
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
