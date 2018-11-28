package services_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/onego-project/onego/requests"
	"github.com/onego-project/onego/services"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/onego-project/onego"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/errors"
	"github.com/onego-project/onego/resources"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	virtualNetworkAllocate               = "records/virtualNetwork/allocate"
	virtualNetworkAllocateClusterWrongID = "records/virtualNetwork/allocateClusterWrongID"
	virtualNetworkAllocateClusterNoID    = "records/virtualNetwork/allocateClusterNoID"

	virtualNetworkDelete        = "records/virtualNetwork/delete"
	virtualNetworkDeleteWrongID = "records/virtualNetwork/deleteWrongID"
	virtualNetworkDeleteNoID    = "records/virtualNetwork/deleteNoID"

	virtualNetworkUpdateMerge        = "records/virtualNetwork/updateMerge"
	virtualNetworkUpdateReplace      = "records/virtualNetwork/updateReplace"
	virtualNetworkUpdateEmptyMerge   = "records/virtualNetwork/updateEmptyMerge"
	virtualNetworkUpdateEmptyReplace = "records/virtualNetwork/updateEmptyReplace"
	virtualNetworkUpdateNoUser       = "records/virtualNetwork/updateNoUser"
	virtualNetworkUpdateUnknown      = "records/virtualNetwork/updateUnknown"

	virtualNetworkRename                 = "records/virtualNetwork/rename"
	virtualNetworkRenameEmpty            = "records/virtualNetwork/renameEmpty"
	virtualNetworkRenameUnknown          = "records/virtualNetwork/renameUnknown"
	virtualNetworkRenameNoVirtualNetwork = "records/virtualNetwork/renameNoVirtualNetwork"

	virtualNetworkChmod                 = "records/virtualNetwork/chmod"
	virtualNetworkPermRequestDefault    = "records/virtualNetwork/chmodPermReqDefault"
	virtualNetworkChmodUnknown          = "records/virtualNetwork/chmodUnknown"
	virtualNetworkChmodNoVirtualNetwork = "records/virtualNetwork/chmodNoVirtualNetwork"

	virtualNetworkChown                 = "records/virtualNetwork/chown"
	virtualNetworkOwnershipReqDefault   = "records/virtualNetwork/chownDefault"
	virtualNetworkChownUnknown          = "records/virtualNetwork/chownUnknown"
	virtualNetworkChownNoVirtualNetwork = "records/virtualNetwork/chownNoVirtualNetwork"

	virtualNetworkRetrieveInfo        = "records/virtualNetwork/retrieveInfo"
	virtualNetworkRetrieveInfoUnknown = "records/virtualNetwork/retrieveInfoUnknown"

	virtualNetworkListAllPrimaryGroup = "records/virtualNetwork/listAllPrimaryGroup"
	virtualNetworkListAllUser         = "records/virtualNetwork/listAllUser"
	virtualNetworkListAllAll          = "records/virtualNetwork/listAllAll"
	virtualNetworkListAllUserGroup    = "records/virtualNetwork/listAllUserGroup"

	virtualNetworkListAllForUser        = "records/virtualNetwork/listAllForUser"
	virtualNetworkListAllForUserUnknown = "records/virtualNetwork/listAllForUserUnknown"
	virtualNetworkListAllForUserEmpty   = "records/virtualNetwork/listAllForUserEmpty"

	virtualNetworkListPagination      = "records/virtualNetwork/listPagination"
	virtualNetworkListPaginationWrong = "records/virtualNetwork/listPaginationWrong"

	virtualNetworkListForUser        = "records/virtualNetwork/listForUser"
	virtualNetworkListForUserUnknown = "records/virtualNetwork/listForUserUnknown"
	virtualNetworkListForUserEmpty   = "records/virtualNetwork/listForUserEmpty"
)

