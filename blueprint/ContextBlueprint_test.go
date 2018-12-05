package blueprint

import (
	"net"
	"strconv"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("ContextBlueprint", func() {
	var blueprint *ContextBlueprint

	ginkgo.Describe("CreateContextBlueprint", func() {
		ginkgo.BeforeEach(func() {
			blueprint = CreateContextBlueprint()
		})

		ginkgo.It("should create a blueprint with CONTEXT element", func() {
			gomega.Expect(blueprint).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root()).ShouldNot(gomega.BeNil())
			gomega.Expect(blueprint.XMLData.Root().Tag).To(gomega.Equal("CONTEXT"))
		})
	})

	ginkgo.Describe("SetDiskID", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
			value = 41
		})

		ginkgo.It("should set DISK_ID tag to specified value", func() {
			blueprint.SetDiskID(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("CONTEXT/DISK_ID").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetEmail", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
			value = "test-value"
		})

		ginkgo.It("should set EMAIL tag to specified value", func() {
			blueprint.SetEmail(value)

			gomega.Expect(blueprint.XMLData.FindElement("CONTEXT/EMAIL").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetPublicIP", func() {
		var ip net.IP

		ginkgo.BeforeEach(func() {
			blueprint = &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
			ip = net.ParseIP("10.0.0.10")
		})

		ginkgo.It("should set PUBLIC_IP tag to specified value", func() {
			blueprint.SetPublicIP(ip)

			gomega.Expect(blueprint.XMLData.FindElement("CONTEXT/PUBLIC_IP").Text()).To(gomega.Equal(ip.String()))
		})
	})

	ginkgo.Describe("SetSSHKey", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
			value = "test-value"
		})

		ginkgo.It("should set SSH_KEY tag to specified value", func() {
			blueprint.SetSSHKey(value)

			gomega.Expect(blueprint.XMLData.FindElement("CONTEXT/SSH_KEY").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetTarget", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
			value = "test-value"
		})

		ginkgo.It("should set TARGET tag to specified value", func() {
			blueprint.SetTarget(value)

			gomega.Expect(blueprint.XMLData.FindElement("CONTEXT/TARGET").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetToken", func() {
		var value bool

		ginkgo.BeforeEach(func() {
			blueprint = &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
			value = true
		})

		ginkgo.It("should set TOKEN tag to specified value", func() {
			blueprint.SetToken(value)

			gomega.Expect(blueprint.XMLData.FindElement("CONTEXT/TOKEN").Text()).To(
				gomega.Equal(boolToString(value)))
		})
	})

	ginkgo.Describe("SetUserData", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
			value = "test-value"
		})

		ginkgo.It("should set USER_DATA tag to specified value", func() {
			blueprint.SetUserData(value)

			gomega.Expect(blueprint.XMLData.FindElement("CONTEXT/USER_DATA").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetVirtualMachineID", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
			value = 41
		})

		ginkgo.It("should set VMID tag to specified value", func() {
			blueprint.SetVirtualMachineID(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("CONTEXT/VMID").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetVirtualMachineGroupID", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
			value = 41
		})

		ginkgo.It("should set VM_GID tag to specified value", func() {
			blueprint.SetVirtualMachineGroupID(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("CONTEXT/VM_GID").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetVirtualMachineGroupName", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
			value = "test-value"
		})

		ginkgo.It("should set VM_GNAME tag to specified value", func() {
			blueprint.SetVirtualMachineGroupName(value)

			gomega.Expect(blueprint.XMLData.FindElement("CONTEXT/VM_GNAME").Text()).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetVirtualMachineUserID", func() {
		var value int

		ginkgo.BeforeEach(func() {
			blueprint = &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
			value = 41
		})

		ginkgo.It("should set VM_UID tag to specified value", func() {
			blueprint.SetVirtualMachineUserID(value)

			i, err := strconv.Atoi(blueprint.XMLData.FindElement("CONTEXT/VM_UID").Text())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(i).To(gomega.Equal(value))
		})
	})

	ginkgo.Describe("SetVirtualMachineUserName", func() {
		var value string

		ginkgo.BeforeEach(func() {
			blueprint = &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
			value = "test-value"
		})

		ginkgo.It("should set VM_UNAME tag to specified value", func() {
			blueprint.SetVirtualMachineUserName(value)

			gomega.Expect(blueprint.XMLData.FindElement("CONTEXT/VM_UNAME").Text()).To(gomega.Equal(value))
		})
	})
})
