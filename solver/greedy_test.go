package solver

import (
	"log"
	"sort"

	"github.com/greenmonn/tsp-go/graph"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thoas/go-funk"
)

var _ = Describe("Greedy", func() {
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

	Describe("SolveGreedy", func() {
		It("returns greedy solution", func() {
			tour := SolveGreedy()

			idPath := funk.Map(tour.Path, func(node *graph.Node) int {
				return node.ID
			}).([]int)

			Expect(len(idPath)).To(Equal(N))

			sort.Ints(idPath)
			Expect(idPath).To(Equal(makeRange(1, N)))

			log.Println("Distance: ", tour.Distance)

			n := graph.WritePathToFile(tour.Path, filename)

			log.Printf("%d Bytes Wrote\n", n)
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
