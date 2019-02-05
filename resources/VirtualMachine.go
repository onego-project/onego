package resources

import (
	"encoding/xml"
	"net"
	"strconv"
	"time"

	"github.com/onego-project/onego/errors"

	"github.com/beevik/etree"
)

// VirtualMachine structure to manage OpenNebula virtual machines.
type VirtualMachine struct {
	Resource
}

// VirtualMachineStateMap contains string representation of Virtual Machine State.
var VirtualMachineStateMap = map[VirtualMachineState]string{
	VirtualMachineStateInit:           "INIT",
	VirtualMachineStatePending:        "PENDING",
	VirtualMachineStateHold:           "HOLD",
	VirtualMachineStateActive:         "ACTIVE",
	VirtualMachineStateStopped:        "STOPPED",
	VirtualMachineStateSuspended:      "SUSPENDED",
	VirtualMachineStateDone:           "DONE",
	VirtualMachineStatePowerOff:       "POWEROFF",
	VirtualMachineStateUndeployed:     "UNDEPLOYED",
	VirtualMachineStateCloning:        "CLONING",
	VirtualMachineStateCloningFailure: "CLONING_FAILURE",
}

// VirtualMachineState - state of Virtual Machine.
type VirtualMachineState int

const (
	// VirtualMachineStateInit - state Init
	VirtualMachineStateInit VirtualMachineState = iota
	// VirtualMachineStatePending - state Pending
	VirtualMachineStatePending
	// VirtualMachineStateHold - state Hold
	VirtualMachineStateHold
	// VirtualMachineStateActive - state Active
	VirtualMachineStateActive
	// VirtualMachineStateStopped - state Stopped
	VirtualMachineStateStopped
	// VirtualMachineStateSuspended - state Suspended
	VirtualMachineStateSuspended
	// VirtualMachineStateDone - state Done
	VirtualMachineStateDone
	// VirtualMachineStateFailed - state Failed
	VirtualMachineStateFailed
	// VirtualMachineStatePowerOff - state Power off
	VirtualMachineStatePowerOff
	// VirtualMachineStateUndeployed - state Undeployed
	VirtualMachineStateUndeployed
	// VirtualMachineStateCloning - state Cloning
	VirtualMachineStateCloning
	// VirtualMachineStateCloningFailure - state Cloning Failure
	VirtualMachineStateCloningFailure
)

// VirtualMachineLcmState - lcm state of Virtual Machine.
type VirtualMachineLcmState int

