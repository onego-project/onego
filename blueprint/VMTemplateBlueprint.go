package blueprint

// VMTemplateBlueprint to set VMTemplate elements.
type VMTemplateBlueprint struct {
	Blueprint
}

// CreateAllocateVMTemplateBlueprint creates empty VMTemplateBlueprint.
func CreateAllocateVMTemplateBlueprint() *VMTemplateBlueprint {
	return &VMTemplateBlueprint{Blueprint: *CreateBlueprint("VMTEMPLATE")}
}

// CreateUpdateVMTemplateBlueprint creates empty VMTemplateBlueprint.
func CreateUpdateVMTemplateBlueprint() *VMTemplateBlueprint {
	return &VMTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
}

// SetCPU sets CPU of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetCPU(cpu string) {
	tb.SetElement("CPU", cpu)
}

// SetDescription sets DESCRIPTION of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetDescription(description string) {
	tb.SetElement("DESCRIPTION", description)
}

// SetLogo sets LOGO of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetLogo(logoPath string) {
	tb.SetElement("LOGO", logoPath)
}

// SetMemory sets MEMORY of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetMemory(memory string) {
	tb.SetElement("MEMORY", memory)
}

// SetSchedRequirements sets SCHED_REQUIREMENTS of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetSchedRequirements(schedReqs string) {
	tb.SetElement("SCHED_REQUIREMENTS", schedReqs)
}

// SetVCPU sets VCPU of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetVCPU(vcpu string) {
	tb.SetElement("VCPU", vcpu)
}
