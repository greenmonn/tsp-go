package operator

import (
	"container/heap"

	"github.com/greenmonn/tsp-go/container"
	"github.com/greenmonn/tsp-go/graph"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Greedy", func() {
	var (
		nodes []*graph.Node
		N     int
		D     func(i, j int) float64
	)

	const (
		filename = "burma14"
	)

	BeforeEach(func() {
		graph.SetGraphFromFile("problems/" + filename + ".tsp")

		N = graph.GetNodesCount()
		D = graph.GetDistanceByIndex
		nodes = graph.GetNodes()
	})

	Describe("GreedyConnect", func() {
		It("connects nodes in greedy way", func() {
			edges := &container.PriorityQueue{}
			heap.Init(edges)

			for i := 0; i < N; i++ {
				for j := 0; j < i; j++ {
					edge := &graph.Edge{From: nodes[i], To: nodes[j],
						Distance: D(i, j)}
					heap.Push(edges, edge)
				}
			}

			GreedyConnect(edges, nodes, nil, -1, nil)

			tour := graph.NewTour()
			tour.FromNodes(nodes)

			Expect(tour.Distance - 36).Should(BeNumerically("<", 1.0))
		})
	})
})
