package solver

import (
	"container/heap"

	"github.com/greenmonn/tsp-go/graph"
)

type priorityQueue []*graph.Edge

func (q priorityQueue) Len() int { return len(q) }

func (q priorityQueue) Less(i, j int) bool { return q[i].Distance < q[j].Distance }

func (q priorityQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }

func (q *priorityQueue) Push(x interface{}) {
	*q = append(*q, x.(*graph.Edge))
}

func (q *priorityQueue) Pop() interface{} {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[0 : n-1]

	return x
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
