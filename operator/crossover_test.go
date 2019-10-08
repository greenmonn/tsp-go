package operator

import (
	"fmt"

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

			children := OrderCrossover(p1, p2)

			for _, c := range children {
				Expect(len(c.Path)).To(Equal(N))
			}

			fmt.Println("Finish")
		})
	})

	Describe("EdgeRecombinationCrossover", func() {
		It("crossover two parent tours", func() {
			p1 := graph.NewRandomTour()
			p2 := graph.NewRandomTour()

			fmt.Println("Parent 1 Distance: ", p1.Distance)
			fmt.Println("Parent 2 Distance: ", p2.Distance)

			children := EdgeRecombinationCrossover(p1, p2)

			for _, c := range children {
				Expect(len(c.Path)).To(Equal(N))
				fmt.Println("Distance: ", c.Distance)
			}

		})
	})
})
