package resources

import (
	"net"

	"github.com/beevik/etree"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	vnetXML = "xml/virtualNetwork.xml"
)

var _ = ginkgo.Describe("VirtualNetwork", func() {
	var (
		doc            *etree.Document
		virtualNetwork *VirtualNetwork
		err            error
	)

	ginkgo.Describe("test getters", func() {
		ginkgo.Context("when virtual network has all attributes", func() {
			ginkgo.BeforeEach(func() {
				// create user with data
				doc = etree.NewDocument()
				err = doc.ReadFromFile(vnetXML)
				virtualNetwork = CreateVirtualNetworkFromXML(doc.Root())
			})

			ginkgo.It("should find all virtual network attributes", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach
				gomega.Expect(virtualNetwork).ShouldNot(gomega.BeNil())

				gomega.Expect(virtualNetwork.ID()).To(gomega.Equal(738))
				gomega.Expect(virtualNetwork.Name()).To(gomega.Equal("metacloud"))

				var userID int
				userID, err = virtualNetwork.User()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(userID).To(gomega.Equal(10))

				var groupID int
				groupID, err = virtualNetwork.Group()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(groupID).To(gomega.Equal(101))

				var permissions *Permissions
				permissions, err = virtualNetwork.Permissions()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(permissions.User.Use).To(gomega.Equal(true))
				gomega.Expect(permissions.User.Manage).To(gomega.Equal(false))
				gomega.Expect(permissions.User.Admin).To(gomega.Equal(false))
				gomega.Expect(permissions.Group.Use).To(gomega.Equal(true))
				gomega.Expect(permissions.Group.Manage).To(gomega.Equal(false))
				gomega.Expect(permissions.Group.Admin).To(gomega.Equal(false))
				gomega.Expect(permissions.Other.Use).To(gomega.Equal(false))
				gomega.Expect(permissions.Other.Manage).To(gomega.Equal(false))
				gomega.Expect(permissions.Other.Admin).To(gomega.Equal(false))

				var clusterIDs []int
				clusterIDs, err = virtualNetwork.Clusters()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(clusterIDs[0]).To(gomega.Equal(112))
				gomega.Expect(clusterIDs).To(gomega.HaveLen(4))

				gomega.Expect(virtualNetwork.Bridge()).To(gomega.Equal("onebr0"))
				gomega.Expect(virtualNetwork.ParentNetworkID()).To(gomega.Equal(11))
				gomega.Expect(virtualNetwork.VnMad()).To(gomega.Equal("fw"))
				gomega.Expect(virtualNetwork.PhysicalDevice()).To(gomega.Equal("asdf"))
				gomega.Expect(virtualNetwork.VirtualLanID()).To(gomega.Equal(22))
				gomega.Expect(virtualNetwork.VirtualLanIDAutomatic()).To(gomega.Equal(0))
				gomega.Expect(virtualNetwork.UsedLeases()).To(gomega.Equal(108))

				var objIDs []int
				objIDs, err = virtualNetwork.VirtualRouters()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(objIDs[0]).To(gomega.Equal(111))
				gomega.Expect(objIDs).To(gomega.HaveLen(2))
			})

			ginkgo.It("should find virtual network address ranges", func() {
				var objects []*AddressRange
				objects, err = virtualNetwork.AddressRanges()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				id := 5
				ip := net.ParseIP("10.10.10.10")
				size := 1
				usedLeases := 1

				ar0 := objects[0]
				gomega.Expect(ar0.ID).To(gomega.Equal(&id))
				gomega.Expect(ar0.IP).To(gomega.Equal(&ip))
				gomega.Expect(ar0.Mac).To(gomega.Equal("11:22:33:44:55:66"))
				gomega.Expect(ar0.Size).To(gomega.Equal(&size))
				gomega.Expect(ar0.Type).To(gomega.Equal("IP4"))
				gomega.Expect(ar0.MacEnd).To(gomega.Equal("11:22:33:44:55:66"))
				gomega.Expect(ar0.IPEnd).To(gomega.Equal(net.ParseIP("10.10.10.10")))
				gomega.Expect(ar0.UsedLeases).To(gomega.Equal(&usedLeases))

				gomega.Expect(objects).To(gomega.HaveLen(5))
			})
		})

		ginkgo.Context("when virtualNetwork has only ID", func() {
			var virtualNetwork *VirtualNetwork

			ginkgo.It("should create virtualNetwork", func() {
				virtualNetwork = CreateVirtualNetworkWithID(42)

				gomega.Expect(virtualNetwork.ID()).To(gomega.Equal(42))
			})

			ginkgo.Context("when virtualNetwork doesn't have given attribute", func() {
				ginkgo.BeforeEach(func() {
					virtualNetwork = CreateVirtualNetworkWithID(42)
					err = nil
				})

				ginkgo.It("should return that virtualNetwork doesn't have name", func() {
					_, err = virtualNetwork.Name()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that virtualNetwork doesn't have user", func() {
					_, err = virtualNetwork.User()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that virtualNetwork doesn't have group", func() {
					_, err = virtualNetwork.Group()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that virtualNetwork doesn't have permissions", func() {
					_, err = virtualNetwork.Permissions()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that virtualNetwork doesn't have CLUSTERS", func() {
					var objectIDs []int

					objectIDs, err = virtualNetwork.Clusters()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(objectIDs).Should(gomega.HaveLen(0))
				})

				ginkgo.It("should return that virtualNetwork doesn't have BRIDGE", func() {
					_, err = virtualNetwork.Bridge()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that virtualNetwork doesn't have PARENT_NETWORK_ID", func() {
					_, err = virtualNetwork.ParentNetworkID()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that virtualNetwork doesn't have VN_MAD", func() {
					_, err = virtualNetwork.VnMad()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that virtualNetwork doesn't have PHYDEV", func() {
					_, err = virtualNetwork.PhysicalDevice()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that virtualNetwork doesn't have VLAN_ID", func() {
					_, err = virtualNetwork.VirtualLanID()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that virtualNetwork doesn't haveVLAN_ID_AUTOMATIC", func() {
					_, err = virtualNetwork.VirtualLanIDAutomatic()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that virtualNetwork doesn't have USED_LEASES", func() {
					_, err = virtualNetwork.UsedLeases()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that virtualNetwork doesn't have VROUTERS", func() {
					var objectIDs []int

					objectIDs, err = virtualNetwork.VirtualRouters()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(objectIDs).Should(gomega.HaveLen(0))
				})

				ginkgo.It("should return that virtualNetwork doesn't have AddressRanges", func() {
					var objects []*AddressRange

					objects, err = virtualNetwork.AddressRanges()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(objects).Should(gomega.HaveLen(0))
				})
			})
		})
	})
})
