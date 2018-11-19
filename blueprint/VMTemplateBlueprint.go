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

// SetDisk sets DISK of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetDisk(diskBlueprint VMTemplateDiskBlueprint) {
	tb.XMLData.Root().AddChild(diskBlueprint.XMLData.Root())
}

// SetFeatures sets FEATURES of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetFeatures(featuresBlueprint VMTemplateFeaturesBlueprint) {
	tb.XMLData.Root().AddChild(featuresBlueprint.XMLData.Root())
}

// SetGraphics sets GRAPHICS of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetGraphics(graphicsBlueprint VMTemplateGraphicsBlueprint) {
	tb.XMLData.Root().AddChild(graphicsBlueprint.XMLData.Root())
}

// SetLogo sets LOGO of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetLogo(logoPath string) {
	tb.SetElement("LOGO", logoPath)
}

// SetMemory sets MEMORY of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetMemory(memory string) {
	tb.SetElement("MEMORY", memory)
}

// SetNIC sets NIC of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetNIC(blueprint VMTemplateNICBlueprint) {
	tb.XMLData.Root().AddChild(blueprint.XMLData.Root())
}

// SetOS sets OS of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetOS(blueprint VMTemplateOSBlueprint) {
	tb.XMLData.Root().AddChild(blueprint.XMLData.Root())
}

// SetRAW sets RAW of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetRAW(blueprint VMTemplateRAWBlueprint) {
	tb.XMLData.Root().AddChild(blueprint.XMLData.Root())
}

// SetSchedRequirements sets SCHED_REQUIREMENTS of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetSchedRequirements(schedReqs string) {
	tb.SetElement("SCHED_REQUIREMENTS", schedReqs)
}

// SetVCPU sets VCPU of a given VMTemplate.
func (tb *VMTemplateBlueprint) SetVCPU(vcpu string) {
	tb.SetElement("VCPU", vcpu)
}