const (
	// VirtualMachineLcmInit - state LcmInit
	VirtualMachineLcmInit VirtualMachineLcmState = iota
	// VirtualMachineProlog - state Prolog
	VirtualMachineProlog
	// VirtualMachineBoot - state Boot
	VirtualMachineBoot
	// VirtualMachineRunning - state Running
	VirtualMachineRunning
	// VirtualMachineMigrate - state Migrate
	VirtualMachineMigrate
	// VirtualMachineSaveStop - state SaveStop
	VirtualMachineSaveStop
	// VirtualMachineSaveSuspend - state SaveSuspend
	VirtualMachineSaveSuspend
	// VirtualMachineSaveMigrate - state SaveMigrate
	VirtualMachineSaveMigrate
	// VirtualMachinePrologMigrate - state PrologMigrate
	VirtualMachinePrologMigrate
	// VirtualMachinePrologResume - state PrologResume
	VirtualMachinePrologResume
	// VirtualMachineEpilogStop - state EpilogStop
	VirtualMachineEpilogStop
	// VirtualMachineEpilog - state Epilog
	VirtualMachineEpilog
	// VirtualMachineShutdown - state Shutdown
	VirtualMachineShutdown
	// VirtualMachineCleanupResubmit - state CleanupResubmit
	VirtualMachineCleanupResubmit
	// VirtualMachineUnknown - state Unknown
	VirtualMachineUnknown
	// VirtualMachineHotplug - state Hotplug
	VirtualMachineHotplug
	// VirtualMachineShutdownPoweroff - state ShutdownPoweroff
	VirtualMachineShutdownPoweroff
	// VirtualMachineBootUnknown - state BootUnknown
	VirtualMachineBootUnknown
	// VirtualMachineBootPoweroff - state BootPoweroff
	VirtualMachineBootPoweroff
	// VirtualMachineBootSuspended - state BootSuspended
	VirtualMachineBootSuspended
	// VirtualMachineBootStopped - state BootStopped
	VirtualMachineBootStopped
	// VirtualMachineCleanupDelete - state CleanupDelete
	VirtualMachineCleanupDelete
	// VirtualMachineHotplugSnapshot - state HotplugSnapshot
	VirtualMachineHotplugSnapshot
	// VirtualMachineHotplugNic - state HotplugNic
	VirtualMachineHotplugNic
	// VirtualMachineHotplugSaveas - state HotplugSaveas
	VirtualMachineHotplugSaveas
	// VirtualMachineHotplugSaveasPoweroff - state HotplugSaveasPoweroff
	VirtualMachineHotplugSaveasPoweroff
	// VirtualMachineHotplugSaveasSuspended - state HotplugSaveasSuspended
	VirtualMachineHotplugSaveasSuspended
	// VirtualMachineShutdownUndeploy - state ShutdownUndeploy
	VirtualMachineShutdownUndeploy
	// VirtualMachineEpilogUndeploy - state EpilogUndeploy
	VirtualMachineEpilogUndeploy
	// VirtualMachinePrologUndeploy - state PrologUndeploy
	VirtualMachinePrologUndeploy
	// VirtualMachineBootUndeploy - state BootUndeploy
	VirtualMachineBootUndeploy
	// VirtualMachineHotplugPrologPoweroff - state HotplugPrologPoweroff
	VirtualMachineHotplugPrologPoweroff
	// VirtualMachineHotplugEpilogPoweroff - state HotplugEpilogPoweroff
	VirtualMachineHotplugEpilogPoweroff
	// VirtualMachineBootMigrate - state BootMigrate
	VirtualMachineBootMigrate
	// VirtualMachineBootFailure - state BootFailure
	VirtualMachineBootFailure
	// VirtualMachineBootMigrateFailure - state BootMigrateFailure
	VirtualMachineBootMigrateFailure
	// VirtualMachinePrologMigrateFailure - state PrologMigrateFailure
	VirtualMachinePrologMigrateFailure
	// VirtualMachinePrologFailure - state PrologFailure
	VirtualMachinePrologFailure
	// VirtualMachineEpilogFailure - state EpilogFailure
	VirtualMachineEpilogFailure
	// VirtualMachineEpilogStopFailure - state EpilogStopFailure
	VirtualMachineEpilogStopFailure
	// VirtualMachineEpilogUndeployFailure - state EpilogUndeployFailure
	VirtualMachineEpilogUndeployFailure
	// VirtualMachinePrologMigratePoweroff - state PrologMigratePoweroff
	VirtualMachinePrologMigratePoweroff
	// VirtualMachinePrologMigratePoweroffFailure - state PrologMigratePoweroffFailure
	VirtualMachinePrologMigratePoweroffFailure
	// VirtualMachinePrologMigrateSuspend - state PrologMigrateSuspend
	VirtualMachinePrologMigrateSuspend
	// VirtualMachinePrologMigrateSuspendFailure - state PrologMigrateSuspendFailure
	VirtualMachinePrologMigrateSuspendFailure
	// VirtualMachineBootStoppedFailure - state BootStoppedFailure
	VirtualMachineBootStoppedFailure
	// VirtualMachineBootUndeployFailure - state BootUndeployFailure
	VirtualMachineBootUndeployFailure
	// VirtualMachinePrologResumeFailure - state PrologResumeFailure
	VirtualMachinePrologResumeFailure
	// VirtualMachinePrologUndeployFailure - state PrologUndeployFailure
	VirtualMachinePrologUndeployFailure
	// VirtualMachineDiskSnapshotPoweroff - state DiskSnapshotPoweroff
	VirtualMachineDiskSnapshotPoweroff
	// VirtualMachineDiskSnapshotRevertPoweroff - state DiskSnapshotRevertPoweroff
	VirtualMachineDiskSnapshotRevertPoweroff
	// VirtualMachineDiskSnapshotDeletePoweroff - state DiskSnapshotDeletePoweroff
	VirtualMachineDiskSnapshotDeletePoweroff
	// VirtualMachineDiskSnapshotSuspended - state DiskSnapshotSuspended
	VirtualMachineDiskSnapshotSuspended
	// VirtualMachineDiskSnapshotRevertSuspended - state DiskSnapshotRevertSuspended
	VirtualMachineDiskSnapshotRevertSuspended
	// VirtualMachineDiskSnapshotDeleteSuspended - state DiskSnapshotDeleteSuspended
	VirtualMachineDiskSnapshotDeleteSuspended
	// VirtualMachineDiskSnapshot - state DiskSnapshot
	VirtualMachineDiskSnapshot
	// VirtualMachineDiskSnapshotDelete - state DiskSnapshotDelete
	VirtualMachineDiskSnapshotDelete
	// VirtualMachinePrologMigrateUnknown - state PrologMigrateUnknown
	VirtualMachinePrologMigrateUnknown
	// VirtualMachinePrologMigrateUnknownFailure - state PrologMigrateUnknownFailure
	VirtualMachinePrologMigrateUnknownFailure
	// VirtualMachineDiskResize - state DiskResize
	VirtualMachineDiskResize
	// VirtualMachineDiskResizePoweroff - state DiskResizePoweroff
	VirtualMachineDiskResizePoweroff
	// VirtualMachineDiskResizeUndeployed - state DiskResizeUndeployed
	VirtualMachineDiskResizeUndeployed
)

