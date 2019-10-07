package graph

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Graph", func() {
	Describe("NewGraph", func() {
		It("creates graph from nodes", func() {
			SetGraphFromFile("problems/burma14.tsp")

			Expect(graph.N).To(Equal(14))
			Expect(len(graph.Nodes)).To(Equal(14))

			Expect(GetDistance(graph.Nodes[1], graph.Nodes[2])).To(Equal(graph.Nodes[1].Distance(graph.Nodes[2])))

		})
	})
})
