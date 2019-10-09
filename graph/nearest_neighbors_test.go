package graph

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NearestNeighbors", func() {
	var (
		N int
	)

	const (
		filename = "burma14"
	)

	BeforeEach(func() {
		SetGraphFromFile("problems/" + filename + ".tsp")

		N = GetNodesCount()
	})

	Describe("SetNearestNeighbors", func() {
		It("saves nearest neighbors of each node", func() {
			K := 5
			SetNearestNeighbors(K)

			Expect(len(NearestNeighbors(N))).To(Equal(K))
		})
	})

	Describe("NearestNeighbors", func() {
		It("returns nearest neighbors of a node with given ID", func() {
			K := 5
			SetNearestNeighbors(K)

			for id := 1; id <= N; id++ {
				Expect(len(NearestNeighbors(id))).To(Equal(K))
			}

		})
	})

	Describe("findNearest", func() {
		It("returns K nearest neighbors of nodes[i]", func() {
			K := 5
			nodes := graph.Nodes

			neighbors := findNearests(nodes, K, 0)

			Expect(PathToIDs(neighbors)).To(Equal([]int{1, 8, 11, 9, 2}))
		})
	})
})
