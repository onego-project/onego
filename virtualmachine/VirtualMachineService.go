package virtualmachine

import (
	"context"
	"fmt"
	"github.com/beevik/etree"
	"github.com/onego-project/xmlrpc"
	"github.com/owlet123/onego/blueprint"
	"github.com/owlet123/onego/datastore"
	"github.com/owlet123/onego/host"
	"github.com/owlet123/onego/ownershiprequest"
	"github.com/owlet123/onego/permissionrequest"
)

// Service struct
type Service struct {
	RPC *RPC
}

const (
	vmDelete       = "terminate"
	vmForceDelete  = "terminate-hard"
	vmUndeploy     = "undeploy"
	vmUndeployHard = "undeploy-hard"
	vmPoweroff     = "poweroff"
	vmPoweroffHard = "poweroff-hard"
	vmReboot       = "reboot"
	vmRebootHard   = "reboot-hard"
	vmHold         = "hold"
	vmRelease      = "release"
	vmStop         = "stop"
	vmSuspend      = "suspend"
	vmResume       = "resume"
	vmReschedule   = "resched"
	vmUnreschedule = "unresched"
)

// UpdateType type
type UpdateType int

const (
	// Replace const
	Replace UpdateType = iota
	// Merge const
	Merge
)

// RecoverOperation type
type RecoverOperation int

const (
	// Failure const
	Failure RecoverOperation = iota
	// Success const
	Success
	// Retry const
	Retry
	// Delete const
	Delete
	// DeleteRecreate const
	DeleteRecreate
)

// OwnershipFilter type
type OwnershipFilter int

const (
	// PrimaryGroup const
	PrimaryGroup OwnershipFilter = iota - 4
	// User const
	User
	// All const
	All
	// UserGroups const
	UserGroups
)

// StateFilter type
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

func (s Service) call(methodName string, args ...interface{}) ([]*xmlrpc.Result, error) {
	ctx := context.TODO()

	result, err := s.RPC.Client.Call(ctx, methodName, args...)
	if err != nil {
		return nil, err
	}

	resArr := result.ResultArray()
	if !resArr[0].ResultBoolean() {
		return nil, fmt.Errorf("%s, code: %d", resArr[1].ResultString(), resArr[2].ResultInt())
	}

	return resArr, nil
}

// Deploy method
func (s Service) Deploy(vm VirtualMachine, host host.Host, overCommit bool, datastore datastore.DataStore) error {
	args := []interface{}{s.RPC.Key, vm.GetID(), host.GetID(), overCommit, datastore.GetID()}
	_, err := s.call("one.vm.deploy", args...)
	return err
}

func (s Service) actions(vm VirtualMachine, action string) error {
	args := []interface{}{s.RPC.Key, action, vm.GetID()}
	_, err := s.call("one.vm.action", args...)
	return err
}

// Terminate method
func (s Service) Terminate(vm VirtualMachine, hard bool) error {
	if hard {
		return s.actions(vm, vmForceDelete)
	}
	return s.actions(vm, vmDelete)
}

// Undeploy method
func (s Service) Undeploy(vm VirtualMachine, hard bool) error {
	if hard {
		return s.actions(vm, vmUndeployHard)
	}
	return s.actions(vm, vmUndeploy)
}

// Poweroff method
func (s Service) Poweroff(vm VirtualMachine, hard bool) error {
	if hard {
		return s.actions(vm, vmPoweroffHard)
	}
	return s.actions(vm, vmPoweroff)
}

// Reboot method
func (s Service) Reboot(vm VirtualMachine, hard bool) error {
	if hard {
		return s.actions(vm, vmRebootHard)
	}
	return s.actions(vm, vmReboot)
}

// Hold method
func (s Service) Hold(vm VirtualMachine) error {
	return s.actions(vm, vmHold)
}

// Release method
func (s Service) Release(vm VirtualMachine) error {
	return s.actions(vm, vmRelease)
}

// Stop method
func (s Service) Stop(vm VirtualMachine) error {
	return s.actions(vm, vmStop)
}

// Suspend method
func (s Service) Suspend(vm VirtualMachine) error {
	return s.actions(vm, vmSuspend)
}

// Resume method
func (s Service) Resume(vm VirtualMachine) error {
	return s.actions(vm, vmResume)
}

