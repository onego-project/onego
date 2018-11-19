package blueprint

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("VMTemplateBlueprint", func() {
	var blueprint *VMTemplateBlueprint
	value := "thisIsTheValue"

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

	ginkgo.Describe("SetDisk", func() {
		var disk *VMTemplateDiskBlueprint

		ginkgo.BeforeEach(func() {
			disk = CreateVMTemplateDiskBlueprint()
			disk.SetImage(value)
			disk.SetDevicePrefix("123456")

			blueprint = &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
		})

		ginkgo.It("should set a DISK tag to specified value", func() {
			blueprint.SetDisk(*disk)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DISK/IMAGE").Text()).To(gomega.Equal(value))
			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DISK/DEV_PREFIX").Text()).To(gomega.Equal("123456"))
		})
	})

	ginkgo.Describe("SetFeatures", func() {
		var features *VMTemplateFeaturesBlueprint

		ginkgo.BeforeEach(func() {
			features = CreateVMTemplateFeaturesBlueprint()
			features.SetGuestAgent(value)

			blueprint = &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
		})

		ginkgo.It("should set a FEATURES tag to specified value", func() {
			blueprint.SetFeatures(*features)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/FEATURES/GUEST_AGENT").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetGraphics", func() {
		var graphics *VMTemplateGraphicsBlueprint

		ginkgo.BeforeEach(func() {
			graphics = CreateVMTemplateGraphicsBlueprint()
			graphics.SetListen(value)
			graphics.SetType("asdf")

			blueprint = &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
		})

		ginkgo.It("should set a GRAPHICS tag to specified value", func() {
			blueprint.SetGraphics(*graphics)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/GRAPHICS/LISTEN").Text()).To(gomega.Equal(value))
			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/GRAPHICS/TYPE").Text()).To(gomega.Equal("asdf"))
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

	ginkgo.Describe("SetNIC", func() {
		var nic *VMTemplateNICBlueprint

		ginkgo.BeforeEach(func() {
			nic = CreateVMTemplateNICBlueprint()
			nic.SetNetwork(value)

			blueprint = &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
		})

		ginkgo.It("should set a NIC tag to specified value", func() {
			blueprint.SetNIC(*nic)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/NIC/NETWORK").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetOS", func() {
		var os *VMTemplateOSBlueprint

		ginkgo.BeforeEach(func() {
			os = CreateVMTemplateOSBlueprint()
			os.SetArch(value)

			blueprint = &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
		})

		ginkgo.It("should set a OS tag to specified value", func() {
			blueprint.SetOS(*os)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/OS/ARCH").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetRAW", func() {
		var raw *VMTemplateRAWBlueprint

		ginkgo.BeforeEach(func() {
			raw = CreateVMTemplateRAWBlueprint()
			raw.SetData(value)
			raw.SetType("21258")

			blueprint = &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
		})

		ginkgo.It("should set a RAW tag to specified value", func() {
			blueprint.SetRAW(*raw)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/RAW/DATA").Text()).To(gomega.Equal(value))
			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/RAW/TYPE").Text()).To(gomega.Equal("21258"))
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
