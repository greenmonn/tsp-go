package solver

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/greenmonn/tsp-go/graph"
)

var _ = Describe("GA", func() {
	var (
		N int
	)

	const (
		filename = "burma14"
	)

	BeforeEach(func() {
		graph.SetGraphFromFile("problems/" + filename + ".tsp")

		N = graph.GetNodesCount()
	})

	Describe("bestTour", func() {
		It("returns best tour among population", func() {
			initialTours := make([]*graph.Tour, 10)

			for i := 0; i < 10; i++ {
				idPath := makeRange(1, N)
				path := graph.IDsToPath(idPath)

				tour := graph.NewTour()
				tour.FromPath(path)

				initialTours[i] = tour
			}
			population := NewPopulation(10, initialTours)

			Expect(population.BestTour().Distance).To(Equal(initialTours[0].Distance))
		})
	})

	Describe("EvolvePopulation", func() {
		It("returns evolved population", func() {
			initialTours := make([]*graph.Tour, 10)

			for i := 0; i < 10; i++ {
				idPath := makeRange(1, N)
				path := graph.IDsToPath(idPath)

				tour := graph.NewTour()
				tour.FromPath(path)

				initialTours[i] = tour
			}
			population := NewPopulation(10, initialTours)
			nextPopulation := EvolvePopulation(population)

			Expect(len(population.Tours)).To(Equal(len(nextPopulation.Tours)))

			Expect(population.BestTour().Distance >= nextPopulation.BestTour().Distance)
		})
	})

	Describe("SolveGA", func() {
		Context("with list of tours as initial population", func() {

			It("returns optimized tour", func() {
				initialTours := make([]*graph.Tour, 10)

				for i := 0; i < 10; i++ {
					idPath := makeRange(1, N)
					path := graph.IDsToPath(idPath)

					tour := graph.NewTour()
					tour.FromPath(path)

					initialTours[i] = tour
				}
				optimizedTour := SolveGA(initialTours, 10, 10)

				Expect(len(optimizedTour.Path)).To(Equal(N))

				Expect(optimizedTour.Path[0]).NotTo(BeNil())

				idPath := graph.PathToIDs(optimizedTour.Path)

				sort.Ints(idPath)
				Expect(idPath).To(Equal(makeRange(1, N)))
			})
		})

		Context("with empty list", func() {
			It("returns optimized tour from a random population", func() {
				optimizedTour := SolveGA([]*graph.Tour{}, 10, 10)

				Expect(len(optimizedTour.Path)).To(Equal(N))

				idPath := graph.PathToIDs(optimizedTour.Path)

				sort.Ints(idPath)
				Expect(idPath).To(Equal(makeRange(1, N)))
			})
		})
	})
})
