package services_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/onego-project/onego/errors"

	"github.com/onego-project/onego/requests"
	"github.com/onego-project/onego/services"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/onego-project/onego"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/resources"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	imageAllocate               = "records/image/allocate"
	imageAllocateWrongDatastore = "records/image/allocateWrongDatastore"
	imageAllocateEmptyDatastore = "records/image/allocateEmptyDatastore"
	imageAllocateEmptyBlueprint = "records/image/allocateEmptyBlueprint"
	imageAllocateExisting       = "records/image/allocateExisting"

	imageClone               = "records/image/clone"
	imageCloneWrongDatastore = "records/image/cloneWrongDatastore"
	imageCloneExisting       = "records/image/cloneExisting"
	imageCloneNoID           = "records/image/cloneNoID"

	imageDelete        = "records/image/delete"
	imageDeleteWrongID = "records/image/deleteWrongID"
	imageDeleteNoID    = "records/image/deleteNoID"

	imageUpdateMerge        = "records/image/updateMerge"
	imageUpdateReplace      = "records/image/updateReplace"
	imageUpdateEmptyMerge   = "records/image/updateEmptyMerge"
	imageUpdateEmptyReplace = "records/image/updateEmptyReplace"
	imageUpdateNoUser       = "records/image/updateNoUser"
	imageUpdateUnknown      = "records/image/updateUnknown"

	imageChangeType        = "records/image/changeType"
	imageChangeTypeUnknown = "records/image/changeTypeUnknown"
	imageChangeTypeNoImage = "records/image/changeTypeNoImage"

	imageChmod              = "records/image/chmod"
	imagePermRequestDefault = "records/image/chmodPermReqDefault"
	imageChmodUnknown       = "records/image/chmodUnknown"
	imageChmodNoImage       = "records/image/chmodNoImage"

	imageChown               = "records/image/chown"
	imageOwnershipReqDefault = "records/image/chownDefault"
	imageChownUnknown        = "records/image/chownUnknown"
	imageChownNoImage        = "records/image/chownNoImage"

	imageRename        = "records/image/rename"
	imageRenameEmpty   = "records/image/renameEmpty"
	imageRenameUnknown = "records/image/renameUnknown"
	imageRenameNoImage = "records/image/renameNoImage"

	imageEnable         = "records/image/enable"
	imageDisable        = "records/image/disable"
	imageDisableUnknown = "records/image/disableUnknown"
	imageDisableNoImage = "records/image/disableNoImage"

	imagePersistent           = "records/image/persistent"
	imageNonpersistent        = "records/image/nonpersistent"
	imageNonpersistentUnknown = "records/image/nonpersistentUnknown"
	imageNonpersistentNoImage = "records/image/nonpersistentNoImage"

	imageRetrieveInfo        = "records/image/retrieveInfo"
	imageRetrieveInfoUnknown = "records/image/retrieveInfoUnknown"

	imageListAllPrimaryGroup = "records/image/listAllPrimaryGroup"
	imageListAllUser         = "records/image/listAllUser"
	imageListAllAll          = "records/image/listAllAll"
	imageListAllUserGroup    = "records/image/listAllUserGroup"

	imageListAllForUser        = "records/image/listAllForUser"
	imageListAllForUserUnknown = "records/image/listAllForUserUnknown"
	imageListAllForUserEmpty   = "records/image/listAllForUserEmpty"

	imageListPagination      = "records/image/listPagination"
	imageListPaginationWrong = "records/image/listPaginationWrong"

	imageListForUser        = "records/image/listForUser"
	imageListForUserUnknown = "records/image/listForUserUnknown"
	imageListForUserEmpty   = "records/image/listForUserEmpty"
)

