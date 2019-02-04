package services

import (
	"context"

	"github.com/onego-project/onego/blueprint"

	"github.com/beevik/etree"
	"github.com/onego-project/onego/requests"
	"github.com/onego-project/onego/resources"
)

// VirtualMachineService structure to manage OpenNebula virtual machine.
type VirtualMachineService struct {
	Service
}

// VM constants to perform an action on a virtual machine.
const (
	vmTerminate     = "terminate"
	vmTerminateHard = "terminate-hard"
	vmUndeploy      = "undeploy"
	vmUndeployHard  = "undeploy-hard"
	vmPoweroff      = "poweroff"
	vmPoweroffHard  = "poweroff-hard"
	vmReboot        = "reboot"
	vmRebootHard    = "reboot-hard"
	vmHold          = "hold"
	vmRelease       = "release"
	vmStop          = "stop"
	vmSuspend       = "suspend"
	vmResume        = "resume"
	vmReschedule    = "resched"
	vmUnreschedule  = "unresched"
)

// RecoverOperation to recover a stuck VM that is waiting for a driver operation.
type RecoverOperation int

const (
	// Failure to failure.
	Failure RecoverOperation = iota
	// Success to success.
	Success
	// Retry to retry.
	Retry
	// Delete to delete.
	Delete
	// DeleteRecreate to delete-recreate.
	DeleteRecreate
)

// StateFilter - VM state to filter by.
type StateFilter int

const (
	// AnyStateIncludingDone const
	AnyStateIncludingDone StateFilter = iota - 2
	// AnyStateExceptDone const
	AnyStateExceptDone
	// Init const
	Init
	// Pending const
	Pending
	// Hold const
	Hold
	// Active const
	Active
	// Stopped const
	Stopped
	// Suspended const
	Suspended
	// Done const
	Done
	// Failed const
	Failed
	// PowerOff const
	PowerOff
	// Undeployed const
	Undeployed
	// Cloning const
	Cloning
	// CloningFailure const
	CloningFailure
)

//// Snapshot structure represents VM Snapshot.
//type Snapshot struct {
//	Name string
//	ID   int
//}