// VMTemplate represents XML data of VM template.
type VMTemplate *etree.Element

// VMContext represents XML data of VM template context.
type VMContext *etree.Element

// Disk structure represents disk of Virtual Machine.
type Disk struct {
	ClusterID   int
	DatastoreID int
	DevPrefix   string
	DiskID      int
	DiskType    DiskType
	Driver      string
	ImageID     int
	ImageState  ImageState
	ReadOnly    bool
	Size        int
	Target      string
	TmMad       string
	Type        DiskType
}

// GraphicsTypeMap to convert GraphicsType to string.
var GraphicsTypeMap = map[GraphicsType]string{
	GraphicsTypeVNC:   "VNC",
	GraphicsTypeSpice: "SPICE",
	GraphicsTypeSDL:   "SDL",
	GraphicsTypeNone:  "NONE",
}

// GraphicsType to set graphics type.
type GraphicsType int

const (
	// GraphicsTypeVNC graphics type.
	GraphicsTypeVNC GraphicsType = iota
	// GraphicsTypeSpice graphics type.
	GraphicsTypeSpice
	// GraphicsTypeSDL graphics type.
	GraphicsTypeSDL
	// GraphicsTypeNone graphics type.
	GraphicsTypeNone
)

// Graphics represents graphics of Virtual Machine.
type Graphics struct {
	Listen         net.IP
	Password       string
	Port           *int
	RandomPassword bool
	Type           GraphicsType
	KeyMap         string
}

// NIC represents Network Interface of VM.
type NIC struct {
	XMLName        xml.Name `xml:"NIC,omitempty"`
	AddressRangeID int      `xml:"AR_ID,omitempty"`
	Bridge         string   `xml:"BRIDGE,omitempty"`
	ClusterIDs     []int    `xml:"CLUSTER_ID,omitempty"`
	IP             net.IP   `xml:"IP,omitempty"`
	Mac            string   `xml:"MAC,omitempty"`
	MTU            *int     `xml:"MTU,omitempty"`
	Network        string   `xml:"NETWORK,omitempty"`
	NetworkID      int      `xml:"NETWORK_ID,omitempty"`
	NicID          int      `xml:"NIC_ID,omitempty"`
	Target         string   `xml:"TARGET,omitempty"`
	VnMad          string   `xml:"VN_MAD,omitempty"`
}

// OperatingSystem represents OS of VM.
type OperatingSystem struct {
	Architecture ArchitectureType
	Boot         string
	Bootloader   string
	KernelCMD    string
	Machine      string
	Root         string
}

