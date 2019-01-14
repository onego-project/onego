package blueprint

import (
	"net"

	"github.com/onego-project/onego/resources"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GraphicsBlueprint", func() {
	var blueprint *GraphicsBlueprint

	ginkgo.Describe("CreateGraphicsBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateGraphicsBlueprint()
		})

		ginkgo.It("should create a blueprint with GRAPHICS element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("GRAPHICS"))
		})
	})

	ginkgo.Describe("SetListen", func() {
		var ip net.IP

		ginkgo.BeforeEach(func() {
			blueprint = &GraphicsBlueprint{Blueprint: *CreateBlueprint("GRAPHICS")}
			ip = net.ParseIP("10.0.0.10")
		})

		ginkgo.It("should set LISTEN tag to specified value", func() {
			blueprint.SetListen(ip)

			gomega.Expect(blueprint.XMLData.FindElement("GRAPHICS/LISTEN").Text()).To(
				gomega.Equal(ip.String()))
		})
	})

	ginkgo.Describe("SetType", func() {
		var value resources.GraphicsType

		ginkgo.BeforeEach(func() {
			blueprint = &GraphicsBlueprint{Blueprint: *CreateBlueprint("GRAPHICS")}
			value = resources.GraphicsTypeSDL
		})

		ginkgo.It("should set TYPE tag to specified value", func() {
			blueprint.SetType(value)

			gomega.Expect(blueprint.XMLData.FindElement("GRAPHICS/TYPE").Text()).To(
				gomega.Equal(resources.GraphicsTypeMap[value]))
		})
	})

	ginkgo.Describe("SetPort", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &GraphicsBlueprint{Blueprint: *CreateBlueprint("GRAPHICS")}
			value = "test-value"
		})

		ginkgo.It("should set PORT tag to specified value", func() {
			blueprint.SetPort(value)

			gomega.Expect(blueprint.XMLData.FindElement("GRAPHICS/PORT").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetKeyMap", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &GraphicsBlueprint{Blueprint: *CreateBlueprint("GRAPHICS")}
			value = "test-value"
		})

		ginkgo.It("should set KEYMAP tag to specified value", func() {
			blueprint.SetKeyMap(value)

			gomega.Expect(blueprint.XMLData.FindElement("GRAPHICS/KEYMAP").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetPassword", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &GraphicsBlueprint{Blueprint: *CreateBlueprint("GRAPHICS")}
			value = "test-value"
		})

		ginkgo.It("should set PASSWD tag to specified value", func() {
			blueprint.SetPassword(value)

			gomega.Expect(blueprint.XMLData.FindElement("GRAPHICS/PASSWD").Text()).To(gomega.Equal(value))
		})
	})
})