// Allocate allocates a new virtual machine in OpenNebula.
// Blueprint should contain at least VM name, memory and cpu.
func (vms *VirtualMachineService) Allocate(ctx context.Context, blueprintInterface blueprint.Interface,
	onHold bool) (*resources.VirtualMachine, error) {
	blueprintText, err := blueprintInterface.Render()
	if err != nil {
		return nil, err
	}

	resArr, err := vms.call(ctx, "one.vm.allocate", blueprintText, onHold)
	if err != nil {
		return nil, err
	}

	return vms.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Deploy initiates the instance of the given virtual machine on the target host.
func (vms *VirtualMachineService) Deploy(ctx context.Context, vm resources.VirtualMachine,
	host resources.Host, overCommit bool, dataStore resources.Datastore) error {
	vmID, err := vm.ID()
	if err != nil {
		return err
	}

	hostID, err := host.ID()
	if err != nil {
		return err
	}

	dataStoreID, err := dataStore.ID()
	if err != nil {
		return err
	}

	_, err = vms.call(ctx, "one.vm.deploy", vmID, hostID, overCommit, dataStoreID)

	return err
}

// action submits an action to be performed on a virtual machine.
func (vms *VirtualMachineService) action(ctx context.Context, vm resources.VirtualMachine, action string) error {
	vmID, err := vm.ID()
	if err != nil {
		return err
	}

	_, err = vms.call(ctx, "one.vm.action", action, vmID)

	return err
}

// Terminate performs a terminate on a virtual machine.
func (vms *VirtualMachineService) Terminate(ctx context.Context, vm resources.VirtualMachine, hard bool) error {
	if hard {
		return vms.action(ctx, vm, vmTerminateHard)
	}
	return vms.action(ctx, vm, vmTerminate)
}

// Undeploy performs an undeploy on a virtual machine.
func (vms *VirtualMachineService) Undeploy(ctx context.Context, vm resources.VirtualMachine, hard bool) error {
	if hard {
		return vms.action(ctx, vm, vmUndeployHard)
	}
	return vms.action(ctx, vm, vmUndeploy)
}

// PowerOff performs a power off on a virtual machine.
func (vms *VirtualMachineService) PowerOff(ctx context.Context, vm resources.VirtualMachine, hard bool) error {
	if hard {
		return vms.action(ctx, vm, vmPoweroffHard)
	}
	return vms.action(ctx, vm, vmPoweroff)
}

// Reboot performs a reboot on a virtual machine.
func (vms *VirtualMachineService) Reboot(ctx context.Context, vm resources.VirtualMachine, hard bool) error {
	if hard {
		return vms.action(ctx, vm, vmRebootHard)
	}
	return vms.action(ctx, vm, vmReboot)
}

// Hold performs a hold on a virtual machine.
func (vms *VirtualMachineService) Hold(ctx context.Context, vm resources.VirtualMachine) error {
	return vms.action(ctx, vm, vmHold)
}

// Release performs a release on a virtual machine.
func (vms *VirtualMachineService) Release(ctx context.Context, vm resources.VirtualMachine) error {
	return vms.action(ctx, vm, vmRelease)
}

// Stop performs a stop on a virtual machine.
func (vms *VirtualMachineService) Stop(ctx context.Context, vm resources.VirtualMachine) error {
	return vms.action(ctx, vm, vmStop)
}

// Suspend performs a suspend on a virtual machine.
func (vms *VirtualMachineService) Suspend(ctx context.Context, vm resources.VirtualMachine) error {
	return vms.action(ctx, vm, vmSuspend)
}

// Resume performs a resume on a virtual machine.
func (vms *VirtualMachineService) Resume(ctx context.Context, vm resources.VirtualMachine) error {
	return vms.action(ctx, vm, vmResume)
}

// Reschedule performs a reschedule on a virtual machine.
func (vms *VirtualMachineService) Reschedule(ctx context.Context, vm resources.VirtualMachine) error {
	return vms.action(ctx, vm, vmReschedule)
}

// Unreschedule performs an unreschedule on a virtual machine.
func (vms *VirtualMachineService) Unreschedule(ctx context.Context, vm resources.VirtualMachine) error {
	return vms.action(ctx, vm, vmUnreschedule)
}

// Migrate migrates one virtual machine to the target host.
func (vms *VirtualMachineService) Migrate(ctx context.Context, vm resources.VirtualMachine,
	host resources.Host, dataStore resources.Datastore, liveMigration bool, overCommit bool) error {
	vmID, err := vm.ID()
	if err != nil {
		return err
	}

	hostID, err := host.ID()
	if err != nil {
		return err
	}

	dataStoreID, err := dataStore.ID()
	if err != nil {
		return err
	}

	_, err = vms.call(ctx, "one.vm.migrate", vmID, hostID, liveMigration,
		overCommit, dataStoreID)
	return err
}

// Chmod changes the permission bits of a virtual machine.
func (vms *VirtualMachineService) Chmod(ctx context.Context, vm resources.VirtualMachine,
	request requests.PermissionRequest) error {
	vmID, err := vm.ID()
	if err != nil {
		return err
	}

	return vms.chmod(ctx, "one.vm.chmod", vmID, request)
}

// Chown changes the ownership of a virtual machine.
func (vms *VirtualMachineService) Chown(ctx context.Context, vm resources.VirtualMachine,
	request requests.OwnershipRequest) error {
	vmID, err := vm.ID()
	if err != nil {
		return err
	}

	return vms.chown(ctx, "one.vm.chown", vmID, request)
}

// Rename renames a virtual machine.
func (vms *VirtualMachineService) Rename(ctx context.Context, vm resources.VirtualMachine, name string) error {
	vmID, err := vm.ID()
	if err != nil {
		return err
	}

	_, err = vms.call(ctx, "one.vm.rename", vmID, name)

	return err
}

//// CreateSnapshot creates a new virtual machine snapshot.
//func (vms *VirtualMachineService) CreateSnapshot(ctx context.Context, vm resources.VirtualMachine,
//	name string) (*Snapshot, error) {
//	vmID, err := vm.ID()
//	if err != nil {
//		return nil, err
//	}
//
//	resArr, err := vms.call(ctx, "one.vm.snapshotcreate", vmID, name)
//	if err != nil {
//		return nil, err
//	}
//
//	snapshot := Snapshot{ID: int(resArr[resultIndex].ResultInt()), Name: name}
//
//	return &snapshot, nil
//}
//
//// RevertSnapshot reverts a virtual machine to a snapshot.
//func (vms *VirtualMachineService) RevertSnapshot(ctx context.Context, vm resources.VirtualMachine,
//	snapshot Snapshot) error {
//	vmID, err := vm.ID()
//	if err != nil {
//		return err
//	}
//
//	_, err = vms.call(ctx, "one.vm.snapshotrevert", vmID, snapshot.ID)
//
//	return err
//}
//
//// DeleteSnapshot deletes a virtual machine snapshot.
//func (vms *VirtualMachineService) DeleteSnapshot(ctx context.Context, vm resources.VirtualMachine,
//	snapshot Snapshot) error {
//	vmID, err := vm.ID()
//	if err != nil {
//		return err
//	}
//
//	_, err = vms.call(ctx, "one.vm.snapshotdelete", vmID, snapshot.ID)
//
//	return err
//}

// Resize changes the capacity of the virtual machine.
func (vms *VirtualMachineService) Resize(ctx context.Context, vm resources.VirtualMachine,
	request requests.ResizeRequest, overCommit bool) error {
	vmID, err := vm.ID()
	if err != nil {
		return err
	}

	text, err := resources.RenderInterfaceToXMLString(request)
	if err != nil {
		// error should never occur
		return err
	}

	_, err = vms.call(ctx, "one.vm.resize", vmID, text, overCommit)
	return err
}

// UpdateUserTemplate merges/replaces the user template contents.
func (vms *VirtualMachineService) UpdateUserTemplate(ctx context.Context, vm resources.VirtualMachine,
	blueprint blueprint.Interface, updateType UpdateType) error {
	vmID, err := vm.ID()
	if err != nil {
		return err
	}

	blueprintText, err := blueprint.Render()
	if err != nil {
		return err
	}

	_, err = vms.call(ctx, "one.vm.update", vmID, blueprintText, updateType)
	return err
}

// UpdateTemplate updates (appends) a set of supported configuration attributes in the VM template.
func (vms *VirtualMachineService) UpdateTemplate(ctx context.Context, vm resources.VirtualMachine,
	blueprint blueprint.Interface) error {
	vmID, err := vm.ID()
	if err != nil {
		return err
	}

	blueprintText, err := blueprint.Render()
	if err != nil {
		return err
	}

	_, err = vms.call(ctx, "one.vm.updateconf", vmID, blueprintText)
	return err
}

// Recover recovers a stuck VM that is waiting for a driver operation.
// The recovery may be done by failing or succeeding the pending operation.
// You need to manually check the vm status on the host, to decide if the operation was successful or not.
func (vms *VirtualMachineService) Recover(ctx context.Context, vm resources.VirtualMachine,
	operation RecoverOperation) error {
	vmID, err := vm.ID()
	if err != nil {
		return err
	}

	_, err = vms.call(ctx, "one.vm.recover", vmID, operation)

	return err
}

// RetrieveInfo retrieves information for the virtual machine.
func (vms *VirtualMachineService) RetrieveInfo(ctx context.Context,
	vmID int) (*resources.VirtualMachine, error) {
	doc, err := vms.retrieveInfo(ctx, "one.vm.info", vmID)
	if err != nil {
		return nil, err
	}

	return resources.CreateVirtualMachineFromXML(doc.Root()), nil
}

func (vms *VirtualMachineService) list(ctx context.Context, filterFlag, pageOffset,
	pageSize int, filter StateFilter) ([]*resources.VirtualMachine, error) {
	resArr, err := vms.call(ctx, "one.vmpool.info", filterFlag, pageOffset, pageSize, filter)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromString(resArr[resultIndex].ResultString()); err != nil {
		return nil, err
	}

	elements := doc.FindElements("VM_POOL/VM")

	virtualMachines := make([]*resources.VirtualMachine, len(elements))
	for i, e := range elements {
		virtualMachines[i] = resources.CreateVirtualMachineFromXML(e)
	}

	return virtualMachines, nil
}

// ListAll retrieves information for part of the vms in the pool which belong to given owner(s) in ownership filter.
func (vms *VirtualMachineService) ListAll(ctx context.Context, ownershipFilter OwnershipFilter,
	stateFilter StateFilter) ([]*resources.VirtualMachine, error) {
	return vms.list(ctx, int(ownershipFilter), pageOffsetDefault, pageSizeDefault, stateFilter)
}

// ListAllForUser retrieves information for part of the vms in the pool.
func (vms *VirtualMachineService) ListAllForUser(ctx context.Context, user resources.User,
	stateFilter StateFilter) ([]*resources.VirtualMachine, error) {
	userID, err := user.ID()
	if err != nil {
		return nil, err
	}

	return vms.list(ctx, userID, pageOffsetDefault, pageSizeDefault, stateFilter)
}

// List retrieves information for all the vms in the pool.
func (vms *VirtualMachineService) List(ctx context.Context, pageOffset int,
	pageSize int, ownershipFilter OwnershipFilter, stateFilter StateFilter) ([]*resources.VirtualMachine, error) {
	return vms.list(ctx, int(ownershipFilter), (pageOffset-1)*pageSize, -pageSize, stateFilter)
}

// ListForUser retrieves information for part of the vms in the pool.
func (vms *VirtualMachineService) ListForUser(ctx context.Context, user resources.User, pageOffset int,
	pageSize int, stateFilter StateFilter) ([]*resources.VirtualMachine, error) {
	userID, err := user.ID()
	if err != nil {
		return nil, err
	}

	return vms.list(ctx, userID, (pageOffset-1)*pageSize, -pageSize, stateFilter)
}
