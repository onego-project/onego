package services_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/onego-project/onego/errors"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/onego-project/onego"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/requests"
	"github.com/onego-project/onego/resources"
	"github.com/onego-project/onego/services"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	datastoreAllocate             = "records/datastore/allocate"
	datastoreAllocateWrongCluster = "records/datastore/allocateWrongCluster"
	datastoreAllocateExisting     = "records/datastore/allocateExisting"

	datastoreDelete        = "records/datastore/delete"
	datastoreDeleteWrongID = "records/datastore/deleteWrongID"
	datastoreDeleteNoID    = "records/datastore/deleteNoID"

	datastoreUpdateMerge        = "records/datastore/updateMerge"
	datastoreUpdateReplace      = "records/datastore/updateReplace"
	datastoreUpdateEmptyMerge   = "records/datastore/updateEmptyMerge"
	datastoreUpdateEmptyReplace = "records/datastore/updateEmptyReplace"
	datastoreUpdateNoUser       = "records/datastore/updateNoUser"
	datastoreUpdateUnknown      = "records/datastore/updateUnknown"

	datastoreChmod              = "records/datastore/chmod"
	datastorePermRequestDefault = "records/datastore/chmodPermReqDefault"
	datastoreChmodUnknown       = "records/datastore/chmodUnknown"
	datastoreChmodNoDatastore   = "records/datastore/chmodNoDatastore"

	datastoreChown               = "records/datastore/chown"
	datastoreOwnershipReqDefault = "records/datastore/chownDefault"
	datastoreChownUnknown        = "records/datastore/chownUnknown"
	datastoreChownNoDatastore    = "records/datastore/chownNoDatastore"

	datastoreRename            = "records/datastore/rename"
	datastoreRenameEmpty       = "records/datastore/renameEmpty"
	datastoreRenameUnknown     = "records/datastore/renameUnknown"
	datastoreRenameNoDatastore = "records/datastore/renameNoDatastore"

	datastoreEnable             = "records/datastore/enable"
	datastoreDisable            = "records/datastore/disable"
	datastoreDisableUnknown     = "records/datastore/disableUnknown"
	datastoreDisableNoDatastore = "records/datastore/disableNoDatastore"

	datastoreRetrieveInfo        = "records/datastore/retrieveInfo"
	datastoreRetrieveInfoUnknown = "records/datastore/retrieveInfoUnknown"

	datastoreList = "records/datastore/list"
)

