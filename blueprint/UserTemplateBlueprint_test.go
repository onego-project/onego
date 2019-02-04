package blueprint

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("UserTemplateBlueprint", func() {
	var blueprint *UserTemplateBlueprint

	ginkgo.Describe("CreateUpdateUserTemplateBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateUpdateUserTemplateBlueprint()
		})

		ginkgo.It("should create a blueprint with TEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("TEMPLATE"))
		})
	})

	ginkgo.Describe("SetDescription", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &UserTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = "test-value"
		})

		ginkgo.It("should set DESCRIPTION tag to specified value", func() {
			blueprint.SetDescription(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DESCRIPTION").Text()).To(gomega.Equal(value))
		})
	})
})
