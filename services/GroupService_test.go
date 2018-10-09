package services_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/onego-project/onego"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/resources"
	"github.com/onego-project/onego/services"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	existingGroupID    = 120
	nonExistingGroupID = 158
)

const (
	groupAllocate       = "records/group/allocate"
	groupAllocationFail = "records/group/allocationFail"

	groupDelete      = "records/group/delete"
	groupDeleteFail  = "records/group/deleteFail"
	groupDeleteEmpty = "records/group/deleteEmpty"

	groupUpdateMerge      = "records/group/updateMerge"
	groupUpdateReplace    = "records/group/updateReplace"
	groupUpdateEmpty      = "records/group/updateEmpty"
	groupUpdateFail       = "records/group/updateFail"
	groupUpdateGroupEmpty = "records/group/updateGroupEmpty"

	groupRetrieveInfoGroup        = "records/group/retrieveInfoGroup"
	groupRetrieveInfoGroupUnknown = "records/group/retrieveInfoGroupUnknown"

	groupList = "records/group/list"

	groupAddAdmin           = "records/group/addAdmin"
	groupAddAdminUnknown    = "records/group/addAdminUnknown"
	groupAddAdminFail       = "records/group/addAdminFail"
	groupAddAdminGroupEmpty = "records/group/addAdminGroupEmpty"
	groupAddAdminUserEmpty  = "records/group/addAdminUserEmpty"

	groupDelAdmin           = "records/group/delAdmin"
	groupDelAdminUnknown    = "records/group/delAdminUnknown"
	groupDelAdminFail       = "records/group/delAdminFail"
	groupDelAdminGroupEmpty = "records/group/delAdminGroupEmpty"
	groupDelAdminUserEmpty  = "records/group/delAdminUserEmpty"
)

