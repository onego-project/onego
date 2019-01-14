package blueprint

import (
	"strconv"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("VirtualMachineBlueprint", func() {
	var blueprint *VirtualMachineBlueprint

	ginkgo.Describe("CreateUpdateVirtualMachineBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateUpdateVirtualMachineBlueprint()
		})

		ginkgo.It("should create a blueprint with TEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("TEMPLATE"))
		})
	})

	ginkgo.Describe("CreateAllocateVirtualMachineBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateAllocateVirtualMachineBlueprint()
		})

		ginkgo.It("should create a blueprint with VM element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("VM"))
		})
	})

	ginkgo.Describe("SetAutomaticDSRequirements", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = "test-value"
		})

		ginkgo.It("should set AUTOMATIC_DS_REQUIREMENTS tag to specified value", func() {
			blueprint.SetAutomaticDSRequirements(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/AUTOMATIC_DS_REQUIREMENTS").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetAutomaticRequirements", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = "test-value"
		})

		ginkgo.It("should set AUTOMATIC_REQUIREMENTS tag to specified value", func() {
			blueprint.SetAutomaticRequirements(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/AUTOMATIC_REQUIREMENTS").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetContext", func() {
		var value *ContextBlueprint

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = CreateContextBlueprint()
		})

		ginkgo.It("should set CONTEXT tag to specified value", func() {
			blueprint.SetContext(*value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/CONTEXT").Text()).To(gomega.Equal(value.XMLData.Tag))
		})
	})

	ginkgo.Describe("SetCPU", func() {
		var value float64

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = 0.55
		})

		ginkgo.It("should set CPU tag to specified value", func() {
			blueprint.SetCPU(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/CPU").Text()).To(gomega.Equal(
				strconv.FormatFloat(value, 'f', -1, 64)))
		})
	})

	ginkgo.Describe("SetMemory", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = 2048
		})

		ginkgo.It("should set MEMORY tag to specified value", func() {
			blueprint.SetMemory(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("TEMPLATE/MEMORY").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetDisk", func() {
		var value *DiskBlueprint

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = CreateDiskBlueprint()
		})

		ginkgo.It("should set DISK tag to specified value", func() {
			blueprint.SetDisk(*value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DISK").Text()).To(gomega.Equal(value.XMLData.Tag))
		})
	})

	ginkgo.Describe("SetFeatures", func() {
		var value *FeaturesBlueprint

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = CreateFeaturesBlueprint()
		})

		ginkgo.It("should set FEATURES tag to specified value", func() {
			blueprint.SetFeatures(*value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/FEATURES").Text()).To(gomega.Equal(value.XMLData.Tag))
		})
	})

	ginkgo.Describe("SetGraphics", func() {
		var value *GraphicsBlueprint

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = CreateGraphicsBlueprint()
		})

		ginkgo.It("should set GRAPHICS tag to specified value", func() {
			blueprint.SetGraphics(*value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/GRAPHICS").Text()).To(gomega.Equal(value.XMLData.Tag))
		})
	})

	ginkgo.Describe("SetNIC", func() {
		var value *NICBlueprint

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = CreateNICBlueprint()
		})

		ginkgo.It("should set NIC tag to specified value", func() {
			blueprint.SetNIC(*value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/NIC").Text()).To(gomega.Equal(value.XMLData.Tag))
		})
	})

	ginkgo.Describe("SetOS", func() {
		var value *OSBlueprint

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = CreateOSBlueprint()
		})

		ginkgo.It("should set OS tag to specified value", func() {
			blueprint.SetOS(*value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/OS").Text()).To(gomega.Equal(value.XMLData.Tag))
		})
	})

	ginkgo.Describe("SetRaw", func() {
		var value *RawBlueprint

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = CreateRawBlueprint()
		})

		ginkgo.It("should set RAW tag to specified value", func() {
			blueprint.SetRaw(*value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/RAW").Text()).To(gomega.Equal(value.XMLData.Tag))
		})
	})

	ginkgo.Describe("SetTemplateID", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = 2048
		})

		ginkgo.It("should set TEMPLATE_ID tag to specified value", func() {
			blueprint.SetTemplateID(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("TEMPLATE/TEMPLATE_ID").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetVCPU", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = 2048
		})

		ginkgo.It("should set VCPU tag to specified value", func() {
			blueprint.SetVCPU(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("TEMPLATE/VCPU").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})
})
