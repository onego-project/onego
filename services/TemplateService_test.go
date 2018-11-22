package services_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/onego-project/onego/services"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/onego-project/onego"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/errors"
	"github.com/onego-project/onego/requests"
	"github.com/onego-project/onego/resources"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	templateAllocate               = "records/template/allocate"
	templateAllocateEmptyBlueprint = "records/template/allocateEmptyBlueprint"
	templateAllocateExisting       = "records/template/allocateExisting"

	templateClone         = "records/template/clone"
	templateCloneExisting = "records/template/cloneExisting"
	templateCloneNoID     = "records/template/cloneNoID"

	templateDelete        = "records/template/delete"
	templateDeleteWrongID = "records/template/deleteWrongID"
	templateDeleteNoID    = "records/template/deleteNoID"

	templateUpdateMerge        = "records/template/updateMerge"
	templateUpdateReplace      = "records/template/updateReplace"
	templateUpdateEmptyMerge   = "records/template/updateEmptyMerge"
	templateUpdateEmptyReplace = "records/template/updateEmptyReplace"
	templateUpdateNoUser       = "records/template/updateNoUser"
	templateUpdateUnknown      = "records/template/updateUnknown"

	templateChmod              = "records/template/chmod"
	templatePermRequestDefault = "records/template/chmodPermReqDefault"
	templateChmodUnknown       = "records/template/chmodUnknown"
	templateChmodNoTemplate    = "records/template/chmodNoTemplate"

	templateChown               = "records/template/chown"
	templateOwnershipReqDefault = "records/template/chownDefault"
	templateChownUnknown        = "records/template/chownUnknown"
	templateChownNoTemplate     = "records/template/chownNoTemplate"

	templateRename           = "records/template/rename"
	templateRenameEmpty      = "records/template/renameEmpty"
	templateRenameUnknown    = "records/template/renameUnknown"
	templateRenameNoTemplate = "records/template/renameNoTemplate"

	templateRetrieveInfo        = "records/template/retrieveInfo"
	templateRetrieveInfoUnknown = "records/template/retrieveInfoUnknown"

	templateListAllPrimaryGroup = "records/template/listAllPrimaryGroup"
	templateListAllUser         = "records/template/listAllUser"
	templateListAllAll          = "records/template/listAllAll"
	templateListAllUserGroup    = "records/template/listAllUserGroup"

	templateListAllForUser        = "records/template/listAllForUser"
	templateListAllForUserUnknown = "records/template/listAllForUserUnknown"
	templateListAllForUserEmpty   = "records/template/listAllForUserEmpty"

	templateListPagination      = "records/template/listPagination"
	templateListPaginationWrong = "records/template/listPaginationWrong"

	templateListForUser        = "records/template/listForUser"
	templateListForUserUnknown = "records/template/listForUserUnknown"
	templateListForUserEmpty   = "records/template/listForUserEmpty"
)

