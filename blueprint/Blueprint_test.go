package blueprint

import (
	"github.com/beevik/etree"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Blueprint", func() {
	var (
		blueprint *Blueprint
		err       error
	)

	ginkgo.BeforeEach(func() {
		doc := etree.NewDocument()
		root := doc.CreateElement("ROOT")
		node := root.CreateElement("NODE")
		node.SetText("TEXT")

		blueprint = &Blueprint{XMLData: doc}
	})

	ginkgo.Describe("CreateBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateBlueprint("ROOT")
		})

		ginkgo.It("should create a blueprint with root element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("ROOT"))
		})
	})

	ginkgo.Describe("Render", func() {
		ginkgo.Context("without XML data", func() {
			ginkgo.BeforeEach(func() {
				blueprint = &Blueprint{XMLData: nil}
			})

			ginkgo.It("returns an error", func() {
				_, err = blueprint.Render()
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("with XML data", func() {
			ginkgo.It("renders its content", func() {
				xml, err := blueprint.Render()

				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(xml).To(gomega.Equal("<ROOT><NODE>TEXT</NODE></ROOT>"))
			})
		})
	})

	ginkgo.Describe("SetElement", func() {
		ginkgo.Context("with already existing element", func() {
			ginkgo.It("sets element to specified value", func() {
				blueprint.SetElement("NODE", "BLOB")

				gomega.Expect(blueprint.XMLData.FindElement("ROOT/NODE").Text()).To(gomega.Equal("BLOB"))
			})
		})

		ginkgo.Context("without already existing element", func() {
			ginkgo.BeforeEach(func() {
				doc := etree.NewDocument()
				doc.CreateElement("ROOT")

				blueprint = &Blueprint{XMLData: doc}
			})

			ginkgo.It("creates element and sets it to specified value", func() {
				blueprint.SetElement("NODE", "BLOB")

				gomega.Expect(blueprint.XMLData.FindElement("ROOT/NODE").Text()).To(gomega.Equal("BLOB"))
			})
		})
	})

	ginkgo.Describe("SetName", func() {
		ginkgo.Context("with already existing element", func() {
			ginkgo.It("sets name to specified value", func() {
				blueprint.SetName("BLOB")

				gomega.Expect(blueprint.XMLData.FindElement("ROOT/NAME").Text()).To(gomega.Equal("BLOB"))
			})
		})

		ginkgo.Context("without already existing element", func() {
			ginkgo.BeforeEach(func() {
				doc := etree.NewDocument()
				doc.CreateElement("ROOT")

				blueprint = &Blueprint{XMLData: doc}
			})

			ginkgo.It("creates name and sets it to specified value", func() {
				blueprint.SetName("BLOB")

				gomega.Expect(blueprint.XMLData.FindElement("ROOT/NAME").Text()).To(gomega.Equal("BLOB"))
			})
		})
	})
})
