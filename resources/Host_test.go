package resources

import (
	"github.com/beevik/etree"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	hostXML = "xml/host.xml"
)

var _ = ginkgo.Describe("Host", func() {
	var (
		doc  *etree.Document
		host *Host
		err  error
	)

	ginkgo.Describe("test getters", func() {
		ginkgo.Context("when host has all attributes", func() {
			ginkgo.BeforeEach(func() {
				// create host with data
				doc = etree.NewDocument()
				err = doc.ReadFromFile(hostXML)
				host = CreateHostFromXML(doc.Root())
			})

			ginkgo.It("should find all host attributes", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach
				gomega.Expect(host).ShouldNot(gomega.BeNil())

				gomega.Expect(host.ID()).To(gomega.Equal(934))
				gomega.Expect(host.Name()).To(gomega.Equal("gorbag.ics"))

				gomega.Expect(host.State()).To(gomega.Equal(HostState(2)))
				gomega.Expect(host.IMMad()).To(gomega.Equal("kvm"))
				gomega.Expect(host.VMMad()).To(gomega.Equal("kvm"))
				gomega.Expect(host.LastMonitoringTime()).To(gomega.Equal(1539689293))

				var clusterID int
				clusterID, err = host.Cluster()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(clusterID).To(gomega.Equal(118))

				gomega.Expect(host.DiskUsage()).To(gomega.Equal(0))
				gomega.Expect(host.MemoryUsage()).To(gomega.Equal(125829120))
				gomega.Expect(host.CPUUsage()).To(gomega.Equal(3500))
				gomega.Expect(host.MaxDisk()).To(gomega.Equal(368))
				gomega.Expect(host.MaxMemory()).To(gomega.Equal(127608968))
				gomega.Expect(host.MaxCPU()).To(gomega.Equal(4000))
				gomega.Expect(host.FreeDisk()).To(gomega.Equal(239))
				gomega.Expect(host.FreeMemory()).To(gomega.Equal(45884464))
				gomega.Expect(host.FreeCPU()).To(gomega.Equal(3880))
				gomega.Expect(host.UsedDisk()).To(gomega.Equal(110))
				gomega.Expect(host.UsedMemory()).To(gomega.Equal(85918808))
				gomega.Expect(host.UsedCPU()).To(gomega.Equal(120))
				gomega.Expect(host.RunningVMs()).To(gomega.Equal(18))

				var datastores []int
				datastores, err = host.Datastores()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(datastores).To(gomega.HaveLen(1))
				gomega.Expect(datastores[0]).To(gomega.Equal(151))

				var pcis []*PCI
				pcis, err = host.PCIDevices()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(pcis).To(gomega.HaveLen(2))
				gomega.Expect(pcis[0].Address).To(gomega.Equal("0000:02:00:0"))
				gomega.Expect(pcis[0].Bus).To(gomega.Equal("02"))
				gomega.Expect(pcis[0].Class).To(gomega.Equal("0300"))
				gomega.Expect(pcis[0].ClassName).To(gomega.Equal("VGA compatible controller"))
				gomega.Expect(pcis[0].Device).To(gomega.Equal("100c"))
				gomega.Expect(pcis[0].DeviceName).To(gomega.Equal("GK110B [GeForce GTX TITAN Black]"))
				gomega.Expect(pcis[0].Domain).To(gomega.Equal("0000"))
				gomega.Expect(pcis[0].Function).To(gomega.Equal("0"))
				gomega.Expect(pcis[0].ShortAddress).To(gomega.Equal("02:00.0"))
				gomega.Expect(pcis[0].Slot).To(gomega.Equal("00"))
				gomega.Expect(pcis[0].Type).To(gomega.Equal("10de:100c:0300"))
				gomega.Expect(pcis[0].Vendor).To(gomega.Equal("10de"))
				gomega.Expect(pcis[0].VendorName).To(gomega.Equal("NVIDIA Corporation"))
				gomega.Expect(pcis[0].VMID).To(gomega.Equal("-1"))

				var vmIDs []int
				vmIDs, err = host.VirtualMachines()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(vmIDs).To(gomega.HaveLen(2))
				gomega.Expect(vmIDs[0]).To(gomega.Equal(42810))
				gomega.Expect(vmIDs[1]).To(gomega.Equal(42868))
			})
		})

		ginkgo.Context("when host has only ID", func() {
			var host *Host

			ginkgo.It("should create host", func() {
				host = CreateHostWithID(42)

				gomega.Expect(host.ID()).To(gomega.Equal(42))
			})

			ginkgo.Context("when host doesn't have given attribute", func() {
				ginkgo.BeforeEach(func() {
					host = CreateHostWithID(42)
					err = nil
				})

				ginkgo.It("should return that host doesn't have name", func() {
					_, err = host.Name()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have State", func() {
					_, err = host.State()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have IMMad", func() {
					_, err = host.IMMad()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have VMMad", func() {
					_, err = host.VMMad()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have LastMonitoringTime", func() {
					_, err = host.LastMonitoringTime()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have Cluster", func() {
					_, err = host.Cluster()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have DiskUsage", func() {
					_, err = host.DiskUsage()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have MemoryUsage", func() {
					_, err = host.MemoryUsage()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have CPUUsage", func() {
					_, err = host.CPUUsage()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have MaxDisk", func() {
					_, err = host.MaxDisk()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have MaxMemory", func() {
					_, err = host.MaxMemory()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have MaxCPU", func() {
					_, err = host.MaxCPU()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have FreeDisk", func() {
					_, err = host.FreeDisk()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have FreeMemory", func() {
					_, err = host.FreeMemory()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have FreeCPU", func() {
					_, err = host.FreeCPU()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have UsedDisk", func() {
					_, err = host.UsedDisk()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have UsedMemory", func() {
					_, err = host.UsedMemory()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have UsedCPU", func() {
					_, err = host.UsedCPU()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have RunningVMs", func() {
					_, err = host.RunningVMs()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that host doesn't have Datastores", func() {
					var datastores []int

					datastores, err = host.Datastores()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(datastores).Should(gomega.HaveLen(0))
				})

				ginkgo.It("should return that host doesn't have PCIDevices", func() {
					var pcis []*PCI

					pcis, err = host.PCIDevices()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(pcis).Should(gomega.HaveLen(0))
				})

				ginkgo.It("should return that host doesn't have VMS", func() {
					var vms []int

					vms, err = host.VirtualMachines()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(vms).Should(gomega.HaveLen(0))
				})
			})
		})
	})
})