// ArchitectureTypeMap to convert architecture type to string.
var ArchitectureTypeMap = map[ArchitectureType]string{
	ArchitectureTypeI686:   "i686",
	ArchitectureTypeX86_64: "x86_64",
}

// ArchitectureType - type of architecture.
type ArchitectureType int

const (
	// ArchitectureTypeI686 to set architecture type to i686.
	ArchitectureTypeI686 ArchitectureType = iota
	// ArchitectureTypeX86_64 to set architecture type to x86_64.
	ArchitectureTypeX86_64
)

// Raw represents Raw of Virtual Machine.
type Raw struct {
	Data string
	Type string
}

// UserTemplate represents XML data of VM User Template.
type UserTemplate *etree.Element

// History structure represents history record of Virtual Machine.
type History struct {
	OID         int
	Seq         int
	Hostname    string
	HID         *int
	CID         *int
	STime       *time.Time
	ETime       *time.Time
	VMMad       string
	TmMad       string
	DatastoreID *int
	PSTime      *time.Time
	PETime      *time.Time
	RSTime      *time.Time
	REtime      *time.Time
	ESTime      *time.Time
	EETime      *time.Time
	Action      Action
}

// ActionMap to convert Action to string representation.
var ActionMap = map[Action]string{
	ActionMigrate:             "migrate",
	ActionPoweroffMigrate:     "poweroff_migrate",
	ActionPoweroffHardMigrate: "poweroff_hard_migrate",
	ActionLiveMigrate:         "live_migrate",
	ActionTerminate:           "terminate",
	ActionTerminateHard:       "terminate_hard",
	ActionUndeploy:            "undeploy",
	ActionUndeployHard:        "undeploy_hard",
	ActionHold:                "hold",
	ActionRelease:             "release",
	ActionStop:                "stop",
	ActionSuspend:             "suspend",
	ActionResume:              "resume",
	ActionDelete:              "delete",
	ActionDeleteRecreate:      "delete_recreate",
	ActionReboot:              "reboot",
	ActionRebootHard:          "reboot_hard",
	ActionResched:             "resched",
	ActionUnresched:           "unresched",
	ActionPoweroff:            "poweroff",
	ActionPoweroffHard:        "poweroff_hard",
	ActionDiskAttach:          "disk_attach",
	ActionDiskDetach:          "disk_detach",
	ActionNicAttach:           "nic_attach",
	ActionNicDetach:           "nic_detach",
	ActionAliasAttach:         "alias_attach",
	ActionAliasDetach:         "alias_detach",
	ActionDiskSnapshotCreate:  "disk_snapshot_create",
	ActionDiskSnapshotDelete:  "disk_snapshot_delete",
	ActionDiskSnapshotRename:  "disk_snapshot_rename",
	ActionDiskResize:          "disk_resize",
	ActionDeploy:              "deploy",
	ActionChown:               "chown",
	ActionChmod:               "chmod",
	ActionUpdateconf:          "updateconf",
	ActionRename:              "rename",
	ActionResize:              "resize",
	ActionUpdate:              "update",
	ActionSnapshotCreate:      "snapshot_create",
	ActionSnapshotDelete:      "snapshot_delete",
	ActionSnapshotRevert:      "snapshot_revert",
	ActionDiskSaveas:          "disk_saveas",
	ActionDiskSnapshotRevert:  "disk_snapshot_revert",
	ActionRecover:             "recover",
	ActionRetry:               "retry",
	ActionMonitor:             "monitor",
	ActionNone:                "none",
}

// Action type to enumerate VM Action constants.
type Action int

