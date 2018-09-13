package services_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/onego-project/onego"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/resources"
	"github.com/onego-project/onego/services"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
)

const (
	endpoint = "http://localhost:2633/RPC2"
	name     = "oneadmin"
	password = "qwerty123"
	token    = name + ":" + password
)

const (
	idExistingUser    = 22
	idNonExistingUser = 25

	idExistingNotMainGroup = 118
	idExistingGroup        = 120
	idNonExistingGroup     = 122
)

const (
	existingAuthDriver    = "server_cipher"
	nonExistingAuthDriver = "non-existing-driver"
)

const (
	userAllocate           = "records/user/allocate"
	userAllocateWrongGroup = "records/user/allocateWrongGroup"
	userAllocateExisting   = "records/user/allocateExisting"

	userDelete        = "records/user/delete"
	userDeleteWrongID = "records/user/deleteWrongID"
	userDeleteNoID    = "records/user/deleteNoID"

	userPassword            = "records/user/password"
	userPasswordEmpty       = "records/user/passwordEmpty"
	userPasswordNoUser      = "records/user/passwdNoUser"
	userPasswordUnknownUser = "records/user/passwdUnknownUser"

	userUpdateMerge        = "records/user/updateMerge"
	userUpdateReplace      = "records/user/updateReplace"
	userUpdateEmptyMerge   = "records/user/updateEmptyMerge"
	userUpdateEmptyReplace = "records/user/updateEmptyReplace"
	userUpdateNoUser       = "records/user/updateNoUser"
	userUpdateUnknownUser  = "records/user/updateUnknownUser"

	userAuthDriver            = "records/user/authDriver"
	userAuthDriverNonExisting = "records/user/authDriverNonExisting"
	userAuthDriverEmpty       = "records/user/authDriverEmpty"
	userAuthDriverNoUser      = "records/user/authDriverNoUser"
	userAuthDriverUnknownUser = "records/user/authDriverUnknownUser"

	userMainGroup            = "records/user/mainGroup"
	userMainGroupNonExisting = "records/user/mainGroupNonExisting"
	userMainGroupEmpty       = "records/user/mainGroupEmpty"
	userMainGroupNoUser      = "records/user/mainGroupNoUser"
	userMainGroupUnknownUser = "records/user/mainGroupUnknownUser"

	userSecGroupAdd = "records/user/secGroupAdd"
	userSecGroupDel = "records/user/secGroupDel"

	userRetrieveInfoUser        = "records/user/retrieveInfoUser"
	userRetrieveInfoUserUnknown = "records/user/retrieveInfoUserUnknown"

	userRetrieveInfoConnectedUser = "records/user/retrieveInfoConnectedUser"
	userRetrieveInfoAll           = "records/user/retrieveInfoAll"
)

