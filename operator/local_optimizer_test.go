package operator

import (
	"github.com/greenmonn/tsp-go/graph"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Optimize", func() {
	BeforeEach(func() {
		N := 5
		distances := make([][]float64, N)
		for i := 0; i < N; i++ {
			distances[i] = make([]float64, i)
			for j := 0; j < i; j++ {
				distances[i][j] = 1.0
			}
		}

		graph.SetGraphByMatrix(distances)
	})
	Describe("SwapTwoEdges", func() {
		It("swap two edgese", func() {
			tour := graph.NewTour()

			tour.FromPath(graph.IDsToPath([]int{1, 3, 2, 4, 5}))

			SwapTwoEdges(tour, 0, 3, false)

			Expect(graph.PathToIDs(tour.Path)).To(Equal([]int{1, 5, 3, 2, 4}))

		})
	})
})
