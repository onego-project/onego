package blueprint

import (
	"net"
	"strconv"
)

// ContextBlueprint to set Context elements.
type ContextBlueprint struct {
	Blueprint
}

// CreateContextBlueprint creates empty ContextBlueprint.
func CreateContextBlueprint() *ContextBlueprint {
	return &ContextBlueprint{Blueprint: *CreateBlueprint("CONTEXT")}
}

// SetDiskID sets DISK_ID in a Context.
func (cb *ContextBlueprint) SetDiskID(value int) {
	cb.SetElement("DISK_ID", strconv.Itoa(value))
}

// SetEmail sets EMAIL in a Context.
func (cb *ContextBlueprint) SetEmail(value string) {
	cb.SetElement("EMAIL", value)
}

// SetOnegateEndpoint sets ONEGATE_ENDPOINT in a Context.
func (cb *ContextBlueprint) SetOnegateEndpoint(value string) {
	cb.SetElement("ONEGATE_ENDPOINT", value)
}

// SetPublicIP sets PUBLIC_IP in a Context.
func (cb *ContextBlueprint) SetPublicIP(value net.IP) {
	cb.SetElement("PUBLIC_IP", value.String())
}

// SetSSHKey sets SSH_KEY in a Context.
func (cb *ContextBlueprint) SetSSHKey(value string) {
	cb.SetElement("SSH_KEY", value)
}

// SetTarget sets TARGET in a Context.
func (cb *ContextBlueprint) SetTarget(value string) {
	cb.SetElement("TARGET", value)
}

// SetToken sets TOKEN (YES/NO) in a Context.
func (cb *ContextBlueprint) SetToken(value bool) {
	cb.SetElement("TOKEN", boolToString(value))
}

// SetUserData sets USER_DATA in a Context.
func (cb *ContextBlueprint) SetUserData(value string) {
	cb.SetElement("USER_DATA", value)
}

// SetVirtualMachineID sets VMID in a Context.
func (cb *ContextBlueprint) SetVirtualMachineID(value int) {
	cb.SetElement("VMID", strconv.Itoa(value))
}

// SetVirtualMachineGroupID sets VM_GID in a Context.
func (cb *ContextBlueprint) SetVirtualMachineGroupID(value int) {
	cb.SetElement("VM_GID", strconv.Itoa(value))
}

// SetVirtualMachineGroupName sets VM_GNAME in a Context.
func (cb *ContextBlueprint) SetVirtualMachineGroupName(value string) {
	cb.SetElement("VM_GNAME", value)
}

// SetVirtualMachineUserID sets VM_UID in a Context.
func (cb *ContextBlueprint) SetVirtualMachineUserID(value int) {
	cb.SetElement("VM_UID", strconv.Itoa(value))
}

// SetVirtualMachineUserName sets VM_UNAME in a Context.
func (cb *ContextBlueprint) SetVirtualMachineUserName(value string) {
	cb.SetElement("VM_UNAME", value)
}