var _ = ginkgo.Describe("User Service", func() {
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
			err = fmt.Errorf("no client")
			return
		}
	})

	ginkgo.AfterEach(func() {
		rec.Stop()
	})

	ginkgo.Describe("allocate user", func() {
		var user *resources.User
		var oneUser *resources.User
		var userID int

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.When("when attributes are set correctly", func() {
				ginkgo.BeforeEach(func() {
					recName = userAllocate
				})

				ginkgo.It("should create new user", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					secGroups := make([]resources.Group, 2)
					secGroups[0] = *resources.CreateGroup(116)
					secGroups[1] = *resources.CreateGroup(120)

					user, err = client.UserService.Allocate(context.TODO(), "Dusan", "password", "core", *resources.CreateGroup(118), secGroups)
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(user).ShouldNot(gomega.BeNil())

					// check whether User really exists in OpenNebula
					userID, err = user.ID()
					gomega.Expect(err).Should(gomega.BeNil())

					oneUser, err = client.UserService.RetrieveInfo(context.TODO(), userID)
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(oneUser.Name()).To(gomega.Equal("Dusan"))
				})
			})

			ginkgo.When("when group is wrong", func() {
				ginkgo.BeforeEach(func() {
					recName = userAllocateWrongGroup
				})

				ginkgo.It("shouldn't create new user", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					user, err = client.UserService.Allocate(context.TODO(), "Miso", "password", "core", *resources.CreateGroup(-1), nil)
					gomega.Expect(err).ShouldNot(gomega.BeNil())
					gomega.Expect(user).Should(gomega.BeNil())
				})
			})
		})

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = userAllocateExisting
			})

			ginkgo.It("should return that user with name Dusan already exists", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				user, err = client.UserService.Allocate(context.TODO(), "Dusan", "password", "core", *resources.CreateGroup(118), nil)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
				gomega.Expect(user).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("delete user", func() {
		var user *resources.User
		var oneUser *resources.User
		var userID int

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = userDelete

				user = resources.CreateUserWithID(idExistingUser)
				if user == nil {
					err = fmt.Errorf("no user to finish test")
				}
			})

			ginkgo.It("should delete user", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.UserService.Delete(context.TODO(), *user)
				gomega.Expect(err).Should(gomega.BeNil())

				// check whether User was really deleted in OpenNebula
				userID, err = user.ID()
				gomega.Expect(err).Should(gomega.BeNil())

				oneUser, err = client.UserService.RetrieveInfo(context.TODO(), userID)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
				gomega.Expect(oneUser).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = userDeleteWrongID

				user = resources.CreateUserWithID(idNonExistingUser)
				if user == nil {
					err = fmt.Errorf("no user to finish test")
				}
			})

			ginkgo.It("should return that user with given ID doesn't exist", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.UserService.Delete(context.TODO(), *user)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = userDeleteNoID

				user = &resources.User{}
			})

			ginkgo.It("should return that user has no ID", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.UserService.Delete(context.TODO(), *user)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("change user password", func() {
		var user *resources.User

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				user = resources.CreateUserWithID(33)
				if user == nil {
					err = fmt.Errorf("no user to finish test")
				}
			})

			ginkgo.When("when password is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = userPassword
				})

				ginkgo.It("should change password of given user", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.UserService.ChangePassword(context.TODO(), *user, "helloworld")
					gomega.Expect(err).Should(gomega.BeNil())
				})
			})

			ginkgo.When("password is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = userPasswordEmpty
				})

				ginkgo.It("should return error that password cannot be empty", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.UserService.ChangePassword(context.TODO(), *user, "")
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = userPasswordUnknownUser

				user = resources.CreateUserWithID(idNonExistingUser)
				if user == nil {
					err = fmt.Errorf("no user to finish test")
				}
			})

			ginkgo.It("should return that user with given ID doesn't exist", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.UserService.ChangePassword(context.TODO(), *user, "helloworld")
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = userPasswordNoUser
			})

			ginkgo.It("should return that user has no ID", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.UserService.ChangePassword(context.TODO(), resources.User{}, "helloworld")
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("update user", func() {
		var (
			user          *resources.User
			userBlueprint *blueprint.UserBlueprint
			oneUser       *resources.User
			userID        int
		)

		ginkgo.Context("when user exists", func() {
			ginkgo.Context("when update data is not empty", func() {
				ginkgo.BeforeEach(func() {
					user = resources.CreateUserWithID(33)
					if user == nil {
						err = fmt.Errorf("no user to finish test")
						return
					}

					userBlueprint = blueprint.CreateUserBlueprint()
					if userBlueprint == nil {
						err = fmt.Errorf("no user blueprint to finish test")
						return
					}
				})

				ginkgo.When("when merge data of given user", func() {
					ginkgo.BeforeEach(func() {
						recName = userUpdateMerge
					})

					ginkgo.It("should merge data of given user", func() {
						gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

						userBlueprint.SetEmail("pancake@pizza.com")

						err = client.UserService.Update(context.TODO(), *user, userBlueprint, services.Merge)
						gomega.Expect(err).Should(gomega.BeNil())

						// check whether User data was really updated in OpenNebula
						userID, err = user.ID()
						gomega.Expect(err).Should(gomega.BeNil())

						oneUser, err = client.UserService.RetrieveInfo(context.TODO(), userID)
						gomega.Expect(err).Should(gomega.BeNil())
						gomega.Expect(oneUser).ShouldNot(gomega.BeNil())
						gomega.Expect(oneUser.Attribute("TEMPLATE/EMAIL")).To(gomega.Equal("pancake@pizza.com"))
					})
				})

				ginkgo.When("when replace data of given user", func() {
					ginkgo.BeforeEach(func() {
						recName = userUpdateReplace
					})

					ginkgo.It("should replace data of given user", func() {
						gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

						userBlueprint.SetEmail("lasagne@pizza.com")
						userBlueprint.SetFullName("Frantisek Slovak")

						err = client.UserService.Update(context.TODO(), *user, userBlueprint, services.Replace)
						gomega.Expect(err).Should(gomega.BeNil())

						// check whether User data was really replaced in OpenNebula
						userID, err = user.ID()
						gomega.Expect(err).Should(gomega.BeNil())

						oneUser, err = client.UserService.RetrieveInfo(context.TODO(), userID)
						gomega.Expect(err).Should(gomega.BeNil())
						gomega.Expect(oneUser).ShouldNot(gomega.BeNil())
						gomega.Expect(oneUser.Attribute("TEMPLATE/EMAIL")).To(gomega.Equal("lasagne@pizza.com"))
						gomega.Expect(oneUser.Attribute("TEMPLATE/NAME")).To(gomega.Equal("Frantisek Slovak"))
					})
				})
			})

			ginkgo.Context("when update data is empty", func() {
				ginkgo.BeforeEach(func() {
					user = resources.CreateUserWithID(idExistingUser)
					if user == nil {
						err = fmt.Errorf("no user to finish test")
						return
					}

					userBlueprint = &blueprint.UserBlueprint{}
				})

				ginkgo.When("when merge data of given user", func() {
					ginkgo.BeforeEach(func() {
						recName = userUpdateEmptyMerge
					})

					ginkgo.It("should merge data of given user", func() {
						gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

						err = client.UserService.Update(context.TODO(), *user, userBlueprint, services.Merge)
						gomega.Expect(err).ShouldNot(gomega.BeNil())
					})
				})

				ginkgo.When("when replace data of given user", func() {
					ginkgo.BeforeEach(func() {
						recName = userUpdateEmptyReplace
					})

					ginkgo.It("should replace data of given user", func() {
						gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

						err = client.UserService.Update(context.TODO(), *user, userBlueprint, services.Replace)
						gomega.Expect(err).ShouldNot(gomega.BeNil())
					})
				})
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = userUpdateUnknownUser

				user = resources.CreateUserWithID(idNonExistingUser)
				if user == nil {
					err = fmt.Errorf("no user to finish test")
				}

				userBlueprint = blueprint.CreateUserBlueprint()
				if userBlueprint == nil {
					err = fmt.Errorf("no user blueprint to finish test")
					return
				}
				userBlueprint.SetEmail("pancake@pizza.com")
			})

			ginkgo.It("should return that user with given ID doesn't exist", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.UserService.Update(context.TODO(), *user, userBlueprint, services.Merge)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = userUpdateNoUser

				userBlueprint = blueprint.CreateUserBlueprint()
				if userBlueprint == nil {
					err = fmt.Errorf("no user blueprint to finish test")
					return
				}
				userBlueprint.SetEmail("pancake@pizza.com")
			})

			ginkgo.It("should return that user has no ID", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.UserService.Update(context.TODO(), resources.User{}, userBlueprint, services.Merge)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("change user authentication driver", func() {
		var user *resources.User
		var oneUser *resources.User
		var userID int

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				user = resources.CreateUserWithID(33)
				if user == nil {
					err = fmt.Errorf("no user to finish test")
				}
			})

			ginkgo.When("when authentication driver is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = userAuthDriver
				})

				ginkgo.It("should change authentication driver of given user", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.UserService.ChangeAuthDriver(context.TODO(), *user, existingAuthDriver)
					gomega.Expect(err).Should(gomega.BeNil())

					// check whether User auth. driver was really changed in OpenNebula
					userID, err = user.ID()
					gomega.Expect(err).Should(gomega.BeNil())

					oneUser, err = client.UserService.RetrieveInfo(context.TODO(), userID)
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(oneUser).ShouldNot(gomega.BeNil())
					gomega.Expect(oneUser.AuthDriver()).To(gomega.Equal(existingAuthDriver))
				})
			})

			ginkgo.When("when authentication driver is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = userAuthDriverEmpty
				})

				ginkgo.It("should change authentication driver of given user to empty", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.UserService.ChangeAuthDriver(context.TODO(), *user, "")
					gomega.Expect(err).Should(gomega.BeNil())
				})
			})

			ginkgo.When("when authentication driver is non-existing", func() {
				ginkgo.BeforeEach(func() {
					recName = userAuthDriverNonExisting
				})

				ginkgo.It("should change authentication driver of given user to non-existing auth. driver", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.UserService.ChangeAuthDriver(context.TODO(), *user, nonExistingAuthDriver)
					gomega.Expect(err).Should(gomega.BeNil())
				})
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = userAuthDriverUnknownUser

				user = resources.CreateUserWithID(idNonExistingUser)
				if user == nil {
					err = fmt.Errorf("no user to finish test")
				}
			})

			ginkgo.It("should return that user with given ID doesn't exist", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.UserService.ChangeAuthDriver(context.TODO(), *user, existingAuthDriver)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = userAuthDriverNoUser
			})

			ginkgo.It("should return that user has no ID", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.UserService.ChangeAuthDriver(context.TODO(), resources.User{}, existingAuthDriver)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("change main group", func() {
		var user *resources.User
		var group *resources.Group
		var oneUser *resources.User
		var userID int

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				user = resources.CreateUserWithID(33)
				if user == nil {
					err = fmt.Errorf("no user to finish test")
				}
			})

			ginkgo.When("when group exists", func() {
				ginkgo.BeforeEach(func() {
					recName = userMainGroup

					group = resources.CreateGroup(idExistingGroup)
					if group == nil {
						err = fmt.Errorf("no group to finish test")
					}
				})

				ginkgo.It("should change main group of given user", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.UserService.ChangeMainGroup(context.TODO(), *user, *group)
					gomega.Expect(err).Should(gomega.BeNil())

					// check whether User main group was really changed in OpenNebula
					userID, err = user.ID()
					gomega.Expect(err).Should(gomega.BeNil())

					oneUser, err = client.UserService.RetrieveInfo(context.TODO(), userID)
					gomega.Expect(err).Should(gomega.BeNil())
					gomega.Expect(oneUser).ShouldNot(gomega.BeNil())
					gomega.Expect(oneUser.MainGroup()).To(gomega.Equal(idExistingGroup))
				})
			})

			ginkgo.When("when group doesn't exist", func() {
				ginkgo.BeforeEach(func() {
					recName = userMainGroupNonExisting

					group = resources.CreateGroup(idNonExistingGroup)
					if group == nil {
						err = fmt.Errorf("no group to finish test")
					}
				})

				ginkgo.It("should return error that main group cannot be non-existing", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.UserService.ChangeMainGroup(context.TODO(), *user, *group)
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})
			})

			ginkgo.When("group is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = userMainGroupEmpty
				})

				ginkgo.It("should return error that main group cannot be empty", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.UserService.ChangeMainGroup(context.TODO(), *user, resources.Group{})
					gomega.Expect(err).ShouldNot(gomega.BeNil())
				})
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = userMainGroupUnknownUser

				user = resources.CreateUserWithID(idNonExistingUser)
				if user == nil {
					err = fmt.Errorf("no user to finish test")
				}

				group = resources.CreateGroup(idExistingGroup)
				if group == nil {
					err = fmt.Errorf("no group to finish test")
				}
			})

			ginkgo.It("should return that user with given ID doesn't exist", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.UserService.ChangeMainGroup(context.TODO(), *user, *group)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = userMainGroupNoUser
			})

			ginkgo.It("should return that user has no ID", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				err = client.UserService.ChangeMainGroup(context.TODO(), resources.User{}, *group)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("manage secondary group", func() {
		var user *resources.User
		var group *resources.Group

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				user = resources.CreateUserWithID(idExistingUser)
				if user == nil {
					err = fmt.Errorf("no user to finish test")
				}
			})

			ginkgo.When("when add new secondary group", func() {
				ginkgo.BeforeEach(func() {
					recName = userSecGroupAdd

					group = resources.CreateGroup(idExistingNotMainGroup)
					if group == nil {
						err = fmt.Errorf("no group to finish test")
					}
				})

				ginkgo.It("should add secondary group of given user", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.UserService.AddSecondaryGroup(context.TODO(), *user, *group)
					gomega.Expect(err).Should(gomega.BeNil())
				})
			})

			ginkgo.When("when delete existing group", func() {
				ginkgo.BeforeEach(func() {
					recName = userSecGroupDel

					group = resources.CreateGroup(idExistingNotMainGroup)
					if group == nil {
						err = fmt.Errorf("no group to finish test")
					}
				})

				ginkgo.It("should remove given group of given user", func() {
					gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

					err = client.UserService.RemoveSecondaryGroup(context.TODO(), *user, *group)
					gomega.Expect(err).Should(gomega.BeNil())
				})
			})
		})
	})

	ginkgo.Describe("retrieve info", func() {
		var user *resources.User

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = userRetrieveInfoUser
			})

			ginkgo.It("should return user with full info", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				user, err = client.UserService.RetrieveInfo(context.TODO(), idExistingUser)
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(user).ShouldNot(gomega.BeNil())
				gomega.Expect(user.ID()).To(gomega.Equal(idExistingUser))
				gomega.Expect(user.Name()).To(gomega.Equal("Karol"))
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = userRetrieveInfoUserUnknown
			})

			ginkgo.It("should return that given user doesn't exist", func() {
				gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

				user, err = client.UserService.RetrieveInfo(context.TODO(), idNonExistingUser)
				gomega.Expect(err).ShouldNot(gomega.BeNil())
				gomega.Expect(user).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("retrieve connected user info", func() {
		var user *resources.User

		ginkgo.BeforeEach(func() {
			recName = userRetrieveInfoConnectedUser
		})

		ginkgo.It("should return an info of connected user", func() {
			gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

			user, err = client.UserService.RetrieveConnectedUserInfo(context.TODO())
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(user).ShouldNot(gomega.BeNil())
		})
	})

	ginkgo.Describe("retrieve users info", func() {
		ginkgo.BeforeEach(func() {
			recName = userRetrieveInfoAll
		})

		ginkgo.It("should return an array of users with full info", func() {
			gomega.Expect(err).Should(gomega.BeNil()) // no error during BeforeEach

			var users []*resources.User

			users, err = client.UserService.List(context.TODO())
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(users).ShouldNot(gomega.BeNil())
		})
	})
})
