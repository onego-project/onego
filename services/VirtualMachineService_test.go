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
	virtualMachineAllocate               = "records/virtualMachine/allocate"
	virtualMachineAllocateWrongBlueprint = "records/virtualMachine/allocateWrongBlueprint"

	virtualMachineDeploy                 = "records/virtualMachine/deploy"
	virtualMachineDeployWrongID          = "records/virtualMachine/deployWrongID"
	virtualMachineDeployWrongHostID      = "records/virtualMachine/deployWrongHostID"
	virtualMachineDeployWrongDatastoreID = "records/virtualMachine/deployWrongDatastoreID"

	virtualMachineActionWrongID = "records/virtualMachine/actionWrongID"
	virtualMachineTerminate     = "records/virtualMachine/terminate"
	virtualMachineTerminateHard = "records/virtualMachine/terminateHard"
	virtualMachineUndeploy      = "records/virtualMachine/undeploy"
	virtualMachineUndeployHard  = "records/virtualMachine/undeployHard"
	virtualMachinePowerOff      = "records/virtualMachine/powerOff"
	virtualMachinePowerOffHard  = "records/virtualMachine/powerOffHard"
	virtualMachineReboot        = "records/virtualMachine/reboot"
	virtualMachineRebootHard    = "records/virtualMachine/rebootHard"
	virtualMachineHold          = "records/virtualMachine/hold"
	virtualMachineRelease       = "records/virtualMachine/release"
	virtualMachineStop          = "records/virtualMachine/stop"
	virtualMachineSuspend       = "records/virtualMachine/suspend"
	virtualMachineResume        = "records/virtualMachine/resume"
	virtualMachineReschedule    = "records/virtualMachine/reschedule"
	virtualMachineUnreschedule  = "records/virtualMachine/unreschedule"

	virtualMachineMigrate                 = "records/virtualMachine/migrate"
	virtualMachineMigrateWrongID          = "records/virtualMachine/migrateWrongID"
	virtualMachineMigrateWrongHostID      = "records/virtualMachine/migrateWrongHostID"
	virtualMachineMigrateWrongDatastoreID = "records/virtualMachine/migrateWrongDatastoreID"

	virtualMachineChmod              = "records/virtualMachine/chmod"
	virtualMachinePermRequestDefault = "records/virtualMachine/chmodPermReqDefault"
	virtualMachineChmodWrongID       = "records/virtualMachine/chmodUnknown"

	virtualMachineChown               = "records/virtualMachine/chown"
	virtualMachineOwnershipReqDefault = "records/virtualMachine/chownDefault"
	virtualMachineChownWrongID        = "records/virtualMachine/chownUnknown"

	virtualMachineRename        = "records/virtualMachine/rename"
	virtualMachineRenameEmpty   = "records/virtualMachine/renameEmpty"
	virtualMachineRenameUnknown = "records/virtualMachine/renameUnknown"

	//virtualMachineCreateSnapshot        = "records/virtualMachine/createSnapshot"
	//virtualMachineCreateSnapshotWrongID = "records/virtualMachine/createSnapshotWrongID"
	//
	//virtualMachineRevertSnapshot        = "records/virtualMachine/revertSnapshot"
	//virtualMachineRevertSnapshotWrongID = "records/virtualMachine/revertSnapshotWrongID"
	//
	//virtualMachineDeleteSnapshot        = "records/virtualMachine/deleteSnapshot"
	//virtualMachineDeleteSnapshotWrongID = "records/virtualMachine/deleteSnapshotWrongID"

	virtualMachineResize        = "records/virtualMachine/resize"
	virtualMachineResizeWrongID = "records/virtualMachine/resizeWrongID"

	virtualMachineUpdateUserTemplateMerge   = "records/virtualMachine/updateUserTemplateMerge"
	virtualMachineUpdateUserTemplateReplace = "records/virtualMachine/updateUserTemplateReplace"
	virtualMachineUpdateUserTemplateWrongID = "records/virtualMachine/updateUserTemplateWrongID"

	virtualMachineUpdateTemplate        = "records/virtualMachine/updateTemplate"
	virtualMachineUpdateTemplateWrongID = "records/virtualMachine/updateTemplateWrongID"

	virtualMachineRecover        = "records/virtualMachine/recover"
	virtualMachineRecoverWrongID = "records/virtualMachine/recoverWrongID"

	virtualMachineRetrieveInfo        = "records/virtualMachine/retrieveInfo"
	virtualMachineRetrieveInfoWrongID = "records/virtualMachine/retrieveInfoWrongID"

	virtualMachineListAllPrimaryGroup = "records/virtualMachine/listAllPrimaryGroup"
	virtualMachineListAllUser         = "records/virtualMachine/listAllUser"
	virtualMachineListAllAll          = "records/virtualMachine/listAllAll"
	virtualMachineListAllUserGroup    = "records/virtualMachine/listAllUserGroup"

	virtualMachineListAllForUser        = "records/virtualMachine/listAllForUser"
	virtualMachineListAllForUserWrongID = "records/virtualMachine/listAllForUserWrongID"

	virtualMachineListPagination = "records/virtualMachine/listPagination"

	virtualMachineListForUser        = "records/virtualMachine/listForUser"
	virtualMachineListForUserWrongID = "records/virtualMachine/listForUserWrongID"
)

