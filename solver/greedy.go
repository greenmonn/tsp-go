package solver

import (
	"container/heap"
	"log"

	"github.com/greenmonn/tsp-go/container"
	"github.com/greenmonn/tsp-go/graph"
	"github.com/greenmonn/tsp-go/operator"
)

func SolveGreedy() (tour *graph.Tour) {
	log.Println("Start Solving Greedy")

	nodes := graph.CopyNodesFromGraph()

	edges := &container.PriorityQueue{}
	heap.Init(edges)

	container.InitEdges(edges, nodes)

	log.Println("1: Make Edges Priority Queue")

	operator.GreedyConnect(edges, nodes, nil, -1, nil)

	log.Println("2: Connect possible short edges")

	tour = graph.NewTour()
	tour.FromNodes(nodes)

	log.Println("3: Construct tour from connected nodes")

	return

}
