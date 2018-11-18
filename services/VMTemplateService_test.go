package services_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/onego-project/onego"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/errors"
	"github.com/onego-project/onego/requests"
	"github.com/onego-project/onego/resources"
	"github.com/onego-project/onego/services"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	vmTemplateAllocate               = "records/vmTemplate/allocate"
	vmTemplateAllocateEmptyBlueprint = "records/vmTemplate/allocateEmptyBlueprint"
	vmTemplateAllocateExisting       = "records/vmTemplate/allocateExisting"

	vmTemplateClone         = "records/vmTemplate/clone"
	vmTemplateCloneExisting = "records/vmTemplate/cloneExisting"
	vmTemplateCloneNoID     = "records/vmTemplate/cloneNoID"

	vmTemplateDelete        = "records/vmTemplate/delete"
	vmTemplateDeleteWrongID = "records/vmTemplate/deleteWrongID"
	vmTemplateDeleteNoID    = "records/vmTemplate/deleteNoID"

	vmTemplateUpdateMerge        = "records/vmTemplate/updateMerge"
	vmTemplateUpdateReplace      = "records/vmTemplate/updateReplace"
	vmTemplateUpdateEmptyMerge   = "records/vmTemplate/updateEmptyMerge"
	vmTemplateUpdateEmptyReplace = "records/vmTemplate/updateEmptyReplace"
	vmTemplateUpdateNoUser       = "records/vmTemplate/updateNoUser"
	vmTemplateUpdateUnknown      = "records/vmTemplate/updateUnknown"

	vmTemplateChmod              = "records/vmTemplate/chmod"
	vmTemplatePermRequestDefault = "records/vmTemplate/chmodPermReqDefault"
	vmTemplateChmodUnknown       = "records/vmTemplate/chmodUnknown"
	vmTemplateChmodNoVMTemplate  = "records/vmTemplate/chmodNoVMTemplate"

	vmTemplateChown               = "records/vmTemplate/chown"
	vmTemplateOwnershipReqDefault = "records/vmTemplate/chownDefault"
	vmTemplateChownUnknown        = "records/vmTemplate/chownUnknown"
	vmTemplateChownNoVMTemplate   = "records/vmTemplate/chownNoVMTemplate"

	vmTemplateRename             = "records/vmTemplate/rename"
	vmTemplateRenameEmpty        = "records/vmTemplate/renameEmpty"
	vmTemplateRenameUnknown      = "records/vmTemplate/renameUnknown"
	vmTemplateRenameNoVMTemplate = "records/vmTemplate/renameNoVMTemplate"

	vmTemplateRetrieveInfo        = "records/vmTemplate/retrieveInfo"
	vmTemplateRetrieveInfoUnknown = "records/vmTemplate/retrieveInfoUnknown"

	vmTemplateListAllPrimaryGroup = "records/vmTemplate/listAllPrimaryGroup"
	vmTemplateListAllUser         = "records/vmTemplate/listAllUser"
	vmTemplateListAllAll          = "records/vmTemplate/listAllAll"
	vmTemplateListAllUserGroup    = "records/vmTemplate/listAllUserGroup"

	vmTemplateListAllForUser        = "records/vmTemplate/listAllForUser"
	vmTemplateListAllForUserUnknown = "records/vmTemplate/listAllForUserUnknown"
	vmTemplateListAllForUserEmpty   = "records/vmTemplate/listAllForUserEmpty"

	vmTemplateListPagination      = "records/vmTemplate/listPagination"
	vmTemplateListPaginationWrong = "records/vmTemplate/listPaginationWrong"

	vmTemplateListForUser        = "records/vmTemplate/listForUser"
	vmTemplateListForUserUnknown = "records/vmTemplate/listForUserUnknown"
	vmTemplateListForUserEmpty   = "records/vmTemplate/listForUserEmpty"
)

