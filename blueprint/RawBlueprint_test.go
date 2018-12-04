package blueprint

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("RawBlueprint", func() {
	var blueprint *RawBlueprint

	ginkgo.Describe("CreateRawBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateRawBlueprint()
		})

		ginkgo.It("should create a blueprint with RAW element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("RAW"))
		})
	})

	ginkgo.Describe("SetData", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &RawBlueprint{Blueprint: *CreateBlueprint("RAW")}
			value = "test-value"
		})

		ginkgo.It("should set DATA tag to specified value", func() {
			blueprint.SetData(value)

			gomega.Expect(blueprint.XMLData.FindElement("RAW/DATA").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetType", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &RawBlueprint{Blueprint: *CreateBlueprint("RAW")}
			value = "test-value"
		})

		ginkgo.It("should set TYPE tag to specified value", func() {
			blueprint.SetType(value)

			gomega.Expect(blueprint.XMLData.FindElement("RAW/TYPE").Text()).To(gomega.Equal(value))
		})
	})
})
