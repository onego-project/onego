package resources

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/onego-project/xmlrpc"
	"strconv"
	"github.com/onego-project/onego/requests"
)

// VirtualMachine struct
type VirtualMachine struct {
	XMLData *etree.Element
}

// RPC struct
type RPC struct {
	Client *xmlrpc.Client
	Key    string
}

// Snapshot struct
type Snapshot struct {
	//name       string
	SnapshotID int
}

// Monitoring struct
type Monitoring struct {
	VMid              int
	MonitoringRecords []MonitoringData
}

// MonitoringData struct
type MonitoringData struct {
	XMLData *etree.Element
}

// History struct
type History struct {
	XMLData *etree.Element
}

// CreateVM constructs VirtualMachine
func CreateVM(id int) *VirtualMachine {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0"`)

	el := doc.CreateElement("VM").CreateElement("ID")
	el.SetText(fmt.Sprintf("%d", id))

	return &VirtualMachine{doc.Root()}
}

// GetAttribute method
func (vm VirtualMachine) GetAttribute(path string) string {
	elements := vm.XMLData.FindElements(path)
	if elements == nil {
		return ""
	}
	return elements[0].Text()
}

// GetID method
func (vm VirtualMachine) GetID() int {
	i, err := strconv.Atoi(vm.GetAttribute("ID"))
	if err != nil {
		return -1
	}
	return i
}

// GetUID method
func (vm VirtualMachine) GetUID() int {
	i, err := strconv.Atoi(vm.GetAttribute("UID"))
	if err != nil {
		return -1
	}
	return i
}

// GetGID method
func (vm VirtualMachine) GetGID() int {
	i, err := strconv.Atoi(vm.GetAttribute("GID"))
	if err != nil {
		return -1
	}
	return i
}

// GetUName method
func (vm VirtualMachine) GetUName() string {
	return vm.GetAttribute("UNAME")
}

// GetGName method
func (vm VirtualMachine) GetGName() string {
	return vm.GetAttribute("GNAME")
}

// GetName method
func (vm VirtualMachine) GetName() string {
	return vm.GetAttribute("NAME")
}

// GetPermissionUserUse method
func (vm VirtualMachine) GetPermissionUserUse() int {
	i, err := strconv.Atoi(vm.GetAttribute("PERMISSIONS/OWNER_U"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionUserManage method
func (vm VirtualMachine) GetPermissionUserManage() int {
	i, err := strconv.Atoi(vm.GetAttribute("PERMISSIONS/OWNER_M"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionUserAdmin method
func (vm VirtualMachine) GetPermissionUserAdmin() int {
	i, err := strconv.Atoi(vm.GetAttribute("PERMISSIONS/OWNER_A"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionGroupUse method
func (vm VirtualMachine) GetPermissionGroupUse() int {
	i, err := strconv.Atoi(vm.GetAttribute("PERMISSIONS/GROUP_U"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionGroupManage method
func (vm VirtualMachine) GetPermissionGroupManage() int {
	i, err := strconv.Atoi(vm.GetAttribute("PERMISSIONS/GROUP_M"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionGroupAdmin method
func (vm VirtualMachine) GetPermissionGroupAdmin() int {
	i, err := strconv.Atoi(vm.GetAttribute("PERMISSIONS/GROUP_A"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionOtherUse method
func (vm VirtualMachine) GetPermissionOtherUse() int {
	i, err := strconv.Atoi(vm.GetAttribute("PERMISSIONS/OTHER_U"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionOtherManage method
func (vm VirtualMachine) GetPermissionOtherManage() int {
	i, err := strconv.Atoi(vm.GetAttribute("PERMISSIONS/OTHER_M"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionOtherAdmin method
func (vm VirtualMachine) GetPermissionOtherAdmin() int {
	i, err := strconv.Atoi(vm.GetAttribute("PERMISSIONS/OTHER_A"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionRequest method
func (vm VirtualMachine) GetPermissionRequest() *requests.PermissionRequest {
	return &requests.PermissionRequest{
		Permissions: [][]int{{vm.GetPermissionUserUse(), vm.GetPermissionUserManage(), vm.GetPermissionUserAdmin()},
			{vm.GetPermissionGroupUse(), vm.GetPermissionGroupManage(), vm.GetPermissionGroupAdmin()},
			{vm.GetPermissionOtherUse(), vm.GetPermissionOtherManage(), vm.GetPermissionOtherAdmin()}}}
}

// GetLastPoll method
func (vm VirtualMachine) GetLastPoll() int {
	i, err := strconv.Atoi(vm.GetAttribute("LAST_POLL"))
	if err != nil {
		return -1
	}
	return i
}

// GetState method
func (vm VirtualMachine) GetState() int {
	i, err := strconv.Atoi(vm.GetAttribute("STATE"))
	if err != nil {
		return -1
	}
	return i
}

// GetLCMState method
func (vm VirtualMachine) GetLCMState() int {
	i, err := strconv.Atoi(vm.GetAttribute("LCM_STATE"))
	if err != nil {
		return -1
	}
	return i
}

// GetPrevState method
func (vm VirtualMachine) GetPrevState() int {
	i, err := strconv.Atoi(vm.GetAttribute("PREV_STATE"))
	if err != nil {
		return -1
	}
	return i
}

// GetPrevLCMState method
func (vm VirtualMachine) GetPrevLCMState() int {
	i, err := strconv.Atoi(vm.GetAttribute("PREV_LCM_STATE"))
	if err != nil {
		return -1
	}
	return i
}

// GetResched method
func (vm VirtualMachine) GetResched() int {
	i, err := strconv.Atoi(vm.GetAttribute("RESCHED"))
	if err != nil {
		return -1
	}
	return i
}

// GetSTime method
func (vm VirtualMachine) GetSTime() int64 {
	i, err := strconv.ParseInt(vm.GetAttribute("STIME"), 10, 64)
	if err != nil {
		return -1
	}
	return i
}

// GetETime method
func (vm VirtualMachine) GetETime() int64 {
	i, err := strconv.ParseInt(vm.GetAttribute("ETIME"), 10, 64)
	if err != nil {
		return -1
	}
	return i
}

// GetDeployID method
func (vm VirtualMachine) GetDeployID() string {
	return vm.GetAttribute("DEPLOY_ID")
}

// GetMonitoring method
func (vm VirtualMachine) GetMonitoring() *Monitoring {
	return &Monitoring{
		VMid:              vm.GetID(),
		MonitoringRecords: []MonitoringData{{vm.XMLData.FindElements("MONITORING")[0]}}}
}

// GetTemplate method
func (vm VirtualMachine) GetTemplate() *etree.Element {
	return vm.XMLData.FindElements("TEMPLATE")[0]
}

// GetUserTemplate method
func (vm VirtualMachine) GetUserTemplate() *etree.Element {
	return vm.XMLData.FindElements("USERTEMPLATE")[0]
}

// GetHistoryRecords method
func (vm VirtualMachine) GetHistoryRecords() []History {
	elements := vm.XMLData.FindElements("HISTORY_RECORDS/HISTORY")
	history := make([]History, len(elements))
	for i, e := range elements {
		history[i] = History{e}
	}
	return history
}
