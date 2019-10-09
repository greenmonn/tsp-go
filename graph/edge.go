package graph

import "strconv"

type Edge struct {
	From     *Node
	To       *Node
	Distance float64
}

func NewEdge(from *Node, to *Node) *Edge {
	return &Edge{
		From:     from,
		To:       to,
		Distance: 0.,
	}
}

func NewEdges(nodes []*Node) map[string]*Edge {
	N := GetNodesCount()
	D := GetDistanceByIndex
	edges := make(map[string]*Edge)

	for i := 0; i < N; i++ {
		for j := 0; j < i; j++ {
			edge := &Edge{From: nodes[i], To: nodes[j], Distance: D(i, j)}

			edges[edge.Hash()] = edge
		}
	}
	return edges

}

func (e *Edge) Hash() string {
	x := e.From.ID
	y := e.To.ID

	return EdgeID(x, y)
}

func EdgeID(x int, y int) string {
	if x > y {
		x, y = y, x
	}
	return strconv.Itoa(x) + "#" + strconv.Itoa(y)
}

func (e *Edge) UpdateNodes() {
	e.From.Degree++
	e.To.Degree++

	e.From.Connected = append(e.From.Connected, e.To)
	e.To.Connected = append(e.To.Connected, e.From)
}
