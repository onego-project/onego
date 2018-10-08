package resources

import (
	"time"

	"github.com/beevik/etree"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	imageXML = "xml/image.xml"
)

var _ = ginkgo.Describe("Image", func() {
	var (
		doc   *etree.Document
		image *Image
		err   error
	)

	ginkgo.Describe("test getters", func() {
		ginkgo.Context("when image has all attributes", func() {
			ginkgo.BeforeEach(func() {
				// create user with data
				doc = etree.NewDocument()
				err = doc.ReadFromFile(imageXML)
				image = CreateImageFromXML(doc.Root())
			})

			ginkgo.It("should find all image attributes", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach
				gomega.Expect(image).ShouldNot(gomega.BeNil())

				gomega.Expect(image.ID()).To(gomega.Equal(123))
				gomega.Expect(image.Name()).To(gomega.Equal("tty-linux-local"))

				var userID int
				userID, err = image.User()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(userID).To(gomega.Equal(1881))

				var groupID int
				groupID, err = image.Group()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(groupID).To(gomega.Equal(8118))

				var permissions *Permissions
				permissions, err = image.Permissions()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(permissions.User.Use).To(gomega.Equal(true))
				gomega.Expect(permissions.User.Manage).To(gomega.Equal(true))
				gomega.Expect(permissions.User.Admin).To(gomega.Equal(false))
				gomega.Expect(permissions.Group.Use).To(gomega.Equal(true))
				gomega.Expect(permissions.Group.Manage).To(gomega.Equal(false))
				gomega.Expect(permissions.Group.Admin).To(gomega.Equal(false))
				gomega.Expect(permissions.Other.Use).To(gomega.Equal(true))
				gomega.Expect(permissions.Other.Manage).To(gomega.Equal(false))
				gomega.Expect(permissions.Other.Admin).To(gomega.Equal(false))

				gomega.Expect(image.Type()).To(gomega.Equal(ImageType(0)))
				gomega.Expect(image.DiskType()).To(gomega.Equal(DiskType(0)))
				gomega.Expect(image.Persistent()).To(gomega.Equal(false))

				regTime := time.Unix(int64(1490793710), 0)
				gomega.Expect(image.RegistrationTime()).To(gomega.Equal(&regTime))
				gomega.Expect(image.Source()).To(gomega.Equal(
					"/var/lib/one//datastores/101/5995f01c3c35883b297eed95a12a271b"))
				gomega.Expect(image.Path()).To(gomega.Equal("/var/tmp/tty-linux.img"))
				gomega.Expect(image.FileSystemType()).To(gomega.Equal("raw"))
				gomega.Expect(image.Size()).To(gomega.Equal(40))
				gomega.Expect(image.State()).To(gomega.Equal(ImageState(2)))
				gomega.Expect(image.RunningVMs()).To(gomega.Equal(1))

				var datastoreID int
				datastoreID, err = image.Datastore()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(datastoreID).To(gomega.Equal(101))

				var vmIDs []int
				vmIDs, err = image.VirtualMachines()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(vmIDs[0]).To(gomega.Equal(84))
				gomega.Expect(vmIDs).To(gomega.HaveLen(2))

				var cloneIDs []int
				cloneIDs, err = image.Clones()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(cloneIDs[0]).To(gomega.Equal(55))
				gomega.Expect(cloneIDs).To(gomega.HaveLen(3))

				var appCloneIDs []int
				appCloneIDs, err = image.AppClones()
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(appCloneIDs[0]).To(gomega.Equal(22))
				gomega.Expect(appCloneIDs).To(gomega.HaveLen(4))
			})

			ginkgo.It("should find image snapshots", func() {
				var snapshots []*ImageSnapshot
				snapshots, err = image.Snapshots()
				gomega.Expect(err).Should(gomega.BeNil())

				snapshot1 := snapshots[0]
				gomega.Expect(snapshot1.ID).To(gomega.Equal(42))
				gomega.Expect(snapshot1.Active).To(gomega.Equal("How"))
				gomega.Expect(snapshot1.Children).To(gomega.Equal("are"))
				gomega.Expect(snapshot1.Name).To(gomega.Equal("you"))

				time1 := time.Unix(int64(123456), 0)
				gomega.Expect(snapshot1.Date).To(gomega.Equal(&time1))
				gomega.Expect(snapshot1.Parent).To(gomega.Equal(1))
				gomega.Expect(snapshot1.Size).To(gomega.Equal(22))

				snapshot2 := snapshots[1]
				gomega.Expect(snapshot2.ID).To(gomega.Equal(23))
				gomega.Expect(snapshot2.Active).To(gomega.Equal(""))
				gomega.Expect(snapshot2.Children).To(gomega.Equal(""))
				gomega.Expect(snapshot2.Name).To(gomega.Equal(""))

				time2 := time.Unix(int64(26485), 0)
				gomega.Expect(snapshot2.Date).To(gomega.Equal(&time2))
				gomega.Expect(snapshot2.Parent).To(gomega.Equal(6))
				gomega.Expect(snapshot2.Size).To(gomega.Equal(228))

				gomega.Expect(snapshots).To(gomega.HaveLen(2))
			})
		})

		ginkgo.Context("when image has only ID", func() {
			var image *Image

			ginkgo.It("should create image", func() {
				image = CreateImageWithID(42)

				gomega.Expect(image.ID()).To(gomega.Equal(42))
			})

			ginkgo.Context("when image doesn't have given attribute", func() {
				ginkgo.BeforeEach(func() {
					image = CreateImageWithID(42)
					err = nil
				})

				ginkgo.It("should return that image doesn't have name", func() {
					_, err = image.Name()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have user", func() {
					_, err = image.User()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have group", func() {
					_, err = image.Group()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have permissions", func() {
					_, err = image.Permissions()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have TYPE", func() {
					_, err = image.Type()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have DISK_TYPE", func() {
					_, err = image.DiskType()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have PERSISTENT", func() {
					_, err = image.Persistent()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have REGTIME", func() {
					_, err = image.RegistrationTime()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have SOURCE", func() {
					_, err = image.Source()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have PATH", func() {
					_, err = image.Path()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have FSTYPE", func() {
					_, err = image.FileSystemType()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have SIZE", func() {
					_, err = image.Size()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have STATE", func() {
					_, err = image.State()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have RUNNING_VMS", func() {
					_, err = image.RunningVMs()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have DATASTORE", func() {
					_, err = image.Datastore()
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})

				ginkgo.It("should return that image doesn't have VMS", func() {
					var vms []int

					vms, err = image.VirtualMachines()
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(vms).Should(gomega.HaveLen(0))
				})

				ginkgo.It("should return that image doesn't have CLONES", func() {
					var clones []int

					clones, err = image.Clones()
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(clones).Should(gomega.HaveLen(0))
				})

				ginkgo.It("should return that image doesn't have APP_CLONES", func() {
					var appClones []int

					appClones, err = image.AppClones()
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(appClones).Should(gomega.HaveLen(0))
				})

				ginkgo.It("should return that image doesn't have SNAPSHOTS", func() {
					var snapshots []*ImageSnapshot

					snapshots, err = image.Snapshots()
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(snapshots).Should(gomega.HaveLen(0))
				})
			})
		})
	})
})
