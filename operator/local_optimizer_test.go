package operator

import (
	"sort"

	"github.com/greenmonn/tsp-go/graph"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	N int
)

const filename = "fl1400"

var _ = Describe("Optimize", func() {
	BeforeEach(func() {
		graph.SetGraphFromFile("problems/" + filename + ".tsp")

		N = graph.GetNodesCount()

		graph.SetNearestNeighbors(5)
	})

	Describe("FastLocalSearchOptimize", func() {
		It("optimizes to the better neighbor", func() {
			parent1 := PartialRandomGreedy()
			parent2 := PartialRandomGreedy()

			tours := GXCrossover(parent1, parent2, 1.0, 0.25, 0.75)
			tour := tours[0]

			FastLocalSearchOptimize(tour)
			tour.FromNodes(tour.Path)

			idPath := graph.PathToIDs(tour.Path)
			sort.Ints(idPath)
			Expect(idPath).To(Equal(makeRange(1, 1400)))

		})
	})

	Describe("find2OptBetterMoveFromEdges", func() {
		FIt("optimizes to the better neighbor", func() {
			parent1 := PartialRandomGreedy()
			parent2 := PartialRandomGreedy()

			tours := GXCrossover(parent1, parent2, 1.0, 0.25, 0.75)
			tour := tours[0]
			EdgeExchangeMutateForGX(tour, 1.0)

			found := true
			for found == true {
				found = find2OptBetterMoveFromEdges(tour)
				tour.FromNodes(tour.Path)

				idPath := graph.PathToIDs(tour.Path)
				sort.Ints(idPath)
				Expect(idPath).To(Equal(makeRange(1, 1400)))
			}

		})
	})

	Describe("findOrderedVertices", func() {
		It("finds two vertices of given edge with same direction of given nodes", func() {
			tour := PartialRandomGreedy()

			edges := make([]*graph.Edge, N)
			index := 0
			for _, e := range tour.Edges {
				edges[index] = e
				index++
			}

			for i := 0; i < N; i++ {
				for j := i + 1; j < N; j++ {
					e1 := edges[i]
					e2 := edges[j]

					a, b := e1.From, e2.To

					c, d := findOrderedVertices(a, b, e2)

					Expect(c.ID).Should(Or(Equal(e2.From.ID), Equal(e2.To.ID)))

					Expect(d.ID).Should(Or(Equal(e2.From.ID), Equal(e2.To.ID)))

					Expect(c.ID).NotTo(Equal(d.ID))

				}
			}

		})
	})
})
