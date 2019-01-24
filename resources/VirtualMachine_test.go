package resources

import (
	"net"
	"time"

	"github.com/beevik/etree"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	vmXML = "xml/virtualMachine.xml"
)

var _ = ginkgo.Describe("VirtualMachine", func() {
	var (
		doc            *etree.Document
		virtualMachine *VirtualMachine
		err            error
	)

	ginkgo.Describe("test getters", func() {
		ginkgo.Context("when virtual machine has all attributes", func() {
			ginkgo.BeforeEach(func() {
				// create virtual machine from XML data
				doc = etree.NewDocument()
				err = doc.ReadFromFile(vmXML)
				virtualMachine = CreateVirtualMachineFromXML(doc.Root())
			})

			ginkgo.It("should find all virtual machine attributes (excepting Template)", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach
				gomega.Expect(virtualMachine).ShouldNot(gomega.BeNil())

				gomega.Expect(virtualMachine.ID()).To(gomega.Equal(57502))
				gomega.Expect(virtualMachine.Name()).To(gomega.Equal("METACLOUD"))

				var userID int
				userID, err = virtualMachine.User()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(userID).To(gomega.Equal(46))

				var groupID int
				groupID, err = virtualMachine.Group()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(groupID).To(gomega.Equal(113))

				var permissions *Permissions
				permissions, err = virtualMachine.Permissions()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(permissions.User.Use).To(gomega.Equal(true))
				gomega.Expect(permissions.User.Manage).To(gomega.Equal(true))
				gomega.Expect(permissions.User.Admin).To(gomega.Equal(false))
				gomega.Expect(permissions.Group.Use).To(gomega.Equal(false))
				gomega.Expect(permissions.Group.Manage).To(gomega.Equal(false))
				gomega.Expect(permissions.Group.Admin).To(gomega.Equal(false))
				gomega.Expect(permissions.Other.Use).To(gomega.Equal(false))
				gomega.Expect(permissions.Other.Manage).To(gomega.Equal(false))
				gomega.Expect(permissions.Other.Admin).To(gomega.Equal(false))

				lastPoll := time.Unix(int64(1543406223), 0)
				gomega.Expect(virtualMachine.LastPoll()).To(gomega.Equal(&lastPoll))

				gomega.Expect(virtualMachine.State()).To(gomega.Equal(VirtualMachineState(3)))
				gomega.Expect(virtualMachine.LCMState()).To(gomega.Equal(VirtualMachineLcmState(3)))
				gomega.Expect(virtualMachine.PrevState()).To(gomega.Equal(VirtualMachineState(3)))
				gomega.Expect(virtualMachine.PrevLCMState()).To(gomega.Equal(VirtualMachineLcmState(3)))

				gomega.Expect(virtualMachine.Reschedule()).To(gomega.BeFalse())

				sTime := time.Unix(int64(1540931164), 0)
				gomega.Expect(virtualMachine.STime()).To(gomega.Equal(&sTime))
				gomega.Expect(virtualMachine.ETime()).To(gomega.BeNil()) // ETIME = 0

				gomega.Expect(virtualMachine.DeployID()).To(gomega.Equal("one-57502"))
			})

			ginkgo.It("should find VM Template Disk attributes", func() {
				var disks []*Disk
				disks, err = virtualMachine.Disks()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(disks).To(gomega.HaveLen(2))
				gomega.Expect(disks[0].ClusterID).To(gomega.Equal(119))
				gomega.Expect(disks[0].DevPrefix).To(gomega.Equal("vd"))
				gomega.Expect(disks[0].DiskType).To(gomega.Equal(DiskTypeBlock))
				gomega.Expect(disks[0].ReadOnly).To(gomega.Equal(false))

				gomega.Expect(disks[1].DatastoreID).To(gomega.Equal(152))
				gomega.Expect(disks[1].Driver).To(gomega.Equal("raw"))
				gomega.Expect(disks[1].DiskType).To(gomega.Equal(DiskTypeRbd))
				gomega.Expect(disks[1].ReadOnly).To(gomega.Equal(false))
			})

			ginkgo.It("should find VM Template Graphics attributes", func() {
				var graphics *Graphics
				graphics, err = virtualMachine.Graphics()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(graphics.Type).To(gomega.Equal(GraphicsTypeVNC))

				gomega.Expect(graphics.Listen).To(gomega.Equal(net.ParseIP("0.0.0.0")))
				port := 63402
				gomega.Expect(graphics.Port).To(gomega.Equal(&port))
				gomega.Expect(graphics.RandomPassword).To(gomega.Equal(true))
			})

			ginkgo.It("should find VM Template NIC attributes", func() {
				var nics []*NIC
				nics, err = virtualMachine.NICs()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(nics).To(gomega.HaveLen(1))
				gomega.Expect(nics[0].AddressRangeID).To(gomega.Equal(83))

				gomega.Expect(nics[0].IP).To(gomega.Equal(net.ParseIP("123.123.123.85")))
				gomega.Expect(nics[0].Mac).To(gomega.Equal("02:00:00:f4:fd:33"))
			})

			ginkgo.It("should find VM Template OperatingSystem attributes", func() {
				var os *OperatingSystem
				os, err = virtualMachine.OperatingSystem()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(os.Architecture).To(gomega.Equal(ArchitectureTypeX86_64))
				gomega.Expect(os.Bootloader).To(gomega.Equal(""))
			})

			ginkgo.It("should find VM Template Raw attributes", func() {
				var raw *Raw
				raw, err = virtualMachine.Raw()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(raw.Data).To(gomega.Equal("<cpu mode=’host-model’></cpu>"))
				gomega.Expect(raw.Type).To(gomega.Equal("kvm"))
			})

			ginkgo.It("should find VM Template HistoryRecords attributes", func() {
				var historyRecords []*History
				historyRecords, err = virtualMachine.HistoryRecords()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(historyRecords).To(gomega.HaveLen(2))
				gomega.Expect(historyRecords[0].OID).To(gomega.Equal(49277))
				gomega.Expect(historyRecords[0].Hostname).To(gomega.Equal("dukan24.ics.muni.cz"))
				hid := 932
				gomega.Expect(historyRecords[0].HID).To(gomega.Equal(&hid))
				cid := 117
				gomega.Expect(historyRecords[0].CID).To(gomega.Equal(&cid))

				psTime := time.Unix(int64(1519209121), 0)
				gomega.Expect(historyRecords[0].PSTime).To(gomega.Equal(&psTime))
				gomega.Expect(historyRecords[0].ESTime).To(gomega.BeNil()) // ESTIME = 0
				gomega.Expect(historyRecords[0].Action).To(gomega.Equal(Action(19)))
			})

			ginkgo.It("should find other VM Template attributes", func() {
				var template *etree.Element
				template, err = virtualMachine.Template()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(template).NotTo(gomega.BeNil())
				gomega.Expect(template.Tag).To(gomega.Equal("TEMPLATE"))

				var context *etree.Element
				context, err = virtualMachine.Context()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(context).NotTo(gomega.BeNil())
				gomega.Expect(context.Tag).To(gomega.Equal("CONTEXT"))

				var userTemplate *etree.Element
				userTemplate, err = virtualMachine.UserTemplate()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(userTemplate).NotTo(gomega.BeNil())
				gomega.Expect(userTemplate.Tag).To(gomega.Equal("USER_TEMPLATE"))

				gomega.Expect(virtualMachine.CPU()).To(gomega.Equal(0.25))
				gomega.Expect(virtualMachine.CreatedBy()).To(gomega.Equal(46))
				gomega.Expect(virtualMachine.Memory()).To(gomega.Equal(2048))
				gomega.Expect(virtualMachine.TemplateID()).To(gomega.Equal(4572))
				gomega.Expect(virtualMachine.VCPU()).To(gomega.Equal(1))
			})

			ginkgo.Context("when virtualMachine has only ID", func() {
				var virtualMachine *VirtualMachine

				ginkgo.It("should create virtualMachine", func() {
					virtualMachine = CreateVirtualMachineWithID(42)

					gomega.Expect(virtualMachine.ID()).To(gomega.Equal(42))
				})

				ginkgo.Context("when virtualMachine doesn't have given attribute", func() {
					ginkgo.BeforeEach(func() {
						virtualMachine = CreateVirtualMachineWithID(42)
						err = nil
					})

					ginkgo.It("should return that virtualMachine doesn't have name", func() {
						_, err = virtualMachine.Name()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have user", func() {
						_, err = virtualMachine.User()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have group", func() {
						_, err = virtualMachine.Group()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have permissions", func() {
						_, err = virtualMachine.Permissions()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have last poll", func() {
						_, err = virtualMachine.LastPoll()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have state", func() {
						_, err = virtualMachine.State()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have lcm state", func() {
						_, err = virtualMachine.LCMState()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have prev state", func() {
						_, err = virtualMachine.PrevState()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have prev lcm state", func() {
						_, err = virtualMachine.PrevLCMState()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have reschedule", func() {
						_, err = virtualMachine.Reschedule()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have start time", func() {
						_, err = virtualMachine.STime()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have end time", func() {
						_, err = virtualMachine.ETime()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have deploy ID", func() {
						_, err = virtualMachine.DeployID()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have template", func() {
						_, err = virtualMachine.Template()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have context", func() {
						_, err = virtualMachine.Context()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have cpu", func() {
						_, err = virtualMachine.CPU()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have created by", func() {
						_, err = virtualMachine.CreatedBy()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have disks", func() {
						var disks []*Disk

						disks, err = virtualMachine.Disks()
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(disks).Should(gomega.HaveLen(0))
					})

					ginkgo.It("should return that virtualMachine doesn't have graphics", func() {
						_, err = virtualMachine.Graphics()
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have memory", func() {
						_, err = virtualMachine.Memory()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have nic", func() {
						var nics []*NIC

						nics, err = virtualMachine.NICs()
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(nics).Should(gomega.HaveLen(0))
					})

					ginkgo.It("should return that virtualMachine doesn't have OS", func() {
						_, err = virtualMachine.OperatingSystem()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have raw", func() {
						_, err = virtualMachine.Raw()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have template ID", func() {
						_, err = virtualMachine.TemplateID()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have vcpu", func() {
						_, err = virtualMachine.VCPU()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have user template", func() {
						_, err = virtualMachine.UserTemplate()
						gomega.Expect(err).To(gomega.HaveOccurred())
					})

					ginkgo.It("should return that virtualMachine doesn't have history records", func() {
						var historyRecords []*History

						historyRecords, err = virtualMachine.HistoryRecords()
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(historyRecords).Should(gomega.HaveLen(0))
					})
				})
			})
		})
	})
})