// Reschedule method
func (s Service) Reschedule(vm VirtualMachine) error {
	return s.actions(vm, vmReschedule)
}

// Unreschedule method
func (s Service) Unreschedule(vm VirtualMachine) error {
	return s.actions(vm, vmUnreschedule)
}

// Migrate method
func (s Service) Migrate(vm VirtualMachine, host host.Host, datastore datastore.DataStore, liveMigration bool, overcommit bool) error {
	args := []interface{}{s.RPC.Key, vm.GetID(), host.GetID(), liveMigration, overcommit, datastore.GetID()}
	_, err := s.call("one.vm.migrate", args...)
	return err
}

// Chmod method
func (s Service) Chmod(vm VirtualMachine, request permissionrequest.PermissionRequest) error {
	args := []interface{}{s.RPC.Key, vm.GetID()}
	for pGroup := 0; pGroup < 3; pGroup++ {
		for pType := 0; pType < 3; pType++ {
			args = append(args, request.Permissions[pGroup][pType])
		}
	}

	_, err := s.call("one.vm.chmod", args...)
	return err
}

// Chown method
func (s Service) Chown(vm VirtualMachine, request ownershiprequest.OwnershipRequest) error {
	args := []interface{}{s.RPC.Key, vm.GetID(), request.User, request.Group}
	_, err := s.call("one.vm.chown", args...)
	return err
}

// Rename method
func (s Service) Rename(vm VirtualMachine, name string) error {
	args := []interface{}{s.RPC.Key, vm.GetID(), name}
	_, err := s.call("one.vm.rename", args...)
	return err
}

// CreateSnapshot method
func (s Service) CreateSnapshot(vm VirtualMachine, name string) (*Snapshot, error) {
	args := []interface{}{s.RPC.Key, vm.GetID(), name}

	resArr, err := s.call("one.vm.snapshotcreate", args...)
	if err != nil {
		return nil, err
	}

	snapshot := Snapshot{SnapshotID: int(resArr[1].ResultInt())}

	return &snapshot, nil
}

// RevertSnapshot method
func (s Service) RevertSnapshot(vm VirtualMachine, snapshot Snapshot) error {
	args := []interface{}{s.RPC.Key, vm.GetID(), snapshot.SnapshotID}
	_, err := s.call("one.vm.snapshotrevert", args...)
	return err
}

// DeleteSnapshot method
func (s Service) DeleteSnapshot(vm VirtualMachine, snapshot Snapshot) error {
	args := []interface{}{s.RPC.Key, vm.GetID(), snapshot.SnapshotID}
	_, err := s.call("one.vm.snapshotdelete", args...)
	return err
}

// Resize method
func (s Service) Resize(vm VirtualMachine, request blueprint.Interface, overCommit bool) error {
	args := []interface{}{s.RPC.Key, vm.GetID(), request.Render(), overCommit}
	_, err := s.call("one.vm.resize", args...)
	return err
}

// UpdateUserTemplate method
func (s Service) UpdateUserTemplate(vm VirtualMachine, blueprint blueprint.Interface, updateType UpdateType) error {
	args := []interface{}{s.RPC.Key, vm.GetID(), blueprint.Render(), updateType}
	_, err := s.call("one.vm.update", args...)
	return err
}

// UpdateTemplate method
func (s Service) UpdateTemplate(vm VirtualMachine, blueprint blueprint.Interface) error {
	args := []interface{}{s.RPC.Key, vm.GetID(), blueprint.Render()}
	_, err := s.call("one.vm.updateconf", args...)
	return err
}

// Recover method
func (s Service) Recover(vm VirtualMachine, operation RecoverOperation) error {
	args := []interface{}{s.RPC.Key, vm.GetID(), operation}
	_, err := s.call("one.vm.recover", args...)
	return err
}

// RetrieveInfo method
func (s Service) RetrieveInfo(vm VirtualMachine) (*VirtualMachine, error) {
	args := []interface{}{s.RPC.Key, vm.GetID()}
	resArr, err := s.call("one.vm.info", args...)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromString(resArr[1].ResultString()); err != nil {
		return nil, err
	}

	vminfo := VirtualMachine{doc.Root()}

	return &vminfo, nil
}