const (
	// ActionMigrate - action Migrate
	ActionMigrate Action = iota
	// ActionPoweroffMigrate - action PoweroffMigrate
	ActionPoweroffMigrate
	// ActionPoweroffHardMigrate - action PoweroffHardMigrate
	ActionPoweroffHardMigrate
	// ActionLiveMigrate - action LiveMigrate
	ActionLiveMigrate
	// ActionTerminate - action Terminate
	ActionTerminate
	// ActionTerminateHard - action TerminateHard
	ActionTerminateHard
	// ActionUndeploy - action Undeploy
	ActionUndeploy
	// ActionUndeployHard - action UndeployHard
	ActionUndeployHard
	// ActionHold - action Hold
	ActionHold
	// ActionRelease - action Release
	ActionRelease
	// ActionStop - action Stop
	ActionStop
	// ActionSuspend - action Suspend
	ActionSuspend
	// ActionResume - action Resume
	ActionResume
	// ActionDelete - action Delete
	ActionDelete
	// ActionDeleteRecreate - action DeleteRecreate
	ActionDeleteRecreate
	// ActionReboot - action Reboot
	ActionReboot
	// ActionRebootHard - action RebootHard
	ActionRebootHard
	// ActionResched - action Resched
	ActionResched
	// ActionUnresched - action Unresched
	ActionUnresched
	// ActionPoweroff - action Poweroff
	ActionPoweroff
	// ActionPoweroffHard - action PoweroffHard
	ActionPoweroffHard
	// ActionDiskAttach - action DiskAttach
	ActionDiskAttach
	// ActionDiskDetach - action DiskDetach
	ActionDiskDetach
	// ActionNicAttach - action NicAttach
	ActionNicAttach
	// ActionNicDetach - action NicDetach
	ActionNicDetach
	// ActionAliasAttach - action AliasAttach
	ActionAliasAttach
	// ActionAliasDetach - action AliasDetach
	ActionAliasDetach
	// ActionDiskSnapshotCreate - action DiskSnapshotCreate
	ActionDiskSnapshotCreate
	// ActionDiskSnapshotDelete - action DiskSnapshotDelete
	ActionDiskSnapshotDelete
	// ActionDiskSnapshotRename - action DiskSnapshotRename
	ActionDiskSnapshotRename
	// ActionDiskResize - action DiskResize
	ActionDiskResize
	// ActionDeploy - action Deploy
	ActionDeploy
	// ActionChown - action Chown
	ActionChown
	// ActionChmod - action Chmod
	ActionChmod
	// ActionUpdateconf - action Updateconf
	ActionUpdateconf
	// ActionRename - action Rename
	ActionRename
	// ActionResize - action Resize
	ActionResize
	// ActionUpdate - action Update
	ActionUpdate
	// ActionSnapshotCreate - action SnapshotCreate
	ActionSnapshotCreate
	// ActionSnapshotDelete - action SnapshotDelete
	ActionSnapshotDelete
	// ActionSnapshotRevert - action SnapshotRevert
	ActionSnapshotRevert
	// ActionDiskSaveas - action DiskSaveas
	ActionDiskSaveas
	// ActionDiskSnapshotRevert - action DiskSnapshotRevert
	ActionDiskSnapshotRevert
	// ActionRecover - action Recover
	ActionRecover
	// ActionRetry - action Retry
	ActionRetry
	// ActionMonitor - action Monitor
	ActionMonitor
	// ActionNone - action None
	ActionNone
)

// CreateVirtualMachineWithID constructs virtual machine with id.
func CreateVirtualMachineWithID(id int) *VirtualMachine {
	return &VirtualMachine{*CreateResource("VM", id)}
}

// CreateVirtualMachineFromXML constructs virtual machine with full xml data.
func CreateVirtualMachineFromXML(XMLdata *etree.Element) *VirtualMachine {
	return &VirtualMachine{Resource: Resource{XMLData: XMLdata}}
}

// User gets User ID of given VM.
func (vm *VirtualMachine) User() (int, error) {
	return vm.intAttribute("UID")
}

// Group gets Group ID of given VM.
func (vm *VirtualMachine) Group() (int, error) {
	return vm.intAttribute("GID")
}

// Permissions gets VM permissions.
func (vm *VirtualMachine) Permissions() (*Permissions, error) {
	return vm.permissions()
}

// LastPoll gets last poll of given VM.
func (vm *VirtualMachine) LastPoll() (*time.Time, error) {
	return vm.parseTime("LAST_POLL")
}

// State gets state of given VM.
func (vm *VirtualMachine) State() (VirtualMachineState, error) {
	s, err := vm.intAttribute("STATE")
	return VirtualMachineState(s), err
}

// LCMState gets LCM state of given VM.
func (vm *VirtualMachine) LCMState() (VirtualMachineLcmState, error) {
	s, err := vm.intAttribute("LCM_STATE")
	return VirtualMachineLcmState(s), err
}

