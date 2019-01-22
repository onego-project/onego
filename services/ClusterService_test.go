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
	clusterAllocate     = "records/cluster/allocate"
	clusterAllocateFail = "records/cluster/allocateFail"

	clusterDelete        = "records/cluster/delete"
	clusterDeleteWrongID = "records/cluster/deleteWrongID"
	clusterDeleteNoID    = "records/cluster/deleteNoID"

	clusterUpdateMerge        = "records/cluster/updateMerge"
	clusterUpdateReplace      = "records/cluster/updateReplace"
	clusterUpdateEmptyMerge   = "records/cluster/updateEmptyMerge"
	clusterUpdateEmptyReplace = "records/cluster/updateEmptyReplace"
	clusterUpdateNoUser       = "records/cluster/updateNoUser"
	clusterUpdateUnknown      = "records/cluster/updateUnknown"

	clusterAddHost          = "records/cluster/addHost"
	clusterAddNoneHost      = "records/cluster/addNoneHost"
	clusterAddHostUnknown   = "records/cluster/addHostUnknown"
	clusterAddHostNoCluster = "records/cluster/addHostNoCluster"

	clusterDelHost          = "records/cluster/delHost"
	clusterDelNoneHost      = "records/cluster/delNoneHost"
	clusterDelHostUnknown   = "records/cluster/delHostUnknown"
	clusterDelHostNoCluster = "records/cluster/delHostNoCluster"

	clusterAddDatastore          = "records/cluster/addDatastore"
	clusterAddNoneDatastore      = "records/cluster/addNoneDatastore"
	clusterAddDatastoreUnknown   = "records/cluster/addDatastoreUnknown"
	clusterAddDatastoreNoCluster = "records/cluster/addDatastoreNoCluster"

	clusterDelDatastore          = "records/cluster/delDatastore"
	clusterDelNoneDatastore      = "records/cluster/delNoneDatastore"
	clusterDelDatastoreUnknown   = "records/cluster/delDatastoreUnknown"
	clusterDelDatastoreNoCluster = "records/cluster/delDatastoreNoCluster"

	clusterAddVirtualNetwork          = "records/cluster/addVirtualNetwork"
	clusterAddNoneVirtualNetwork      = "records/cluster/addNoneVirtualNetwork"
	clusterAddVirtualNetworkUnknown   = "records/cluster/addVirtualNetworkUnknown"
	clusterAddVirtualNetworkNoCluster = "records/cluster/addVirtualNetworkNoCluster"

	clusterDelVirtualNetwork          = "records/cluster/delVirtualNetwork"
	clusterDelNoneVirtualNetwork      = "records/cluster/delNoneVirtualNetwork"
	clusterDelVirtualNetworkUnknown   = "records/cluster/delVirtualNetworkUnknown"
	clusterDelVirtualNetworkNoCluster = "records/cluster/delVirtualNetworkNoCluster"

	clusterRename          = "records/cluster/rename"
	clusterRenameEmpty     = "records/cluster/renameEmpty"
	clusterRenameUnknown   = "records/cluster/renameUnknown"
	clusterRenameNoCluster = "records/cluster/renameNoCluster"

	clusterRetrieveInfo        = "records/cluster/retrieveInfo"
	clusterRetrieveInfoUnknown = "records/cluster/retrieveInfoUnknown"

	clusterList = "records/cluster/list"
)

