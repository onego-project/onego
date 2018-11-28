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
	"github.com/onego-project/onego/resources"
	"github.com/onego-project/onego/services"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	hostAllocate               = "records/host/allocate"
	hostAllocateFail           = "records/host/allocateFail"
	hostAllocateClusterWrongID = "records/host/allocateClusterWrongID"
	hostAllocateClusterNoID    = "records/host/allocateClusterNoID"

	hostDelete        = "records/host/delete"
	hostDeleteWrongID = "records/host/deleteWrongID"
	hostDeleteNoID    = "records/host/deleteNoID"

	hostUpdateMerge        = "records/host/updateMerge"
	hostUpdateReplace      = "records/host/updateReplace"
	hostUpdateEmptyMerge   = "records/host/updateEmptyMerge"
	hostUpdateEmptyReplace = "records/host/updateEmptyReplace"
	hostUpdateNoUser       = "records/host/updateNoUser"
	hostUpdateUnknown      = "records/host/updateUnknown"

	hostRename        = "records/host/rename"
	hostRenameEmpty   = "records/host/renameEmpty"
	hostRenameUnknown = "records/host/renameUnknown"
	hostRenameNoHost  = "records/host/renameNoHost"

	hostEnable         = "records/host/enable"
	hostDisable        = "records/host/disable"
	hostOffline        = "records/host/offline"
	hostDisableUnknown = "records/host/disableUnknown"
	hostDisableNoHost  = "records/host/disableNoHost"

	hostRetrieveInfo        = "records/host/retrieveInfo"
	hostRetrieveInfoUnknown = "records/host/retrieveInfoUnknown"

	hostList = "records/host/list"
)

