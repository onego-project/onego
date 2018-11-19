package blueprint

// VMTemplateOSBlueprint to set VMTemplate elements.
type VMTemplateOSBlueprint struct {
	Blueprint
}

// CreateVMTemplateOSBlueprint creates empty VMTemplateOSBlueprint.
func CreateVMTemplateOSBlueprint() *VMTemplateOSBlueprint {
	return &VMTemplateOSBlueprint{Blueprint: *CreateBlueprint("OS")}
}

// SetArch sets ARCH of a given VMTemplate.
func (vmtob *VMTemplateOSBlueprint) SetArch(value string) {
	vmtob.SetElement("ARCH", value)
}
