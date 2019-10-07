package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileUtil", func() {
	Describe("parseLine", func() {
		It("parse file's each line to id, positions", func() {
			line := "   1  16.47       96.10"

			id, posX, posY := ParseLine(line)

			Expect(id).To(Equal(1))
			Expect(posX).To(Equal(16.47))
			Expect(posY).To(Equal(96.10))
		})
		Context("without node's information", func() {
			It("returns id -1", func() {
				line := "EDGE_WEIGHT_TYPE: GEO"

				id, _, _ := ParseLine(line)

				Expect(id).To(Equal(-1))
			})
		})

	})

})
