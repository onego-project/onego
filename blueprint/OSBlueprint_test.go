package blueprint

import (
	"github.com/onego-project/onego/resources"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("OSBlueprint", func() {
	var blueprint *OSBlueprint

	ginkgo.Describe("CreateOSBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateOSBlueprint()
		})

		ginkgo.It("should create a blueprint with OS element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("OS"))
		})
	})

	ginkgo.Describe("SetArchitecture", func() {
		var value resources.ArchitectureType

		ginkgo.BeforeEach(func() {
			blueprint = &OSBlueprint{Blueprint: *CreateBlueprint("OS")}
			value = resources.ArchitectureTypeX86_64
		})

		ginkgo.It("should set ARCH tag to specified value", func() {
			blueprint.SetArchitecture(value)

			gomega.Expect(blueprint.XMLData.FindElement("OS/ARCH").Text()).To(gomega.Equal(resources.ArchitectureTypeMap[value]))
		})
	})
})
