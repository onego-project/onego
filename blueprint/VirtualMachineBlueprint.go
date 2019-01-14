package blueprint

import "strconv"

// VirtualMachineBlueprint to allocate and update OpenNebula virtual machine.
type VirtualMachineBlueprint struct {
	Blueprint
}

// CreateAllocateVirtualMachineBlueprint creates empty VirtualMachineBlueprint.
func CreateAllocateVirtualMachineBlueprint() *VirtualMachineBlueprint {
	return &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("VM")}
}

// CreateUpdateVirtualMachineBlueprint creates empty VirtualMachineBlueprint.
func CreateUpdateVirtualMachineBlueprint() *VirtualMachineBlueprint {
	return &VirtualMachineBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
}

// SetAutomaticDSRequirements sets automatic DS requirements of the given virtual machine.
func (vmb *VirtualMachineBlueprint) SetAutomaticDSRequirements(dsReq string) {
	vmb.SetElement("AUTOMATIC_DS_REQUIREMENTS", dsReq)
}

// SetAutomaticRequirements sets automatic requirements of the given virtual machine.
func (vmb *VirtualMachineBlueprint) SetAutomaticRequirements(req string) {
	vmb.SetElement("AUTOMATIC_REQUIREMENTS", req)
}

// SetContext sets CONTEXT of a given virtual machine.
func (vmb *VirtualMachineBlueprint) SetContext(blueprint ContextBlueprint) {
	vmb.AddElement(*blueprint.XMLData)
}

// SetCPU sets CPU of the given virtual machine.
func (vmb *VirtualMachineBlueprint) SetCPU(cpu float64) {
	vmb.SetElement("CPU", strconv.FormatFloat(cpu, 'f', -1, 64))
}

// SetMemory sets MEMORY of a given virtual machine.
func (vmb *VirtualMachineBlueprint) SetMemory(memory int) {
	vmb.SetElement("MEMORY", strconv.Itoa(memory))
}

// SetDisk sets DISK of a given virtual machine.
func (vmb *VirtualMachineBlueprint) SetDisk(blueprint DiskBlueprint) {
	vmb.AddElement(*blueprint.XMLData)
}

// SetFeatures sets FEATURES of a given virtual machine.
func (vmb *VirtualMachineBlueprint) SetFeatures(blueprint FeaturesBlueprint) {
	vmb.AddElement(*blueprint.XMLData)
}

// SetGraphics sets GRAPHICS of a given virtual machine.
func (vmb *VirtualMachineBlueprint) SetGraphics(blueprint GraphicsBlueprint) {
	vmb.AddElement(*blueprint.XMLData)
}

// SetNIC sets NIC of a given virtual machine.
func (vmb *VirtualMachineBlueprint) SetNIC(blueprint NICBlueprint) {
	vmb.AddElement(*blueprint.XMLData)
}

// SetOS sets OS of a given virtual machine.
func (vmb *VirtualMachineBlueprint) SetOS(blueprint OSBlueprint) {
	vmb.AddElement(*blueprint.XMLData)
}

// SetRaw sets RAW of a given virtual machine.
func (vmb *VirtualMachineBlueprint) SetRaw(blueprint RawBlueprint) {
	vmb.AddElement(*blueprint.XMLData)
}

// SetTemplateID sets template id of a given virtual machine.
func (vmb *VirtualMachineBlueprint) SetTemplateID(id int) {
	vmb.SetElement("TEMPLATE_ID", strconv.Itoa(id))
}

// SetVCPU sets VCPU of a given virtual machine.
func (vmb *VirtualMachineBlueprint) SetVCPU(vcpu int) {
	vmb.SetElement("VCPU", strconv.Itoa(vcpu))
}
