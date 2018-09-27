package resources

import (
	"github.com/beevik/etree"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	datastoreXML = "xml/datastore.xml"
)

var _ = ginkgo.Describe("Datastore", func() {
	var (
		doc       *etree.Document
		datastore *Datastore
		err       error
	)

	ginkgo.Describe("test getters", func() {
		ginkgo.Context("when datastore has all attributes", func() {
			ginkgo.BeforeEach(func() {
				// create user with data
				doc = etree.NewDocument()
				err = doc.ReadFromFile(datastoreXML)
				datastore = CreateDatastoreFromXML(doc.Root())
			})

			ginkgo.It("should find all datastore attributes", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach
				gomega.Expect(datastore).ShouldNot(gomega.BeNil())

				gomega.Expect(datastore.ID()).To(gomega.Equal(104))
				gomega.Expect(datastore.Name()).To(gomega.Equal("kingslanding"))

				var userID int
				userID, err = datastore.User()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(userID).To(gomega.Equal(0))

				var groupID int
				groupID, err = datastore.Group()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(groupID).To(gomega.Equal(111))

				var permissions *Permissions
				permissions, err = datastore.Permissions()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(permissions.User.Use).To(gomega.Equal(true))
				gomega.Expect(permissions.User.Manage).To(gomega.Equal(true))
				gomega.Expect(permissions.User.Admin).To(gomega.Equal(false))
				gomega.Expect(permissions.Group.Use).To(gomega.Equal(true))
				gomega.Expect(permissions.Group.Manage).To(gomega.Equal(false))
				gomega.Expect(permissions.Group.Admin).To(gomega.Equal(false))
				gomega.Expect(permissions.Other.Use).To(gomega.Equal(false))
				gomega.Expect(permissions.Other.Manage).To(gomega.Equal(false))
				gomega.Expect(permissions.Other.Admin).To(gomega.Equal(false))

				gomega.Expect(datastore.DsMad()).To(gomega.Equal("fs"))
				gomega.Expect(datastore.TmMad()).To(gomega.Equal("shared"))
				gomega.Expect(datastore.BasePath()).To(gomega.Equal("/var/lib/one//datastores/104"))
				gomega.Expect(datastore.Type()).To(gomega.Equal(ImageDs))
				gomega.Expect(datastore.DiskType()).To(gomega.Equal(DatastoreDiskType(0)))
				gomega.Expect(datastore.State()).To(gomega.Equal(Enabled))

				var clusterIDs []int
				clusterIDs, err = datastore.Clusters()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(clusterIDs[0]).To(gomega.Equal(0))

				gomega.Expect(datastore.TotalMB()).To(gomega.Equal(9714))
				gomega.Expect(datastore.FreeMB()).To(gomega.Equal(5704))
				gomega.Expect(datastore.UsedMB()).To(gomega.Equal(3495))

				gomega.Expect(datastore.Images()).To(gomega.HaveLen(4))
			})
		})

		ginkgo.Context("when datastore has only ID", func() {
			var datastore *Datastore

			ginkgo.It("should create datastore", func() {
				datastore = CreateDatastoreWithID(42)

				gomega.Expect(datastore.ID()).To(gomega.Equal(42))
			})

			ginkgo.Context("when datastore doesn't have given attribute", func() {
				ginkgo.BeforeEach(func() {
					datastore = CreateDatastoreWithID(42)
					err = nil
				})

				ginkgo.It("should return that datastore doesn't have name", func() {
					_, err = datastore.Name()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have user", func() {
					_, err = datastore.User()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have group", func() {
					_, err = datastore.Group()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have permissions", func() {
					_, err = datastore.Permissions()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have DS_MAD", func() {
					_, err = datastore.DsMad()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have TM_MAD", func() {
					_, err = datastore.TmMad()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have base path", func() {
					_, err = datastore.BasePath()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have type", func() {
					_, err = datastore.Type()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have disk type", func() {
					_, err = datastore.DiskType()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have state", func() {
					_, err = datastore.State()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have clusters", func() {
					var clusters []int

					clusters, err = datastore.Clusters()
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(clusters).Should(gomega.HaveLen(0))
				})

				ginkgo.It("should return that datastore doesn't have total MB", func() {
					_, err = datastore.TotalMB()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have free MB", func() {
					_, err = datastore.FreeMB()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have used MB", func() {
					_, err = datastore.UsedMB()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that datastore doesn't have images", func() {
					var images []int

					images, err = datastore.Images()
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(images).Should(gomega.HaveLen(0))
				})
			})
		})
	})
})
