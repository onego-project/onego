package blueprint

import (
	"net"
	"strconv"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("TemplateBlueprint", func() {
	var blueprint *TemplateBlueprint
	value := "thisIsTheValue"

	ginkgo.Describe("CreateAllocateTemplateBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateAllocateTemplateBlueprint()
		})

		ginkgo.It("should create a blueprint with VMTEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("VMTEMPLATE"))
		})
	})

	ginkgo.Describe("CreateUpdateTemplateBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateUpdateTemplateBlueprint()
		})

		ginkgo.It("should create a blueprint with TEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("TEMPLATE"))
		})
	})

	ginkgo.Describe("SetCPU", func() {
		var element float64

		ginkgo.BeforeEach(func() {
			blueprint = &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			element = 456.78
		})

		ginkgo.It("should set a CPU tag to specified value", func() {
			blueprint.SetCPU(element)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/CPU").Text()).To(gomega.Equal(
				strconv.FormatFloat(element, 'f', -1, 64)))
		})
	})

	ginkgo.Describe("SetDescription", func() {
		var element string

		ginkgo.BeforeEach(func() {
			blueprint = &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			element = "xyz"
		})

		ginkgo.It("should set a DESCRIPTION tag to specified value", func() {
			blueprint.SetDescription(element)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DESCRIPTION").Text()).To(gomega.Equal(element))
		})
	})

	ginkgo.Describe("SetDisk", func() {
		var disk *DiskBlueprint

		ginkgo.BeforeEach(func() {
			disk = CreateDiskBlueprint()
			disk.SetImage(value)
			disk.SetDevicePrefix("123456")

			blueprint = &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
		})

		ginkgo.It("should set a DISK tag to specified value", func() {
			blueprint.SetDisk(*disk)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DISK/IMAGE").Text()).To(gomega.Equal(value))
			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DISK/DEV_PREFIX").Text()).To(gomega.Equal("123456"))
		})
	})

	ginkgo.Describe("SetFeatures", func() {
		var features *FeaturesBlueprint

		ginkgo.BeforeEach(func() {
			features = CreateFeaturesBlueprint()
			features.SetGuestAgent(true)

			blueprint = &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
		})

		ginkgo.It("should set a FEATURES tag to specified value", func() {
			blueprint.SetFeatures(*features)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/FEATURES/GUEST_AGENT").Text()).To(gomega.Equal("YES"))
		})
	})

	ginkgo.Describe("SetGraphics", func() {
		var graphics *GraphicsBlueprint
		var ip net.IP

		ginkgo.BeforeEach(func() {
			graphics = CreateGraphicsBlueprint()

			ip = net.ParseIP("10.0.0.10")
			graphics.SetListen(ip)

			graphics.SetType(VNC)

			blueprint = &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
		})

		ginkgo.It("should set a GRAPHICS tag to specified value", func() {
			blueprint.SetGraphics(*graphics)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/GRAPHICS/LISTEN").Text()).To(gomega.Equal(ip.String()))
			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/GRAPHICS/TYPE").Text()).To(gomega.Equal(GraphicsTypeMap[VNC]))
		})
	})

	ginkgo.Describe("SetLogo", func() {
		var element string

		ginkgo.BeforeEach(func() {
			blueprint = &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			element = "xyz"
		})

		ginkgo.It("should set a LOGO tag to specified value", func() {
			blueprint.SetLogo(element)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/LOGO").Text()).To(gomega.Equal(element))
		})
	})

	ginkgo.Describe("SetMemory", func() {
		var element int

		ginkgo.BeforeEach(func() {
			blueprint = &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			element = 45
		})

		ginkgo.It("should set a MEMORY tag to specified value", func() {
			blueprint.SetMemory(element)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("TEMPLATE/MEMORY").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(element))
		})
	})

	ginkgo.Describe("SetNIC", func() {
		var nic *NICBlueprint

		ginkgo.BeforeEach(func() {
			nic = CreateNICBlueprint()
			nic.SetNetworkName(value)

			blueprint = &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
		})

		ginkgo.It("should set a NIC tag to specified value", func() {
			blueprint.SetNIC(*nic)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/NIC/NETWORK").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetOS", func() {
		var os *OSBlueprint
		var arch ArchitectureType

		ginkgo.BeforeEach(func() {
			os = CreateOSBlueprint()

			arch = ArchitectureTypeX86_64
			os.SetArchitecture(arch)

			blueprint = &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
		})

		ginkgo.It("should set a OS tag to specified value", func() {
			blueprint.SetOS(*os)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/OS/ARCH").Text()).To(gomega.Equal(ArchitectureTypeMap[arch]))
		})
	})

	ginkgo.Describe("SetRaw", func() {
		var raw *RawBlueprint

		ginkgo.BeforeEach(func() {
			raw = CreateRAWBlueprint()
			raw.SetData(value)
			raw.SetType("21258")

			blueprint = &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
		})

		ginkgo.It("should set a RAW tag to specified value", func() {
			blueprint.SetRaw(*raw)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/RAW/DATA").Text()).To(gomega.Equal(value))
			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/RAW/TYPE").Text()).To(gomega.Equal("21258"))
		})
	})

	ginkgo.Describe("SetSchedulingRequirements", func() {
		var element string

		ginkgo.BeforeEach(func() {
			blueprint = &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			element = "xyz"
		})

		ginkgo.It("should set a SCHED_REQUIREMENTS tag to specified value", func() {
			blueprint.SetSchedulingRequirements(element)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/SCHED_REQUIREMENTS").Text()).To(
				gomega.Equal(element))
		})
	})

	ginkgo.Describe("SetVCPU", func() {
		var element int

		ginkgo.BeforeEach(func() {
			blueprint = &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			element = 77
		})

		ginkgo.It("should set a VCPU tag to specified value", func() {
			blueprint.SetVCPU(element)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("TEMPLATE/VCPU").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(element))
		})
	})
})