var _ = ginkgo.Describe("Virtual Network Service", func() {
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

	ginkgo.Describe("allocate virtual network", func() {
		var virtualNetwork *resources.VirtualNetwork
		var cluster *resources.Cluster
		var virtualNetworkID int
		var oneVirtualNetwork *resources.VirtualNetwork
		var virtualNetworkBlueprint *blueprint.VirtualNetworkBlueprint

		ginkgo.Context("when virtualNetwork doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkAllocate

				virtualNetworkBlueprint = blueprint.CreateAllocateVirtualNetworkBlueprint()
				virtualNetworkBlueprint.SetElement("BRIDGE", "asdf")
				virtualNetworkBlueprint.SetElement("NAME", "my_new_vnet")
				virtualNetworkBlueprint.SetElement("VN_MAD", "test")
			})

			ginkgo.It("should create new virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				cluster = resources.CreateClusterWithID(101)

				virtualNetwork, err = client.VirtualNetworkService.Allocate(context.TODO(), virtualNetworkBlueprint,
					*cluster)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualNetwork).ShouldNot(gomega.BeNil())

				// check whether VirtualNetwork really exists in OpenNebula
				virtualNetworkID, err = virtualNetwork.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneVirtualNetwork, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), virtualNetworkID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualNetwork.Name()).To(gomega.Equal("my_new_vnet"))
			})
		})

		ginkgo.Context("when cluster has wrong ID", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkAllocateClusterWrongID
			})

			ginkgo.It("shouldn't create new virtualNetwork", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				cluster = resources.CreateClusterWithID(1552)

				virtualNetwork, err = client.VirtualNetworkService.Allocate(context.TODO(), virtualNetworkBlueprint,
					*cluster)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(virtualNetwork).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when cluster has no ID", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkAllocateClusterNoID
			})

			ginkgo.It("shouldn't create new virtualNetwork", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetwork, err = client.VirtualNetworkService.Allocate(context.TODO(), virtualNetworkBlueprint,
					resources.Cluster{})
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(virtualNetwork).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("delete virtual network", func() {
		var (
			virtualNetwork    *resources.VirtualNetwork
			oneVirtualNetwork *resources.VirtualNetwork
			virtualNetworkID  int
		)

		ginkgo.Context("when virtualNetwork exists", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkDelete

				virtualNetwork = resources.CreateVirtualNetworkWithID(11)
				if virtualNetwork == nil {
					err = errors.ErrNoVirtualNetwork
				}
			})

			ginkgo.It("should delete virtual network", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualNetworkService.Delete(context.TODO(), *virtualNetwork)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether virtualNetwork was really deleted in OpenNebula
				virtualNetworkID, err = virtualNetwork.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneVirtualNetwork, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), virtualNetworkID)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(oneVirtualNetwork).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when virtual network doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkDeleteWrongID

				virtualNetwork = resources.CreateVirtualNetworkWithID(110)
				if virtualNetwork == nil {
					err = errors.ErrNoVirtualNetwork
				}
			})

			ginkgo.It("should return that virtual network with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualNetworkService.Delete(context.TODO(), *virtualNetwork)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtual network is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkDeleteNoID

				virtualNetwork = &resources.VirtualNetwork{}
			})

			ginkgo.It("should return that virtualNetwork has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualNetworkService.Delete(context.TODO(), *virtualNetwork)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("update virtual network", func() {
		var (
			virtualNetwork          *resources.VirtualNetwork
			virtualNetworkBlueprint *blueprint.VirtualNetworkBlueprint
			retVN                   *resources.VirtualNetwork
		)

		ginkgo.Context("when virtual network exists", func() {
			ginkgo.Context("when update data is not empty", func() {
				ginkgo.BeforeEach(func() {
					virtualNetwork = resources.CreateVirtualNetworkWithID(12)
					if virtualNetwork == nil {
						err = errors.ErrNoVirtualNetwork
						return
					}
				})

				ginkgo.When("when merge data of given virtualNetwork", func() {
					ginkgo.BeforeEach(func() {
						recName = virtualNetworkUpdateMerge
					})

					ginkgo.It("should merge data of given virtual network", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						virtualNetworkBlueprint = blueprint.CreateUpdateVirtualNetworkBlueprint()
						if virtualNetworkBlueprint == nil {
							err = errors.ErrNoVirtualNetworkBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						virtualNetworkBlueprint.SetGateway("dummy")

						retVN, err = client.VirtualNetworkService.Update(context.TODO(), *virtualNetwork,
							virtualNetworkBlueprint, services.Merge)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						gomega.Expect(retVN).ShouldNot(gomega.BeNil())
						gomega.Expect(retVN.Attribute("TEMPLATE/GATEWAY")).To(gomega.Equal(
							"dummy"))
					})
				})

				ginkgo.When("when replace data of given virtual network", func() {
					ginkgo.BeforeEach(func() {
						recName = virtualNetworkUpdateReplace
					})

					ginkgo.It("should replace data of given virtual network", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						virtualNetworkBlueprint = blueprint.CreateUpdateVirtualNetworkBlueprint()
						if virtualNetworkBlueprint == nil {
							err = errors.ErrNoVirtualNetworkBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						virtualNetworkBlueprint.SetVnMad("blabla")

						retVN, err = client.VirtualNetworkService.Update(context.TODO(), *virtualNetwork,
							virtualNetworkBlueprint, services.Replace)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						gomega.Expect(retVN).ShouldNot(gomega.BeNil())
						gomega.Expect(retVN.Attribute("TEMPLATE/VN_MAD")).To(gomega.Equal(
							"blabla"))
					})
				})
			})

			ginkgo.Context("when update data is empty", func() {
				ginkgo.BeforeEach(func() {
					virtualNetwork = resources.CreateVirtualNetworkWithID(12)
					if virtualNetwork == nil {
						err = errors.ErrNoVirtualNetwork
						return
					}

					virtualNetworkBlueprint = &blueprint.VirtualNetworkBlueprint{}
				})

				ginkgo.When("when merge data of given virtual network", func() {
					ginkgo.BeforeEach(func() {
						recName = virtualNetworkUpdateEmptyMerge
					})

					ginkgo.It("should return an error", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						retVN, err = client.VirtualNetworkService.Update(context.TODO(), *virtualNetwork,
							virtualNetworkBlueprint, services.Merge)
						gomega.Expect(err).To(gomega.HaveOccurred())
						gomega.Expect(retVN).Should(gomega.BeNil())
					})
				})

				ginkgo.When("when replace data of given virtualnNetwork", func() {
					ginkgo.BeforeEach(func() {
						recName = virtualNetworkUpdateEmptyReplace
					})

					ginkgo.It("should return an error", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						retVN, err = client.VirtualNetworkService.Update(context.TODO(), *virtualNetwork,
							virtualNetworkBlueprint, services.Replace)
						gomega.Expect(err).To(gomega.HaveOccurred())
						gomega.Expect(retVN).Should(gomega.BeNil())
					})
				})
			})
		})

		ginkgo.Context("when virtualNetwork doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkUpdateUnknown

				virtualNetwork = resources.CreateVirtualNetworkWithID(110)
				if virtualNetwork == nil {
					err = errors.ErrNoVirtualNetwork
				}

				virtualNetworkBlueprint = blueprint.CreateUpdateVirtualNetworkBlueprint()
				if virtualNetworkBlueprint == nil {
					err = errors.ErrNoVirtualNetworkBlueprint
					return
				}
				virtualNetworkBlueprint.SetGateway("dummy")
			})

			ginkgo.It("should return that virtualNetwork with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				retVN, err = client.VirtualNetworkService.Update(context.TODO(), *virtualNetwork, virtualNetworkBlueprint,
					services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(retVN).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when virtualNetwork is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkUpdateNoUser

				virtualNetworkBlueprint = blueprint.CreateUpdateVirtualNetworkBlueprint()
				if virtualNetworkBlueprint == nil {
					err = errors.ErrNoVirtualNetworkBlueprint
					return
				}
				virtualNetworkBlueprint.SetGateway("dummy")
			})

			ginkgo.It("should return that virtualNetwork has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				retVN, err = client.VirtualNetworkService.Update(context.TODO(), resources.VirtualNetwork{},
					virtualNetworkBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(retVN).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("virtualNetwork rename", func() {
		var (
			virtualNetwork    *resources.VirtualNetwork
			oneVirtualNetwork *resources.VirtualNetwork
			virtualNetworkID  int
		)

		ginkgo.Context("when virtual network exists", func() {
			ginkgo.BeforeEach(func() {
				virtualNetwork = resources.CreateVirtualNetworkWithID(12)
				if virtualNetwork == nil {
					err = errors.ErrNoVirtualNetwork
				}
			})

			ginkgo.When("when new name is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualNetworkRename
				})

				ginkgo.It("should change name of given virtual network", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					// get virtualNetwork name
					oneVirtualNetwork, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), 12)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVirtualNetwork).ShouldNot(gomega.BeNil())

					newName := "my_vnet"
					gomega.Expect(oneVirtualNetwork.Name()).NotTo(gomega.Equal(newName))

					// change name
					err = client.VirtualNetworkService.Rename(context.TODO(), *virtualNetwork, newName)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether name was really changed in OpenNebula
					virtualNetworkID, err = virtualNetwork.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneVirtualNetwork, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), virtualNetworkID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVirtualNetwork).ShouldNot(gomega.BeNil())

					gomega.Expect(oneVirtualNetwork.Name()).To(gomega.Equal("the_best_virtualNetwork_ever"))
				})
			})

			ginkgo.When("when new name is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualNetworkRenameEmpty
				})

				ginkgo.It("should not change name of given virtualNetwork", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.VirtualNetworkService.Rename(context.TODO(), *virtualNetwork, "")
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when virtualNetwork doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkRenameUnknown

				virtualNetwork = resources.CreateVirtualNetworkWithID(110)
				if virtualNetwork == nil {
					err = errors.ErrNoVirtualNetwork
				}
			})

			ginkgo.It("should return that virtualNetwork with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualNetworkService.Rename(context.TODO(), *virtualNetwork, "mastodont")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtualNetwork is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkRenameNoVirtualNetwork

				virtualNetwork = &resources.VirtualNetwork{}
			})

			ginkgo.It("should return that virtualNetwork has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualNetworkService.Rename(context.TODO(), *virtualNetwork, "rex")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("virtual network chmod", func() {
		var (
			virtualNetwork    *resources.VirtualNetwork
			oneVirtualNetwork *resources.VirtualNetwork
			virtualNetworkID  int
			permRequest       requests.PermissionRequest
		)

		ginkgo.Context("when virtual network exists", func() {
			ginkgo.BeforeEach(func() {
				virtualNetwork = resources.CreateVirtualNetworkWithID(12)
				if virtualNetwork == nil {
					err = errors.ErrNoVirtualNetwork
				}
			})

			ginkgo.When("when permission request is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualNetworkChmod
				})

				ginkgo.It("should change permission of given virtualNetwork", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					permRequest = requests.CreatePermissionRequestBuilder().Deny(requests.User,
						requests.Manage).Allow(requests.Other, requests.Admin).Build()

					err = client.VirtualNetworkService.Chmod(context.TODO(), *virtualNetwork, permRequest)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chmod was really changed in OpenNebula
					virtualNetworkID, err = virtualNetwork.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneVirtualNetwork, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), virtualNetworkID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVirtualNetwork).ShouldNot(gomega.BeNil())

					var perm *resources.Permissions
					perm, err = oneVirtualNetwork.Permissions()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					gomega.Expect(perm.User.Manage).To(gomega.Equal(false))
					gomega.Expect(perm.Other.Admin).To(gomega.Equal(true))
				})
			})

			ginkgo.When("when permission request is default", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualNetworkPermRequestDefault
				})

				ginkgo.It("should not change permissions of given virtual network", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.VirtualNetworkService.Chmod(context.TODO(), *virtualNetwork,
						requests.CreatePermissionRequestBuilder().Build())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when virtual network doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkChmodUnknown

				virtualNetwork = resources.CreateVirtualNetworkWithID(110)
				if virtualNetwork == nil {
					err = errors.ErrNoVirtualNetwork
				}
			})

			ginkgo.It("should return that virtual network with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				permRequest = requests.CreatePermissionRequestBuilder().Allow(requests.User,
					requests.Manage).Deny(requests.Other, requests.Admin).Build()

				err = client.VirtualNetworkService.Chmod(context.TODO(), *virtualNetwork, permRequest)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtual network is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkChmodNoVirtualNetwork

				virtualNetwork = &resources.VirtualNetwork{}
			})

			ginkgo.It("should return that virtual network has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualNetworkService.Chmod(context.TODO(), *virtualNetwork, requests.PermissionRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("virtual network chown", func() {
		var (
			virtualNetwork    *resources.VirtualNetwork
			oneVirtualNetwork *resources.VirtualNetwork
			virtualNetworkID  int

			user         *resources.User
			group        *resources.Group
			ownershipReq requests.OwnershipRequest
		)

		ginkgo.Context("when virtual network exists", func() {
			ginkgo.BeforeEach(func() {
				virtualNetwork = resources.CreateVirtualNetworkWithID(12)
				if virtualNetwork == nil {
					err = errors.ErrNoVirtualNetwork
				}
			})

			ginkgo.When("when ownership request is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualNetworkChown
				})

				ginkgo.It("should change owner of given virtual network", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					user = resources.CreateUserWithID(31)
					group = resources.CreateGroupWithID(120)

					ownershipReq = requests.CreateOwnershipRequestBuilder().User(*user).Group(*group).Build()

					err = client.VirtualNetworkService.Chown(context.TODO(), *virtualNetwork, ownershipReq)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chown was really changed in OpenNebula
					virtualNetworkID, err = virtualNetwork.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneVirtualNetwork, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), virtualNetworkID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVirtualNetwork).ShouldNot(gomega.BeNil())

					gomega.Expect(oneVirtualNetwork.User()).To(gomega.Equal(31))
					gomega.Expect(oneVirtualNetwork.Group()).To(gomega.Equal(120))
				})
			})

			ginkgo.When("when ownership request is default", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualNetworkOwnershipReqDefault
				})

				ginkgo.It("should not change permissions of given virtualNetwork", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.VirtualNetworkService.Chown(context.TODO(), *virtualNetwork,
						requests.CreateOwnershipRequestBuilder().Build())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when virtual network doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkChownUnknown

				virtualNetwork = resources.CreateVirtualNetworkWithID(110)
				if virtualNetwork == nil {
					err = errors.ErrNoVirtualNetwork
				}
			})

			ginkgo.It("should return that virtual network with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualNetworkService.Chown(context.TODO(), *virtualNetwork, requests.OwnershipRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtualNetwork is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkChownNoVirtualNetwork

				virtualNetwork = &resources.VirtualNetwork{}
			})

			ginkgo.It("should return that virtualNetwork has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualNetworkService.Chown(context.TODO(), *virtualNetwork, requests.OwnershipRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("virtual network retrieve info", func() {
		var virtualNetwork *resources.VirtualNetwork

		ginkgo.Context("when virtual network exists", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkRetrieveInfo
			})

			ginkgo.It("should return virtual network with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetwork, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), 12)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualNetwork).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualNetwork.ID()).To(gomega.Equal(12))
				gomega.Expect(virtualNetwork.Name()).To(gomega.Equal("my_vnet"))
			})
		})

		ginkgo.Context("when virtual network doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkRetrieveInfoUnknown
			})

			ginkgo.It("should return that given virtual network doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetwork, err = client.VirtualNetworkService.RetrieveInfo(context.TODO(), 110)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(virtualNetwork).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("virtual network list all", func() {
		var virtualNetworks []*resources.VirtualNetwork

		ginkgo.Context("when ownership filter is not empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkListAllPrimaryGroup
			})

			ginkgo.It("should return list of all virtual network with full info belongs to primary group", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetworks, err = client.VirtualNetworkService.ListAll(context.TODO(),
					services.OwnershipFilterPrimaryGroup)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualNetworks).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualNetworks).To(gomega.HaveLen(6))
			})
		})

		ginkgo.Context("when ownership filter is not empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkListAllUser
			})

			ginkgo.It("should return list of all virtual network with full info belongs to the user", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetworks, err = client.VirtualNetworkService.ListAll(context.TODO(),
					services.OwnershipFilterUser)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualNetworks).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualNetworks).To(gomega.HaveLen(6))
				gomega.Expect(virtualNetworks[0].ID()).To(gomega.Equal(0))
				gomega.Expect(virtualNetworks[1].ID()).To(gomega.Equal(2))
				gomega.Expect(virtualNetworks[2].ID()).To(gomega.Equal(3))
				gomega.Expect(virtualNetworks[3].ID()).To(gomega.Equal(7))
				gomega.Expect(virtualNetworks[4].ID()).To(gomega.Equal(9))
				gomega.Expect(virtualNetworks[5].ID()).To(gomega.Equal(10))
			})
		})

		ginkgo.Context("when ownership filter is set to all", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkListAllAll
			})

			ginkgo.It("should return list of all virtual network with full info belongs to all", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetworks, err = client.VirtualNetworkService.ListAll(context.TODO(), services.OwnershipFilterAll)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualNetworks).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualNetworks).To(gomega.HaveLen(7))
			})
		})

		ginkgo.Context("when ownership filter is set to UserGroup", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkListAllUserGroup
			})

			ginkgo.It("should return list of all virtual network with full info belongs to the user and any "+
				"of his groups", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetworks, err = client.VirtualNetworkService.ListAll(context.TODO(),
					services.OwnershipFilterUserGroup)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualNetworks).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualNetworks).To(gomega.HaveLen(7))
			})
		})
	})

	ginkgo.Describe("virtual network list all for user", func() {
		var virtualNetworks []*resources.VirtualNetwork

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkListAllForUser
			})

			ginkgo.It("should return virtual network with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetworks, err = client.VirtualNetworkService.ListAllForUser(context.TODO(),
					*resources.CreateUserWithID(31))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualNetworks).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualNetworks).To(gomega.HaveLen(1))
				gomega.Expect(virtualNetworks[0].ID()).To(gomega.Equal(12))
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkListAllForUserUnknown
			})

			ginkgo.It("should return empty list of virtual network (length 0)", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetworks, err = client.VirtualNetworkService.ListAllForUser(context.TODO(),
					*resources.CreateUserWithID(42))
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualNetworks).Should(gomega.Equal(make([]*resources.VirtualNetwork, 0)))
				gomega.Expect(virtualNetworks).Should(gomega.HaveLen(0))
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkListAllForUserEmpty
			})

			ginkgo.It("should return that user doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetworks, err = client.VirtualNetworkService.ListAllForUser(context.TODO(), resources.User{})
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(virtualNetworks).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("List methods with pagination", func() {
		var (
			virtualNetworks []*resources.VirtualNetwork
		)

		ginkgo.Context("pagination ok", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkListPagination
			})

			ginkgo.It("should return virtual network with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetworks, err = client.VirtualNetworkService.List(context.TODO(), 3,
					2, services.OwnershipFilterUserGroup)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualNetworks).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualNetworks).To(gomega.HaveLen(2))
				gomega.Expect(virtualNetworks[0].ID()).To(gomega.Equal(9))
				gomega.Expect(virtualNetworks[1].ID()).To(gomega.Equal(10))
			})
		})

		ginkgo.Context("pagination wrong", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkListPaginationWrong
			})
		})

		ginkgo.It("should return that pagination is wrong", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

			virtualNetworks, err = client.VirtualNetworkService.List(context.TODO(), -2,
				-2, services.OwnershipFilterPrimaryGroup)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(virtualNetworks).Should(gomega.BeNil())
		})
	})

	ginkgo.Describe("virtual network list for user", func() {
		var virtualNetworks []*resources.VirtualNetwork

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkListForUser
			})

			ginkgo.It("should return virtual network with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetworks, err = client.VirtualNetworkService.ListForUser(context.TODO(),
					*resources.CreateUserWithID(0), 2, 2)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualNetworks).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualNetworks).To(gomega.HaveLen(2))
				gomega.Expect(virtualNetworks[0].ID()).To(gomega.Equal(3))
				gomega.Expect(virtualNetworks[1].ID()).To(gomega.Equal(7))
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkListForUserUnknown
			})

			ginkgo.It("should return empty list of virtual network (length 0)", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetworks, err = client.VirtualNetworkService.ListForUser(context.TODO(),
					*resources.CreateUserWithID(88), 2, 2)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualNetworks).Should(gomega.Equal(make([]*resources.VirtualNetwork, 0)))
				gomega.Expect(virtualNetworks).Should(gomega.HaveLen(0))
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualNetworkListForUserEmpty
			})

			ginkgo.It("should return that user doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualNetworks, err = client.VirtualNetworkService.ListForUser(context.TODO(),
					resources.User{}, 2, 2)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(virtualNetworks).Should(gomega.BeNil())
			})
		})
	})
})