var _ = ginkgo.Describe("Cluster Service", func() {
	var (
		recName string
		rec     *recorder.Recorder
		client  *onego.Client
		err     error
	)

	var existingClusterID = 101
	var nonExistingClusterID = 110
	var allocatedDeletedClusterID = 102

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

	ginkgo.Describe("allocate cluster", func() {
		var cluster *resources.Cluster
		var clusterID int
		var oneCluster *resources.Cluster

		clusterName := "my_cluster"

		ginkgo.Context("when cluster doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterAllocate
			})

			ginkgo.It("should create new cluster", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				cluster, err = client.ClusterService.Allocate(context.TODO(), clusterName)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(cluster).ShouldNot(gomega.BeNil())

				// check whether Cluster really exists in OpenNebula
				clusterID, err = cluster.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneCluster, err = client.ClusterService.RetrieveInfo(context.TODO(), clusterID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneCluster.Name()).To(gomega.Equal(clusterName))
			})
		})

		ginkgo.Context("when cluster exists", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterAllocateFail
			})

			ginkgo.It("shouldn't create new cluster", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				cluster, err = client.ClusterService.Allocate(context.TODO(), clusterName)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(cluster).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("delete cluster", func() {
		var (
			cluster    *resources.Cluster
			oneCluster *resources.Cluster
			clusterID  int
		)

		ginkgo.Context("when cluster exists", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterDelete

				cluster = resources.CreateClusterWithID(allocatedDeletedClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.It("should delete cluster", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.Delete(context.TODO(), *cluster)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether cluster was really deleted in OpenNebula
				clusterID, err = cluster.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneCluster, err = client.ClusterService.RetrieveInfo(context.TODO(), clusterID)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(oneCluster).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("when cluster doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterDeleteWrongID

				cluster = resources.CreateClusterWithID(nonExistingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.It("should return that cluster with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.Delete(context.TODO(), *cluster)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when cluster is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterDeleteNoID

				cluster = &resources.Cluster{}
			})

			ginkgo.It("should return that cluster has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.Delete(context.TODO(), *cluster)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("update cluster", func() {
		var (
			cluster          *resources.Cluster
			clusterBlueprint *blueprint.ClusterBlueprint
			retCluster       *resources.Cluster
		)

		ginkgo.Context("when cluster exists", func() {
			ginkgo.Context("when update data is not empty", func() {
				ginkgo.BeforeEach(func() {
					cluster = resources.CreateClusterWithID(existingClusterID)
					if cluster == nil {
						err = errors.ErrNoCluster
						return
					}
				})

				ginkgo.When("when merge data of given cluster", func() {
					ginkgo.BeforeEach(func() {
						recName = clusterUpdateMerge
					})

					ginkgo.It("should merge data of given cluster", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						clusterBlueprint = blueprint.CreateUpdateClusterBlueprint()
						if clusterBlueprint == nil {
							err = errors.ErrNoClusterBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						clusterBlueprint.SetReservedMemory(123)
						clusterBlueprint.SetReservedCPU(321)

						retCluster, err = client.ClusterService.Update(context.TODO(), *cluster, clusterBlueprint,
							services.Merge)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						gomega.Expect(retCluster).ShouldNot(gomega.BeNil())
						gomega.Expect(retCluster.Attribute("TEMPLATE/RESERVED_MEM")).To(
							gomega.Equal("123"))
						gomega.Expect(retCluster.Attribute("TEMPLATE/RESERVED_CPU")).To(
							gomega.Equal("321"))
					})
				})

				ginkgo.When("when replace data of given cluster", func() {
					ginkgo.BeforeEach(func() {
						recName = clusterUpdateReplace
					})

					ginkgo.It("should replace data of given cluster", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						clusterBlueprint = blueprint.CreateUpdateClusterBlueprint()
						if clusterBlueprint == nil {
							err = errors.ErrNoClusterBlueprint
							gomega.Expect(err).NotTo(gomega.HaveOccurred())
						}
						clusterBlueprint.SetReservedCPU(123)

						retCluster, err = client.ClusterService.Update(context.TODO(), *cluster, clusterBlueprint,
							services.Replace)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						gomega.Expect(retCluster).ShouldNot(gomega.BeNil())
						gomega.Expect(retCluster.Attribute("TEMPLATE/RESERVED_CPU")).To(
							gomega.Equal("123"))
					})
				})
			})

			ginkgo.Context("when update data is empty", func() {
				ginkgo.BeforeEach(func() {
					cluster = resources.CreateClusterWithID(existingClusterID)
					if cluster == nil {
						err = errors.ErrNoCluster
						return
					}

					clusterBlueprint = &blueprint.ClusterBlueprint{}
				})

				ginkgo.When("when merge data of given cluster", func() {
					ginkgo.BeforeEach(func() {
						recName = clusterUpdateEmptyMerge
					})

					ginkgo.It("should return an error", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						retCluster, err = client.ClusterService.Update(context.TODO(), *cluster, clusterBlueprint,
							services.Merge)
						gomega.Expect(err).To(gomega.HaveOccurred())
					})
				})

				ginkgo.When("when replace data of given cluster", func() {
					ginkgo.BeforeEach(func() {
						recName = clusterUpdateEmptyReplace
					})

					ginkgo.It("should return an error", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						retCluster, err = client.ClusterService.Update(context.TODO(), *cluster, clusterBlueprint,
							services.Replace)
						gomega.Expect(err).To(gomega.HaveOccurred())
					})
				})
			})
		})

		ginkgo.Context("when cluster doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterUpdateUnknown

				cluster = resources.CreateClusterWithID(nonExistingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}

				clusterBlueprint = blueprint.CreateUpdateClusterBlueprint()
				if clusterBlueprint == nil {
					err = errors.ErrNoClusterBlueprint
					return
				}
				clusterBlueprint.SetReservedCPU(321)
			})

			ginkgo.It("should return that cluster with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				retCluster, err = client.ClusterService.Update(context.TODO(), *cluster, clusterBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when cluster is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterUpdateNoUser

				clusterBlueprint = blueprint.CreateUpdateClusterBlueprint()
				if clusterBlueprint == nil {
					err = errors.ErrNoClusterBlueprint
					return
				}
				clusterBlueprint.SetReservedCPU(321)
			})

			ginkgo.It("should return that cluster has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				retCluster, err = client.ClusterService.Update(context.TODO(), resources.Cluster{},
					clusterBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("cluster add host", func() {
		var (
			cluster           *resources.Cluster
			oneCluster        *resources.Cluster
			clusterID         int
			host              *resources.Host
			existingHostID    = 4
			nonExistingHostID = 110
		)

		ginkgo.Context("when cluster exists", func() {
			ginkgo.BeforeEach(func() {
				cluster = resources.CreateClusterWithID(existingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.When("when new host exists", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterAddHost

					host = resources.CreateHostWithID(existingHostID)
				})

				ginkgo.It("should add host to the given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.AddHost(context.TODO(), *cluster, *host)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether host was really added in OpenNebula
					clusterID, err = cluster.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneCluster, err = client.ClusterService.RetrieveInfo(context.TODO(), clusterID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneCluster).ShouldNot(gomega.BeNil())

					var hosts []int
					hosts, err = oneCluster.Hosts()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(hosts).Should(gomega.ContainElement(existingHostID))
				})
			})

			ginkgo.When("when host does not exists", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterAddNoneHost

					host = resources.CreateHostWithID(nonExistingHostID)
				})

				ginkgo.It("should not add host to given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.AddHost(context.TODO(), *cluster, *host)
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when cluster doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterAddHostUnknown

				cluster = resources.CreateClusterWithID(nonExistingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.It("should return that cluster with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.AddHost(context.TODO(), *cluster, *host)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when cluster is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterAddHostNoCluster

				cluster = &resources.Cluster{}
			})

			ginkgo.It("should return that cluster has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.AddHost(context.TODO(), *cluster, *host)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("cluster delete host", func() {
		var (
			cluster           *resources.Cluster
			oneCluster        *resources.Cluster
			clusterID         int
			host              *resources.Host
			existingHostID    = 4
			nonExistingHostID = 110
		)

		ginkgo.Context("when cluster exists", func() {
			ginkgo.BeforeEach(func() {
				cluster = resources.CreateClusterWithID(existingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.When("when new host exists", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterDelHost

					host = resources.CreateHostWithID(existingHostID)
				})

				ginkgo.It("should delete host to the given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.DeleteHost(context.TODO(), *cluster, *host)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether host was really added in OpenNebula
					clusterID, err = cluster.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneCluster, err = client.ClusterService.RetrieveInfo(context.TODO(), clusterID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneCluster).ShouldNot(gomega.BeNil())

					var hosts []int
					hosts, err = oneCluster.Hosts()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(hosts).ShouldNot(gomega.ContainElement(existingHostID))
				})
			})

			ginkgo.When("when host does not exists", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterDelNoneHost

					host = resources.CreateHostWithID(nonExistingHostID)
				})

				ginkgo.It("should not delete host to given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.DeleteHost(context.TODO(), *cluster, *host)
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when cluster doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterDelHostUnknown

				cluster = resources.CreateClusterWithID(nonExistingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.It("should return that cluster with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.DeleteHost(context.TODO(), *cluster, *host)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when cluster is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterDelHostNoCluster

				cluster = &resources.Cluster{}
			})

			ginkgo.It("should return that cluster has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.DeleteHost(context.TODO(), *cluster, *host)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("cluster add datastore", func() {
		var (
			cluster                *resources.Cluster
			oneCluster             *resources.Cluster
			clusterID              int
			datastore              *resources.Datastore
			existingDatastoreID    = 110
			nonExistingDatastoreID = 120
		)

		ginkgo.Context("when cluster exists", func() {
			ginkgo.BeforeEach(func() {
				cluster = resources.CreateClusterWithID(existingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.When("when new datastore exists", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterAddDatastore

					datastore = resources.CreateDatastoreWithID(existingDatastoreID)
				})

				ginkgo.It("should add datastore to the given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.AddDatastore(context.TODO(), *cluster, *datastore)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether datastore was really added in OpenNebula
					clusterID, err = cluster.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneCluster, err = client.ClusterService.RetrieveInfo(context.TODO(), clusterID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneCluster).ShouldNot(gomega.BeNil())

					var datastores []int
					datastores, err = oneCluster.Datastores()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					founded := false
					for i := range datastores {
						if datastores[i] == existingDatastoreID {
							founded = true
						}
					}
					gomega.Expect(founded).To(gomega.Equal(true))
				})
			})

			ginkgo.When("when datastore does not exists", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterAddNoneDatastore

					datastore = resources.CreateDatastoreWithID(nonExistingDatastoreID)
				})

				ginkgo.It("should not add datastore to given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.AddDatastore(context.TODO(), *cluster, *datastore)
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when cluster doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterAddDatastoreUnknown

				cluster = resources.CreateClusterWithID(nonExistingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.It("should return that cluster with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.AddDatastore(context.TODO(), *cluster, *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when cluster is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterAddDatastoreNoCluster

				cluster = &resources.Cluster{}
			})

			ginkgo.It("should return that cluster has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.AddDatastore(context.TODO(), *cluster, *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("cluster delete datastore", func() {
		var (
			cluster                *resources.Cluster
			oneCluster             *resources.Cluster
			clusterID              int
			datastore              *resources.Datastore
			existingDatastoreID    = 110
			nonExistingDatastoreID = 120
		)

		ginkgo.Context("when cluster exists", func() {
			ginkgo.BeforeEach(func() {
				cluster = resources.CreateClusterWithID(existingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.When("when new datastore exists", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterDelDatastore

					datastore = resources.CreateDatastoreWithID(existingDatastoreID)
				})

				ginkgo.It("should delete datastore to the given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.DeleteDatastore(context.TODO(), *cluster, *datastore)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether datastore was really added in OpenNebula
					clusterID, err = cluster.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneCluster, err = client.ClusterService.RetrieveInfo(context.TODO(), clusterID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneCluster).ShouldNot(gomega.BeNil())

					var datastores []int
					datastores, err = oneCluster.Datastores()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					founded := false
					for i := range datastores {
						if datastores[i] == existingDatastoreID {
							founded = true
						}
					}
					gomega.Expect(founded).To(gomega.Equal(false))
				})
			})

			ginkgo.When("when datastore does not exists", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterDelNoneDatastore

					datastore = resources.CreateDatastoreWithID(nonExistingDatastoreID)
				})

				ginkgo.It("should not delete datastore to given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.DeleteDatastore(context.TODO(), *cluster, *datastore)
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when cluster doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterDelDatastoreUnknown

				cluster = resources.CreateClusterWithID(nonExistingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.It("should return that cluster with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.DeleteDatastore(context.TODO(), *cluster, *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when cluster is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterDelDatastoreNoCluster

				cluster = &resources.Cluster{}
			})

			ginkgo.It("should return that cluster has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.DeleteDatastore(context.TODO(), *cluster, *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("cluster add virtualNetwork", func() {
		var (
			cluster                     *resources.Cluster
			oneCluster                  *resources.Cluster
			clusterID                   int
			virtualNetwork              *resources.VirtualNetwork
			existingVirtualNetworkID    = 0
			nonExistingVirtualNetworkID = 110
		)

		ginkgo.Context("when cluster exists", func() {
			ginkgo.BeforeEach(func() {
				cluster = resources.CreateClusterWithID(existingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.When("when new virtualNetwork exists", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterAddVirtualNetwork

					virtualNetwork = resources.CreateVirtualNetworkWithID(existingVirtualNetworkID)
				})

				ginkgo.It("should add virtualNetwork to the given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.AddVirtualNetwork(context.TODO(), *cluster, *virtualNetwork)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether virtualNetwork was really added in OpenNebula
					clusterID, err = cluster.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneCluster, err = client.ClusterService.RetrieveInfo(context.TODO(), clusterID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneCluster).ShouldNot(gomega.BeNil())

					var virtualNetworks []int
					virtualNetworks, err = oneCluster.VirtualNetworks()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					founded := false
					for i := range virtualNetworks {
						if virtualNetworks[i] == existingVirtualNetworkID {
							founded = true
						}
					}
					gomega.Expect(founded).To(gomega.Equal(true))
				})
			})

			ginkgo.When("when virtualNetwork does not exists", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterAddNoneVirtualNetwork

					virtualNetwork = resources.CreateVirtualNetworkWithID(nonExistingVirtualNetworkID)
				})

				ginkgo.It("should not add virtualNetwork to given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.AddVirtualNetwork(context.TODO(), *cluster, *virtualNetwork)
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when cluster doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterAddVirtualNetworkUnknown

				cluster = resources.CreateClusterWithID(nonExistingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.It("should return that cluster with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.AddVirtualNetwork(context.TODO(), *cluster, *virtualNetwork)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when cluster is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterAddVirtualNetworkNoCluster

				cluster = &resources.Cluster{}
			})

			ginkgo.It("should return that cluster has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.AddVirtualNetwork(context.TODO(), *cluster, *virtualNetwork)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("cluster delete virtualNetwork", func() {
		var (
			cluster                     *resources.Cluster
			oneCluster                  *resources.Cluster
			clusterID                   int
			virtualNetwork              *resources.VirtualNetwork
			existingVirtualNetworkID    = 0
			nonExistingVirtualNetworkID = 110
		)

		ginkgo.Context("when cluster exists", func() {
			ginkgo.BeforeEach(func() {
				cluster = resources.CreateClusterWithID(existingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.When("when new virtualNetwork exists", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterDelVirtualNetwork

					virtualNetwork = resources.CreateVirtualNetworkWithID(existingVirtualNetworkID)
				})

				ginkgo.It("should delete virtualNetwork to the given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.DeleteVirtualNetwork(context.TODO(), *cluster, *virtualNetwork)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether virtualNetwork was really added in OpenNebula
					clusterID, err = cluster.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneCluster, err = client.ClusterService.RetrieveInfo(context.TODO(), clusterID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneCluster).ShouldNot(gomega.BeNil())

					var virtualNetworks []int
					virtualNetworks, err = oneCluster.VirtualNetworks()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					founded := false
					for i := range virtualNetworks {
						if virtualNetworks[i] == existingVirtualNetworkID {
							founded = true
						}
					}
					gomega.Expect(founded).To(gomega.Equal(false))
				})
			})

			ginkgo.When("when virtualNetwork does not exists", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterDelNoneVirtualNetwork

					virtualNetwork = resources.CreateVirtualNetworkWithID(nonExistingVirtualNetworkID)
				})

				ginkgo.It("should not delete virtualNetwork to given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.DeleteVirtualNetwork(context.TODO(), *cluster, *virtualNetwork)
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when cluster doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterDelVirtualNetworkUnknown

				cluster = resources.CreateClusterWithID(nonExistingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.It("should return that cluster with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.DeleteVirtualNetwork(context.TODO(), *cluster, *virtualNetwork)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when cluster is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterDelVirtualNetworkNoCluster

				cluster = &resources.Cluster{}
			})

			ginkgo.It("should return that cluster has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.DeleteVirtualNetwork(context.TODO(), *cluster, *virtualNetwork)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("cluster rename", func() {
		var (
			cluster    *resources.Cluster
			oneCluster *resources.Cluster
			clusterID  int
		)

		ginkgo.Context("when cluster exists", func() {
			ginkgo.BeforeEach(func() {
				cluster = resources.CreateClusterWithID(existingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.When("when new name is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterRename
				})

				ginkgo.It("should change name of given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					// get cluster name
					oneCluster, err = client.ClusterService.RetrieveInfo(context.TODO(), existingClusterID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneCluster).ShouldNot(gomega.BeNil())

					newName := "the_best_cluster_ever"
					gomega.Expect(oneCluster.Name()).NotTo(gomega.Equal(newName))

					// change name
					err = client.ClusterService.Rename(context.TODO(), *cluster, newName)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether name was really changed in OpenNebula
					clusterID, err = cluster.ID()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					oneCluster, err = client.ClusterService.RetrieveInfo(context.TODO(), clusterID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneCluster).ShouldNot(gomega.BeNil())

					gomega.Expect(oneCluster.Name()).To(gomega.Equal(newName))
				})
			})

			ginkgo.When("when new name is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = clusterRenameEmpty
				})

				ginkgo.It("should not change name of given cluster", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.ClusterService.Rename(context.TODO(), *cluster, "")
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when cluster doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterRenameUnknown

				cluster = resources.CreateClusterWithID(nonExistingClusterID)
				if cluster == nil {
					err = errors.ErrNoCluster
				}
			})

			ginkgo.It("should return that cluster with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.Rename(context.TODO(), *cluster, "mastodont")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when cluster is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterRenameNoCluster

				cluster = &resources.Cluster{}
			})

			ginkgo.It("should return that cluster has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.ClusterService.Rename(context.TODO(), *cluster, "rex")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("cluster retrieve info", func() {
		var cluster *resources.Cluster

		ginkgo.Context("when cluster exists", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterRetrieveInfo
			})

			ginkgo.It("should return cluster with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				cluster, err = client.ClusterService.RetrieveInfo(context.TODO(), existingClusterID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(cluster).ShouldNot(gomega.BeNil())
				gomega.Expect(cluster.ID()).To(gomega.Equal(existingClusterID))
				gomega.Expect(cluster.Name()).To(gomega.Equal("the_best_cluster_ever"))
			})
		})

		ginkgo.Context("when cluster doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = clusterRetrieveInfoUnknown
			})

			ginkgo.It("should return that given cluster doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				cluster, err = client.ClusterService.RetrieveInfo(context.TODO(), nonExistingClusterID)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(cluster).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("cluster list", func() {
		ginkgo.BeforeEach(func() {
			recName = clusterList
		})

		ginkgo.It("should return an array of clusters with full info", func() {
			gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

			var clusters []*resources.Cluster

			clusters, err = client.ClusterService.List(context.TODO())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(clusters).ShouldNot(gomega.BeNil())
		})
	})
})
