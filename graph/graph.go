package graph

type Graph struct {
	N                int
	Nodes            []*Node
	Distances        [][]float64
	NearestNeighbors map[int][]*Node
}

var (
	graph *Graph // Singleton
)

func SetGraphFromFile(filepath string) {
	SetGraph(ReadNodesFromFile(filepath))
}

func SetGraphByMatrix(distances [][]float64) {
	N := len(distances)
	nodes := make([]*Node, N)
	for i := 0; i < N; i++ {
		nodes[i] = NewNode(i+1, 0., 0.)
	}
	graph = &Graph{N: N, Nodes: nodes}
	graph.Distances = distances
}

func SetGraph(nodes []*Node) {
	graph = &Graph{N: len(nodes), Nodes: nodes}
	graph.Distances = make([][]float64, graph.N)

	for i := 0; i < graph.N; i++ {
		graph.Distances[i] = make([]float64, graph.N)
		for j := 0; j < i; j++ {
			graph.Distances[i][j] = nodes[i].Distance(nodes[j])
		}
	}
}

func CopyNodesFromGraph() []*Node {
	N := GetNodesCount()

	nodes := make([]*Node, N)
	for i := 0; i < N; i++ {
		id := GetNode(i).ID
		x := GetNode(i).X
		y := GetNode(i).Y
		nodes[i] = NewNode(id, x, y)
	}

	return nodes
}

func GetNodes() []*Node {
	return graph.Nodes
}

func GetNode(i int) *Node {
	return graph.Nodes[i]
}

func GetNodesCount() int {
	return graph.N
}

func GetDistanceByIndex(i int, j int) float64 {
	if i < j {
		i, j = j, i
	}
	return graph.Distances[i][j]
}

func GetDistance(n1 *Node, n2 *Node) float64 {
	// Precondition: nodes[i] has ID=i+1

	i := n1.ID - 1
	j := n2.ID - 1

	if i < j {
		i, j = j, i
	}

	return graph.Distances[i][j]
}
