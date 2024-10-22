package graph

import (
	"math/rand"

	"github.com/thoas/go-funk"
)

type Tour struct {
	Path      []*Node
	Distance  float64
	Edges     map[string]*Edge
	FlexEdges []*Edge
}

func NewTour() *Tour {
	return &Tour{
		Path:     make([]*Node, GetNodesCount()),
		Distance: 0.0,
	}
}

func NewRandomTour() *Tour {
	N := GetNodesCount()

	ids := makeRange(1, N)
	rand.Shuffle(N, func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })

	t := &Tour{Path: IDsToPath(ids)}

	t.UpdateDistance()

	return t
}

func (t *Tour) FromNodes(connectedNodes []*Node) {
	if len(connectedNodes) == 0 {
		return
	}

	// log.Println("Constructing path from connected nodes. . .")

	N := GetNodesCount()
	var prevNode *Node
	node := connectedNodes[0]

	t.Path[0] = node
	index := 1

	t.Distance = 0.0

	for index < N {
		for _, nextNode := range node.Connected {
			if prevNode != nil && nextNode.ID == prevNode.ID {
				continue
			}

			t.Distance += GetDistance(node, nextNode)

			t.Path[index] = nextNode
			index++

			prevNode = node
			node = nextNode

			break
		}
	}

	t.Distance += GetDistance(t.Path[len(t.Path)-1], t.Path[0])
}

func (t *Tour) FromPath(path []*Node) {
	t.Path = path
	t.UpdateConnections()
	t.UpdateDistance()
}

func (t *Tour) UpdateDistance() {
	node := t.Path[0]

	t.Distance = 0.0

	// Calculate distance from path
	for _, nextNode := range t.Path[1:] {
		t.Distance += GetDistance(node, nextNode)
		node = nextNode
	}

	t.Distance += GetDistance(t.Path[len(t.Path)-1], t.Path[0])
}

func (t *Tour) UpdateConnections() {
	// Update each node's connection information
	for i := 0; i < GetNodesCount(); i++ {
		t.Path[i].Connected =
			[]*Node{t.GetNode(i - 1), t.GetNode(i + 1)}
	}
}

func (t *Tour) UpdateEdges() {
	node := t.Path[0]

	t.Edges = make(map[string]*Edge)

	// Find edges from path
	for _, nextNode := range t.Path[1:] {
		edge := NewEdge(node, nextNode)
		t.Edges[edge.Hash()] = edge
		node = nextNode
	}

	edge := NewEdge(t.Path[len(t.Path)-1], t.Path[0])
	t.Edges[edge.Hash()] = edge
}

func (t *Tour) GetNode(index int) *Node {
	N := GetNodesCount()
	for index >= N {
		index -= N
	}

	for index < 0 {
		index += N
	}

	return t.Path[index]
}

func (t *Tour) SetNode(index int, node *Node) {
	N := GetNodesCount()
	for index >= N {
		index -= N
	}

	for index < 0 {
		index += N
	}

	t.Path[index] = node
}

func (t *Tour) WritePathToFile(name string) int {
	return WritePathToFile(t.Path, name)
}

func IDsToPath(idList []int) []*Node {
	return funk.Map(idList, func(id int) *Node {
		return NewNode(id, 0., 0.)
	}).([]*Node)
}

func PathToIDs(path []*Node) []int {
	return funk.Map(path, func(node *Node) int {
		if node == nil {
			return 0
		}
		return node.ID
	}).([]int)
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func areEqualIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
