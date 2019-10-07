package solver

import (
	"container/heap"
	"fmt"

	"github.com/greenmonn/tsp-go/graph"
)

func SolveGreedy() (tour *graph.Tour) {
	fmt.Println("Start Solving Greedy")

	N := graph.GetNodesCount()
	D := graph.GetDistanceByIndex

	nodes := make([]*graph.Node, N)
	for i := 0; i < N; i++ {
		id := graph.GetNode(i).ID
		x := graph.GetNode(i).X
		y := graph.GetNode(i).Y
		nodes[i] = graph.NewNode(id, x, y)
	}

	edges := &priorityQueue{}
	heap.Init(edges)

	for i := 0; i < N; i++ {
		for j := 0; j < i; j++ {
			edge := &graph.Edge{From: nodes[i], To: nodes[j], Distance: D(i, j)}
			heap.Push(edges, edge)
		}
	}

	fmt.Println("1: Make Edges Priority Queue")

	_ = connect(edges, nodes, N)

	fmt.Println("2: Connect possible short edges")

	tour = graph.NewTour()
	tour.FromNodes(nodes)

	fmt.Println("3: Construct tour from connected nodes")

	return

}

func connect(edges *priorityQueue, nodes []*graph.Node, N int) (distance float64) {
	sets := make(map[int]*[]*graph.Node)
	setsCount := N

	distance = 0

	for i := 0; i < N; i++ {
		node := nodes[i]
		sets[node.ID] = &[]*graph.Node{node}
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

		distance += e.Distance

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
