package solver

import "github.com/greenmonn/tsp-go/graph"

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
