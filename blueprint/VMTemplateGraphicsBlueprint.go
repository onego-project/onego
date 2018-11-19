package blueprint

// VMTemplateGraphicsBlueprint to set VMTemplate elements.
type VMTemplateGraphicsBlueprint struct {
	Blueprint
}

// CreateVMTemplateGraphicsBlueprint creates empty VMTemplateGraphicsBlueprint.
func CreateVMTemplateGraphicsBlueprint() *VMTemplateGraphicsBlueprint {
	return &VMTemplateGraphicsBlueprint{Blueprint: *CreateBlueprint("GRAPHICS")}
}

// SetListen sets LISTEN of a given VMTemplate.
func (vmtgb *VMTemplateGraphicsBlueprint) SetListen(value string) {
	vmtgb.SetElement("LISTEN", value)
}

// SetType sets TYPE of a given VMTemplate.
func (vmtgb *VMTemplateGraphicsBlueprint) SetType(value string) {
	vmtgb.SetElement("TYPE", value)
}