var _ = ginkgo.Describe("Image Service", func() {
	var (
		recName string
		rec     *recorder.Recorder
		client  *onego.Client
		err     error
	)

	ginkgo.JustBeforeEach(func() {
		// Start recorder
		rec, err = recorder.New(recName)
		if err != nil {
			return
		}

		rec.SetMatcher(func(r *http.Request, i cassette.Request) bool {
			var b bytes.Buffer
			if _, err = b.ReadFrom(r.Body); err != nil {
				return false
			}
			r.Body = ioutil.NopCloser(&b)
			return cassette.DefaultMatcher(r, i) && (b.String() == "" || b.String() == i.Body)
		})

		// Create an HTTP client and inject our transport
		clientHTTP := &http.Client{
			Transport: rec, // Inject as transport!
		}

		// create onego client
		client = onego.CreateClient(endpoint, token, clientHTTP)
		if client == nil {
			err = errors.ErrNoClient
			return
		}
	})

	ginkgo.AfterEach(func() {
		rec.Stop()
	})

	ginkgo.Describe("allocate image", func() {
		var (
			image          *resources.Image
			datastore      *resources.Datastore
			imageID        int
			oneImage       *resources.Image
			imageBlueprint *blueprint.ImageBlueprint
		)

		ginkgo.BeforeEach(func() {
			imageBlueprint = blueprint.CreateAllocateImageBlueprint()
			imageBlueprint.SetElement("SOURCE", "/var/lib/one//datastores/101/5995f01c3c35883b297eed95a12a271b")
			imageBlueprint.SetElement("SIZE", "2252")
			imageBlueprint.SetElement("NAME", "myFirstImage")
		})

		ginkgo.Context("when image doesn't exist", func() {
			ginkgo.When("when attributes are set correctly", func() {
				ginkgo.BeforeEach(func() {
					recName = imageAllocate
				})

				ginkgo.It("should create new image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					datastore = resources.CreateDatastoreWithID(110)

					image, err = client.ImageService.Allocate(context.TODO(), imageBlueprint, *datastore)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(image).ShouldNot(gomega.BeNil())

					// check whether Image really exists in OpenNebula
					imageID, err = image.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneImage.Name()).To(gomega.Equal("myFirstImage"))
				})
			})

			ginkgo.When("when datastore is wrong", func() {
				ginkgo.BeforeEach(func() {
					recName = imageAllocateWrongDatastore
				})

				ginkgo.It("shouldn't create new image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					datastore = resources.CreateDatastoreWithID(33)

					image, err = client.ImageService.Allocate(context.TODO(), imageBlueprint, *datastore)
					gomega.Expect(err).To(gomega.HaveOccurred())
					gomega.Expect(image).Should(gomega.BeNil())
				})
			})

			ginkgo.When("when datastore is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = imageAllocateEmptyDatastore
				})

				ginkgo.It("shouldn't create new image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					datastore = &resources.Datastore{}

					image, err = client.ImageService.Allocate(context.TODO(), imageBlueprint, *datastore)
					gomega.Expect(err).To(gomega.HaveOccurred())
					gomega.Expect(image).Should(gomega.BeNil())
				})
			})

			ginkgo.When("when blueprint is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = imageAllocateEmptyBlueprint
				})

				ginkgo.It("shouldn't create new image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					datastore = resources.CreateDatastoreWithID(33)

					imageBlueprint = &blueprint.ImageBlueprint{}

					image, err = client.ImageService.Allocate(context.TODO(), imageBlueprint, *datastore)
					gomega.Expect(err).To(gomega.HaveOccurred())
					gomega.Expect(image).Should(gomega.BeNil())
				})
			})
		})

		ginkgo.Context("when image exists", func() {
			ginkgo.BeforeEach(func() {
				recName = imageAllocateExisting
			})

			ginkgo.It("should return that image already exists", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				image, err = client.ImageService.Allocate(context.TODO(), imageBlueprint, *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(image).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("clone image", func() {
		var (
			image     *resources.Image
			clone     *resources.Image
			oneImage  *resources.Image
			datastore *resources.Datastore
			cloneID   int
		)

		ginkgo.Context("when image doesn't exist", func() {
			ginkgo.When("when attributes are set correctly", func() {
				ginkgo.BeforeEach(func() {
					recName = imageClone

					image = resources.CreateImageWithID(378)
					if image == nil {
						err = errors.ErrNoImage
					}
				})

				ginkgo.It("should create new image clone", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					datastore = resources.CreateDatastoreWithID(110)

					clone, err = client.ImageService.Clone(context.TODO(), *image, "Golias", *datastore)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(clone).ShouldNot(gomega.BeNil())

					// check whether Image really exists in OpenNebula
					cloneID, err = clone.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), cloneID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneImage.Name()).To(gomega.Equal("Golias"))
				})
			})

			ginkgo.When("when datastore is wrong", func() {
				ginkgo.BeforeEach(func() {
					recName = imageCloneWrongDatastore

					image = resources.CreateImageWithID(378)
					if image == nil {
						err = errors.ErrNoImage
					}
				})

				ginkgo.It("shouldn't create new image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					datastore = resources.CreateDatastoreWithID(33)

					clone, err = client.ImageService.Clone(context.TODO(), *image, "asdf", *datastore)
					gomega.Expect(err).To(gomega.HaveOccurred())
					gomega.Expect(clone).Should(gomega.BeNil())
				})
			})
		})

		ginkgo.Context("when image exists", func() {
			ginkgo.BeforeEach(func() {
				recName = imageCloneExisting
			})

			ginkgo.It("should return that image already exists", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				clone, err = client.ImageService.Clone(context.TODO(), *image, "asdfgr", *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(clone).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when image is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageCloneNoID

				image = &resources.Image{}
			})

			ginkgo.It("should return that image has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				clone, err = client.ImageService.Clone(context.TODO(), *image, "asdfg", *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(clone).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("delete image", func() {
		var (
			image    *resources.Image
			oneImage *resources.Image
			imageID  int
		)

		ginkgo.Context("when image exists", func() {
			ginkgo.BeforeEach(func() {
				recName = imageDelete

				image = resources.CreateImageWithID(377)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.It("should delete image", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.Delete(context.TODO(), *image)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether image was really deleted in OpenNebula
				imageID, err = image.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(oneImage).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when image doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = imageDeleteWrongID

				image = resources.CreateImageWithID(110)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.It("should return that image with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.Delete(context.TODO(), *image)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when image is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageDeleteNoID

				image = &resources.Image{}
			})

			ginkgo.It("should return that image has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.Delete(context.TODO(), *image)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("update image", func() {
		var (
			image          *resources.Image
			oneImage       *resources.Image
			imageID        int
			imageBlueprint *blueprint.ImageBlueprint
		)

		ginkgo.Context("when image exists", func() {
			ginkgo.Context("when update data is not empty", func() {
				ginkgo.BeforeEach(func() {
					image = resources.CreateImageWithID(378)
					if image == nil {
						err = errors.ErrNoImage
						return
					}
				})

				ginkgo.When("when merge data of given image", func() {
					ginkgo.BeforeEach(func() {
						recName = imageUpdateMerge
					})

					ginkgo.It("should merge data of given image", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						imageBlueprint = blueprint.CreateUpdateImageBlueprint()
						if imageBlueprint == nil {
							err = errors.ErrNoImageBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						imageBlueprint.SetDescription("dummy")

						err = client.ImageService.Update(context.TODO(), *image, imageBlueprint, services.Merge)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						// check whether image data was really updated in OpenNebula
						imageID, err = image.ID()
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(oneImage).ShouldNot(gomega.BeNil())
						gomega.Expect(oneImage.Attribute("TEMPLATE/DESCRIPTION")).To(gomega.Equal("dummy"))
					})
				})

				ginkgo.When("when replace data of given image", func() {
					ginkgo.BeforeEach(func() {
						recName = imageUpdateReplace
					})

					ginkgo.It("should replace data of given image", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						imageBlueprint = blueprint.CreateUpdateImageBlueprint()
						if imageBlueprint == nil {
							err = errors.ErrNoImageBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						imageBlueprint.SetDiskType(resources.DiskTypeBlock)

						err = client.ImageService.Update(context.TODO(), *image, imageBlueprint, services.Replace)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						// check whether image data was really replaced in OpenNebula
						imageID, err = image.ID()
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(oneImage).ShouldNot(gomega.BeNil())
						gomega.Expect(oneImage.Attribute("TEMPLATE/DISK_TYPE")).To(gomega.Equal("BLOCK"))
					})
				})
			})

			ginkgo.Context("when update data is empty", func() {
				ginkgo.BeforeEach(func() {
					image = resources.CreateImageWithID(378)
					if image == nil {
						err = errors.ErrNoImage
						return
					}

					imageBlueprint = &blueprint.ImageBlueprint{}
				})

				ginkgo.When("when merge data of given image", func() {
					ginkgo.BeforeEach(func() {
						recName = imageUpdateEmptyMerge
					})

					ginkgo.It("should merge data of given image", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						err = client.ImageService.Update(context.TODO(), *image, imageBlueprint, services.Merge)
						gomega.Expect(err).To(gomega.HaveOccurred())
					})
				})

				ginkgo.When("when replace data of given image", func() {
					ginkgo.BeforeEach(func() {
						recName = imageUpdateEmptyReplace
					})

					ginkgo.It("should replace data of given image", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						err = client.ImageService.Update(context.TODO(), *image, imageBlueprint, services.Replace)
						gomega.Expect(err).To(gomega.HaveOccurred())
					})
				})
			})
		})

		ginkgo.Context("when image doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = imageUpdateUnknown

				image = resources.CreateImageWithID(110)
				if image == nil {
					err = errors.ErrNoImage
				}

				imageBlueprint = blueprint.CreateUpdateImageBlueprint()
				if imageBlueprint == nil {
					err = errors.ErrNoImageBlueprint
					return
				}
				imageBlueprint.SetDescription("dummy")
			})

			ginkgo.It("should return that image with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.Update(context.TODO(), *image, imageBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when image is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageUpdateNoUser

				imageBlueprint = blueprint.CreateUpdateImageBlueprint()
				if imageBlueprint == nil {
					err = errors.ErrNoImageBlueprint
					return
				}
				imageBlueprint.SetDescription("dummy")
			})

			ginkgo.It("should return that image has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.Update(context.TODO(), resources.Image{},
					imageBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("image change type", func() {
		var (
			image    *resources.Image
			oneImage *resources.Image
			imageID  int
		)

		ginkgo.Context("when image exists", func() {
			ginkgo.BeforeEach(func() {
				recName = imageChangeType

				image = resources.CreateImageWithID(378)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.It("should change type of given image", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.ChangeType(context.TODO(), *image, resources.ImageTypeOs)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether image type was really changed in OpenNebula
				imageID, err = image.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneImage).ShouldNot(gomega.BeNil())
				gomega.Expect(oneImage.Type()).To(gomega.Equal(resources.ImageTypeOs))
			})
		})

		ginkgo.Context("when image doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = imageChangeTypeUnknown

				image = resources.CreateImageWithID(110)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.It("should return that image with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.ChangeType(context.TODO(), *image, resources.ImageTypeDataBlock)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when image is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageChangeTypeNoImage

				image = &resources.Image{}
			})

			ginkgo.It("should return that image has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.ChangeType(context.TODO(), *image, resources.ImageTypeKernel)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("image chmod", func() {
		var (
			image       *resources.Image
			oneImage    *resources.Image
			imageID     int
			permRequest requests.PermissionRequest
		)

		ginkgo.Context("when image exists", func() {
			ginkgo.BeforeEach(func() {
				image = resources.CreateImageWithID(378)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.When("when permission request is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = imageChmod
				})

				ginkgo.It("should change permission of given image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					permRequest = requests.CreatePermissionRequestBuilder().Deny(requests.User,
						requests.Manage).Allow(requests.Other, requests.Admin).Build()

					err = client.ImageService.Chmod(context.TODO(), *image, permRequest)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chmod was really changed in OpenNebula
					imageID, err = image.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneImage).ShouldNot(gomega.BeNil())

					var perm *resources.Permissions
					perm, err = oneImage.Permissions()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					gomega.Expect(perm.User.Manage).To(gomega.Equal(false))
					gomega.Expect(perm.Other.Admin).To(gomega.Equal(true))
				})
			})

			ginkgo.When("when permission request is default", func() {
				ginkgo.BeforeEach(func() {
					recName = imagePermRequestDefault
				})

				ginkgo.It("should not change permissions of given image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ImageService.Chmod(context.TODO(), *image, requests.CreatePermissionRequestBuilder().Build())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when image doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = imageChmodUnknown

				image = resources.CreateImageWithID(110)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.It("should return that image with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				permRequest = requests.CreatePermissionRequestBuilder().Allow(requests.User,
					requests.Manage).Deny(requests.Other, requests.Admin).Build()

				err = client.ImageService.Chmod(context.TODO(), *image, permRequest)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when image is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageChmodNoImage

				image = &resources.Image{}
			})

			ginkgo.It("should return that image has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.Chmod(context.TODO(), *image, requests.PermissionRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("image chown", func() {
		var (
			image    *resources.Image
			oneImage *resources.Image
			imageID  int

			user         *resources.User
			group        *resources.Group
			ownershipReq requests.OwnershipRequest
		)

		ginkgo.Context("when image exists", func() {
			ginkgo.BeforeEach(func() {
				image = resources.CreateImageWithID(378)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.When("when ownership request is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = imageChown
				})

				ginkgo.It("should change owner of given image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					user = resources.CreateUserWithID(31)
					group = resources.CreateGroupWithID(120)

					ownershipReq = requests.CreateOwnershipRequestBuilder().User(*user).Group(*group).Build()

					err = client.ImageService.Chown(context.TODO(), *image, ownershipReq)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chown was really changed in OpenNebula
					imageID, err = image.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneImage).ShouldNot(gomega.BeNil())

					gomega.Expect(oneImage.User()).To(gomega.Equal(31))
					gomega.Expect(oneImage.Group()).To(gomega.Equal(120))
				})
			})

			ginkgo.When("when ownership request is default", func() {
				ginkgo.BeforeEach(func() {
					recName = imageOwnershipReqDefault
				})

				ginkgo.It("should not change permissions of given image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ImageService.Chown(context.TODO(), *image,
						requests.CreateOwnershipRequestBuilder().Build())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when image doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = imageChownUnknown

				image = resources.CreateImageWithID(110)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.It("should return that image with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.Chown(context.TODO(), *image, requests.OwnershipRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when image is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageChownNoImage

				image = &resources.Image{}
			})

			ginkgo.It("should return that image has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.Chown(context.TODO(), *image, requests.OwnershipRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("image rename", func() {
		var (
			image    *resources.Image
			oneImage *resources.Image
			imageID  int
		)

		ginkgo.Context("when image exists", func() {
			ginkgo.BeforeEach(func() {
				image = resources.CreateImageWithID(378)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.When("when new name is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = imageRename
				})

				ginkgo.It("should change name of given image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					// get image name
					oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneImage).ShouldNot(gomega.BeNil())

					newName := "monkey2"
					gomega.Expect(oneImage.Name()).NotTo(gomega.Equal(newName))

					// change name
					err = client.ImageService.Rename(context.TODO(), *image, newName)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether name was really changed in OpenNebula
					imageID, err = image.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneImage).ShouldNot(gomega.BeNil())

					gomega.Expect(oneImage.Name()).To(gomega.Equal("monkey2"))
				})
			})

			ginkgo.When("when new name is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = imageRenameEmpty
				})

				ginkgo.It("should not change name of given image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ImageService.Rename(context.TODO(), *image, "")
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when image doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = imageRenameUnknown

				image = resources.CreateImageWithID(110)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.It("should return that image with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.Rename(context.TODO(), *image, "mastodont")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when image is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageRenameNoImage

				image = &resources.Image{}
			})

			ginkgo.It("should return that image has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.Rename(context.TODO(), *image, "rex")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("image enable", func() {
		var (
			image    *resources.Image
			oneImage *resources.Image
			imageID  int
		)

		ginkgo.Context("when image exists", func() {
			ginkgo.BeforeEach(func() {
				image = resources.CreateImageWithID(378)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.When("when enable", func() {
				ginkgo.BeforeEach(func() {
					recName = imageEnable
				})

				ginkgo.It("should enable given image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ImageService.Enable(context.TODO(), *image)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chmod was really changed in OpenNebula
					imageID, err = image.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneImage).ShouldNot(gomega.BeNil())

					gomega.Expect(oneImage.State()).To(gomega.Equal(resources.ImageStateReady))
				})
			})

			ginkgo.When("when disable", func() {
				ginkgo.BeforeEach(func() {
					recName = imageDisable
				})

				ginkgo.It("should disable given image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ImageService.Disable(context.TODO(), *image)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chmod was really changed in OpenNebula
					imageID, err = image.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneImage).ShouldNot(gomega.BeNil())

					gomega.Expect(oneImage.State()).To(gomega.Equal(resources.ImageStateDisabled))
				})
			})
		})

		ginkgo.Context("when image doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = imageDisableUnknown

				image = resources.CreateImageWithID(110)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.It("should return that image with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.Disable(context.TODO(), *image)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when image is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageDisableNoImage

				image = &resources.Image{}
			})

			ginkgo.It("should return that image has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.Disable(context.TODO(), *image)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("image make persistent", func() {
		var (
			image    *resources.Image
			oneImage *resources.Image
			imageID  int
		)

		ginkgo.Context("when image exists", func() {
			ginkgo.BeforeEach(func() {
				image = resources.CreateImageWithID(378)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.When("when persistent", func() {
				ginkgo.BeforeEach(func() {
					recName = imagePersistent
				})

				ginkgo.It("should make persistent given image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ImageService.MakePersistent(context.TODO(), *image)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether persistence was really changed in OpenNebula
					imageID, err = image.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneImage).ShouldNot(gomega.BeNil())

					gomega.Expect(oneImage.Persistent()).To(gomega.Equal(true))
				})
			})

			ginkgo.When("when non-persistent", func() {
				ginkgo.BeforeEach(func() {
					recName = imageNonpersistent
				})

				ginkgo.It("should make non-persistent given image", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ImageService.MakeNonPersistent(context.TODO(), *image)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether persistence was really changed in OpenNebula
					imageID, err = image.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneImage, err = client.ImageService.RetrieveInfo(context.TODO(), imageID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneImage).ShouldNot(gomega.BeNil())

					gomega.Expect(oneImage.Persistent()).To(gomega.Equal(false))
				})
			})
		})

		ginkgo.Context("when image doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = imageNonpersistentUnknown

				image = resources.CreateImageWithID(110)
				if image == nil {
					err = errors.ErrNoImage
				}
			})

			ginkgo.It("should return that image with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.MakeNonPersistent(context.TODO(), *image)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when image is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageNonpersistentNoImage

				image = &resources.Image{}
			})

			ginkgo.It("should return that image has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ImageService.MakeNonPersistent(context.TODO(), *image)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("image retrieve info", func() {
		var image *resources.Image

		ginkgo.Context("when image exists", func() {
			ginkgo.BeforeEach(func() {
				recName = imageRetrieveInfo
			})

			ginkgo.It("should return image with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				image, err = client.ImageService.RetrieveInfo(context.TODO(), 378)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(image).ShouldNot(gomega.BeNil())
				gomega.Expect(image.ID()).To(gomega.Equal(378))
				gomega.Expect(image.Name()).To(gomega.Equal("monkey"))
			})
		})

		ginkgo.Context("when image doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = imageRetrieveInfoUnknown
			})

			ginkgo.It("should return that given image doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				image, err = client.ImageService.RetrieveInfo(context.TODO(), 110)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(image).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("image list all", func() {
		var images []*resources.Image

		ginkgo.Context("when ownership filter is not empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageListAllPrimaryGroup
			})

			ginkgo.It("should return list of all images with full info belongs to primary group", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				images, err = client.ImageService.ListAll(context.TODO(), services.OwnershipFilterPrimaryGroup)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(images).ShouldNot(gomega.BeNil())
				gomega.Expect(images).To(gomega.HaveLen(2))
			})
		})

		ginkgo.Context("when ownership filter is not empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageListAllUser
			})

			ginkgo.It("should return list of all images with full info belongs to the user", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				images, err = client.ImageService.ListAll(context.TODO(), services.OwnershipFilterUser)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(images).ShouldNot(gomega.BeNil())
				gomega.Expect(images).To(gomega.HaveLen(4))
				gomega.Expect(images[0].ID()).To(gomega.Equal(0))
				gomega.Expect(images[1].ID()).To(gomega.Equal(375))
				gomega.Expect(images[2].ID()).To(gomega.Equal(376))
				gomega.Expect(images[3].ID()).To(gomega.Equal(379))
			})
		})

		ginkgo.Context("when ownership filter is set to ALl", func() {
			ginkgo.BeforeEach(func() {
				recName = imageListAllAll
			})

			ginkgo.It("should return list of all images with full info belongs to all", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				images, err = client.ImageService.ListAll(context.TODO(), services.OwnershipFilterAll)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(images).ShouldNot(gomega.BeNil())
				gomega.Expect(images).To(gomega.HaveLen(5))
			})
		})

		ginkgo.Context("when ownership filter is set to UserGroup", func() {
			ginkgo.BeforeEach(func() {
				recName = imageListAllUserGroup
			})

			ginkgo.It("should return list of all images with full info belongs to the user and any of his groups", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				images, err = client.ImageService.ListAll(context.TODO(), services.OwnershipFilterUserGroup)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(images).ShouldNot(gomega.BeNil())
				gomega.Expect(images).To(gomega.HaveLen(5))
			})
		})
	})

	ginkgo.Describe("image list all for user", func() {
		var images []*resources.Image

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = imageListAllForUser
			})

			ginkgo.It("should return images with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				images, err = client.ImageService.ListAllForUser(context.TODO(), *resources.CreateUserWithID(31))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(images).ShouldNot(gomega.BeNil())
				gomega.Expect(images).To(gomega.HaveLen(1))
				gomega.Expect(images[0].ID()).To(gomega.Equal(378))
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = imageListAllForUserUnknown
			})

			ginkgo.It("should return empty list of images (length 0)", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				images, err = client.ImageService.ListAllForUser(context.TODO(), *resources.CreateUserWithID(42))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(images).Should(gomega.Equal(make([]*resources.Image, 0)))
				gomega.Expect(images).Should(gomega.HaveLen(0))
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageListAllForUserEmpty
			})

			ginkgo.It("should return that user doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				images, err = client.ImageService.ListAllForUser(context.TODO(), resources.User{})
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(images).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("List methods with pagination", func() {
		var (
			images []*resources.Image
		)

		ginkgo.Context("pagination ok", func() {
			ginkgo.BeforeEach(func() {
				recName = imageListPagination
			})

			ginkgo.It("should return images with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				images, err = client.ImageService.List(context.TODO(), 3, 2, services.OwnershipFilterUserGroup)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(images).ShouldNot(gomega.BeNil())
				gomega.Expect(images).To(gomega.HaveLen(2))
				gomega.Expect(images[0].ID()).To(gomega.Equal(380))
				gomega.Expect(images[1].ID()).To(gomega.Equal(381))
			})
		})

		ginkgo.Context("pagination wrong", func() {
			ginkgo.BeforeEach(func() {
				recName = imageListPaginationWrong
			})
		})

		ginkgo.It("should return that pagination is wrong", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

			images, err = client.ImageService.List(context.TODO(), -2, -2, services.OwnershipFilterPrimaryGroup)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(images).Should(gomega.BeNil())
		})
	})

	ginkgo.Describe("image list for user", func() {
		var images []*resources.Image

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = imageListForUser
			})

			ginkgo.It("should return images with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				images, err = client.ImageService.ListForUser(context.TODO(), *resources.CreateUserWithID(0), 2, 2)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(images).ShouldNot(gomega.BeNil())
				gomega.Expect(images).To(gomega.HaveLen(2))
				gomega.Expect(images[0].ID()).To(gomega.Equal(379))
				gomega.Expect(images[1].ID()).To(gomega.Equal(380))
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = imageListForUserUnknown
			})

			ginkgo.It("should return empty list of images (length 0)", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				images, err = client.ImageService.ListForUser(context.TODO(), *resources.CreateUserWithID(88), 2, 2)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(images).Should(gomega.Equal(make([]*resources.Image, 0)))
				gomega.Expect(images).Should(gomega.HaveLen(0))
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = imageListForUserEmpty
			})

			ginkgo.It("should return that user doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				images, err = client.ImageService.ListForUser(context.TODO(), resources.User{}, 2, 2)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(images).Should(gomega.BeNil())
			})
		})
	})
})
