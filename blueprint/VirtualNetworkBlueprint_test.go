package blueprint

import (
	"strconv"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("VirtualNetworkBlueprint", func() {
	var blueprint *VirtualNetworkBlueprint

	ginkgo.Describe("CreateUpdateVirtualNetworkBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateUpdateVirtualNetworkBlueprint()
		})

		ginkgo.It("should create a blueprint with TEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("TEMPLATE"))
		})
	})

	ginkgo.Describe("CreateAllocateVirtualNetworkBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateAllocateVirtualNetworkBlueprint()
		})

		ginkgo.It("should create a blueprint with VNET element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("VNET"))
		})
	})

	ginkgo.Describe("SetBridge", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualNetworkBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = "test-value"
		})

		ginkgo.It("should set BRIDGE tag to specified value", func() {
			blueprint.SetBridge(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/BRIDGE").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetFilterIPSpoofing", func() {
		var value bool

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualNetworkBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = true
		})

		ginkgo.It("should set FILTER_IP_SPOOFING tag to specified value", func() {
			blueprint.SetFilterIPSpoofing(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/FILTER_IP_SPOOFING").Text()).To(
				gomega.Equal(boolToString(value)))
		})
	})

	ginkgo.Describe("SetFilterMacSpoofing", func() {
		var value bool

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualNetworkBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = false
		})

		ginkgo.It("should set FILTER_MAC_SPOOFING tag to specified value", func() {
			blueprint.SetFilterMacSpoofing(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/FILTER_MAC_SPOOFING").Text()).To(
				gomega.Equal(boolToString(value)))
		})
	})

	ginkgo.Describe("SetGateway", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualNetworkBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = "test-value"
		})

		ginkgo.It("should set GATEWAY tag to specified value", func() {
			blueprint.SetGateway(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/GATEWAY").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetMTU", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualNetworkBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = 41
		})

		ginkgo.It("should set MTU tag to specified value", func() {
			blueprint.SetMTU(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("TEMPLATE/MTU").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetNetworkAddress", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualNetworkBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = "test-value"
		})

		ginkgo.It("should set NETWORK_ADDRESS tag to specified value", func() {
			blueprint.SetNetworkAddress(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/NETWORK_ADDRESS").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetNetworkMask", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualNetworkBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = "test-value"
		})

		ginkgo.It("should set NETWORK_MASK tag to specified value", func() {
			blueprint.SetNetworkMask(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/NETWORK_MASK").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetPhysicalDevice", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualNetworkBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = "test-value"
		})

		ginkgo.It("should set PHYDEV tag to specified value", func() {
			blueprint.SetPhysicalDevice(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/PHYDEV").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetSecurityGroups", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualNetworkBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = "test-value"
		})

		ginkgo.It("should set SECURITY_GROUPS tag to specified value", func() {
			blueprint.SetSecurityGroups(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/SECURITY_GROUPS").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetVirtualLanID", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualNetworkBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = 11
		})

		ginkgo.It("should set VLAN_ID tag to specified value", func() {
			blueprint.SetVirtualLanID(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("TEMPLATE/VLAN_ID").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetVnMad", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &VirtualNetworkBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			value = "test-value"
		})

		ginkgo.It("should set VN_MAD tag to specified value", func() {
			blueprint.SetVnMad(value)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/VN_MAD").Text()).To(gomega.Equal(value))
		})
	})
})
