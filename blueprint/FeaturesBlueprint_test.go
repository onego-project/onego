package blueprint

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("FeaturesBlueprint", func() {
	var blueprint *FeaturesBlueprint

	ginkgo.Describe("CreateFeaturesBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateFeaturesBlueprint()
		})

		ginkgo.It("should create a blueprint with FEATURES element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("FEATURES"))
		})
	})

	ginkgo.Describe("SetACPI", func() {
		var value bool

		ginkgo.BeforeEach(func() {
			blueprint = &FeaturesBlueprint{Blueprint: *CreateBlueprint("FEATURES")}
			value = true
		})

		ginkgo.It("should set ACPI tag to specified value", func() {
			blueprint.SetACPI(value)

			gomega.Expect(blueprint.XMLData.FindElement("FEATURES/ACPI").Text()).To(
				gomega.Equal(boolToString(value)))
		})
	})

	ginkgo.Describe("SetAPIC", func() {
		var value bool

		ginkgo.BeforeEach(func() {
			blueprint = &FeaturesBlueprint{Blueprint: *CreateBlueprint("FEATURES")}
			value = true
		})

		ginkgo.It("should set APIC tag to specified value", func() {
			blueprint.SetAPIC(value)

			gomega.Expect(blueprint.XMLData.FindElement("FEATURES/APIC").Text()).To(
				gomega.Equal(boolToString(value)))
		})
	})

	ginkgo.Describe("SetLocalTime", func() {
		var value bool

		ginkgo.BeforeEach(func() {
			blueprint = &FeaturesBlueprint{Blueprint: *CreateBlueprint("FEATURES")}
			value = true
		})

		ginkgo.It("should set LOCALTIME tag to specified value", func() {
			blueprint.SetLocalTime(value)

			gomega.Expect(blueprint.XMLData.FindElement("FEATURES/LOCALTIME").Text()).To(
				gomega.Equal(boolToString(value)))
		})
	})

	ginkgo.Describe("SetPAE", func() {
		var value bool

		ginkgo.BeforeEach(func() {
			blueprint = &FeaturesBlueprint{Blueprint: *CreateBlueprint("FEATURES")}
			value = true
		})

		ginkgo.It("should set PAE tag to specified value", func() {
			blueprint.SetPAE(value)

			gomega.Expect(blueprint.XMLData.FindElement("FEATURES/PAE").Text()).To(
				gomega.Equal(boolToString(value)))
		})
	})

	ginkgo.Describe("SetHyperV", func() {
		var value bool

		ginkgo.BeforeEach(func() {
			blueprint = &FeaturesBlueprint{Blueprint: *CreateBlueprint("FEATURES")}
			value = true
		})

		ginkgo.It("should set HYPERV tag to specified value", func() {
			blueprint.SetHyperV(value)

			gomega.Expect(blueprint.XMLData.FindElement("FEATURES/HYPERV").Text()).To(
				gomega.Equal(boolToString(value)))
		})
	})

	ginkgo.Describe("SetGuestAgent", func() {
		var value bool

		ginkgo.BeforeEach(func() {
			blueprint = &FeaturesBlueprint{Blueprint: *CreateBlueprint("FEATURES")}
			value = true
		})

		ginkgo.It("should set GUEST_AGENT tag to specified value", func() {
			blueprint.SetGuestAgent(value)

			gomega.Expect(blueprint.XMLData.FindElement("FEATURES/GUEST_AGENT").Text()).To(
				gomega.Equal(boolToString(value)))
		})
	})
})