// PrevState gets previous state of given VM.
func (vm *VirtualMachine) PrevState() (VirtualMachineState, error) {
	s, err := vm.intAttribute("PREV_STATE")
	return VirtualMachineState(s), err
}

// PrevLCMState gets previous LCM state of given VM.
func (vm *VirtualMachine) PrevLCMState() (VirtualMachineLcmState, error) {
	s, err := vm.intAttribute("PREV_LCM_STATE")
	return VirtualMachineLcmState(s), err
}

// Reschedule gets true for reschedule false otherwise.
func (vm *VirtualMachine) Reschedule() (bool, error) {
	ret, err := vm.intAttribute("RESCHED")
	if err != nil {
		return false, err
	}

	return intToBool(ret), nil
}

// STime gets start time of given VM.
func (vm *VirtualMachine) STime() (*time.Time, error) {
	return vm.parseTime("STIME")
}

// ETime gets end time of given VM.
func (vm *VirtualMachine) ETime() (*time.Time, error) {
	return vm.parseTime("ETIME")
}

// DeployID gets (string) ID of given VM.
func (vm *VirtualMachine) DeployID() (string, error) {
	return vm.Attribute("DEPLOY_ID")
}

// Template gets template of given VM as an element.
func (vm *VirtualMachine) Template() (VMTemplate, error) {
	template := vm.XMLData.FindElement("TEMPLATE")
	if template == nil {
		return nil, &errors.XMLElementError{Path: "TEMPLATE"}
	}
	return template, nil
}

// Context gets Context of given VM as an element.
func (vm *VirtualMachine) Context() (VMContext, error) {
	context := vm.XMLData.FindElement("TEMPLATE/CONTEXT")
	if context == nil {
		return nil, &errors.XMLElementError{Path: "TEMPLATE/CONTEXT"}
	}
	return context, nil
}

// CPU gets CPU of given VM.
func (vm *VirtualMachine) CPU() (float64, error) {
	f, err := vm.Attribute("TEMPLATE/CPU")
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(f, 64)
}

// CreatedBy gets a ID who create given VM.
func (vm *VirtualMachine) CreatedBy() (int, error) {
	return vm.intAttribute("TEMPLATE/CREATED_BY")
}

// Disks gets an array of Disks of given VM.
func (vm *VirtualMachine) Disks() ([]*Disk, error) {
	elements := vm.XMLData.FindElements("TEMPLATE/DISK")
	if len(elements) == 0 {
		return make([]*Disk, 0), nil
	}

	array := make([]*Disk, len(elements))
	var err error

	for i, e := range elements {
		array[i], err = createDiskFromElement(e)
		if err != nil {
			return nil, err
		}
	}
	return array, nil
}

func createDiskFromElement(element *etree.Element) (*Disk, error) {
	if element == nil {
		return nil, &errors.XMLElementError{Path: "TEMPLATE/DISK"}
	}

	parsedStrings, err := parseStringsFromElement(element, []string{"DEV_PREFIX", "DISK_TYPE", "DRIVER", "READONLY",
		"TARGET", "TM_MAD", "TYPE"})
	if err != nil {
		return nil, err
	}

	parsedInts, err := parseIntsFromElement(element, []string{"CLUSTER_ID", "DATASTORE_ID", "DISK_ID", "IMAGE_ID",
		"IMAGE_STATE", "SIZE"})
	if err != nil {
		return nil, err
	}

	diskType, err := findDiskTypeByValue(parsedStrings[1])
	if err != nil {
		return nil, err
	}

	ttype, err := findDiskTypeByValue(parsedStrings[6])
	if err != nil {
		return nil, err
	}

	return &Disk{
		ClusterID:   parsedInts[0],
		DatastoreID: parsedInts[1],
		DevPrefix:   parsedStrings[0],
		DiskID:      parsedInts[2],
		DiskType:    *diskType,
		Driver:      parsedStrings[2],
		ImageID:     parsedInts[3],
		ImageState:  ImageState(parsedInts[4]),
		ReadOnly:    stringToBool(parsedStrings[3]),
		Size:        parsedInts[5],
		Target:      parsedStrings[4],
		TmMad:       parsedStrings[5],
		Type:        *ttype,
	}, nil
}

