package operator

import (
	"fmt"
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/greenmonn/tsp-go/graph"
)

var _ = Describe("Crossover", func() {
	var (
		N int
	)

	const (
		filename = "fl1400"
	)

	BeforeEach(func() {
		graph.SetGraphFromFile("problems/" + filename + ".tsp")

		N = graph.GetNodesCount()
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

	Describe("GXCrossover", func() {
		// Parents must have edges map
		BeforeEach(func() {
			graph.SetNearestNeighbors(5)
		})

		It("crossover two parent tours", func() {
			p1 := PartialRandomGreedy()
			p2 := PartialRandomGreedy()

			fmt.Println("Parent 1 Distance: ", p1.Distance)
			fmt.Println("Parent 2 Distance: ", p2.Distance)

			children := GXCrossover(p1, p2, 1., 0.25, 0.75)

			for _, c := range children {
				Expect(len(c.Path)).To(Equal(N))

				idPath := graph.PathToIDs(c.Path)
				sort.Ints(idPath)

				Expect(idPath).To(Equal(makeRange(1, N)))
				fmt.Println("Distance: ", c.Distance)
			}

		})
	})
})