var _ = ginkgo.Describe("Datastore Service", func() {
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

	ginkgo.Describe("allocate datastore", func() {
		var (
			datastore          *resources.Datastore
			oneDatastore       *resources.Datastore
			cluster            *resources.Cluster
			datastoreID        int
			datastoreBlueprint *blueprint.DatastoreBlueprint
		)

		ginkgo.BeforeEach(func() {
			datastoreBlueprint = blueprint.CreateAllocateDatastoreBlueprint()
			datastoreBlueprint.SetElement("NAME", "hot_dog")
			datastoreBlueprint.SetElement("DS_MAD", "dummy")
			datastoreBlueprint.SetElement("TM_MAD", "dummy")
		})

		ginkgo.Context("when datastore doesn't exist", func() {
			ginkgo.When("when attributes are set correctly", func() {
				ginkgo.BeforeEach(func() {
					recName = datastoreAllocate
				})

				ginkgo.It("should create new datastore", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					cluster = resources.CreateClusterWithID(0)

					datastore, err = client.DatastoreService.Allocate(context.TODO(), datastoreBlueprint, *cluster)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(datastore).ShouldNot(gomega.BeNil())

					// check whether Datastore really exists in OpenNebula
					datastoreID, err = datastore.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneDatastore, err = client.DatastoreService.RetrieveInfo(context.TODO(), datastoreID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneDatastore.Name()).To(gomega.Equal("hot_dog"))
				})
			})

			ginkgo.When("when cluster is wrong", func() {
				ginkgo.BeforeEach(func() {
					recName = datastoreAllocateWrongCluster
				})

				ginkgo.It("shouldn't create new user", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					cluster = resources.CreateClusterWithID(33)

					datastore, err = client.DatastoreService.Allocate(context.TODO(), datastoreBlueprint, *cluster)
					gomega.Expect(err).To(gomega.HaveOccurred())
					gomega.Expect(datastore).Should(gomega.BeNil())
				})
			})
		})

		ginkgo.Context("when datastore exists", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreAllocateExisting
			})

			ginkgo.It("should return that datastore with name hot_dog already exists", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				datastore, err = client.DatastoreService.Allocate(context.TODO(), datastoreBlueprint, *cluster)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(datastore).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("delete datastore", func() {
		var (
			datastore    *resources.Datastore
			oneDatastore *resources.Datastore
			datastoreID  int
		)

		ginkgo.Context("when datastore exists", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreDelete

				datastore = resources.CreateDatastoreWithID(105)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should delete datastore", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.DatastoreService.Delete(context.TODO(), *datastore)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether datastore was really deleted in OpenNebula
				datastoreID, err = datastore.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneDatastore, err = client.DatastoreService.RetrieveInfo(context.TODO(), datastoreID)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(oneDatastore).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when datastore doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreDeleteWrongID

				datastore = resources.CreateDatastoreWithID(110)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should return that datastore with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.DatastoreService.Delete(context.TODO(), *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when datastore is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreDeleteNoID

				datastore = &resources.Datastore{}
			})

			ginkgo.It("should return that datastore has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.DatastoreService.Delete(context.TODO(), *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("update datastore", func() {
		var (
			datastore          *resources.Datastore
			oneDatastore       *resources.Datastore
			datastoreID        int
			datastoreBlueprint *blueprint.DatastoreBlueprint
		)

		ginkgo.Context("when datastore exists", func() {
			ginkgo.Context("when update data is not empty", func() {
				ginkgo.BeforeEach(func() {
					datastore = resources.CreateDatastoreWithID(104)
					if datastore == nil {
						err = errors.ErrNoDatastore
						return
					}
				})

				ginkgo.When("when merge data of given datastore", func() {
					ginkgo.BeforeEach(func() {
						recName = datastoreUpdateMerge
					})

					ginkgo.It("should merge data of given datastore", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						datastoreBlueprint = blueprint.CreateUpdateDatastoreBlueprint()
						if datastoreBlueprint == nil {
							err = errors.ErrNoDatastoreBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						datastoreBlueprint.SetDsMad("dummy")

						err = client.DatastoreService.Update(context.TODO(), *datastore, datastoreBlueprint, services.Merge)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						// check whether datastore data was really updated in OpenNebula
						datastoreID, err = datastore.ID()
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						oneDatastore, err = client.DatastoreService.RetrieveInfo(context.TODO(), datastoreID)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(oneDatastore).ShouldNot(gomega.BeNil())
						gomega.Expect(oneDatastore.Attribute("TEMPLATE/DS_MAD")).To(gomega.Equal("dummy"))
					})
				})

				ginkgo.When("when replace data of given datastore", func() {
					ginkgo.BeforeEach(func() {
						recName = datastoreUpdateReplace
					})

					ginkgo.It("should replace data of given datastore", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						datastoreBlueprint = blueprint.CreateUpdateDatastoreBlueprint()
						if datastoreBlueprint == nil {
							err = errors.ErrNoDatastoreBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						datastoreBlueprint.SetType(resources.DatastoreTypeSystem)

						err = client.DatastoreService.Update(context.TODO(), *datastore, datastoreBlueprint, services.Replace)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						// check whether datastore data was really replaced in OpenNebula
						datastoreID, err = datastore.ID()
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						oneDatastore, err = client.DatastoreService.RetrieveInfo(context.TODO(), datastoreID)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(oneDatastore).ShouldNot(gomega.BeNil())
						gomega.Expect(oneDatastore.Attribute("TEMPLATE/TYPE")).To(gomega.Equal("SYSTEM_DS"))
					})
				})
			})

			ginkgo.Context("when update data is empty", func() {
				ginkgo.BeforeEach(func() {
					datastore = resources.CreateDatastoreWithID(104)
					if datastore == nil {
						err = errors.ErrNoDatastore
						return
					}

					datastoreBlueprint = &blueprint.DatastoreBlueprint{}
				})

				ginkgo.When("when merge data of given datastore", func() {
					ginkgo.BeforeEach(func() {
						recName = datastoreUpdateEmptyMerge
					})

					ginkgo.It("should merge data of given datastore", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						err = client.DatastoreService.Update(context.TODO(), *datastore, datastoreBlueprint, services.Merge)
						gomega.Expect(err).To(gomega.HaveOccurred())
					})
				})

				ginkgo.When("when replace data of given datastore", func() {
					ginkgo.BeforeEach(func() {
						recName = datastoreUpdateEmptyReplace
					})

					ginkgo.It("should replace data of given datastore", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						err = client.DatastoreService.Update(context.TODO(), *datastore, datastoreBlueprint, services.Replace)
						gomega.Expect(err).To(gomega.HaveOccurred())
					})
				})
			})
		})

		ginkgo.Context("when datastore doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreUpdateUnknown

				datastore = resources.CreateDatastoreWithID(110)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}

				datastoreBlueprint = blueprint.CreateUpdateDatastoreBlueprint()
				if datastoreBlueprint == nil {
					err = errors.ErrNoDatastoreBlueprint
					return
				}
				datastoreBlueprint.SetType(resources.DatastoreTypeSystem)
			})

			ginkgo.It("should return that datastore with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.DatastoreService.Update(context.TODO(), *datastore, datastoreBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when datastore is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreUpdateNoUser

				datastoreBlueprint = blueprint.CreateUpdateDatastoreBlueprint()
				if datastoreBlueprint == nil {
					err = errors.ErrNoDatastoreBlueprint
					return
				}
				datastoreBlueprint.SetType(resources.DatastoreTypeSystem)
			})

			ginkgo.It("should return that datastore has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.DatastoreService.Update(context.TODO(), resources.Datastore{},
					datastoreBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("datastore chmod", func() {
		var (
			datastore    *resources.Datastore
			oneDatastore *resources.Datastore
			datastoreID  int
		)

		ginkgo.Context("when datastore exists", func() {
			ginkgo.BeforeEach(func() {
				datastore = resources.CreateDatastoreWithID(106)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.When("when permission request is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = datastoreChmod
				})

				ginkgo.It("should change permission of given datastore", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					permRequest := requests.CreatePermissionRequestBuilder().Deny(requests.User,
						requests.Manage).Allow(requests.Other, requests.Admin).Build()

					err = client.DatastoreService.Chmod(context.TODO(), *datastore, permRequest)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chmod was really changed in OpenNebula
					datastoreID, err = datastore.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneDatastore, err = client.DatastoreService.RetrieveInfo(context.TODO(), datastoreID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneDatastore).ShouldNot(gomega.BeNil())

					var perm *resources.Permissions
					perm, err = oneDatastore.Permissions()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					gomega.Expect(perm.User.Manage).To(gomega.Equal(false))
					gomega.Expect(perm.Other.Admin).To(gomega.Equal(true))
				})
			})

			ginkgo.When("when permission request is default", func() {
				ginkgo.BeforeEach(func() {
					recName = datastorePermRequestDefault
				})

				ginkgo.It("should not change permissions of given datastore", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.DatastoreService.Chmod(context.TODO(), *datastore, requests.CreatePermissionRequestBuilder().Build())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when datastore doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreChmodUnknown

				datastore = resources.CreateDatastoreWithID(110)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should return that datastore with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				permRequest := requests.CreatePermissionRequestBuilder().Allow(requests.User,
					requests.Manage).Deny(requests.Other, requests.Admin).Build()

				err = client.DatastoreService.Chmod(context.TODO(), *datastore, permRequest)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when datastore is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreChmodNoDatastore
			})

			ginkgo.It("should return that datastore has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				permRequest := requests.CreatePermissionRequestBuilder().Allow(requests.User,
					requests.Manage).Deny(requests.Other, requests.Admin).Build()

				err = client.DatastoreService.Chmod(context.TODO(), *datastore, permRequest)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("datastore chown", func() {
		var (
			datastore    *resources.Datastore
			oneDatastore *resources.Datastore
			datastoreID  int
		)

		ginkgo.Context("when datastore exists", func() {
			ginkgo.BeforeEach(func() {
				datastore = resources.CreateDatastoreWithID(106)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.When("when ownership request is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = datastoreChown
				})

				ginkgo.It("should change owner of given datastore", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					user := resources.CreateUserWithID(31)
					group := resources.CreateGroupWithID(120)

					ownershipReq := requests.CreateOwnershipRequestBuilder().User(*user).Group(*group).Build()

					err = client.DatastoreService.Chown(context.TODO(), *datastore, ownershipReq)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chown was really changed in OpenNebula
					datastoreID, err = datastore.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneDatastore, err = client.DatastoreService.RetrieveInfo(context.TODO(), datastoreID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneDatastore).ShouldNot(gomega.BeNil())

					gomega.Expect(oneDatastore.User()).To(gomega.Equal(31))
					gomega.Expect(oneDatastore.Group()).To(gomega.Equal(120))
				})
			})

			ginkgo.When("when ownership request is default", func() {
				ginkgo.BeforeEach(func() {
					recName = datastoreOwnershipReqDefault
				})

				ginkgo.It("should not change permissions of given datastore", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.DatastoreService.Chown(context.TODO(), *datastore,
						requests.CreateOwnershipRequestBuilder().Build())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when datastore doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreChownUnknown

				datastore = resources.CreateDatastoreWithID(110)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should return that datastore with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				user := resources.CreateUserWithID(31)
				group := resources.CreateGroupWithID(120)

				ownershipReq := requests.CreateOwnershipRequestBuilder().User(*user).Group(*group).Build()

				err = client.DatastoreService.Chown(context.TODO(), *datastore, ownershipReq)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when datastore is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreChownNoDatastore
			})

			ginkgo.It("should return that datastore has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				user := resources.CreateUserWithID(31)
				group := resources.CreateGroupWithID(120)

				ownershipReq := requests.CreateOwnershipRequestBuilder().User(*user).Group(*group).Build()

				err = client.DatastoreService.Chown(context.TODO(), *datastore, ownershipReq)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("datastore rename", func() {
		var (
			datastore    *resources.Datastore
			oneDatastore *resources.Datastore
			datastoreID  int
		)

		ginkgo.Context("when datastore exists", func() {
			ginkgo.BeforeEach(func() {
				datastore = resources.CreateDatastoreWithID(106)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.When("when new name is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = datastoreRename
				})

				ginkgo.It("should change name of given datastore", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.DatastoreService.Rename(context.TODO(), *datastore, "monkey")
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chmod was really changed in OpenNebula
					datastoreID, err = datastore.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneDatastore, err = client.DatastoreService.RetrieveInfo(context.TODO(), datastoreID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneDatastore).ShouldNot(gomega.BeNil())

					gomega.Expect(oneDatastore.Name()).To(gomega.Equal("monkey"))
				})
			})

			ginkgo.When("when new name is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = datastoreRenameEmpty
				})

				ginkgo.It("should not change name of given datastore", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.DatastoreService.Rename(context.TODO(), *datastore, "")
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when datastore doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreRenameUnknown

				datastore = resources.CreateDatastoreWithID(110)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should return that datastore with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.DatastoreService.Rename(context.TODO(), *datastore, "mastodont")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when datastore is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreRenameNoDatastore
			})

			ginkgo.It("should return that datastore has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.DatastoreService.Rename(context.TODO(), *datastore, "rex")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("datastore enable", func() {
		var (
			datastore    *resources.Datastore
			oneDatastore *resources.Datastore
			datastoreID  int
		)

		ginkgo.Context("when datastore exists", func() {
			ginkgo.BeforeEach(func() {
				datastore = resources.CreateDatastoreWithID(106)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.When("when enable", func() {
				ginkgo.BeforeEach(func() {
					recName = datastoreEnable
				})

				ginkgo.It("should enable given datastore", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.DatastoreService.Enable(context.TODO(), *datastore)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chmod was really changed in OpenNebula
					datastoreID, err = datastore.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneDatastore, err = client.DatastoreService.RetrieveInfo(context.TODO(), datastoreID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneDatastore).ShouldNot(gomega.BeNil())

					gomega.Expect(oneDatastore.State()).To(gomega.Equal(resources.DatastoreStateReady))
				})
			})

			ginkgo.When("when disable", func() {
				ginkgo.BeforeEach(func() {
					recName = datastoreDisable
				})

				ginkgo.It("should disable given datastore", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.DatastoreService.Disable(context.TODO(), *datastore)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chmod was really changed in OpenNebula
					datastoreID, err = datastore.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneDatastore, err = client.DatastoreService.RetrieveInfo(context.TODO(), datastoreID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneDatastore).ShouldNot(gomega.BeNil())

					gomega.Expect(oneDatastore.State()).To(gomega.Equal(resources.DatastoreStateDisabled))
				})
			})
		})

		ginkgo.Context("when datastore doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreDisableUnknown

				datastore = resources.CreateDatastoreWithID(110)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should return that datastore with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.DatastoreService.Disable(context.TODO(), *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when datastore is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreDisableNoDatastore
			})

			ginkgo.It("should return that datastore has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.DatastoreService.Disable(context.TODO(), *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("datastore retrieve info", func() {
		var datastore *resources.Datastore

		ginkgo.Context("when datastore exists", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreRetrieveInfo
			})

			ginkgo.It("should return datastore with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				datastore, err = client.DatastoreService.RetrieveInfo(context.TODO(), 106)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(datastore).ShouldNot(gomega.BeNil())
				gomega.Expect(datastore.ID()).To(gomega.Equal(106))
				gomega.Expect(datastore.Name()).To(gomega.Equal("monkey"))
			})
		})

		ginkgo.Context("when datastore doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = datastoreRetrieveInfoUnknown
			})

			ginkgo.It("should return that given datastore doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				datastore, err = client.DatastoreService.RetrieveInfo(context.TODO(), 110)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(datastore).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("datastore list", func() {
		ginkgo.BeforeEach(func() {
			recName = datastoreList
		})

		ginkgo.It("should return an array of datastores with full info", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

			var datastores []*resources.Datastore

			datastores, err = client.DatastoreService.List(context.TODO())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(datastores).ShouldNot(gomega.BeNil())
		})
	})
})
