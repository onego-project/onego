package blueprint

import (
	"strconv"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("HostBlueprint", func() {
	var blueprint *HostBlueprint

	ginkgo.Describe("CreateUpdateHostBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateUpdateHostBlueprint()
		})

		ginkgo.It("should create a blueprint with TEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("TEMPLATE"))
		})
	})

	ginkgo.Describe("SetArchitecture", func() {
		var arch string

		ginkgo.BeforeEach(func() {
			blueprint = &HostBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			arch = "x86_64"
		})

		ginkgo.It("should set ARCH tag to specified value", func() {
			blueprint.SetArchitecture(arch)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/ARCH").Text()).To(gomega.Equal(arch))
		})
	})

	ginkgo.Describe("SetClusterName", func() {
		var clusterName string

		ginkgo.BeforeEach(func() {
			blueprint = &HostBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			clusterName = "Kontrola"
		})

		ginkgo.It("should set CLUSTER_NAME tag to specified value", func() {
			blueprint.SetClusterName(clusterName)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/CLUSTER_NAME").Text()).To(
				gomega.Equal(clusterName))
		})
	})

	ginkgo.Describe("SetHostName", func() {
		var hostName string

		ginkgo.BeforeEach(func() {
			blueprint = &HostBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			hostName = "blob"
		})

		ginkgo.It("should set HOSTNAME tag to specified value", func() {
			blueprint.SetHostName(hostName)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/HOSTNAME").Text()).To(gomega.Equal(hostName))
		})
	})

	ginkgo.Describe("SetImMad", func() {
		var imMad string

		ginkgo.BeforeEach(func() {
			blueprint = &HostBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			imMad = "HITYY"
		})

		ginkgo.It("should set IM_MAD tag to specified value", func() {
			blueprint.SetImMad(imMad)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/IM_MAD").Text()).To(gomega.Equal(imMad))
		})
	})

	ginkgo.Describe("SetPriority", func() {
		var priority int

		ginkgo.BeforeEach(func() {
			blueprint = &HostBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			priority = 5
		})

		ginkgo.It("should set PRIORITY tag to specified value", func() {
			blueprint.SetPriority(priority)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("TEMPLATE/PRIORITY").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(priority))
		})
	})

	ginkgo.Describe("SetReservedCPU", func() {
		var reservedCPU int

		ginkgo.BeforeEach(func() {
			blueprint = &HostBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			reservedCPU = 88
		})

		ginkgo.It("should set RESERVED_CPU tag to specified value", func() {
			blueprint.SetReservedCPU(reservedCPU)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("TEMPLATE/RESERVED_CPU").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(reservedCPU))
		})
	})

	ginkgo.Describe("SetReservedMemory", func() {
		var reservedMemory int

		ginkgo.BeforeEach(func() {
			blueprint = &HostBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			reservedMemory = 44
		})

		ginkgo.It("should set RESERVED_MEM tag to specified value", func() {
			blueprint.SetReservedMemory(reservedMemory)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("TEMPLATE/RESERVED_MEM").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(reservedMemory))
		})
	})

	ginkgo.Describe("SetStatus", func() {
		var status string

		ginkgo.BeforeEach(func() {
			blueprint = &HostBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			status = "HITYY"
		})

		ginkgo.It("should set STATUS tag to specified value", func() {
			blueprint.SetStatus(status)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/STATUS").Text()).To(gomega.Equal(status))
		})
	})

	ginkgo.Describe("SetVMMad", func() {
		var vmMad string

		ginkgo.BeforeEach(func() {
			blueprint = &HostBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			vmMad = "HITYY"
		})

		ginkgo.It("should set VM_MAD tag to specified value", func() {
			blueprint.SetVMMad(vmMad)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/VM_MAD").Text()).To(gomega.Equal(vmMad))
		})
	})
})