var _ = ginkgo.Describe("VMTemplate Service", func() {
	var (
		recName string
		rec     *recorder.Recorder
		client  *onego.Client
		err     error
	)

	var existingVMTemplateID = 325
	var nonExistingVMTemplateID = 420
	var allocatedDeletedVMTemplateID = 326

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

	ginkgo.Describe("allocate vmTemplate", func() {
		var (
			vmTemplate          *resources.VMTemplate
			vmTemplateID        int
			oneVMTemplate       *resources.VMTemplate
			vmTemplateBlueprint *blueprint.VMTemplateBlueprint
		)

		ginkgo.BeforeEach(func() {
			vmTemplateBlueprint = blueprint.CreateAllocateVMTemplateBlueprint()
			vmTemplateBlueprint.SetElement("NAME", "hello")
		})

		ginkgo.Context("when vmTemplate doesn't exist", func() {
			ginkgo.When("when attributes are set correctly", func() {
				ginkgo.BeforeEach(func() {
					recName = vmTemplateAllocate
				})

				ginkgo.It("should create new vmTemplate", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					vmTemplate, err = client.VMTemplateService.Allocate(context.TODO(), vmTemplateBlueprint)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(vmTemplate).ShouldNot(gomega.BeNil())

					// check whether VMTemplate really exists in OpenNebula
					vmTemplateID, err = vmTemplate.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneVMTemplate, err = client.VMTemplateService.RetrieveInfo(context.TODO(), vmTemplateID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVMTemplate.Name()).To(gomega.Equal("hello"))
				})
			})

			ginkgo.When("when blueprint is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = vmTemplateAllocateEmptyBlueprint
				})

				ginkgo.It("shouldn't create new vmTemplate", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					vmTemplate, err = client.VMTemplateService.Allocate(context.TODO(),
						&blueprint.VMTemplateBlueprint{})
					gomega.Expect(err).To(gomega.HaveOccurred())
					gomega.Expect(vmTemplate).Should(gomega.BeNil())
				})
			})
		})

		ginkgo.Context("when vmTemplate exists", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateAllocateExisting
			})

			ginkgo.It("should return that vmTemplate already exists", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplate, err = client.VMTemplateService.Allocate(context.TODO(), vmTemplateBlueprint)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(vmTemplate).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("clone vmTemplate", func() {
		var (
			vmTemplate    *resources.VMTemplate
			clone         *resources.VMTemplate
			oneVMTemplate *resources.VMTemplate
			cloneID       int
		)

		ginkgo.Context("when vmTemplate doesn't exist", func() {
			ginkgo.When("when attributes are set correctly", func() {
				ginkgo.BeforeEach(func() {
					recName = vmTemplateClone

					vmTemplate = resources.CreateVMTemplateWithID(existingVMTemplateID)
					if vmTemplate == nil {
						err = errors.ErrNoVMTemplate
					}
				})

				ginkgo.It("should create new vmTemplate clone", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					clone, err = client.VMTemplateService.Clone(context.TODO(), *vmTemplate, "Golias", true)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(clone).ShouldNot(gomega.BeNil())

					// check whether VMTemplate really exists in OpenNebula
					cloneID, err = clone.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneVMTemplate, err = client.VMTemplateService.RetrieveInfo(context.TODO(), cloneID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVMTemplate.Name()).To(gomega.Equal("Golias"))
				})
			})
		})

		ginkgo.Context("when vmTemplate exists", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateCloneExisting
			})

			ginkgo.It("should return that vmTemplate already exists", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				clone, err = client.VMTemplateService.Clone(context.TODO(), *vmTemplate, "test-delete", true)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(clone).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when vmTemplate is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateCloneNoID

				vmTemplate = &resources.VMTemplate{}
			})

			ginkgo.It("should return that vmTemplate has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				clone, err = client.VMTemplateService.Clone(context.TODO(), *vmTemplate, "asdfg", true)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(clone).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("delete vmTemplate", func() {
		var (
			vmTemplate    *resources.VMTemplate
			oneVMTemplate *resources.VMTemplate
			vmTemplateID  int
		)

		ginkgo.Context("when vmTemplate exists", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateDelete

				vmTemplate = resources.CreateVMTemplateWithID(allocatedDeletedVMTemplateID)
				if vmTemplate == nil {
					err = errors.ErrNoVMTemplate
				}
			})

			ginkgo.It("should delete vmTemplate", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VMTemplateService.Delete(context.TODO(), *vmTemplate, true)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether vmTemplate was really deleted in OpenNebula
				vmTemplateID, err = vmTemplate.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneVMTemplate, err = client.VMTemplateService.RetrieveInfo(context.TODO(), vmTemplateID)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(oneVMTemplate).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when vmTemplate doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateDeleteWrongID

				vmTemplate = resources.CreateVMTemplateWithID(nonExistingVMTemplateID)
				if vmTemplate == nil {
					err = errors.ErrNoVMTemplate
				}
			})

			ginkgo.It("should return that vmTemplate with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VMTemplateService.Delete(context.TODO(), *vmTemplate, true)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when vmTemplate is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateDeleteNoID
			})

			ginkgo.It("should return that vmTemplate has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VMTemplateService.Delete(context.TODO(), resources.VMTemplate{}, true)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("update vmTemplate", func() {
		var (
			vmTemplate          *resources.VMTemplate
			oneVMTemplate       *resources.VMTemplate
			vmTemplateID        int
			vmTemplateBlueprint *blueprint.VMTemplateBlueprint
		)

		ginkgo.Context("when vmTemplate exists", func() {
			ginkgo.Context("when update data is not empty", func() {
				ginkgo.BeforeEach(func() {
					vmTemplate = resources.CreateVMTemplateWithID(existingVMTemplateID)
					if vmTemplate == nil {
						err = errors.ErrNoVMTemplate
						return
					}
				})

				ginkgo.When("when merge data of given vmTemplate", func() {
					ginkgo.BeforeEach(func() {
						recName = vmTemplateUpdateMerge
					})

					ginkgo.It("should merge data of given vmTemplate", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						vmTemplateBlueprint = blueprint.CreateUpdateVMTemplateBlueprint()
						if vmTemplateBlueprint == nil {
							err = errors.ErrNoVMTemplateBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						vmTemplateBlueprint.SetDescription("dummy")

						err = client.VMTemplateService.Update(context.TODO(), *vmTemplate, vmTemplateBlueprint,
							services.Merge)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						// check whether vmTemplate data was really updated in OpenNebula
						vmTemplateID, err = vmTemplate.ID()
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						oneVMTemplate, err = client.VMTemplateService.RetrieveInfo(context.TODO(), vmTemplateID)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(oneVMTemplate).ShouldNot(gomega.BeNil())
						gomega.Expect(oneVMTemplate.Attribute("TEMPLATE/DESCRIPTION")).To(gomega.Equal(
							"dummy"))
					})
				})

				ginkgo.When("when replace data of given vmTemplate", func() {
					ginkgo.BeforeEach(func() {
						recName = vmTemplateUpdateReplace
					})

					ginkgo.It("should replace data of given vmTemplate", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						vmTemplateBlueprint = blueprint.CreateUpdateVMTemplateBlueprint()
						if vmTemplateBlueprint == nil {
							err = errors.ErrNoVMTemplateBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						vmTemplateBlueprint.SetCPU("dummy")

						err = client.VMTemplateService.Update(context.TODO(), *vmTemplate, vmTemplateBlueprint,
							services.Replace)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						// check whether vmTemplate data was really replaced in OpenNebula
						vmTemplateID, err = vmTemplate.ID()
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						oneVMTemplate, err = client.VMTemplateService.RetrieveInfo(context.TODO(), vmTemplateID)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(oneVMTemplate).ShouldNot(gomega.BeNil())
						gomega.Expect(oneVMTemplate.Attribute("TEMPLATE/CPU")).To(gomega.Equal("dummy"))
					})
				})
			})

			ginkgo.Context("when update data is empty", func() {
				ginkgo.BeforeEach(func() {
					vmTemplate = resources.CreateVMTemplateWithID(existingVMTemplateID)
					if vmTemplate == nil {
						err = errors.ErrNoVMTemplate
						return
					}
				})

				ginkgo.When("when merge data of given vmTemplate", func() {
					ginkgo.BeforeEach(func() {
						recName = vmTemplateUpdateEmptyMerge
					})

					ginkgo.It("should merge data of given vmTemplate", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						err = client.VMTemplateService.Update(context.TODO(), *vmTemplate,
							&blueprint.VMTemplateBlueprint{}, services.Merge)
						gomega.Expect(err).To(gomega.HaveOccurred())
					})
				})

				ginkgo.When("when replace data of given vmTemplate", func() {
					ginkgo.BeforeEach(func() {
						recName = vmTemplateUpdateEmptyReplace
					})

					ginkgo.It("should replace data of given vmTemplate", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						err = client.VMTemplateService.Update(context.TODO(), *vmTemplate,
							&blueprint.VMTemplateBlueprint{}, services.Replace)
						gomega.Expect(err).To(gomega.HaveOccurred())
					})
				})
			})
		})

		ginkgo.Context("when vmTemplate doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateUpdateUnknown

				vmTemplate = resources.CreateVMTemplateWithID(nonExistingVMTemplateID)
				if vmTemplate == nil {
					err = errors.ErrNoVMTemplate
				}

				vmTemplateBlueprint = blueprint.CreateUpdateVMTemplateBlueprint()
				if vmTemplateBlueprint == nil {
					err = errors.ErrNoVMTemplateBlueprint
					return
				}
				vmTemplateBlueprint.SetDescription("dummy")
			})

			ginkgo.It("should return that vmTemplate with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VMTemplateService.Update(context.TODO(), *vmTemplate, vmTemplateBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when vmTemplate is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateUpdateNoUser

				vmTemplateBlueprint = blueprint.CreateUpdateVMTemplateBlueprint()
				if vmTemplateBlueprint == nil {
					err = errors.ErrNoVMTemplateBlueprint
					return
				}
				vmTemplateBlueprint.SetDescription("dummy")
			})

			ginkgo.It("should return that vmTemplate has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VMTemplateService.Update(context.TODO(), resources.VMTemplate{},
					vmTemplateBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("vmTemplate chmod", func() {
		var (
			vmTemplate    *resources.VMTemplate
			oneVMTemplate *resources.VMTemplate
			vmTemplateID  int
			permRequest   requests.PermissionRequest
		)

		ginkgo.Context("when vmTemplate exists", func() {
			ginkgo.BeforeEach(func() {
				vmTemplate = resources.CreateVMTemplateWithID(existingVMTemplateID)
				if vmTemplate == nil {
					err = errors.ErrNoVMTemplate
				}
			})

			ginkgo.When("when permission request is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = vmTemplateChmod
				})

				ginkgo.It("should change permission of given vmTemplate", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					permRequest = requests.CreatePermissionRequestBuilder().Deny(requests.User,
						requests.Manage).Allow(requests.Other, requests.Admin).Build()

					err = client.VMTemplateService.Chmod(context.TODO(), *vmTemplate, permRequest)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chmod was really changed in OpenNebula
					vmTemplateID, err = vmTemplate.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneVMTemplate, err = client.VMTemplateService.RetrieveInfo(context.TODO(), vmTemplateID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVMTemplate).ShouldNot(gomega.BeNil())

					var perm *resources.Permissions
					perm, err = oneVMTemplate.Permissions()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					gomega.Expect(perm.User.Manage).To(gomega.Equal(false))
					gomega.Expect(perm.Other.Admin).To(gomega.Equal(true))
				})
			})

			ginkgo.When("when permission request is default", func() {
				ginkgo.BeforeEach(func() {
					recName = vmTemplatePermRequestDefault
				})

				ginkgo.It("should not change permissions of given vmTemplate", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.VMTemplateService.Chmod(context.TODO(), *vmTemplate,
						requests.CreatePermissionRequestBuilder().Build())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when vmTemplate doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateChmodUnknown

				vmTemplate = resources.CreateVMTemplateWithID(nonExistingVMTemplateID)
				if vmTemplate == nil {
					err = errors.ErrNoVMTemplate
				}
			})

			ginkgo.It("should return that vmTemplate with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				permRequest = requests.CreatePermissionRequestBuilder().Allow(requests.User,
					requests.Manage).Deny(requests.Other, requests.Admin).Build()

				err = client.VMTemplateService.Chmod(context.TODO(), *vmTemplate, permRequest)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when vmTemplate is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateChmodNoVMTemplate

				vmTemplate = &resources.VMTemplate{}
			})

			ginkgo.It("should return that vmTemplate has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VMTemplateService.Chmod(context.TODO(), *vmTemplate, requests.PermissionRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("vmTemplate chown", func() {
		var (
			vmTemplate    *resources.VMTemplate
			oneVMTemplate *resources.VMTemplate
			vmTemplateID  int

			user         *resources.User
			group        *resources.Group
			ownershipReq requests.OwnershipRequest
		)

		ginkgo.Context("when vmTemplate exists", func() {
			ginkgo.BeforeEach(func() {
				vmTemplate = resources.CreateVMTemplateWithID(existingVMTemplateID)
				if vmTemplate == nil {
					err = errors.ErrNoVMTemplate
				}
			})

			ginkgo.When("when ownership request is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = vmTemplateChown
				})

				ginkgo.It("should change owner of given vmTemplate", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					userID := 31
					groupID := 120

					user = resources.CreateUserWithID(userID)
					group = resources.CreateGroupWithID(groupID)

					ownershipReq = requests.CreateOwnershipRequestBuilder().User(*user).Group(*group).Build()

					err = client.VMTemplateService.Chown(context.TODO(), *vmTemplate, ownershipReq)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chown was really changed in OpenNebula
					vmTemplateID, err = vmTemplate.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneVMTemplate, err = client.VMTemplateService.RetrieveInfo(context.TODO(), vmTemplateID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVMTemplate).ShouldNot(gomega.BeNil())

					gomega.Expect(oneVMTemplate.User()).To(gomega.Equal(userID))
					gomega.Expect(oneVMTemplate.Group()).To(gomega.Equal(groupID))
				})
			})

			ginkgo.When("when ownership request is default", func() {
				ginkgo.BeforeEach(func() {
					recName = vmTemplateOwnershipReqDefault
				})

				ginkgo.It("should not change permissions of given vmTemplate", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.VMTemplateService.Chown(context.TODO(), *vmTemplate,
						requests.CreateOwnershipRequestBuilder().Build())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when vmTemplate doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateChownUnknown

				vmTemplate = resources.CreateVMTemplateWithID(nonExistingVMTemplateID)
				if vmTemplate == nil {
					err = errors.ErrNoVMTemplate
				}
			})

			ginkgo.It("should return that vmTemplate with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VMTemplateService.Chown(context.TODO(), *vmTemplate, requests.OwnershipRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when vmTemplate is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateChownNoVMTemplate
			})

			ginkgo.It("should return that vmTemplate has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VMTemplateService.Chown(context.TODO(), resources.VMTemplate{},
					requests.OwnershipRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("vmTemplate rename", func() {
		var (
			vmTemplate    *resources.VMTemplate
			oneVMTemplate *resources.VMTemplate
			vmTemplateID  int
		)

		ginkgo.Context("when vmTemplate exists", func() {
			ginkgo.BeforeEach(func() {
				vmTemplate = resources.CreateVMTemplateWithID(existingVMTemplateID)
				if vmTemplate == nil {
					err = errors.ErrNoVMTemplate
				}
			})

			ginkgo.When("when new name is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = vmTemplateRename
				})

				ginkgo.It("should change name of given vmTemplate", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					// get vmTemplate name
					oneVMTemplate, err = client.VMTemplateService.RetrieveInfo(context.TODO(), vmTemplateID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVMTemplate).ShouldNot(gomega.BeNil())

					newName := "monkey2"
					gomega.Expect(oneVMTemplate.Name()).NotTo(gomega.Equal(newName))

					// change name
					err = client.VMTemplateService.Rename(context.TODO(), *vmTemplate, newName)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether name was really changed in OpenNebula
					vmTemplateID, err = vmTemplate.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneVMTemplate, err = client.VMTemplateService.RetrieveInfo(context.TODO(), vmTemplateID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVMTemplate).ShouldNot(gomega.BeNil())

					gomega.Expect(oneVMTemplate.Name()).To(gomega.Equal("monkey2"))
				})
			})

			ginkgo.When("when new name is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = vmTemplateRenameEmpty
				})

				ginkgo.It("should not change name of given vmTemplate", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.VMTemplateService.Rename(context.TODO(), *vmTemplate, "")
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when vmTemplate doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateRenameUnknown

				vmTemplate = resources.CreateVMTemplateWithID(nonExistingVMTemplateID)
				if vmTemplate == nil {
					err = errors.ErrNoVMTemplate
				}
			})

			ginkgo.It("should return that vmTemplate with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VMTemplateService.Rename(context.TODO(), *vmTemplate, "mastodont")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when vmTemplate is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateRenameNoVMTemplate
			})

			ginkgo.It("should return that vmTemplate has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VMTemplateService.Rename(context.TODO(), resources.VMTemplate{}, "rex")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("vmTemplate retrieve info", func() {
		var vmTemplate *resources.VMTemplate

		ginkgo.Context("when vmTemplate exists", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateRetrieveInfo
			})

			ginkgo.It("should return vmTemplate with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplate, err = client.VMTemplateService.RetrieveInfo(context.TODO(), existingVMTemplateID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(vmTemplate).ShouldNot(gomega.BeNil())
				gomega.Expect(vmTemplate.ID()).To(gomega.Equal(existingVMTemplateID))
				gomega.Expect(vmTemplate.Name()).To(gomega.Equal("monkey2"))
			})
		})

		ginkgo.Context("when vmTemplate doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateRetrieveInfoUnknown
			})

			ginkgo.It("should return that given vmTemplate doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplate, err = client.VMTemplateService.RetrieveInfo(context.TODO(), nonExistingVMTemplateID)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(vmTemplate).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("vmTemplate list all", func() {
		var vmTemplates []*resources.VMTemplate

		ginkgo.Context("when ownership filter is not empty", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateListAllPrimaryGroup
			})

			ginkgo.It("should return list of all vmTemplates with full info belongs to primary group", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplates, err = client.VMTemplateService.ListAll(context.TODO(),
					services.OwnershipFilterPrimaryGroup)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(vmTemplates).ShouldNot(gomega.BeNil())
				gomega.Expect(vmTemplates).To(gomega.HaveLen(4))
			})
		})

		ginkgo.Context("when ownership filter is not empty", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateListAllUser
			})

			ginkgo.It("should return list of all vmTemplates with full info belongs to the user", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplates, err = client.VMTemplateService.ListAll(context.TODO(), services.OwnershipFilterUser)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(vmTemplates).ShouldNot(gomega.BeNil())
				gomega.Expect(vmTemplates).To(gomega.HaveLen(4))
				gomega.Expect(vmTemplates[0].ID()).To(gomega.Equal(0))
				gomega.Expect(vmTemplates[1].ID()).To(gomega.Equal(324))
				gomega.Expect(vmTemplates[2].ID()).To(gomega.Equal(327))
				gomega.Expect(vmTemplates[3].ID()).To(gomega.Equal(328))
			})
		})

		ginkgo.Context("when ownership filter is set to ALl", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateListAllAll
			})

			ginkgo.It("should return list of all vmTemplates with full info belongs to all", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplates, err = client.VMTemplateService.ListAll(context.TODO(), services.OwnershipFilterAll)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(vmTemplates).ShouldNot(gomega.BeNil())
				gomega.Expect(vmTemplates).To(gomega.HaveLen(5))
			})
		})

		ginkgo.Context("when ownership filter is set to UserGroup", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateListAllUserGroup
			})

			ginkgo.It("should return list of all vmTemplates with full info belongs to the "+
				"user and any of his groups", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplates, err = client.VMTemplateService.ListAll(context.TODO(), services.OwnershipFilterUserGroup)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(vmTemplates).ShouldNot(gomega.BeNil())
				gomega.Expect(vmTemplates).To(gomega.HaveLen(5))
			})
		})
	})

	ginkgo.Describe("vmTemplate list all for user", func() {
		var vmTemplates []*resources.VMTemplate

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateListAllForUser
			})

			ginkgo.It("should return vmTemplates with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplates, err = client.VMTemplateService.ListAllForUser(context.TODO(),
					*resources.CreateUserWithID(31))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(vmTemplates).ShouldNot(gomega.BeNil())
				gomega.Expect(vmTemplates).To(gomega.HaveLen(1))
				gomega.Expect(vmTemplates[0].ID()).To(gomega.Equal(325))
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateListAllForUserUnknown
			})

			ginkgo.It("should return empty list of vmTemplates (length 0)", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplates, err = client.VMTemplateService.ListAllForUser(context.TODO(),
					*resources.CreateUserWithID(42))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(vmTemplates).Should(gomega.Equal(make([]*resources.VMTemplate, 0)))
				gomega.Expect(vmTemplates).Should(gomega.HaveLen(0))
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateListAllForUserEmpty
			})

			ginkgo.It("should return that user doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplates, err = client.VMTemplateService.ListAllForUser(context.TODO(), resources.User{})
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(vmTemplates).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("List methods with pagination", func() {
		var (
			vmTemplates []*resources.VMTemplate
		)

		ginkgo.Context("pagination ok", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateListPagination
			})

			ginkgo.It("should return vmTemplates with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplates, err = client.VMTemplateService.List(context.TODO(), 2, 2,
					services.OwnershipFilterUserGroup)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(vmTemplates).ShouldNot(gomega.BeNil())
				gomega.Expect(vmTemplates).To(gomega.HaveLen(2))
				gomega.Expect(vmTemplates[0].ID()).To(gomega.Equal(325))
				gomega.Expect(vmTemplates[1].ID()).To(gomega.Equal(327))
			})
		})

		ginkgo.Context("pagination wrong", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateListPaginationWrong
			})
		})

		ginkgo.It("should return that pagination is wrong", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

			vmTemplates, err = client.VMTemplateService.List(context.TODO(), -2, -2,
				services.OwnershipFilterPrimaryGroup)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(vmTemplates).Should(gomega.BeNil())
		})
	})

	ginkgo.Describe("vmTemplate list for user", func() {
		var vmTemplates []*resources.VMTemplate

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateListForUser
			})

			ginkgo.It("should return vmTemplates with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplates, err = client.VMTemplateService.ListForUser(context.TODO(),
					*resources.CreateUserWithID(0), 2, 2)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(vmTemplates).ShouldNot(gomega.BeNil())
				gomega.Expect(vmTemplates).To(gomega.HaveLen(2))
				gomega.Expect(vmTemplates[0].ID()).To(gomega.Equal(327))
				gomega.Expect(vmTemplates[1].ID()).To(gomega.Equal(328))
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateListForUserUnknown
			})

			ginkgo.It("should return empty list of vmTemplates (length 0)", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplates, err = client.VMTemplateService.ListForUser(context.TODO(),
					*resources.CreateUserWithID(88), 2, 2)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(vmTemplates).Should(gomega.Equal(make([]*resources.VMTemplate, 0)))
				gomega.Expect(vmTemplates).Should(gomega.HaveLen(0))
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = vmTemplateListForUserEmpty
			})

			ginkgo.It("should return that user doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				vmTemplates, err = client.VMTemplateService.ListForUser(context.TODO(), resources.User{}, 2, 2)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(vmTemplates).Should(gomega.BeNil())
			})
		})
	})
})
