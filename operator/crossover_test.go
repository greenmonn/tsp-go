package operator

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/greenmonn/tsp-go/graph"
)

var _ = Describe("Crossover", func() {
	var (
		N int
	)

	BeforeEach(func() {
		N = 5
		distances := make([][]float64, N)
		for i := 0; i < N; i++ {
			distances[i] = make([]float64, i)
			for j := 0; j < i; j++ {
				distances[i][j] = 1.0
			}
		}

		graph.SetGraphByMatrix(distances)
	})

	Describe("OrderCrossover", func() {
		It("crossover two parent tours", func() {
			p1 := graph.NewRandomTour()
			p2 := graph.NewRandomTour()

			// fmt.Println("parent1: ", graph.PathToIDs(p1.Path))
			// fmt.Println("parent2: ", graph.PathToIDs(p2.Path))

			c1, c2 := OrderCrossover(p1, p2)

			Expect(len(c1.Path)).To(Equal(N))
			Expect(len(c2.Path)).To(Equal(N))

			// fmt.Println("child1: ", graph.PathToIDs(c1.Path))
			// fmt.Println("child2: ", graph.PathToIDs(c2.Path))
		})
	})
})