var _ = ginkgo.Describe("Group", func() {
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

		// Create an HTTP client and inject our transport
		clientHTTP := &http.Client{
			Transport: rec, // Inject as transport!
		}

		// create onego client
		client = onego.CreateClient(endpoint, token, clientHTTP)
		if client == nil {
			err = fmt.Errorf("no client")
			return
		}
	})

	ginkgo.AfterEach(func() {
		rec.Stop()
	})

	ginkgo.Describe("allocate group", func() {
		var group *resources.Group

		ginkgo.Context("when group doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = groupAllocate
			})

			ginkgo.It("should create new group", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				group, err = client.GroupService.Allocate(context.TODO(), "the_best_group")
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(group).ShouldNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when group exists", func() {
			ginkgo.BeforeEach(func() {
				recName = groupAllocationFail
			})

			ginkgo.It("shouldn't create new group", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				group, err = client.GroupService.Allocate(context.TODO(), "the_best_group")
				gomega.Expect(err).ShouldNot(gomega.BeNil())
				gomega.Expect(group).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("delete group", func() {
		var group *resources.Group

		ginkgo.Context("when group exists", func() {
			ginkgo.BeforeEach(func() {
				recName = groupDelete

				group = resources.CreateGroupWithID(existingGroupID)
			})

			ginkgo.It("should delete the given group", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.GroupService.Delete(context.TODO(), *group)
				gomega.Expect(err).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when group doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = groupDeleteFail

				group = resources.CreateGroupWithID(nonExistingGroupID)
			})

			ginkgo.It("should return an error that group doesn't exist", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.GroupService.Delete(context.TODO(), *group)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when group is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = groupDeleteEmpty
			})

			ginkgo.It("should return an error that group is empty", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.GroupService.Delete(context.TODO(), resources.Group{})
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("update group", func() {
		var group *resources.Group
		var blueprintGroup *blueprint.GroupBlueprint

		ginkgo.Context("when group exists", func() {
			ginkgo.When("when update type is merge", func() {
				ginkgo.BeforeEach(func() {
					recName = groupUpdateMerge

					group = resources.CreateGroupWithID(existingGroupID)

					blueprintGroup = blueprint.CreateUpdateGroupBlueprint()
					blueprintGroup.SetElement("someTag", "someValue")
				})

				ginkgo.It("should merge data", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.GroupService.Update(context.TODO(), *group, blueprintGroup, services.Merge)
					gomega.Expect(err).Should(gomega.BeNil())
				})
			})

			ginkgo.When("when update type is replace", func() {
				ginkgo.BeforeEach(func() {
					recName = groupUpdateReplace

					group = resources.CreateGroupWithID(existingGroupID)

					blueprintGroup = blueprint.CreateUpdateGroupBlueprint()
					blueprintGroup.SetElement("newTag", "newValue")
				})

				ginkgo.It("should replace data", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.GroupService.Update(context.TODO(), *group, blueprintGroup, services.Replace)
					gomega.Expect(err).Should(gomega.BeNil())
				})
			})

			ginkgo.When("when update data is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = groupUpdateEmpty

					group = resources.CreateGroupWithID(existingGroupID)
				})

				ginkgo.It("should returns that blueprint data is empty", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.GroupService.Update(context.TODO(), *group, &blueprint.GroupBlueprint{}, services.Merge)
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})
			})
		})

		ginkgo.Context("when group doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = groupUpdateFail

				group = resources.CreateGroupWithID(nonExistingGroupID)

				blueprintGroup = blueprint.CreateUpdateGroupBlueprint()
				blueprintGroup.SetElement("someTag", "someValue")
			})

			ginkgo.It("should return an error that group doesn't exist", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.GroupService.Update(context.TODO(), *group, blueprintGroup, services.Merge)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when group is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = groupUpdateGroupEmpty

				blueprintGroup = blueprint.CreateUpdateGroupBlueprint()
				blueprintGroup.SetElement("someTag", "someValue")
			})

			ginkgo.It("should return an error that group is empty", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.GroupService.Update(context.TODO(), resources.Group{}, blueprintGroup, services.Merge)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("retrieve info", func() {
		var group *resources.Group

		ginkgo.Context("when group exists", func() {
			ginkgo.BeforeEach(func() {
				recName = groupRetrieveInfoGroup
			})

			ginkgo.It("should return group with full info", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				group, err = client.GroupService.RetrieveInfo(context.TODO(), existingGroupID)
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(group).ShouldNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when group doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = groupRetrieveInfoGroupUnknown
			})

			ginkgo.It("should return that given group doesn't exist", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				group, err = client.GroupService.RetrieveInfo(context.TODO(), nonExistingGroupID)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
				gomega.Expect(group).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("list groups", func() {
		ginkgo.BeforeEach(func() {
			recName = groupList
		})

		ginkgo.It("should return an array of groups with full info", func() {
			gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

			var groups []*resources.Group

			groups, err = client.GroupService.List(context.TODO())
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(groups).ShouldNot(gomega.BeNil())
		})
	})

	ginkgo.Describe("add admin", func() {
		var group *resources.Group
		var user *resources.User

		ginkgo.Context("when group exists", func() {
			ginkgo.When("when user exists", func() {
				ginkgo.BeforeEach(func() {
					recName = groupAddAdmin

					group = resources.CreateGroupWithID(existingGroupID)
					user = resources.CreateUserWithID(idExistingUser)
				})

				ginkgo.It("should add admin to given group", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.GroupService.AddAdmin(context.TODO(), *group, *user)
					gomega.Expect(err).Should(gomega.BeNil())
				})
			})

			ginkgo.When("when user doesn't exist", func() {
				ginkgo.BeforeEach(func() {
					recName = groupAddAdminUnknown

					group = resources.CreateGroupWithID(existingGroupID)
					user = resources.CreateUserWithID(idNonExistingUser)
				})

				ginkgo.It("should add admin to given group", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.GroupService.AddAdmin(context.TODO(), *group, *user)
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})
			})

			ginkgo.When("when user is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = groupAddAdminUserEmpty

					group = resources.CreateGroupWithID(existingGroupID)
				})

				ginkgo.It("should add admin to given group", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.GroupService.AddAdmin(context.TODO(), *group, resources.User{})
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})
			})
		})

		ginkgo.Context("when group doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = groupAddAdminFail

				group = resources.CreateGroupWithID(nonExistingGroupID)
				user = resources.CreateUserWithID(idExistingUser)
			})

			ginkgo.It("should return an error that group doesn't exist", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.GroupService.AddAdmin(context.TODO(), *group, *user)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when group is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = groupAddAdminGroupEmpty

				user = resources.CreateUserWithID(idExistingUser)
			})

			ginkgo.It("should return an error that group is empty", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.GroupService.AddAdmin(context.TODO(), resources.Group{}, *user)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("remove admin", func() {
		var group *resources.Group
		var user *resources.User

		ginkgo.Context("when group exists", func() {
			ginkgo.When("when user exists", func() {
				ginkgo.BeforeEach(func() {
					recName = groupDelAdmin

					group = resources.CreateGroupWithID(existingGroupID)
					user = resources.CreateUserWithID(idExistingUser)
				})

				ginkgo.It("should remove admin from given group", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.GroupService.RemoveAdmin(context.TODO(), *group, *user)
					gomega.Expect(err).Should(gomega.BeNil())
				})
			})

			ginkgo.When("when user doesn't exist", func() {
				ginkgo.BeforeEach(func() {
					recName = groupDelAdminUnknown

					group = resources.CreateGroupWithID(existingGroupID)
					user = resources.CreateUserWithID(idNonExistingUser)
				})

				ginkgo.It("should remove admin from given group", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.GroupService.RemoveAdmin(context.TODO(), *group, *user)
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})
			})

			ginkgo.When("when user is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = groupDelAdminUserEmpty

					group = resources.CreateGroupWithID(existingGroupID)
				})

				ginkgo.It("should remove admin from given group", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.GroupService.RemoveAdmin(context.TODO(), *group, resources.User{})
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})
			})
		})

		ginkgo.Context("when group doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = groupDelAdminFail

				group = resources.CreateGroupWithID(nonExistingGroupID)
				user = resources.CreateUserWithID(idExistingUser)
			})

			ginkgo.It("should return an error that group doesn't exist", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.GroupService.RemoveAdmin(context.TODO(), *group, *user)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when group is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = groupDelAdminGroupEmpty

				user = resources.CreateUserWithID(idExistingUser)
			})

			ginkgo.It("should return an error that group is empty", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.GroupService.RemoveAdmin(context.TODO(), resources.Group{}, *user)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})
	})
})
