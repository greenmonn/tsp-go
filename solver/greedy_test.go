package solver

import (
	"container/heap"
	"fmt"
	"math"
	"sort"

	"github.com/greenmonn/tsp-go/graph"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thoas/go-funk"
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

	Describe("SolveGreedy", func() {
		It("returns greedy solution", func() {
			tour := SolveGreedy()

			idPath := funk.Map(tour.Path, func(node *graph.Node) int {
				return node.ID
			}).([]int)

			Expect(len(idPath)).To(Equal(N))

			sort.Ints(idPath)
			Expect(idPath).To(Equal(makeRange(1, N)))

			fmt.Println("Distance: ", tour.Distance)

			n := graph.WritePathToFile(tour.Path, filename)

			fmt.Printf("%d Bytes Wrote\n", n)
		})
	})

	Describe("connect", func() {
		It("connects nodes in greedy way", func() {
			edges := &priorityQueue{}
			heap.Init(edges)

			for i := 0; i < N; i++ {
				for j := 0; j < i; j++ {
					edge := &graph.Edge{From: nodes[i], To: nodes[j],
						Distance: D(i, j)}
					heap.Push(edges, edge)
				}
			}

			edgesSum := connect(edges, nodes, N)

			tour := graph.NewTour()
			tour.FromNodes(nodes)

			Expect(math.Abs(tour.Distance-edgesSum) <= 1)
		})
	})
})

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