var _ = ginkgo.Describe("Virtual Machine Service", func() {
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

	ginkgo.Describe("allocate Virtual Machine", func() {
		var (
			virtualMachine          *resources.VirtualMachine
			oneVirtualMachine       *resources.VirtualMachine
			virtualMachineBlueprint *blueprint.VirtualMachineBlueprint
			virtualMachineID        int
			virtualMachineName      string
		)

		ginkgo.Context("when virtual machine doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineAllocate

				virtualMachineName = "test-allocate-vm"
				virtualMachineBlueprint = blueprint.CreateAllocateVirtualMachineBlueprint()
				if virtualMachineBlueprint == nil {
					err = errors.ErrNoVirtualMachineBlueprint
				}

				virtualMachineBlueprint.SetName(virtualMachineName)
				virtualMachineBlueprint.SetMemory(2048)
				virtualMachineBlueprint.SetCPU(4)
			})

			ginkgo.It("should create new Virtual Machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachine, err = client.VirtualMachineService.Allocate(context.TODO(), virtualMachineBlueprint, true)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualMachine).ShouldNot(gomega.BeNil())

				// check whether VirtualMachine really exists in OpenNebula
				virtualMachineID, err = virtualMachine.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.Name()).To(gomega.Equal(virtualMachineName))
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateHold))
			})
		})

		ginkgo.Context("when virtual machine blueprint has no element", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineAllocateWrongBlueprint

				virtualMachineBlueprint = blueprint.CreateAllocateVirtualMachineBlueprint()
				if virtualMachineBlueprint == nil {
					err = errors.ErrNoVirtualMachineBlueprint
				}
			})

			ginkgo.It("shouldn't create new Virtual Machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachine, err = client.VirtualMachineService.Allocate(context.TODO(), virtualMachineBlueprint, true)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(virtualMachine).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("deploy Virtual Machine", func() {
		var (
			virtualMachine    *resources.VirtualMachine
			oneVirtualMachine *resources.VirtualMachine
			virtualMachineID  int

			hostID int
			host   *resources.Host

			datastoreID int
			datastore   *resources.Datastore
		)

		ginkgo.Context("when virtual machine, host and datastore are set correctly", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineDeploy

				virtualMachineID = 89
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				hostID = 4
				host = resources.CreateHostWithID(hostID)
				if host == nil {
					err = errors.ErrNoHost
				}

				datastoreID = 104
				datastore = resources.CreateDatastoreWithID(datastoreID)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should deploy a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Deploy(context.TODO(), *virtualMachine, *host, true, *datastore)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really deployed in OpenNebula
				virtualMachineID, err = virtualMachine.ID()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).NotTo(gomega.Equal(resources.VirtualMachineStateHold))
			})
		})

		ginkgo.Context("when virtual machine is not set correctly", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineDeployWrongID

				hostID = 4
				host = resources.CreateHostWithID(hostID)
				if host == nil {
					err = errors.ErrNoHost
				}

				datastoreID = 104
				datastore = resources.CreateDatastoreWithID(datastoreID)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should not deploy a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Deploy(context.TODO(), resources.VirtualMachine{}, *host, true, *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when host is not set correctly", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineDeployWrongHostID

				virtualMachineID = 89
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				datastoreID = 104
				datastore = resources.CreateDatastoreWithID(datastoreID)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should not deploy a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Deploy(context.TODO(), *virtualMachine, resources.Host{}, true, *datastore)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when datastore is not set correctly", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineDeployWrongDatastoreID

				virtualMachineID = 89
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				hostID = 4
				host = resources.CreateHostWithID(hostID)
				if host == nil {
					err = errors.ErrNoHost
				}
			})

			ginkgo.It("should not deploy a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Deploy(context.TODO(), *virtualMachine, *host, true, resources.Datastore{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("action Virtual Machine", func() {
		var (
			virtualMachine    *resources.VirtualMachine
			oneVirtualMachine *resources.VirtualMachine
			virtualMachineID  int
		)

		ginkgo.Context("when virtual machine is empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineActionWrongID
			})

			ginkgo.It("should not change a state of a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Terminate(context.TODO(), resources.VirtualMachine{}, false)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when terminate", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineTerminate

				virtualMachineID = 111
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should terminate a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Terminate(context.TODO(), *virtualMachine, false)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really terminated in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateDone))
			})
		})

		ginkgo.Context("when terminate hard", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineTerminateHard

				virtualMachineID = 112
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should terminate hard a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Terminate(context.TODO(), *virtualMachine, true)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really terminated hard in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateDone))
			})
		})

		ginkgo.Context("when undeploy", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineUndeploy

				virtualMachineID = 137
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should undeploy a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Undeploy(context.TODO(), *virtualMachine, false)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really undeployed in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateUndeployed))
			})
		})

		ginkgo.Context("when undeploy hard", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineUndeployHard

				virtualMachineID = 153
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should undeploy hard a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Undeploy(context.TODO(), *virtualMachine, true)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really undeployed hard in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateUndeployed))
			})
		})

		ginkgo.Context("when power off", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachinePowerOff

				virtualMachineID = 116
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should power off a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.PowerOff(context.TODO(), *virtualMachine, false)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really powered off in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStatePowerOff))
			})
		})

		ginkgo.Context("when power off hard", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachinePowerOffHard

				virtualMachineID = 117
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should power off hard a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.PowerOff(context.TODO(), *virtualMachine, true)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really powered off hard in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStatePowerOff))
			})
		})

		ginkgo.Context("when reboot", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineReboot

				virtualMachineID = 118
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should reboot a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Reboot(context.TODO(), *virtualMachine, false)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really reboot in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateActive))
			})
		})

		ginkgo.Context("when reboot hard", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineRebootHard

				virtualMachineID = 122
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should reboot hard a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Reboot(context.TODO(), *virtualMachine, true)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really rebooted hard in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateActive))
			})
		})

		ginkgo.Context("when hold", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineHold

				virtualMachineID = 128
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should hold a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Hold(context.TODO(), *virtualMachine)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really hold in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateHold))
			})
		})

		ginkgo.Context("when release", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineRelease

				virtualMachineID = 128
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should release a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Release(context.TODO(), *virtualMachine)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really released in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateActive))
			})
		})

		ginkgo.Context("when stop", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineStop

				virtualMachineID = 122
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should stop a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Stop(context.TODO(), *virtualMachine)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really stopped in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateStopped))
			})
		})

		ginkgo.Context("when suspend", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineSuspend

				virtualMachineID = 123
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should suspend a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Suspend(context.TODO(), *virtualMachine)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really suspended in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateSuspended))
			})
		})

		ginkgo.Context("when resume", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineResume

				virtualMachineID = 123
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should resume a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Resume(context.TODO(), *virtualMachine)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really resumed in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateActive))
			})
		})

		ginkgo.Context("when reschedule", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineReschedule

				virtualMachineID = 123
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should reschedule a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Reschedule(context.TODO(), *virtualMachine)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really rescheduled in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.Reschedule()).To(gomega.BeTrue())
			})
		})

		ginkgo.Context("when unreschedule", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineUnreschedule

				virtualMachineID = 123
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should unreschedule a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Unreschedule(context.TODO(), *virtualMachine)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really unrescheduled in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine.Reschedule()).To(gomega.BeFalse())
			})
		})
	})

	ginkgo.Describe("Migrate Virtual Machine", func() {
		var (
			virtualMachine    *resources.VirtualMachine
			oneVirtualMachine *resources.VirtualMachine
			virtualMachineID  int

			host   *resources.Host
			hostID int

			datastore   *resources.Datastore
			datastoreID int

			historyRecords []*resources.History
		)

		ginkgo.Context("when VM, host and datastore are set correctly", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineMigrate

				virtualMachineID = 123
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				hostID = 10
				host = resources.CreateHostWithID(hostID)
				if host == nil {
					err = errors.ErrNoHost
				}

				datastoreID = 0
				datastore = resources.CreateDatastoreWithID(datastoreID)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should migrate a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Migrate(context.TODO(), *virtualMachine, *host, *datastore, false, false)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really migrated in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				historyRecords, err = oneVirtualMachine.HistoryRecords()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(historyRecords).NotTo(gomega.HaveLen(0))

				gomega.Expect(historyRecords[len(historyRecords)-1].HID).To(gomega.Equal(&hostID))
			})
		})

		ginkgo.Context("when VM is not set correctly", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineMigrateWrongID

				virtualMachineID = 1000
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				hostID = 10
				host = resources.CreateHostWithID(hostID)
				if host == nil {
					err = errors.ErrNoHost
				}

				datastoreID = 0
				datastore = resources.CreateDatastoreWithID(datastoreID)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should not migrate a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Migrate(context.TODO(), *virtualMachine, *host, *datastore, false, false)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when host is not set correctly", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineMigrateWrongHostID

				virtualMachineID = 123
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				hostID = 100
				host = resources.CreateHostWithID(hostID)
				if host == nil {
					err = errors.ErrNoHost
				}

				datastoreID = 0
				datastore = resources.CreateDatastoreWithID(datastoreID)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should not migrate a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Migrate(context.TODO(), *virtualMachine, *host, *datastore, false, false)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when datastore is not set correctly", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineMigrateWrongDatastoreID

				virtualMachineID = 123
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				hostID = 10
				host = resources.CreateHostWithID(hostID)
				if host == nil {
					err = errors.ErrNoHost
				}

				datastoreID = 500
				datastore = resources.CreateDatastoreWithID(datastoreID)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should not migrate a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Migrate(context.TODO(), *virtualMachine, *host, *datastore, false, false)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when VM is empty", func() {
			ginkgo.BeforeEach(func() {
				hostID = 10
				host = resources.CreateHostWithID(hostID)
				if host == nil {
					err = errors.ErrNoHost
				}

				datastoreID = 0
				datastore = resources.CreateDatastoreWithID(datastoreID)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should not migrate a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Migrate(context.TODO(), resources.VirtualMachine{}, *host, *datastore,
					false, false)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when host is empty", func() {
			ginkgo.BeforeEach(func() {
				virtualMachineID = 123
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				datastoreID = 0
				datastore = resources.CreateDatastoreWithID(datastoreID)
				if datastore == nil {
					err = errors.ErrNoDatastore
				}
			})

			ginkgo.It("should not migrate a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Migrate(context.TODO(), *virtualMachine, resources.Host{}, *datastore,
					false, false)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when datastore is empty", func() {
			ginkgo.BeforeEach(func() {
				virtualMachineID = 123
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				hostID = 10
				host = resources.CreateHostWithID(hostID)
				if host == nil {
					err = errors.ErrNoHost
				}
			})

			ginkgo.It("should not migrate a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Migrate(context.TODO(), *virtualMachine, *host, resources.Datastore{},
					false, false)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("virtual machine chmod", func() {
		var (
			virtualMachine    *resources.VirtualMachine
			oneVirtualMachine *resources.VirtualMachine
			virtualMachineID  int

			permRequest requests.PermissionRequest
		)

		ginkgo.Context("when virtual machine exists", func() {
			ginkgo.BeforeEach(func() {
				virtualMachineID = 130
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.When("when permission request is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualMachineChmod
				})

				ginkgo.It("should change permission of given virtual machine", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					permRequest = requests.CreatePermissionRequestBuilder().Deny(requests.User,
						requests.Manage).Allow(requests.Other, requests.Admin).Build()

					err = client.VirtualMachineService.Chmod(context.TODO(), *virtualMachine, permRequest)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chmod was really changed in OpenNebula
					oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVirtualMachine).ShouldNot(gomega.BeNil())

					var perm *resources.Permissions
					perm, err = oneVirtualMachine.Permissions()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					gomega.Expect(perm.User.Manage).To(gomega.Equal(false))
					gomega.Expect(perm.Other.Admin).To(gomega.Equal(true))
				})
			})

			ginkgo.When("when permission request is default", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualMachinePermRequestDefault
				})

				ginkgo.It("should not change permissions of given virtual machine", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.VirtualMachineService.Chmod(context.TODO(), *virtualMachine,
						requests.CreatePermissionRequestBuilder().Build())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when virtual machine doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineChmodWrongID

				virtualMachine = resources.CreateVirtualMachineWithID(110)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should return that virtual machine with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				permRequest = requests.CreatePermissionRequestBuilder().Allow(requests.User,
					requests.Manage).Deny(requests.Other, requests.Admin).Build()

				err = client.VirtualMachineService.Chmod(context.TODO(), *virtualMachine, permRequest)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtual machine is empty", func() {
			ginkgo.It("should return that virtual machine has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Chmod(context.TODO(), resources.VirtualMachine{}, requests.PermissionRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("virtual machine chown", func() {
		var (
			virtualMachine    *resources.VirtualMachine
			oneVirtualMachine *resources.VirtualMachine
			virtualMachineID  int

			userID int
			user   *resources.User

			groupID int
			group   *resources.Group

			ownershipReq requests.OwnershipRequest
		)

		ginkgo.Context("when virtual machine exists", func() {
			ginkgo.BeforeEach(func() {
				virtualMachineID = 128
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.When("when ownership request is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualMachineChown

					userID = 33
					user = resources.CreateUserWithID(userID)
					if user == nil {
						err = errors.ErrNoUser
					}

					groupID = 120
					group = resources.CreateGroupWithID(groupID)
					if group == nil {
						err = errors.ErrNoGroup
					}

					ownershipReq = requests.CreateOwnershipRequestBuilder().User(*user).Group(*group).Build()
				})

				ginkgo.It("should change owner of given virtual machine", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.VirtualMachineService.Chown(context.TODO(), *virtualMachine, ownershipReq)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether chown was really changed in OpenNebula
					oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVirtualMachine).ShouldNot(gomega.BeNil())

					gomega.Expect(oneVirtualMachine.User()).To(gomega.Equal(userID))
					gomega.Expect(oneVirtualMachine.Group()).To(gomega.Equal(groupID))
				})
			})

			ginkgo.When("when ownership request is default", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualMachineOwnershipReqDefault
				})

				ginkgo.It("should not change permissions of given virtual machine", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.VirtualMachineService.Chown(context.TODO(), *virtualMachine,
						requests.CreateOwnershipRequestBuilder().Build())
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when virtual machine doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineChownWrongID

				virtualMachine = resources.CreateVirtualMachineWithID(110)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should return that virtual machine with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Chown(context.TODO(), *virtualMachine, requests.OwnershipRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtual machine is empty", func() {
			ginkgo.It("should return that virtual machine has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Chown(context.TODO(), resources.VirtualMachine{}, requests.OwnershipRequest{})
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("virtual machine rename", func() {
		var (
			virtualMachine    *resources.VirtualMachine
			oneVirtualMachine *resources.VirtualMachine
			virtualMachineID  int
		)

		ginkgo.Context("when virtual machine exists", func() {
			ginkgo.BeforeEach(func() {
				virtualMachineID = 130
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.When("when new name is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualMachineRename
				})

				ginkgo.It("should change name of given virtual machine", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					// get virtualMachine name
					oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVirtualMachine).ShouldNot(gomega.BeNil())

					newName := "vm_XY"
					gomega.Expect(oneVirtualMachine.Name()).NotTo(gomega.Equal(newName))

					// change name
					err = client.VirtualMachineService.Rename(context.TODO(), *virtualMachine, newName)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether name was really changed in OpenNebula
					oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					gomega.Expect(oneVirtualMachine).ShouldNot(gomega.BeNil())

					gomega.Expect(oneVirtualMachine.Name()).To(gomega.Equal(newName))
				})
			})

			ginkgo.When("when new name is empty", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualMachineRenameEmpty
				})

				ginkgo.It("should not change name of given virtual machine", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.VirtualMachineService.Rename(context.TODO(), *virtualMachine, "")
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when virtual machine doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineRenameUnknown

				virtualMachineID = 1000
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should return that virtual machine with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Rename(context.TODO(), *virtualMachine, "mastodont")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtual machine is empty", func() {
			ginkgo.It("should return that virtual machine has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Rename(context.TODO(), resources.VirtualMachine{}, "rex")
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("Resize Virtual Machine", func() {
		var (
			virtualMachine    *resources.VirtualMachine
			oneVirtualMachine *resources.VirtualMachine
			virtualMachineID  int

			resizeRequest *requests.ResizeRequest
			cpu           float64
			memory        int
		)

		ginkgo.Context("when VM and resize request are set correctly", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineResize

				virtualMachineID = 130
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				cpu = 0.5
				memory = 1
				resizeRequest = &requests.ResizeRequest{
					CPU:    cpu,
					Memory: memory,
				}
			})

			ginkgo.It("should resize a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Resize(context.TODO(), *virtualMachine, *resizeRequest, false)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VirtualMachine was really migrated in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				gomega.Expect(oneVirtualMachine.CPU()).To(gomega.Equal(cpu))
				gomega.Expect(oneVirtualMachine.Memory()).To(gomega.Equal(memory))
			})
		})

		ginkgo.Context("when VM is not set correctly", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineResizeWrongID

				virtualMachineID = 1000
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				cpu = 0.5
				memory = 1
				resizeRequest = &requests.ResizeRequest{
					CPU:    cpu,
					Memory: memory,
				}
			})

			ginkgo.It("should not resize a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Resize(context.TODO(), *virtualMachine, *resizeRequest, false)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when VM is empty", func() {
			ginkgo.BeforeEach(func() {
				cpu = 0.5
				memory = 1
				resizeRequest = &requests.ResizeRequest{
					CPU:    cpu,
					Memory: memory,
				}
			})

			ginkgo.It("should not resize a given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Resize(context.TODO(), resources.VirtualMachine{}, *resizeRequest, false)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("update user template in virtual machine", func() {
		var (
			virtualMachine    *resources.VirtualMachine
			oneVirtualMachine *resources.VirtualMachine
			virtualMachineID  int

			userTemplateBlueprint *blueprint.UserTemplateBlueprint
		)

		ginkgo.Context("when virtual machine exists", func() {
			ginkgo.Context("when update data is not empty", func() {
				ginkgo.BeforeEach(func() {
					virtualMachineID = 128
					virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
					if virtualMachine == nil {
						err = errors.ErrNoVirtualMachine
						return
					}

					userTemplateBlueprint = blueprint.CreateUpdateUserTemplateBlueprint()
					if userTemplateBlueprint == nil {
						err = errors.ErrNoUserTemplateBlueprint
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
					}
					userTemplateBlueprint.SetDescription("dummy")
				})

				ginkgo.When("when merge data of given user template", func() {
					ginkgo.BeforeEach(func() {
						recName = virtualMachineUpdateUserTemplateMerge
					})

					ginkgo.It("should merge data of given user template in virtual machine", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						err = client.VirtualMachineService.UpdateUserTemplate(context.TODO(), *virtualMachine,
							userTemplateBlueprint, services.Merge)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						// check whether VirtualMachine was really migrated in OpenNebula
						oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						gomega.Expect(oneVirtualMachine.Attribute("USER_TEMPLATE/DESCRIPTION")).To(gomega.Equal("dummy"))
					})
				})

				ginkgo.When("when replace data of given user template", func() {
					ginkgo.BeforeEach(func() {
						recName = virtualMachineUpdateUserTemplateReplace
					})

					ginkgo.It("should replace data of given user template in virtual machine", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						err = client.VirtualMachineService.UpdateUserTemplate(context.TODO(), *virtualMachine,
							userTemplateBlueprint, services.Replace)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						// check whether VirtualMachine was really migrated in OpenNebula
						oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
						gomega.Expect(err).NotTo(gomega.HaveOccurred())

						gomega.Expect(oneVirtualMachine.Attribute("USER_TEMPLATE/DESCRIPTION")).To(gomega.Equal("dummy"))
					})
				})
			})

			ginkgo.Context("when update data is empty", func() {
				ginkgo.BeforeEach(func() {
					virtualMachineID = 128
					virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
					if virtualMachine == nil {
						err = errors.ErrNoVirtualMachine
						return
					}

					userTemplateBlueprint = &blueprint.UserTemplateBlueprint{}
				})

				ginkgo.When("when merge data of given virtual machine", func() {
					ginkgo.It("should return an error", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						err = client.VirtualMachineService.UpdateUserTemplate(context.TODO(), *virtualMachine,
							userTemplateBlueprint, services.Merge)
						gomega.Expect(err).To(gomega.HaveOccurred())
					})
				})

				ginkgo.When("when replace data of given virtualnNetwork", func() {
					ginkgo.It("should return an error", func() {
						gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

						err = client.VirtualMachineService.UpdateUserTemplate(context.TODO(), *virtualMachine,
							userTemplateBlueprint, services.Replace)
						gomega.Expect(err).To(gomega.HaveOccurred())
					})
				})
			})
		})

		ginkgo.Context("when virtualMachine doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineUpdateUserTemplateWrongID

				virtualMachineID = 1000
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				userTemplateBlueprint = blueprint.CreateUpdateUserTemplateBlueprint()
				if userTemplateBlueprint == nil {
					err = errors.ErrNoUserTemplateBlueprint
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				}
				userTemplateBlueprint.SetDescription("dummy")
			})

			ginkgo.It("should return that virtualMachine with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.UpdateUserTemplate(context.TODO(), *virtualMachine, userTemplateBlueprint,
					services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtualMachine is empty", func() {
			ginkgo.BeforeEach(func() {
				userTemplateBlueprint = blueprint.CreateUpdateUserTemplateBlueprint()
				if userTemplateBlueprint == nil {
					err = errors.ErrNoUserTemplateBlueprint
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				}
				userTemplateBlueprint.SetDescription("dummy")
			})

			ginkgo.It("should return that virtualMachine has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.UpdateUserTemplate(context.TODO(), resources.VirtualMachine{},
					userTemplateBlueprint, services.Merge)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("update template in virtual machine", func() {
		var (
			virtualMachine    *resources.VirtualMachine
			oneVirtualMachine *resources.VirtualMachine
			virtualMachineID  int

			virtualMachineBlueprint *blueprint.VirtualMachineBlueprint
		)

		ginkgo.Context("when virtual machine exists", func() {
			ginkgo.Context("when update data is not empty", func() {
				ginkgo.BeforeEach(func() {
					recName = virtualMachineUpdateTemplate

					virtualMachineID = 130
					virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
					if virtualMachine == nil {
						err = errors.ErrNoVirtualMachine
						return
					}

					virtualMachineBlueprint = blueprint.CreateUpdateVirtualMachineBlueprint()
					if virtualMachineBlueprint == nil {
						err = errors.ErrNoVirtualMachineBlueprint
						gomega.Expect(err).NotTo(gomega.HaveOccurred())
					}
					virtualMachineBlueprint.SetMemory(1)
				})

				ginkgo.It("should update data of given template in virtual machine", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.VirtualMachineService.UpdateTemplate(context.TODO(), *virtualMachine, virtualMachineBlueprint)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					// check whether VirtualMachine was really migrated in OpenNebula
					oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())

					gomega.Expect(oneVirtualMachine.Attribute("TEMPLATE/MEMORY")).To(gomega.Equal("1"))
				})
			})

			ginkgo.Context("when update data is empty", func() {
				ginkgo.BeforeEach(func() {
					virtualMachineID = 128
					virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
					if virtualMachine == nil {
						err = errors.ErrNoVirtualMachine
						return
					}

					virtualMachineBlueprint = &blueprint.VirtualMachineBlueprint{}
				})

				ginkgo.It("should return an error", func() {
					gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

					err = client.VirtualMachineService.UpdateTemplate(context.TODO(), *virtualMachine, virtualMachineBlueprint)
					gomega.Expect(err).To(gomega.HaveOccurred())
				})
			})
		})

		ginkgo.Context("when virtualMachine doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineUpdateTemplateWrongID

				virtualMachineID = 1000
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}

				virtualMachineBlueprint = blueprint.CreateUpdateVirtualMachineBlueprint()
				if virtualMachineBlueprint == nil {
					err = errors.ErrNoVirtualMachineBlueprint
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				}
				virtualMachineBlueprint.SetName("dummy")
			})

			ginkgo.It("should return that virtualMachine with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.UpdateTemplate(context.TODO(), *virtualMachine, virtualMachineBlueprint)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtualMachine is empty", func() {
			ginkgo.BeforeEach(func() {
				virtualMachineBlueprint = blueprint.CreateUpdateVirtualMachineBlueprint()
				if virtualMachineBlueprint == nil {
					err = errors.ErrNoVirtualMachineBlueprint
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				}
				virtualMachineBlueprint.SetName("dummy")
			})

			ginkgo.It("should return that virtualMachine has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.UpdateTemplate(context.TODO(), resources.VirtualMachine{},
					virtualMachineBlueprint)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("virtual machine recover", func() {
		var (
			virtualMachine    *resources.VirtualMachine
			oneVirtualMachine *resources.VirtualMachine
			virtualMachineID  int
		)

		ginkgo.Context("when virtual machine exists", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineRecover

				virtualMachineID = 135
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should recover of given virtual machine", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				// get virtualMachine state
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine).ShouldNot(gomega.BeNil())

				gomega.Expect(oneVirtualMachine.State()).NotTo(gomega.Equal(resources.VirtualMachineStateActive))

				// recover - delete
				err = client.VirtualMachineService.Recover(context.TODO(), *virtualMachine, services.Delete)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// check whether VM was really recovered in OpenNebula
				oneVirtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(oneVirtualMachine).ShouldNot(gomega.BeNil())

				gomega.Expect(oneVirtualMachine.State()).To(gomega.Equal(resources.VirtualMachineStateDone))
			})
		})

		ginkgo.Context("when virtual machine doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineRecoverWrongID

				virtualMachineID = 1000
				virtualMachine = resources.CreateVirtualMachineWithID(virtualMachineID)
				if virtualMachine == nil {
					err = errors.ErrNoVirtualMachine
				}
			})

			ginkgo.It("should return that virtual machine with given ID doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Recover(context.TODO(), *virtualMachine, services.Delete)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("when virtual machine is empty", func() {
			ginkgo.It("should return that virtual machine has no ID", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				err = client.VirtualMachineService.Recover(context.TODO(), resources.VirtualMachine{}, services.Retry)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
	})

	ginkgo.Describe("virtual machine retrieve info", func() {
		var (
			virtualMachine   *resources.VirtualMachine
			virtualMachineID int
		)

		ginkgo.Context("when virtual machine exists", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineRetrieveInfo

				virtualMachineID = 128
			})

			ginkgo.It("should return virtual machine with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), virtualMachineID)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualMachine).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualMachine.ID()).To(gomega.Equal(virtualMachineID))
				gomega.Expect(virtualMachine.Name()).To(gomega.Equal("cirros-128"))
			})
		})

		ginkgo.Context("when virtual machine doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineRetrieveInfoWrongID
			})

			ginkgo.It("should return that given virtual machine doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachine, err = client.VirtualMachineService.RetrieveInfo(context.TODO(), 1000)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(virtualMachine).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("virtual machine list all", func() {
		var virtualMachines []*resources.VirtualMachine

		ginkgo.Context("when ownership filter is not empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineListAllPrimaryGroup
			})

			ginkgo.It("should return list of all virtual machine with full info belongs to primary group", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachines, err = client.VirtualMachineService.ListAll(context.TODO(),
					services.OwnershipFilterPrimaryGroup, services.Active)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualMachines).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualMachines).To(gomega.HaveLen(4))
			})
		})

		ginkgo.Context("when ownership filter is not empty", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineListAllUser
			})

			ginkgo.It("should return list of all virtual machine with full info belongs to the user", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachines, err = client.VirtualMachineService.ListAll(context.TODO(),
					services.OwnershipFilterUser, services.Active)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualMachines).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualMachines).To(gomega.HaveLen(4))
				gomega.Expect(virtualMachines[0].ID()).To(gomega.Equal(131))
				gomega.Expect(virtualMachines[1].ID()).To(gomega.Equal(132))
				gomega.Expect(virtualMachines[2].ID()).To(gomega.Equal(133))
				gomega.Expect(virtualMachines[3].ID()).To(gomega.Equal(134))
			})
		})

		ginkgo.Context("when ownership filter is set to all", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineListAllAll
			})

			ginkgo.It("should return list of all virtual machine with full info belongs to all", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachines, err = client.VirtualMachineService.ListAll(context.TODO(), services.OwnershipFilterAll,
					services.Active)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualMachines).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualMachines).To(gomega.HaveLen(5))
			})
		})

		ginkgo.Context("when ownership filter is set to UserGroup", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineListAllUserGroup
			})

			ginkgo.It("should return list of all virtual machine with full info belongs to the user and any "+
				"of his groups", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachines, err = client.VirtualMachineService.ListAll(context.TODO(),
					services.OwnershipFilterUserGroup, services.Active)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualMachines).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualMachines).To(gomega.HaveLen(5))
			})
		})
	})

	ginkgo.Describe("virtual machine list all for user", func() {
		var virtualMachines []*resources.VirtualMachine

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineListAllForUser
			})

			ginkgo.It("should return virtual machine with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachines, err = client.VirtualMachineService.ListAllForUser(context.TODO(),
					*resources.CreateUserWithID(33), services.Active)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualMachines).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualMachines).To(gomega.HaveLen(1))
				gomega.Expect(virtualMachines[0].ID()).To(gomega.Equal(128))
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineListAllForUserWrongID
			})

			ginkgo.It("should return empty list of virtual machine (length 0)", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachines, err = client.VirtualMachineService.ListAllForUser(context.TODO(),
					*resources.CreateUserWithID(500), services.Active)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualMachines).Should(gomega.Equal(make([]*resources.VirtualMachine, 0)))
				gomega.Expect(virtualMachines).Should(gomega.HaveLen(0))
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.It("should return that user doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachines, err = client.VirtualMachineService.ListAllForUser(context.TODO(), resources.User{},
					services.Active)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(virtualMachines).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("List methods with pagination", func() {
		var (
			virtualMachines []*resources.VirtualMachine
		)

		ginkgo.Context("pagination ok", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineListPagination
			})

			ginkgo.It("should return virtual machine with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachines, err = client.VirtualMachineService.List(context.TODO(), 3,
					2, services.OwnershipFilterUserGroup, services.Active)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualMachines).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualMachines).To(gomega.HaveLen(1))
				gomega.Expect(virtualMachines[0].ID()).To(gomega.Equal(134))
			})
		})

		ginkgo.Context("pagination wrong", func() {
			ginkgo.It("should return that pagination is wrong", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachines, err = client.VirtualMachineService.List(context.TODO(), -2,
					-2, services.OwnershipFilterPrimaryGroup, services.Active)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(virtualMachines).Should(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("virtual machine list for user", func() {
		var virtualMachines []*resources.VirtualMachine

		ginkgo.Context("when user exists", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineListForUser
			})

			ginkgo.It("should return virtual machine with full info", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachines, err = client.VirtualMachineService.ListForUser(context.TODO(),
					*resources.CreateUserWithID(0), 2, 2, services.Active)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualMachines).ShouldNot(gomega.BeNil())
				gomega.Expect(virtualMachines).To(gomega.HaveLen(2))
				gomega.Expect(virtualMachines[0].ID()).To(gomega.Equal(133))
				gomega.Expect(virtualMachines[1].ID()).To(gomega.Equal(134))
			})
		})

		ginkgo.Context("when user doesn't exist", func() {
			ginkgo.BeforeEach(func() {
				recName = virtualMachineListForUserWrongID
			})

			ginkgo.It("should return empty list of virtual machine (length 0)", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachines, err = client.VirtualMachineService.ListForUser(context.TODO(),
					*resources.CreateUserWithID(88), 2, 2, services.Active)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(virtualMachines).Should(gomega.Equal(make([]*resources.VirtualMachine, 0)))
				gomega.Expect(virtualMachines).Should(gomega.HaveLen(0))
			})
		})

		ginkgo.Context("when user is empty", func() {
			ginkgo.It("should return that user doesn't exist", func() {
				gomega.Expect(err).NotTo(gomega.HaveOccurred()) // no error during BeforeEach

				virtualMachines, err = client.VirtualMachineService.ListForUser(context.TODO(),
					resources.User{}, 2, 2, services.Active)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(virtualMachines).Should(gomega.BeNil())
			})
		})
	})
})
