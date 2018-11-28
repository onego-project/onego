package blueprint

import (
	"strconv"
)

// TemplateBlueprint to set Template elements.
type TemplateBlueprint struct {
	Blueprint
}

// CreateAllocateTemplateBlueprint creates empty TemplateBlueprint.
func CreateAllocateTemplateBlueprint() *TemplateBlueprint {
	return &TemplateBlueprint{Blueprint: *CreateBlueprint("VMTEMPLATE")}
}

// CreateUpdateTemplateBlueprint creates empty TemplateBlueprint.
func CreateUpdateTemplateBlueprint() *TemplateBlueprint {
	return &TemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
}

// SetCPU sets CPU of a given Template.
func (tb *TemplateBlueprint) SetCPU(cpu float64) {
	tb.SetElement("CPU", strconv.FormatFloat(cpu, 'f', -1, 64))
}

// SetDescription sets DESCRIPTION of a given Template.
func (tb *TemplateBlueprint) SetDescription(description string) {
	tb.SetElement("DESCRIPTION", description)
}

// SetDisk sets DISK of a given Template.
func (tb *TemplateBlueprint) SetDisk(diskBlueprint DiskBlueprint) {
	tb.XMLData.Root().AddChild(diskBlueprint.XMLData.Root())
}

// SetFeatures sets FEATURES of a given Template.
func (tb *TemplateBlueprint) SetFeatures(featuresBlueprint FeaturesBlueprint) {
	tb.XMLData.Root().AddChild(featuresBlueprint.XMLData.Root())
}

// SetGraphics sets GRAPHICS of a given Template.
func (tb *TemplateBlueprint) SetGraphics(graphicsBlueprint GraphicsBlueprint) {
	tb.XMLData.Root().AddChild(graphicsBlueprint.XMLData.Root())
}

// SetLogo sets LOGO of a given Template.
func (tb *TemplateBlueprint) SetLogo(logoPath string) {
	tb.SetElement("LOGO", logoPath)
}

// SetMemory sets MEMORY of a given Template.
func (tb *TemplateBlueprint) SetMemory(memory int) {
	tb.SetElement("MEMORY", strconv.Itoa(memory))
}

// SetNIC sets NIC of a given Template.
func (tb *TemplateBlueprint) SetNIC(blueprint NICBlueprint) {
	tb.XMLData.Root().AddChild(blueprint.XMLData.Root())
}

// SetOS sets OS of a given Template.
func (tb *TemplateBlueprint) SetOS(blueprint OSBlueprint) {
	tb.XMLData.Root().AddChild(blueprint.XMLData.Root())
}

// SetRaw sets RAW of a given Template.
func (tb *TemplateBlueprint) SetRaw(blueprint RawBlueprint) {
	tb.XMLData.Root().AddChild(blueprint.XMLData.Root())
}

// SetSchedulingRequirements sets SCHEDULING REQUIREMENTS of a given Template.
func (tb *TemplateBlueprint) SetSchedulingRequirements(schedReqs string) {
	tb.SetElement("SCHED_REQUIREMENTS", schedReqs)
}

// SetVCPU sets VCPU of a given Template.
func (tb *TemplateBlueprint) SetVCPU(vcpu int) {
	tb.SetElement("VCPU", strconv.Itoa(vcpu))
}
