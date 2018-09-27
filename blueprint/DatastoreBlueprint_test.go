package blueprint

import (
	"github.com/onego-project/onego/resources"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("DatastoreBlueprint", func() {
	var blueprint *DatastoreBlueprint

	ginkgo.Describe("CreateUpdateDatastoreBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateUpdateDatastoreBlueprint()
		})

		ginkgo.It("should create a blueprint with TEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("TEMPLATE"))
		})
	})

	ginkgo.Describe("CreateAllocateDatastoreBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateAllocateDatastoreBlueprint()
		})

		ginkgo.It("should create a blueprint with DATASTORE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("DATASTORE"))
		})
	})

	ginkgo.Describe("SetDiskType", func() {
		var diskType resources.DatastoreDiskType

		ginkgo.BeforeEach(func() {
			blueprint = &DatastoreBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			diskType = resources.File
		})

		ginkgo.It("should set DISK_TYPE tag to specified value", func() {
			blueprint.SetDiskType(diskType)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DISK_TYPE").Text()).To(gomega.Equal(resources.DatastoreDiskTypeMap[diskType]))
		})
	})

	ginkgo.Describe("SetDsMad", func() {
		var dsMad string

		ginkgo.BeforeEach(func() {
			blueprint = &DatastoreBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			dsMad = "blob"
		})

		ginkgo.It("should set DS_MAD tag to specified value", func() {
			blueprint.SetDsMad(dsMad)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/DS_MAD").Text()).To(gomega.Equal(dsMad))
		})
	})

	ginkgo.Describe("SetTmMad", func() {
		var tmMad string

		ginkgo.BeforeEach(func() {
			blueprint = &DatastoreBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			tmMad = "blob"
		})

		ginkgo.It("should set TM_MAD tag to specified value", func() {
			blueprint.SetTmMad(tmMad)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/TM_MAD").Text()).To(gomega.Equal(tmMad))
		})
	})

	ginkgo.Describe("SetType", func() {
		var datastoreType resources.DatastoreType

		ginkgo.BeforeEach(func() {
			blueprint = &DatastoreBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			datastoreType = resources.ImageDs
		})

		ginkgo.It("should set TYPE tag to specified value", func() {
			blueprint.SetType(datastoreType)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/TYPE").Text()).To(gomega.Equal(resources.DatastoreTypeMap[datastoreType]))
		})
	})
})
