package graph

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
	if e.From.ID > e.To.ID {
		x, y = y, x
	}
	return string(x) + "#" + string(y)
}

func EdgeID(id1 int, id2 int) string {
	if id1 > id2 {
		id1, id2 = id2, id1
	}
	return string(id1) + "#" + string(id2)
}

func (e *Edge) UpdateNodes() {
	e.From.Degree++
	e.To.Degree++

	e.From.Connected = append(e.From.Connected, e.To)
	e.To.Connected = append(e.To.Connected, e.From)
}
