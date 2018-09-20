package blueprint

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GroupBlueprint", func() {
	var (
		blueprint *GroupBlueprint
	)

	ginkgo.Describe("CreateUpdateGroupBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateUpdateGroupBlueprint()
		})

		ginkgo.It("should create a blueprint with TEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("TEMPLATE"))
		})
	})
})
