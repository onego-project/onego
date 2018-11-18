package blueprint

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("VMTemplateBlueprint", func() {
	var blueprint *VMTemplateBlueprint

	ginkgo.Describe("CreateAllocateVMTemplateBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateAllocateVMTemplateBlueprint()
		})

		ginkgo.It("should create a blueprint with VMTEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("VMTEMPLATE"))
		})
	})

	ginkgo.Describe("CreateUpdateVMTemplateBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateUpdateVMTemplateBlueprint()
		})

		ginkgo.It("should create a blueprint with TEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("TEMPLATE"))
		})
	})

	ginkgo.Describe("SetCPU", func() {
		var element string

		ginkgo.BeforeEach(func() {
			blueprint = &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			element = "xyz"
		})

		ginkgo.It("should set a CPU tag to specified value", func() {
			blueprint.SetCPU(element)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/CPU").Text()).To(gomega.Equal(element))
		})
	})

	ginkgo.Describe("SetDescription", func() {
		var element string

		ginkgo.BeforeEach(func() {
			blueprint = &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			element = "xyz"
		})

		ginkgo.It("should set a DESCRIPTION tag to specified value", func() {
			blueprint.SetDescription(element)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DESCRIPTION").Text()).To(gomega.Equal(element))
		})
	})

	ginkgo.Describe("SetLogo", func() {
		var element string

		ginkgo.BeforeEach(func() {
			blueprint = &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			element = "xyz"
		})

		ginkgo.It("should set a LOGO tag to specified value", func() {
			blueprint.SetLogo(element)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/LOGO").Text()).To(gomega.Equal(element))
		})
	})

	ginkgo.Describe("SetMemory", func() {
		var element string

		ginkgo.BeforeEach(func() {
			blueprint = &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			element = "xyz"
		})

		ginkgo.It("should set a MEMORY tag to specified value", func() {
			blueprint.SetMemory(element)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/MEMORY").Text()).To(gomega.Equal(element))
		})
	})

	ginkgo.Describe("SetSchedRequirements", func() {
		var element string

		ginkgo.BeforeEach(func() {
			blueprint = &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			element = "xyz"
		})

		ginkgo.It("should set a SCHED_REQUIREMENTS tag to specified value", func() {
			blueprint.SetSchedRequirements(element)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/SCHED_REQUIREMENTS").Text()).To(
				gomega.Equal(element))
		})
	})

	ginkgo.Describe("SetVCPU", func() {
		var element string

		ginkgo.BeforeEach(func() {
			blueprint = &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			element = "xyz"
		})

		ginkgo.It("should set a VCPU tag to specified value", func() {
			blueprint.SetVCPU(element)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/VCPU").Text()).To(gomega.Equal(element))
		})
	})
})