// Graphics gets a Graphics of given VM.
func (vm *VirtualMachine) Graphics() (*Graphics, error) {
	element := vm.XMLData.FindElement("TEMPLATE/GRAPHICS")
	if element == nil {
		return &Graphics{Type: GraphicsTypeNone}, nil
	}

	parsedStrings := parseStringsFromElementWithoutError(element, []string{"LISTEN", "PASSWD", "RANDOM_PASSWD", "TYPE",
		"KEYMAP"})

	// no error when LISTEN is empty (return nil)
	var ip net.IP
	if parsedStrings[0] != "" {
		ip = net.ParseIP(parsedStrings[0])
	}

	// no error when RANDOM_PASSWD is empty (return false)
	var randomPasswd bool
	if parsedStrings[2] != "" {
		randomPasswd = stringToBool(parsedStrings[2])
	}

	// no error when PORT is empty (set nil)
	var port *int
	p, err := intAttributeFromElement(element, "PORT")
	if err == nil {
		port = &p
	}

	graphicsType := GraphicsTypeNone
	if parsedStrings[3] != "" {
		for key, val := range GraphicsTypeMap {
			if val == parsedStrings[3] {
				graphicsType = key
				break
			}
		}
	}

	return &Graphics{
		Listen:         ip,
		Password:       parsedStrings[1],
		RandomPassword: randomPasswd,
		Port:           port,
		Type:           graphicsType,
		KeyMap:         parsedStrings[4],
	}, nil
}

// Memory gets size of memory of given VM.
func (vm *VirtualMachine) Memory() (int, error) {
	return vm.intAttribute("TEMPLATE/MEMORY")
}

// NICs gets an array of NICs of given VM.
func (vm *VirtualMachine) NICs() ([]*NIC, error) {
	elements := vm.XMLData.FindElements("TEMPLATE/NIC")
	if len(elements) == 0 {
		return make([]*NIC, 0), nil
	}

	array := make([]*NIC, len(elements))
	var err error

	for i, e := range elements {
		array[i], err = createNICFromElement(e)
		if err != nil {
			return nil, err
		}
	}
	return array, nil
}

func createNICFromElement(element *etree.Element) (*NIC, error) {
	if element == nil {
		return nil, &errors.XMLElementError{Path: "TEMPLATE/NIC"}
	}

	parsedStrings, err := parseStringsFromElement(element, []string{"BRIDGE", "CLUSTER_ID", "MAC", "NETWORK",
		"TARGET", "VN_MAD"})
	if err != nil {
		return nil, err
	}

	parsedInts, err := parseIntsFromElement(element, []string{"AR_ID", "NETWORK_ID", "NIC_ID"})
	if err != nil {
		return nil, err
	}

	// ignore error; create nil attribute
	var ip net.IP
	s, err := attributeFromElement(element, "IP")
	if err == nil {
		ip = net.ParseIP(s)
	}

	// ignore error; create nil attribute
	var mtu *int
	i, err := intAttributeFromElement(element, "MTU")
	if err == nil {
		mtu = &i
	}

	// parse cluster IDs
	clusterIDs, err := parseIntsFromString(parsedStrings[1])
	if err != nil {
		return nil, err
	}

	return &NIC{
		AddressRangeID: parsedInts[0],
		Bridge:         parsedStrings[0],
		ClusterIDs:     clusterIDs,
		IP:             ip,
		Mac:            parsedStrings[2],
		MTU:            mtu,
		Network:        parsedStrings[3],
		NetworkID:      parsedInts[1],
		NicID:          parsedInts[2],
		Target:         parsedStrings[4],
		VnMad:          parsedStrings[5],
	}, nil
}

