package graph

import "strconv"

type Edge struct {
	From     *Node
	To       *Node
	Distance float64
}

func NewEdge(from *Node, to *Node) *Edge {
	return &Edge{
		From: from,
		To:   to,
	}
}

func (e *Edge) Hash() string {
	x := e.From.ID
	y := e.To.ID
	if e.From.ID > e.To.ID {
		x, y = y, x
	}
	return strconv.Itoa(x) + "#" + strconv.Itoa(y)
}