////  method
// func (s Service) RetrieveMonitoring(vm VirtualMachine) (*Monitoring, error) {
//	args := []interface{}{s.RPC.Key, vm.GetID()}
//	resArr, err := s.call("one.vm.monitoring", args...)
//	if err != nil {
//		return nil, err
//	}
//
//	doc := etree.NewDocument()
//	if err = doc.ReadFromString(resArr[1].ResultString()); err != nil {
//		return nil, err
//	}
//
//	monitoringRecords := make([]MonitoringData, len(doc.ChildElements()))
//	for i, e := range doc.ChildElements() {
//		monitoringRecords[i] = MonitoringData{e}
//	}
//
//	vminfo := Monitoring{vm.GetID(), monitoringRecords}
//
//	return &vminfo, nil
//}

////  method
// func (s Service) RetrieveAccounting(ownershipFilter OwnershipFilter, startTime time.Time, endTime time.Time) ([]History, error) {
//	args := []interface{}{s.RPC.Key, ownershipFilter, startTime.Unix(), endTime.Unix()}
//	resArr, err := s.call("one.vmpool.accounting", args...)
//	if err != nil {
//		return nil, err
//	}
//
//	doc := etree.NewDocument()
//	if err = doc.ReadFromString(resArr[1].ResultString()); err != nil {
//		return nil, err
//	}
//
//	historyRecords := make([]History, len(doc.ChildElements()))
//	for i, e := range doc.ChildElements() {
//		historyRecords[i] = History{e}
//	}
//
//	return historyRecords, nil
//}

// ListAll method
func (s Service) ListAll(ownershipFilter OwnershipFilter, stateFilter StateFilter) ([]*VirtualMachine, error) {
	args := []interface{}{s.RPC.Key, ownershipFilter, -1, -1, stateFilter}
	resArr, err := s.call("one.vmpool.info", args...)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromString(resArr[1].ResultString()); err != nil {
		return nil, err
	}

	elements := doc.FindElements("VM_POOL/VM")
	virtualMachines := make([]*VirtualMachine, len(elements))
	for i, e := range elements {
		virtualMachines[i] = &VirtualMachine{e}
	}

	return virtualMachines, nil
}

// ListAllForUser method
func (s Service) ListAllForUser(user int, stateFilter StateFilter) ([]*VirtualMachine, error) {
	args := []interface{}{s.RPC.Key, user, -1, -1, stateFilter}
	resArr, err := s.call("one.vmpool.info", args...)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromString(resArr[1].ResultString()); err != nil {
		return nil, err
	}

	elements := doc.FindElements("VM_POOL/VM")
	virtualMachines := make([]*VirtualMachine, len(elements))
	for i, e := range elements {
		virtualMachines[i] = &VirtualMachine{e}
	}

	return virtualMachines, nil
}

// List method
func (s Service) List(pageOffset int, pageSize int, ownershipFilter OwnershipFilter, stateFilter StateFilter) ([]*VirtualMachine, error) {
	args := []interface{}{s.RPC.Key, ownershipFilter, -pageOffset, -pageSize, stateFilter}
	resArr, err := s.call("one.vmpool.info", args...)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromString(resArr[1].ResultString()); err != nil {
		return nil, err
	}

	elements := doc.FindElements("VM_POOL/VM")
	virtualMachines := make([]*VirtualMachine, len(elements))
	for i, e := range elements {
		virtualMachines[i] = &VirtualMachine{e}
	}

	return virtualMachines, nil
}

// ListForUser method
func (s Service) ListForUser(user int, pageOffset int, pageSize int, stateFilter StateFilter) ([]*VirtualMachine, error) {
	args := []interface{}{s.RPC.Key, user, -pageOffset, -pageSize, stateFilter}
	resArr, err := s.call("one.vmpool.info", args...)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromString(resArr[1].ResultString()); err != nil {
		return nil, err
	}

	elements := doc.FindElements("VM_POOL/VM")
	virtualMachines := make([]*VirtualMachine, len(elements))
	for i, e := range elements {
		virtualMachines[i] = &VirtualMachine{e}
	}

	return virtualMachines, nil
}

// Allocate method
func (s Service) Allocate(blueprintInterface blueprint.Interface, onHold bool) (*VirtualMachine, error) {
	args := []interface{}{s.RPC.Key, blueprintInterface.Render(), onHold}

	resArr, err := s.call("one.vm.allocate", args...)
	if err != nil {
		return nil, err
	}

	vm := CreateVM(int(resArr[1].ResultInt()))

	return vm, nil
}
