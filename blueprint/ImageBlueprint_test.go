package blueprint

import (
	"github.com/onego-project/onego/resources"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("ImageBlueprint", func() {
	var blueprint *ImageBlueprint

	ginkgo.Describe("CreateUpdateImageBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateUpdateImageBlueprint()
		})

		ginkgo.It("should create a blueprint with TEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("TEMPLATE"))
		})
	})

	ginkgo.Describe("CreateAllocateImageBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateAllocateImageBlueprint()
		})

		ginkgo.It("should create a blueprint with IMAGE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("IMAGE"))
		})
	})

	ginkgo.Describe("SetDescription", func() {
		var description string

		ginkgo.BeforeEach(func() {
			blueprint = &ImageBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			description = "Some text"
		})

		ginkgo.It("should set DESCRIPTION tag to specified value", func() {
			blueprint.SetDescription(description)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DESCRIPTION").Text()).To(gomega.Equal(description))
		})
	})

	ginkgo.Describe("SetDevPrefix", func() {
		var devPrefix string

		ginkgo.BeforeEach(func() {
			blueprint = &ImageBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			devPrefix = "Kontrola"
		})

		ginkgo.It("should set DEV_PREFIX tag to specified value", func() {
			blueprint.SetDevPrefix(devPrefix)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DEV_PREFIX").Text()).To(gomega.Equal(devPrefix))
		})
	})

	ginkgo.Describe("SetDiskType", func() {
		var diskType resources.DiskType

		ginkgo.BeforeEach(func() {
			blueprint = &ImageBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			diskType = resources.DiskTypeFile
		})

		ginkgo.It("should set DISK_TYPE tag to specified value", func() {
			blueprint.SetDiskType(diskType)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DISK_TYPE").Text()).To(gomega.Equal(
				resources.DiskTypeMap[diskType]))
		})
	})

	ginkgo.Describe("SetDriver", func() {
		var driver string

		ginkgo.BeforeEach(func() {
			blueprint = &ImageBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			driver = "blob"
		})

		ginkgo.It("should set DRIVER tag to specified value", func() {
			blueprint.SetDriver(driver)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DRIVER").Text()).To(gomega.Equal(driver))
		})
	})

	ginkgo.Describe("SetTarget", func() {
		var target string

		ginkgo.BeforeEach(func() {
			blueprint = &ImageBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			target = "HITYY"
		})

		ginkgo.It("should set TARGET tag to specified value", func() {
			blueprint.SetTarget(target)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/TARGET").Text()).To(gomega.Equal(target))
		})
	})
})
