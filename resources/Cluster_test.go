package resources

import (
	"github.com/beevik/etree"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	clusterXML = "xml/cluster.xml"
)

var _ = ginkgo.Describe("Cluster", func() {
	var (
		doc     *etree.Document
		cluster *Cluster
		err     error
	)

	ginkgo.Describe("test getters", func() {
		ginkgo.Context("when cluster has all attributes", func() {
			ginkgo.BeforeEach(func() {
				// create cluster with data
				doc = etree.NewDocument()
				err = doc.ReadFromFile(clusterXML)
				cluster = CreateClusterFromXML(doc.Root())
			})

			ginkgo.It("should find all cluster attributes", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach
				gomega.Expect(cluster).ShouldNot(gomega.BeNil())

				gomega.Expect(cluster.ID()).To(gomega.Equal(120))
				gomega.Expect(cluster.Name()).To(gomega.Equal("duilin"))

				var hosts []int
				hosts, err = cluster.Hosts()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(hosts).To(gomega.HaveLen(4))
				gomega.Expect(hosts[0]).To(gomega.Equal(937))
				gomega.Expect(hosts[1]).To(gomega.Equal(938))
				gomega.Expect(hosts[2]).To(gomega.Equal(939))
				gomega.Expect(hosts[3]).To(gomega.Equal(940))

				var datastores []int
				datastores, err = cluster.Datastores()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(datastores).To(gomega.HaveLen(3))
				gomega.Expect(datastores[0]).To(gomega.Equal(2))
				gomega.Expect(datastores[1]).To(gomega.Equal(105))
				gomega.Expect(datastores[2]).To(gomega.Equal(120))

				var vnets []int
				vnets, err = cluster.VirtualNetworks()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(vnets).To(gomega.HaveLen(2))
				gomega.Expect(vnets[0]).To(gomega.Equal(1382))
				gomega.Expect(vnets[1]).To(gomega.Equal(1383))
			})
		})

		ginkgo.Context("when cluster has only ID", func() {
			var cluster *Cluster

			ginkgo.It("should create cluster", func() {
				cluster = CreateClusterWithID(42)

				gomega.Expect(cluster.ID()).To(gomega.Equal(42))
			})

			ginkgo.Context("when cluster doesn't have given attribute", func() {
				ginkgo.BeforeEach(func() {
					cluster = CreateClusterWithID(42)
					err = nil
				})

				ginkgo.It("should return that cluster doesn't have name", func() {
					_, err = cluster.Name()
					gomega.Expect(err).To(gomega.HaveOccurred())
				})

				ginkgo.It("should return that cluster doesn't have Hosts", func() {
					var hosts []int

					hosts, err = cluster.Hosts()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(hosts).Should(gomega.HaveLen(0))
				})

				ginkgo.It("should return that cluster doesn't have Datastores", func() {
					var datastores []int

					datastores, err = cluster.Datastores()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(datastores).Should(gomega.HaveLen(0))
				})

				ginkgo.It("should return that cluster doesn't have Virtual Networks", func() {
					var vnets []int

					vnets, err = cluster.VirtualNetworks()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(vnets).Should(gomega.HaveLen(0))
				})
			})
		})
	})
})
