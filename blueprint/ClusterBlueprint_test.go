package blueprint

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("ClusterBlueprint", func() {
	var blueprint *ClusterBlueprint

	ginkgo.Describe("CreateUpdateClusterBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateUpdateClusterBlueprint()
		})

		ginkgo.It("should create a blueprint with TEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("TEMPLATE"))
		})
	})

	ginkgo.Describe("SetReservedCPU", func() {
		var reservedCPU string

		ginkgo.BeforeEach(func() {
			blueprint = &ClusterBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			reservedCPU = "12345"
		})

		ginkgo.It("should set RESERVED_CPU tag to specified value", func() {
			blueprint.SetReservedCPU(reservedCPU)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/RESERVED_CPU").Text()).To(
				gomega.Equal(reservedCPU))
		})
	})

	ginkgo.Describe("SetReservedMemory", func() {
		var reservedMemory string

		ginkgo.BeforeEach(func() {
			blueprint = &ClusterBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			reservedMemory = "54321"
		})

		ginkgo.It("should set RESERVED_MEM tag to specified value", func() {
			blueprint.SetReservedMemory(reservedMemory)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/RESERVED_MEM").Text()).To(
				gomega.Equal(reservedMemory))
		})
	})
})
