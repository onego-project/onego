package blueprint

import (
	"strconv"

	"github.com/onego-project/onego/resources"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("DiskBlueprint", func() {
	var blueprint *DiskBlueprint

	ginkgo.Describe("CreateDiskBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateDiskBlueprint()
		})

		ginkgo.It("should create a blueprint with DISK element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("DISK"))
		})
	})

	ginkgo.Describe("SetDevicePrefix", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = "test-value"
		})

		ginkgo.It("should set DEV_PREFIX tag to specified value", func() {
			blueprint.SetDevicePrefix(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/DEV_PREFIX").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetClone", func() {
		var value bool

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = true
		})

		ginkgo.It("should set CLONE tag to specified value", func() {
			blueprint.SetClone(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/CLONE").Text()).To(
				gomega.Equal(boolToString(value)))
		})
	})

	ginkgo.Describe("SetCloneTarget", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = "test-value"
		})

		ginkgo.It("should set CLONE_TARGET tag to specified value", func() {
			blueprint.SetCloneTarget(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/CLONE_TARGET").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetClusterID", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = 41
		})

		ginkgo.It("should set CLUSTER_ID tag to specified value", func() {
			blueprint.SetClusterID(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("DISK/CLUSTER_ID").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetDatastoreID", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = 41
		})

		ginkgo.It("should set DATASTORE_ID tag to specified value", func() {
			blueprint.SetDatastoreID(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("DISK/DATASTORE_ID").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetDatastore", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = "test-value"
		})

		ginkgo.It("should set DATASTORE tag to specified value", func() {
			blueprint.SetDatastore(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/DATASTORE").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetDiskID", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = 11
		})

		ginkgo.It("should set DISK_ID tag to specified value", func() {
			blueprint.SetDiskID(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("DISK/DISK_ID").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetDiskType", func() {
		var value resources.DiskType

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = resources.DiskTypeBlock
		})

		ginkgo.It("should set DISK_TYPE tag to specified value", func() {
			blueprint.SetDiskType(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/DISK_TYPE").Text()).To(gomega.Equal(
				resources.DiskTypeMap[value]))
		})
	})

	ginkgo.Describe("SetDriver", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = "test-value"
		})

		ginkgo.It("should set DRIVER tag to specified value", func() {
			blueprint.SetDriver(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/DRIVER").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetImage", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = "test-value"
		})

		ginkgo.It("should set IMAGE tag to specified value", func() {
			blueprint.SetImage(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/IMAGE").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetImageID", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = 11
		})

		ginkgo.It("should set IMAGE_ID tag to specified value", func() {
			blueprint.SetImageID(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("DISK/IMAGE_ID").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetImageUserName", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = "test-value"
		})

		ginkgo.It("should set IMAGE_UNAME tag to specified value", func() {
			blueprint.SetImageUserName(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/IMAGE_UNAME").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetPoolName", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = "test-value"
		})

		ginkgo.It("should set POOL_NAME tag to specified value", func() {
			blueprint.SetPoolName(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/POOL_NAME").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetReadOnly", func() {
		var value bool

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = true
		})

		ginkgo.It("should set READONLY tag to specified value", func() {
			blueprint.SetReadOnly(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/READONLY").Text()).To(
				gomega.Equal(boolToString(value)))
		})
	})

	ginkgo.Describe("SetSave", func() {
		var value bool

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = true
		})

		ginkgo.It("should set SAVE tag to specified value", func() {
			blueprint.SetSave(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/SAVE").Text()).To(
				gomega.Equal(boolToString(value)))
		})
	})

	ginkgo.Describe("SetSize", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = 11
		})

		ginkgo.It("should set SIZE tag to specified value", func() {
			blueprint.SetSize(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("DISK/SIZE").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetSource", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = "test-value"
		})

		ginkgo.It("should set SOURCE tag to specified value", func() {
			blueprint.SetSource(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/SOURCE").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetTarget", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = "test-value"
		})

		ginkgo.It("should set TARGET tag to specified value", func() {
			blueprint.SetTarget(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/TARGET").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetTmMad", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = "test-value"
		})

		ginkgo.It("should set TM_MAD tag to specified value", func() {
			blueprint.SetTmMad(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/TM_MAD").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetType", func() {
		var value resources.DiskType

		ginkgo.BeforeEach(func() {
			blueprint = &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
			value = resources.DiskTypeBlock
		})

		ginkgo.It("should set TYPE tag to specified value", func() {
			blueprint.SetType(value)

			gomega.Expect(blueprint.XMLData.FindElement("DISK/TYPE").Text()).To(gomega.Equal(
				resources.DiskTypeMap[value]))
		})
	})
})