var _ = ginkgo.Describe("Host Service", func() {
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

	ginkgo.Describe("allocate host", func() {
		var host *resources.Host
		var cluster *resources.Cluster
		var hostID int
		var oneHost *resources.Host

		ginkgo.Context("when host doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = hostAllocate
			})

			ginkgo.It("should create new host", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				cluster = resources.CreateClusterWithID(101)

				host, err = client.HostService.Allocate(context.TODO(), "the_best_host",
					"kvm", "kvm", *cluster)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(host).ShouldNot(gomega.BeNil())

				// check whether Host really exists in OpenNebula
				hostID, err = host.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneHost, err = client.HostService.RetrieveInfo(context.TODO(), hostID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneHost.Name()).To(gomega.Equal("the_best_host"))
			})
		})

		ginkgo.Context("when cluster has wrong ID", func() {
			ginkgo.BeforeEach(func() {
				recName = hostAllocateClusterWrongID
			})

			ginkgo.It("shouldn't create new host", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				cluster = resources.CreateClusterWithID(1552)

				host, err = client.HostService.Allocate(context.TODO(), "the_best_host2",
					"kvm", "kvm", *cluster)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(host).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when cluster has no ID", func() {
			ginkgo.BeforeEach(func() {
				recName = hostAllocateClusterNoID
			})

			ginkgo.It("shouldn't create new host", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				host, err = client.HostService.Allocate(context.TODO(), "the_best_host3",
					"kvm", "kvm", resources.Cluster{})
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(host).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when host exists", func() {
			ginkgo.BeforeEach(func() {
				recName = hostAllocateFail
			})

			ginkgo.It("shouldn't create new host", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				host, err = client.HostService.Allocate(context.TODO(), "the_best_host4",
					"kvm", "kvm", resources.Cluster{})
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(host).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("delete host", func() {
		var (
			host    *resources.Host
			oneHost *resources.Host
			hostID  int
		)

		ginkgo.Context("when host exists", func() {
			ginkgo.BeforeEach(func() {
				recName = hostDelete

				host = resources.CreateHostWithID(5)
				if host == nil {
					err = errors.ErrNoHost
				}
			})

			ginkgo.It("should delete host", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.HostService.Delete(context.TODO(), *host)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether host was really deleted in OpenNebula
				hostID, err = host.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneHost, err = client.HostService.RetrieveInfo(context.TODO(), hostID)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(oneHost).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when host doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = hostDeleteWrongID

				host = resources.CreateHostWithID(110)
				if host == nil {
					err = errors.ErrNoHost
				}
			})

			ginkgo.It("should return that host with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.HostService.Delete(context.TODO(), *host)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when host is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = hostDeleteNoID

				host = &resources.Host{}
			})

			ginkgo.It("should return that host has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.HostService.Delete(context.TODO(), *host)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("update host", func() {
		var (
			host          *resources.Host
			hostBlueprint *blueprint.HostBlueprint
			retHost       *resources.Host
		)

		ginkgo.Context("when host exists", func() {
			ginkgo.Context("when update data is not empty", func() {
				ginkgo.BeforeEach(func() {
					host = resources.CreateHostWithID(6)
					if host == nil {
						err = errors.ErrNoHost
						return
					}
				})

				ginkgo.When("when merge data of given host", func() {
					ginkgo.BeforeEach(func() {
						recName = hostUpdateMerge
					})

					ginkgo.It("should merge data of given host", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						hostBlueprint = blueprint.CreateUpdateHostBlueprint()
						if hostBlueprint == nil {
							err = errors.ErrNoHostBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						hostBlueprint.SetClusterName("dummy")

						retHost, err = client.HostService.Update(context.TODO(), *host, hostBlueprint, services.Merge)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						gomega.Expect(retHost).ShouldNot(gomega.BeNil())
						gomega.Expect(retHost.Attribute("TEMPLATE/CLUSTER_NAME")).To(gomega.Equal("dummy"))
					})
				})

				ginkgo.When("when replace data of given host", func() {
					ginkgo.BeforeEach(func() {
						recName = hostUpdateReplace
					})

					ginkgo.It("should replace data of given host", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						hostBlueprint = blueprint.CreateUpdateHostBlueprint()
						if hostBlueprint == nil {
							err = errors.ErrNoHostBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						hostBlueprint.SetHostName("blabla")

						retHost, err = client.HostService.Update(context.TODO(), *host, hostBlueprint, services.Replace)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						gomega.Expect(retHost).ShouldNot(gomega.BeNil())
						gomega.Expect(retHost.Attribute("TEMPLATE/HOSTNAME")).To(gomega.Equal("blabla"))
					})
				})
			})

			ginkgo.Context("when update data is empty", func() {
				ginkgo.BeforeEach(func() {
					host = resources.CreateHostWithID(6)
					if host == nil {
						err = errors.ErrNoHost
						return
					}

					hostBlueprint = &blueprint.HostBlueprint{}
				})

				ginkgo.When("when merge data of given host", func() {
					ginkgo.BeforeEach(func() {
						recName = hostUpdateEmptyMerge
					})

					ginkgo.It("should return an error", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						retHost, err = client.HostService.Update(context.TODO(), *host, hostBlueprint, services.Merge)
						gomega.Expect(err).To(gomega.HaveOccurred())
						gomega.Expect(retHost).Should(gomega.BeNil())
					})
				})

				ginkgo.When("when replace data of given host", func() {
					ginkgo.BeforeEach(func() {
						recName = hostUpdateEmptyReplace
					})

					ginkgo.It("should return an error", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						retHost, err = client.HostService.Update(context.TODO(), *host, hostBlueprint, services.Replace)
						gomega.Expect(err).To(gomega.HaveOccurred())
						gomega.Expect(retHost).Should(gomega.BeNil())
					})
				})
			})
		})

		ginkgo.Context("when host doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = hostUpdateUnknown

				host = resources.CreateHostWithID(110)
				if host == nil {
					err = errors.ErrNoHost
				}

				hostBlueprint = blueprint.CreateUpdateHostBlueprint()
				if hostBlueprint == nil {
					err = errors.ErrNoHostBlueprint
					return
				}
				hostBlueprint.SetClusterName("dummy")
			})

			ginkgo.It("should return that host with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				retHost, err = client.HostService.Update(context.TODO(), *host, hostBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(retHost).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when host is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = hostUpdateNoUser

				hostBlueprint = blueprint.CreateUpdateHostBlueprint()
				if hostBlueprint == nil {
					err = errors.ErrNoHostBlueprint
					return
				}
				hostBlueprint.SetClusterName("dummy")
			})

			ginkgo.It("should return that host has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				retHost, err = client.HostService.Update(context.TODO(), resources.Host{},
					hostBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(retHost).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("host rename", func() {
		var (
			host    *resources.Host
			oneHost *resources.Host
			hostID  int
		)

		ginkgo.Context("when host exists", func() {
			ginkgo.BeforeEach(func() {
				host = resources.CreateHostWithID(6)
				if host == nil {
					err = errors.ErrNoHost
				}
			})

			ginkgo.When("when new name is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = hostRename
				})

				ginkgo.It("should change name of given host", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					// get host name
					oneHost, err = client.HostService.RetrieveInfo(context.TODO(), 6)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneHost).ShouldNot(gomega.BeNil())

					newName := "the_best_host_ever"
					gomega.Expect(oneHost.Name()).NotTo(gomega.Equal(newName))

					// change name
					err = client.HostService.Rename(context.TODO(), *host, newName)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether name was really changed in OpenNebula
					hostID, err = host.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneHost, err = client.HostService.RetrieveInfo(context.TODO(), hostID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneHost).ShouldNot(gomega.BeNil())

					gomega.Expect(oneHost.Name()).To(gomega.Equal("the_best_host"))
				})
			})

			ginkgo.When("when new name is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = hostRenameEmpty
				})

				ginkgo.It("should not change name of given host", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.HostService.Rename(context.TODO(), *host, "")
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when host doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = hostRenameUnknown

				host = resources.CreateHostWithID(110)
				if host == nil {
					err = errors.ErrNoHost
				}
			})

			ginkgo.It("should return that host with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.HostService.Rename(context.TODO(), *host, "mastodont")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when host is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = hostRenameNoHost

				host = &resources.Host{}
			})

			ginkgo.It("should return that host has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.HostService.Rename(context.TODO(), *host, "rex")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("host enable", func() {
		var (
			host    *resources.Host
			oneHost *resources.Host
			hostID  int
		)

		ginkgo.Context("when host exists", func() {
			ginkgo.BeforeEach(func() {
				host = resources.CreateHostWithID(6)
				if host == nil {
					err = errors.ErrNoHost
				}
			})

			ginkgo.When("when enable", func() {
				ginkgo.BeforeEach(func() {
					recName = hostEnable
				})

				ginkgo.It("should enable given host", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.HostService.Enable(context.TODO(), *host)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether state was really changed in OpenNebula
					hostID, err = host.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneHost, err = client.HostService.RetrieveInfo(context.TODO(), hostID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneHost).ShouldNot(gomega.BeNil())

					gomega.Expect(oneHost.State()).To(gomega.Equal(resources.HostInit))
				})
			})

			ginkgo.When("when offline", func() {
				ginkgo.BeforeEach(func() {
					recName = hostOffline
				})

				ginkgo.It("should power off given host", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.HostService.Offline(context.TODO(), *host)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether state was really changed in OpenNebula
					hostID, err = host.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneHost, err = client.HostService.RetrieveInfo(context.TODO(), hostID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneHost).ShouldNot(gomega.BeNil())

					gomega.Expect(oneHost.State()).To(gomega.Equal(resources.HostOffline))
				})
			})

			ginkgo.When("when disable", func() {
				ginkgo.BeforeEach(func() {
					recName = hostDisable
				})

				ginkgo.It("should disable given host", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.HostService.Disable(context.TODO(), *host)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether state was really changed in OpenNebula
					hostID, err = host.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneHost, err = client.HostService.RetrieveInfo(context.TODO(), hostID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneHost).ShouldNot(gomega.BeNil())

					gomega.Expect(oneHost.State()).To(gomega.Equal(resources.HostDisabled))
				})
			})
		})

		ginkgo.Context("when host doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = hostDisableUnknown

				host = resources.CreateHostWithID(110)
				if host == nil {
					err = errors.ErrNoHost
				}
			})

			ginkgo.It("should return that host with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.HostService.Disable(context.TODO(), *host)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when host is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = hostDisableNoHost

				host = &resources.Host{}
			})

			ginkgo.It("should return that host has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.HostService.Disable(context.TODO(), *host)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("host retrieve info", func() {
		var host *resources.Host

		ginkgo.Context("when host exists", func() {
			ginkgo.BeforeEach(func() {
				recName = hostRetrieveInfo
			})

			ginkgo.It("should return host with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				host, err = client.HostService.RetrieveInfo(context.TODO(), 6)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(host).ShouldNot(gomega.BeNil())
				gomega.Expect(host.ID()).To(gomega.Equal(6))
				gomega.Expect(host.Name()).To(gomega.Equal("the_best_host_ever"))
			})
		})

		ginkgo.Context("when host doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = hostRetrieveInfoUnknown
			})

			ginkgo.It("should return that given host doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				host, err = client.HostService.RetrieveInfo(context.TODO(), 110)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(host).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("host list", func() {
		ginkgo.BeforeEach(func() {
			recName = hostList
		})

		ginkgo.It("should return an array of hosts with full info", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

			var hosts []*resources.Host

			hosts, err = client.HostService.List(context.TODO())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(hosts).ShouldNot(gomega.BeNil())
		})
	})
})
