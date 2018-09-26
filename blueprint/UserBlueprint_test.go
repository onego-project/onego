package blueprint

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("UserBlueprint", func() {
	var blueprint *UserBlueprint

	ginkgo.Describe("CreateUpdateUserBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateUpdateUserBlueprint()
		})

		ginkgo.It("should create a blueprint with TEMPLATE element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("TEMPLATE"))
		})
	})

	ginkgo.Describe("SetEmail", func() {
		var email string

		ginkgo.BeforeEach(func() {
			blueprint = &UserBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			email = "dont@contact.me"
		})

		ginkgo.It("should set an EMAIL tag to specified value", func() {
			blueprint.SetEmail(email)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/EMAIL").Text()).To(gomega.Equal(email))
		})
	})

	ginkgo.Describe("SetFullName", func() {
		var fullName string

		ginkgo.BeforeEach(func() {
			blueprint = &UserBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			fullName = "John Doe"
		})

		ginkgo.It("should set an NAME tag to specified value", func() {
			blueprint.SetFullName(fullName)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/NAME").Text()).To(gomega.Equal(fullName))
		})
	})

	ginkgo.Describe("SetSSHPublicKey", func() {
		var sshKey string

		ginkgo.BeforeEach(func() {
			blueprint = &UserBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
			sshKey = "ssh-rsa AAAAB3NzPS...id2cznXR8X"
		})

		ginkgo.It("should set an SSH_PUBLIC_KEY tag to specified value", func() {
			blueprint.SetSSHPublicKey(sshKey)

			gomega.Expect(blueprint.XMLData.FindElement("TEMPLATE/SSH_PUBLIC_KEY").Text()).To(gomega.Equal(sshKey))
		})
	})
})