var _ = ginkgo.Describe("Template Service", func() {
	var (
		recName string
		rec     *recorder.Recorder
		client  *onego.Client
		err     error
	)

	var existingTemplateID = 325
	var nonExistingTemplateID = 420
	var allocatedDeletedTemplateID = 326

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

	ginkgo.Describe("allocate template", func() {
		var (
			template          *resources.Template
			templateID        int
			oneTemplate       *resources.Template
			templateBlueprint *blueprint.TemplateBlueprint
		)

		ginkgo.BeforeEach(func() {
			templateBlueprint = blueprint.CreateAllocateTemplateBlueprint()
			templateBlueprint.SetElement("NAME", "hello")
		})

		ginkgo.Context("when template doesn't exist", func() {
			ginkgo.When("when attributes are set correctly", func() {
				ginkgo.BeforeEach(func() {
					recName = templateAllocate
				})

				ginkgo.It("should create new template", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					template, err = client.TemplateService.Allocate(context.TODO(), templateBlueprint)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(template).ShouldNot(gomega.BeNil())

					// check whether Template really exists in OpenNebula
					templateID, err = template.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneTemplate, err = client.TemplateService.RetrieveInfo(context.TODO(), templateID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneTemplate.Name()).To(gomega.Equal("hello"))
				})
			})

			ginkgo.When("when blueprint is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = templateAllocateEmptyBlueprint
				})

				ginkgo.It("shouldn't create new template", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					template, err = client.TemplateService.Allocate(context.TODO(),
						&blueprint.TemplateBlueprint{})
					gomega.Expect(err).To(gomega.HaveOccurred())
					gomega.Expect(template).Should(gomega.BeNil())
				})
			})
		})

		ginkgo.Context("when template exists", func() {
			ginkgo.BeforeEach(func() {
				recName = templateAllocateExisting
			})

			ginkgo.It("should return that template already exists", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				template, err = client.TemplateService.Allocate(context.TODO(), templateBlueprint)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(template).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("clone template", func() {
		var (
			template    *resources.Template
			clone       *resources.Template
			oneTemplate *resources.Template
			cloneID     int
		)

		ginkgo.Context("when template doesn't exist", func() {
			ginkgo.When("when attributes are set correctly", func() {
				ginkgo.BeforeEach(func() {
					recName = templateClone

					template = resources.CreateTemplateWithID(existingTemplateID)
					if template == nil {
						err = errors.ErrNoTemplate
					}
				})

				ginkgo.It("should create new template clone", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					clone, err = client.TemplateService.Clone(context.TODO(), *template, "Golias", true)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(clone).ShouldNot(gomega.BeNil())

					// check whether Template really exists in OpenNebula
					cloneID, err = clone.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneTemplate, err = client.TemplateService.RetrieveInfo(context.TODO(), cloneID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneTemplate.Name()).To(gomega.Equal("Golias"))
				})
			})
		})

		ginkgo.Context("when template exists", func() {
			ginkgo.BeforeEach(func() {
				recName = templateCloneExisting
			})

			ginkgo.It("should return that template already exists", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				clone, err = client.TemplateService.Clone(context.TODO(), *template, "test-delete", true)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(clone).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when template is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = templateCloneNoID

				template = &resources.Template{}
			})

			ginkgo.It("should return that template has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				clone, err = client.TemplateService.Clone(context.TODO(), *template, "asdfg", true)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(clone).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("delete template", func() {
		var (
			template    *resources.Template
			oneTemplate *resources.Template
			templateID  int
		)

		ginkgo.Context("when template exists", func() {
			ginkgo.BeforeEach(func() {
				recName = templateDelete

				template = resources.CreateTemplateWithID(allocatedDeletedTemplateID)
				if template == nil {
					err = errors.ErrNoTemplate
				}
			})

			ginkgo.It("should delete template", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.TemplateService.Delete(context.TODO(), *template, true)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether template was really deleted in OpenNebula
				templateID, err = template.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneTemplate, err = client.TemplateService.RetrieveInfo(context.TODO(), templateID)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(oneTemplate).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when template doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = templateDeleteWrongID

				template = resources.CreateTemplateWithID(nonExistingTemplateID)
				if template == nil {
					err = errors.ErrNoTemplate
				}
			})

			ginkgo.It("should return that template with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.TemplateService.Delete(context.TODO(), *template, true)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when template is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = templateDeleteNoID
			})

			ginkgo.It("should return that template has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.TemplateService.Delete(context.TODO(), resources.Template{}, true)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("update template", func() {
		var (
			template          *resources.Template
			oneTemplate       *resources.Template
			templateID        int
			templateBlueprint *blueprint.TemplateBlueprint
			retTemplate       *resources.Template
		)

		ginkgo.Context("when template exists", func() {
			ginkgo.Context("when update data is not empty", func() {
				ginkgo.BeforeEach(func() {
					template = resources.CreateTemplateWithID(existingTemplateID)
					if template == nil {
						err = errors.ErrNoTemplate
						return
					}
				})

				ginkgo.When("when merge data of given template", func() {
					ginkgo.BeforeEach(func() {
						recName = templateUpdateMerge
					})

					ginkgo.It("should merge data of given template", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						templateBlueprint = blueprint.CreateUpdateTemplateBlueprint()
						if templateBlueprint == nil {
							err = errors.ErrNoTemplateBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						templateBlueprint.SetDescription("dummy")

						retTemplate, err = client.TemplateService.Update(context.TODO(), *template, templateBlueprint,
							services.Merge)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(retTemplate).ShouldNot(gomega.BeNil())

						// check whether template data was really updated in OpenNebula
						templateID, err = template.ID()
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						oneTemplate, err = client.TemplateService.RetrieveInfo(context.TODO(), templateID)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(oneTemplate).ShouldNot(gomega.BeNil())
						gomega.Expect(oneTemplate.Attribute("TEMPLATE/DESCRIPTION")).To(gomega.Equal(
							"dummy"))
					})
				})

				ginkgo.When("when replace data of given template", func() {
					ginkgo.BeforeEach(func() {
						recName = templateUpdateReplace
					})

					ginkgo.It("should replace data of given template", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						templateBlueprint = blueprint.CreateUpdateTemplateBlueprint()
						if templateBlueprint == nil {
							err = errors.ErrNoTemplateBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						templateBlueprint.SetCPU("dummy")

						retTemplate, err = client.TemplateService.Update(context.TODO(), *template, templateBlueprint,
							services.Replace)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(retTemplate).ShouldNot(gomega.BeNil())

						// check whether template data was really replaced in OpenNebula
						templateID, err = template.ID()
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						oneTemplate, err = client.TemplateService.RetrieveInfo(context.TODO(), templateID)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
						gomega.Expect(oneTemplate).ShouldNot(gomega.BeNil())
						gomega.Expect(oneTemplate.Attribute("TEMPLATE/CPU")).To(gomega.Equal("dummy"))
					})
				})
			})

			ginkgo.Context("when update data is empty", func() {
				ginkgo.BeforeEach(func() {
					template = resources.CreateTemplateWithID(existingTemplateID)
					if template == nil {
						err = errors.ErrNoTemplate
						return
					}
				})

				ginkgo.When("when merge data of given template", func() {
					ginkgo.BeforeEach(func() {
						recName = templateUpdateEmptyMerge
					})

					ginkgo.It("should merge data of given template", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						retTemplate, err = client.TemplateService.Update(context.TODO(), *template,
							&blueprint.TemplateBlueprint{}, services.Merge)
						gomega.Expect(err).To(gomega.HaveOccurred())
						gomega.Expect(retTemplate).Should(gomega.BeNil())
					})
				})

				ginkgo.When("when replace data of given template", func() {
					ginkgo.BeforeEach(func() {
						recName = templateUpdateEmptyReplace
					})

					ginkgo.It("should replace data of given template", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						retTemplate, err = client.TemplateService.Update(context.TODO(), *template,
							&blueprint.TemplateBlueprint{}, services.Replace)
						gomega.Expect(err).To(gomega.HaveOccurred())
						gomega.Expect(retTemplate).Should(gomega.BeNil())
					})
				})
			})
		})

		ginkgo.Context("when template doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = templateUpdateUnknown

				template = resources.CreateTemplateWithID(nonExistingTemplateID)
				if template == nil {
					err = errors.ErrNoTemplate
				}

				templateBlueprint = blueprint.CreateUpdateTemplateBlueprint()
				if templateBlueprint == nil {
					err = errors.ErrNoTemplateBlueprint
					return
				}
				templateBlueprint.SetDescription("dummy")
			})

			ginkgo.It("should return that template with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				retTemplate, err = client.TemplateService.Update(context.TODO(), *template, templateBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(retTemplate).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when template is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = templateUpdateNoUser

				templateBlueprint = blueprint.CreateUpdateTemplateBlueprint()
				if templateBlueprint == nil {
					err = errors.ErrNoTemplateBlueprint
					return
				}
				templateBlueprint.SetDescription("dummy")
			})

			ginkgo.It("should return that template has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				retTemplate, err = client.TemplateService.Update(context.TODO(), resources.Template{},
					templateBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(retTemplate).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("template chmod", func() {
		var (
			template    *resources.Template
			oneTemplate *resources.Template
			templateID  int
			permRequest requests.PermissionRequest
		)

		ginkgo.Context("when template exists", func() {
			ginkgo.BeforeEach(func() {
				template = resources.CreateTemplateWithID(existingTemplateID)
				if template == nil {
					err = errors.ErrNoTemplate
				}
			})

			ginkgo.When("when permission request is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = templateChmod
				})

				ginkgo.It("should change permission of given template", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					permRequest = requests.CreatePermissionRequestBuilder().Deny(requests.User,
						requests.Manage).Allow(requests.Other, requests.Admin).Build()

					err = client.TemplateService.Chmod(context.TODO(), *template, permRequest)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chmod was really changed in OpenNebula
					templateID, err = template.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneTemplate, err = client.TemplateService.RetrieveInfo(context.TODO(), templateID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneTemplate).ShouldNot(gomega.BeNil())

					var perm *resources.Permissions
					perm, err = oneTemplate.Permissions()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					gomega.Expect(perm.User.Manage).To(gomega.Equal(false))
					gomega.Expect(perm.Other.Admin).To(gomega.Equal(true))
				})
			})

			ginkgo.When("when permission request is default", func() {
				ginkgo.BeforeEach(func() {
					recName = templatePermRequestDefault
				})

				ginkgo.It("should not change permissions of given template", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.TemplateService.Chmod(context.TODO(), *template,
						requests.CreatePermissionRequestBuilder().Build())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when template doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = templateChmodUnknown

				template = resources.CreateTemplateWithID(nonExistingTemplateID)
				if template == nil {
					err = errors.ErrNoTemplate
				}
			})

			ginkgo.It("should return that template with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				permRequest = requests.CreatePermissionRequestBuilder().Allow(requests.User,
					requests.Manage).Deny(requests.Other, requests.Admin).Build()

				err = client.TemplateService.Chmod(context.TODO(), *template, permRequest)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when template is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = templateChmodNoTemplate

				template = &resources.Template{}
			})

			ginkgo.It("should return that template has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.TemplateService.Chmod(context.TODO(), *template, requests.PermissionRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("template chown", func() {
		var (
			template    *resources.Template
			oneTemplate *resources.Template
			templateID  int

			user         *resources.User
			group        *resources.Group
			ownershipReq requests.OwnershipRequest
		)

		ginkgo.Context("when template exists", func() {
			ginkgo.BeforeEach(func() {
				template = resources.CreateTemplateWithID(existingTemplateID)
				if template == nil {
					err = errors.ErrNoTemplate
				}
			})

			ginkgo.When("when ownership request is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = templateChown
				})

				ginkgo.It("should change owner of given template", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					userID := 31
					groupID := 120

					user = resources.CreateUserWithID(userID)
					group = resources.CreateGroupWithID(groupID)

					ownershipReq = requests.CreateOwnershipRequestBuilder().User(*user).Group(*group).Build()

					err = client.TemplateService.Chown(context.TODO(), *template, ownershipReq)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chown was really changed in OpenNebula
					templateID, err = template.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneTemplate, err = client.TemplateService.RetrieveInfo(context.TODO(), templateID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneTemplate).ShouldNot(gomega.BeNil())

					gomega.Expect(oneTemplate.User()).To(gomega.Equal(userID))
					gomega.Expect(oneTemplate.Group()).To(gomega.Equal(groupID))
				})
			})

			ginkgo.When("when ownership request is default", func() {
				ginkgo.BeforeEach(func() {
					recName = templateOwnershipReqDefault
				})

				ginkgo.It("should not change permissions of given template", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.TemplateService.Chown(context.TODO(), *template,
						requests.CreateOwnershipRequestBuilder().Build())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when template doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = templateChownUnknown

				template = resources.CreateTemplateWithID(nonExistingTemplateID)
				if template == nil {
					err = errors.ErrNoTemplate
				}
			})

			ginkgo.It("should return that template with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.TemplateService.Chown(context.TODO(), *template, requests.OwnershipRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when template is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = templateChownNoTemplate
			})

			ginkgo.It("should return that template has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.TemplateService.Chown(context.TODO(), resources.Template{},
					requests.OwnershipRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("template rename", func() {
		var (
			template    *resources.Template
			oneTemplate *resources.Template
			templateID  int
		)

		ginkgo.Context("when template exists", func() {
			ginkgo.BeforeEach(func() {
				template = resources.CreateTemplateWithID(existingTemplateID)
				if template == nil {
					err = errors.ErrNoTemplate
				}
			})

			ginkgo.When("when new name is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = templateRename
				})

				ginkgo.It("should change name of given template", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					// get template name
					oneTemplate, err = client.TemplateService.RetrieveInfo(context.TODO(), templateID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneTemplate).ShouldNot(gomega.BeNil())

					newName := "monkey2"
					gomega.Expect(oneTemplate.Name()).NotTo(gomega.Equal(newName))

					// change name
					err = client.TemplateService.Rename(context.TODO(), *template, newName)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether name was really changed in OpenNebula
					templateID, err = template.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneTemplate, err = client.TemplateService.RetrieveInfo(context.TODO(), templateID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneTemplate).ShouldNot(gomega.BeNil())

					gomega.Expect(oneTemplate.Name()).To(gomega.Equal("monkey2"))
				})
			})

			ginkgo.When("when new name is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = templateRenameEmpty
				})

				ginkgo.It("should not change name of given template", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.TemplateService.Rename(context.TODO(), *template, "")
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when template doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = templateRenameUnknown

				template = resources.CreateTemplateWithID(nonExistingTemplateID)
				if template == nil {
					err = errors.ErrNoTemplate
				}
			})

			ginkgo.It("should return that template with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.TemplateService.Rename(context.TODO(), *template, "mastodont")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when template is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = templateRenameNoTemplate
			})

			ginkgo.It("should return that template has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.TemplateService.Rename(context.TODO(), resources.Template{}, "rex")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("template retrieve info", func() {
		var template *resources.Template

		ginkgo.Context("when template exists", func() {
			ginkgo.BeforeEach(func() {
				recName = templateRetrieveInfo
			})

			ginkgo.It("should return template with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				template, err = client.TemplateService.RetrieveInfo(context.TODO(), existingTemplateID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(template).ShouldNot(gomega.BeNil())
				gomega.Expect(template.ID()).To(gomega.Equal(existingTemplateID))
				gomega.Expect(template.Name()).To(gomega.Equal("monkey2"))
			})
		})

		ginkgo.Context("when template doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = templateRetrieveInfoUnknown
			})

			ginkgo.It("should return that given template doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				template, err = client.TemplateService.RetrieveInfo(context.TODO(), nonExistingTemplateID)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(template).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("template list all", func() {
		var templates []*resources.Template

		ginkgo.Context("when ownership filter is not empty", func() {
			ginkgo.BeforeEach(func() {
				recName = templateListAllPrimaryGroup
			})

			ginkgo.It("should return list of all templates with full info belongs to primary group", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				templates, err = client.TemplateService.ListAll(context.TODO(),
					services.OwnershipFilterPrimaryGroup)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(templates).ShouldNot(gomega.BeNil())
				gomega.Expect(templates).To(gomega.HaveLen(4))
			})
		})

		ginkgo.Context("when ownership filter is not empty", func() {
			ginkgo.BeforeEach(func() {
				recName = templateListAllUser
			})

			ginkgo.It("should return list of all templates with full info belongs to the user", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				templates, err = client.TemplateService.ListAll(context.TODO(), services.OwnershipFilterUser)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(templates).ShouldNot(gomega.BeNil())
				gomega.Expect(templates).To(gomega.HaveLen(4))
				gomega.Expect(templates[0].ID()).To(gomega.Equal(0))
				gomega.Expect(templates[1].ID()).To(gomega.Equal(324))
				gomega.Expect(templates[2].ID()).To(gomega.Equal(327))
				gomega.Expect(templates[3].ID()).To(gomega.Equal(328))
			})
		})

		ginkgo.Context("when ownership filter is set to ALl", func() {
			ginkgo.BeforeEach(func() {
				recName = templateListAllAll
			})

			ginkgo.It("should return list of all templates with full info belongs to all", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				templates, err = client.TemplateService.ListAll(context.TODO(), services.OwnershipFilterAll)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(templates).ShouldNot(gomega.BeNil())
				gomega.Expect(templates).To(gomega.HaveLen(5))
			})
		})

		ginkgo.Context("when ownership filter is set to UserGroup", func() {
			ginkgo.BeforeEach(func() {
				recName = templateListAllUserGroup
			})

			ginkgo.It("should return list of all templates with full info belongs to the "+
				"user and any of his groups", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				templates, err = client.TemplateService.ListAll(context.TODO(), services.OwnershipFilterUserGroup)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(templates).ShouldNot(gomega.BeNil())
				gomega.Expect(templates).To(gomega.HaveLen(5))
			})
		})
	})

	ginkgo.Describe("template list all for user", func() {
		var templates []*resources.Template

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = templateListAllForUser
			})

			ginkgo.It("should return templates with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				templates, err = client.TemplateService.ListAllForUser(context.TODO(),
					*resources.CreateUserWithID(31))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(templates).ShouldNot(gomega.BeNil())
				gomega.Expect(templates).To(gomega.HaveLen(1))
				gomega.Expect(templates[0].ID()).To(gomega.Equal(325))
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = templateListAllForUserUnknown
			})

			ginkgo.It("should return empty list of templates (length 0)", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				templates, err = client.TemplateService.ListAllForUser(context.TODO(),
					*resources.CreateUserWithID(42))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(templates).Should(gomega.Equal(make([]*resources.Template, 0)))
				gomega.Expect(templates).Should(gomega.HaveLen(0))
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = templateListAllForUserEmpty
			})

			ginkgo.It("should return that user doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				templates, err = client.TemplateService.ListAllForUser(context.TODO(), resources.User{})
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(templates).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("List methods with pagination", func() {
		var (
			templates []*resources.Template
		)

		ginkgo.Context("pagination ok", func() {
			ginkgo.BeforeEach(func() {
				recName = templateListPagination
			})

			ginkgo.It("should return templates with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				templates, err = client.TemplateService.List(context.TODO(), 2, 2,
					services.OwnershipFilterUserGroup)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(templates).ShouldNot(gomega.BeNil())
				gomega.Expect(templates).To(gomega.HaveLen(2))
				gomega.Expect(templates[0].ID()).To(gomega.Equal(325))
				gomega.Expect(templates[1].ID()).To(gomega.Equal(327))
			})
		})

		ginkgo.Context("pagination wrong", func() {
			ginkgo.BeforeEach(func() {
				recName = templateListPaginationWrong
			})
		})

		ginkgo.It("should return that pagination is wrong", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

			templates, err = client.TemplateService.List(context.TODO(), -2, -2,
				services.OwnershipFilterPrimaryGroup)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(templates).Should(gomega.BeNil())
		})
	})

	ginkgo.Describe("template list for user", func() {
		var templates []*resources.Template

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = templateListForUser
			})

			ginkgo.It("should return templates with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				templates, err = client.TemplateService.ListForUser(context.TODO(),
					*resources.CreateUserWithID(0), 2, 2)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(templates).ShouldNot(gomega.BeNil())
				gomega.Expect(templates).To(gomega.HaveLen(2))
				gomega.Expect(templates[0].ID()).To(gomega.Equal(327))
				gomega.Expect(templates[1].ID()).To(gomega.Equal(328))
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = templateListForUserUnknown
			})

			ginkgo.It("should return empty list of templates (length 0)", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				templates, err = client.TemplateService.ListForUser(context.TODO(),
					*resources.CreateUserWithID(88), 2, 2)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(templates).Should(gomega.Equal(make([]*resources.Template, 0)))
				gomega.Expect(templates).Should(gomega.HaveLen(0))
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = templateListForUserEmpty
			})

			ginkgo.It("should return that user doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				templates, err = client.TemplateService.ListForUser(context.TODO(), resources.User{}, 2, 2)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(templates).Should(gomega.BeNil())
			})
		})
	})
})
