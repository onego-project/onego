package blueprint

import (
	"strconv"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("NICBlueprint", func() {
	var blueprint *NICBlueprint

	ginkgo.Describe("CreateNICBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateNICBlueprint()
		})

		ginkgo.It("should create a blueprint with NIC element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("NIC"))
		})
	})

	ginkgo.Describe("SetNetworkName", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &NICBlueprint{Blueprint: *CreateBlueprint("NIC")}
			value = "test-value"
		})

		ginkgo.It("should set NETWORK tag to specified value", func() {
			blueprint.SetNetworkName(value)

			gomega.Expect(blueprint.XMLData.FindElement("NIC/NETWORK").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetNetworkOwnerName", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &NICBlueprint{Blueprint: *CreateBlueprint("NIC")}
			value = "test-value"
		})

		ginkgo.It("should set NETWORK_UNAME tag to specified value", func() {
			blueprint.SetNetworkOwnerName(value)

			gomega.Expect(blueprint.XMLData.FindElement("NIC/NETWORK_UNAME").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetNetworkID", func() {
		var id int

		ginkgo.BeforeEach(func() {
			blueprint = &NICBlueprint{Blueprint: *CreateBlueprint("NIC")}
			id = 420
		})

		ginkgo.It("should set NETWORK_ID tag to specified value", func() {
			blueprint.SetNetworkID(id)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("NIC/NETWORK_ID").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(id))
		})
	})

	ginkgo.Describe("SetNetworkOwnerID", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &NICBlueprint{Blueprint: *CreateBlueprint("NIC")}
			value = 42
		})

		ginkgo.It("should set NETWORK_UID tag to specified value", func() {
			blueprint.SetNetworkOwnerID(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("NIC/NETWORK_UID").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})
})
