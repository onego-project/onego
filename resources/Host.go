package resources

import (
	"github.com/onego-project/onego/errors"

	"github.com/beevik/etree"
)

// Host structure represents OpenNebula Host
type Host struct {
	Resource
}

//  HOST STATES                   +----------------+
//                                |  VM DEPLOYMENT |
//  +----------------+------------+--------+-------+
//  | STATE          | MONITORING | MANUAL | SCHED |
//  +----------------+------------+--------+-------+
//  | INIT/MONITORED |    Yes     |       Yes      |
//  +----------------+------------+--------+-------+
//  | DISABLED       |    Yes     | Yes    |  No   |
//  +----------------+------------+----------------+
//  | OFFLINE        |    No      |        No      |
//  +----------------+-----------------------------+

// HostState - State of Hosts
type HostState int

const (
	// HostInit - Initial state for enabled hosts.
	HostInit HostState = iota
	// HostMonitoringMonitored - Monitoring the host (from monitored).
	HostMonitoringMonitored
	// HostMonitored - The host has been monitored.
	HostMonitored
	// HostError - An error occurred in host monitoring.
	HostError
	// HostDisabled - The host is disabled see above.
	HostDisabled
	// HostMonitoringError - Monitoring the host (from error).
	HostMonitoringError
	// HostMonitoringInit - Monitoring the host (from init).
	HostMonitoringInit
	// HostMonitoringDisabled - Monitoring the host (from disabled).
	HostMonitoringDisabled
	// HostOffline - The host is set offline, see above
	HostOffline
)

// HostStateMap contains string representation of HostState
var HostStateMap = map[HostState]string{
	HostInit:                "INIT",
	HostMonitoringMonitored: "MONITORING_MONITORED",
	HostMonitored:           "MONITORED",
	HostError:               "ERROR",
	HostDisabled:            "DISABLED",
	HostMonitoringError:     "MONITORING_ERROR",
	HostMonitoringInit:      "MONITORING_INIT",
	HostMonitoringDisabled:  "MONITORING_DISABLED",
	HostOffline:             "OFFLINE",
}

// PCI structure represents hardware for Host
type PCI struct {
	Address      string
	Bus          string
	Class        string
	ClassName    string
	Device       string
	DeviceName   string
	Domain       string
	Function     string
	ShortAddress string
	Slot         string
	Type         string
	Vendor       string
	VendorName   string
	VMID         string
}

// CreateHostWithID constructs Host with id
func CreateHostWithID(id int) *Host {
	return &Host{*CreateResource("HOST", id)}
}

// CreateHostFromXML constructs Host with full xml data
func CreateHostFromXML(XMLdata *etree.Element) *Host {
	return &Host{Resource: Resource{XMLData: XMLdata}}
}

// State gets Host state.
// To get string representation of host state use HostStateMap.
func (h *Host) State() (HostState, error) {
	state, err := h.intAttribute("STATE")
	return HostState(state), err
}

// IMMad method
func (h *Host) IMMad() (string, error) {
	return h.Attribute("IM_MAD")
}

// VMMad method
func (h *Host) VMMad() (string, error) {
	return h.Attribute("VM_MAD")
}

// LastMonitoringTime gets last monitoring time
func (h *Host) LastMonitoringTime() (int, error) {
	return h.intAttribute("LAST_MON_TIME")
}

// Cluster gets Cluser ID of given Host
func (h *Host) Cluster() (int, error) {
	return h.intAttribute("CLUSTER_ID")
}

// DiskUsage gets disk usage
func (h *Host) DiskUsage() (int, error) {
	return h.intAttribute("HOST_SHARE/DISK_USAGE")
}

// MemoryUsage gets memory usage
func (h *Host) MemoryUsage() (int, error) {
	return h.intAttribute("HOST_SHARE/MEM_USAGE")
}

// CPUUsage gets cpu usage
func (h *Host) CPUUsage() (int, error) {
	return h.intAttribute("HOST_SHARE/CPU_USAGE")
}

// MaxDisk gets maximal disk size
func (h *Host) MaxDisk() (int, error) {
	return h.intAttribute("HOST_SHARE/MAX_DISK")
}

// MaxMemory gets maximal memory size
func (h *Host) MaxMemory() (int, error) {
	return h.intAttribute("HOST_SHARE/MAX_MEM")
}

// MaxCPU gets maximal cpu size
func (h *Host) MaxCPU() (int, error) {
	return h.intAttribute("HOST_SHARE/MAX_CPU")
}

// FreeDisk gets free disk size
func (h *Host) FreeDisk() (int, error) {
	return h.intAttribute("HOST_SHARE/FREE_DISK")
}

// FreeMemory gets free memory size
func (h *Host) FreeMemory() (int, error) {
	return h.intAttribute("HOST_SHARE/FREE_MEM")
}

// FreeCPU gets free cpu size
func (h *Host) FreeCPU() (int, error) {
	return h.intAttribute("HOST_SHARE/FREE_CPU")
}

// UsedDisk gets used disk size
func (h *Host) UsedDisk() (int, error) {
	return h.intAttribute("HOST_SHARE/USED_DISK")
}

// UsedMemory gets used memory size
func (h *Host) UsedMemory() (int, error) {
	return h.intAttribute("HOST_SHARE/USED_MEM")
}

// UsedCPU gets used cpu size
func (h *Host) UsedCPU() (int, error) {
	return h.intAttribute("HOST_SHARE/USED_CPU")
}

// RunningVMs gets number of running VMs
func (h *Host) RunningVMs() (int, error) {
	return h.intAttribute("HOST_SHARE/RUNNING_VMS")
}

// Datastores gets array of datastore IDs important for Host
func (h *Host) Datastores() ([]int, error) {
	return h.arrayOfIDs("HOST_SHARE/DATASTORES/DS")
}

// PCIDevices gets PCI devices of given Host
func (h *Host) PCIDevices() ([]*PCI, error) {
	elements := h.XMLData.FindElements("HOST_SHARE/PCI_DEVICES/PCI")
	if len(elements) == 0 {
		return make([]*PCI, 0), nil
	}

	pcis := make([]*PCI, len(elements))
	var err error

	for i, e := range elements {
		pcis[i], err = createPCIFromElement(e)
		if err != nil {
			return nil, err
		}
	}
	return pcis, nil
}

func createPCIFromElement(element *etree.Element) (*PCI, error) {
	if element == nil {
		return nil, &errors.XMLElementError{Path: "HOST_SHARE/PCI_DEVICES/PCI"}
	}

	parse := []string{"ADDRESS", "BUS", "CLASS", "CLASS_NAME", "DEVICE", "DEVICE_NAME", "DOMAIN",
		"FUNCTION", "SHORT_ADDRESS", "SLOT", "TYPE", "VENDOR", "VENDOR_NAME", "VMID"}
	parsed := make([]string, len(parse))
	var err error

	for i := 0; i < len(parse); i++ {
		parsed[i], err = attributeFromElement(element, parse[i])
		if err != nil {
			return nil, err
		}
	}

	return &PCI{Address: parsed[0], Bus: parsed[1], Class: parsed[2], ClassName: parsed[3], Device: parsed[4],
		DeviceName: parsed[5], Domain: parsed[6], Function: parsed[7], ShortAddress: parsed[8], Slot: parsed[9],
		Type: parsed[10], Vendor: parsed[11], VendorName: parsed[12], VMID: parsed[13]}, nil
}

// VirtualMachines gets array of VM IDs of given Host
func (h *Host) VirtualMachines() ([]int, error) {
	return h.arrayOfIDs("VMS")
}