// OperatingSystem gets OperatingSystem structure of given VM.
func (vm *VirtualMachine) OperatingSystem() (*OperatingSystem, error) {
	element := vm.XMLData.FindElement("TEMPLATE/OS")
	if element == nil {
		return nil, &errors.XMLElementError{Path: "TEMPLATE/OS"}
	}

	stringsWithoutError := parseStringsFromElementWithoutError(element, []string{"BOOT", "BOOTLOADER", "KERNEL_CMD",
		"MACHINE", "ROOT", "ARCH"})

	arch := stringsWithoutError[5]
	var archType *ArchitectureType
	for key, val := range ArchitectureTypeMap {
		if val == arch {
			archType = &key
			break
		}
	}

	return &OperatingSystem{Architecture: *archType,
		Boot:       stringsWithoutError[0],
		Bootloader: stringsWithoutError[1],
		KernelCMD:  stringsWithoutError[2],
		Machine:    stringsWithoutError[3],
		Root:       stringsWithoutError[4],
	}, nil
}

// Raw gets Raw structure of given Virtual Machine.
func (vm *VirtualMachine) Raw() (*Raw, error) {
	element := vm.XMLData.FindElement("TEMPLATE/RAW")
	if element == nil {
		return nil, &errors.XMLElementError{Path: "TEMPLATE/RAW"}
	}

	parsedStrings, err := parseStringsFromElement(element, []string{"DATA", "TYPE"})
	if err != nil {
		return nil, err
	}

	return &Raw{Data: parsedStrings[0], Type: parsedStrings[1]}, nil
}

// TemplateID gets ID of template of given Virtual Machine.
func (vm *VirtualMachine) TemplateID() (int, error) {
	return vm.intAttribute("TEMPLATE/TEMPLATE_ID")
}

// VCPU gets VCPU of given Virtual Machine.
func (vm *VirtualMachine) VCPU() (int, error) {
	return vm.intAttribute("TEMPLATE/VCPU")
}

// UserTemplate gets UserTemplate of given VM as an element.
func (vm *VirtualMachine) UserTemplate() (UserTemplate, error) {
	ut := vm.XMLData.FindElement("USER_TEMPLATE")
	if ut == nil {
		return nil, &errors.XMLElementError{Path: "USER_TEMPLATE"}
	}
	return ut, nil
}

// HistoryRecords gets an array of History structures of given VM.
func (vm *VirtualMachine) HistoryRecords() ([]*History, error) {
	elements := vm.XMLData.FindElements("HISTORY_RECORDS/HISTORY")
	if len(elements) == 0 {
		return make([]*History, 0), nil
	}

	history := make([]*History, len(elements))
	var err error

	for i, e := range elements {
		history[i], err = createHistoryFromElement(e)
		if err != nil {
			return nil, err
		}
	}
	return history, nil
}

func createHistoryFromElement(element *etree.Element) (*History, error) {
	if element == nil {
		return nil, &errors.XMLElementError{Path: "HISTORY_RECORDS/HISTORY"}
	}

	parsedInts, err := parseIntsFromElement(element, []string{"OID", "SEQ"})
	if err != nil {
		return nil, err
	}

	times, err := parseTimesFromElement(element, []string{"STIME", "ETIME"})
	if err != nil {
		return nil, err
	}

	stringsWithoutError := parseStringsFromElementWithoutError(element, []string{"HOSTNAME", "VM_MAD", "TM_MAD"})
	intsWithoutError := parseIntsFromElementWithoutError(element, []string{"HID", "CID", "DS_ID", "ACTION"})
	timesWithoutError := parseTimesFromElementWithoutError(element, []string{"PSTIME", "PETIME", "RSTIME", "RETIME",
		"ESTIME", "EETIME"})

	return &History{
		Hostname:    stringsWithoutError[0],
		VMMad:       stringsWithoutError[1],
		TmMad:       stringsWithoutError[2],
		OID:         parsedInts[0],
		Seq:         parsedInts[1],
		HID:         intsWithoutError[0],
		CID:         intsWithoutError[1],
		DatastoreID: intsWithoutError[2],
		Action:      Action(*intsWithoutError[3]),
		STime:       times[0],
		ETime:       times[1],
		PSTime:      timesWithoutError[0],
		PETime:      timesWithoutError[1],
		RSTime:      timesWithoutError[2],
		REtime:      timesWithoutError[3],
		ESTime:      timesWithoutError[4],
		EETime:      timesWithoutError[5],
	}, nil
}
