package blueprint

// VMTemplateRAWBlueprint to set VMTemplate elements.
type VMTemplateRAWBlueprint struct {
	Blueprint
}

// CreateVMTemplateRAWBlueprint creates empty VMTemplateRAWBlueprint.
func CreateVMTemplateRAWBlueprint() *VMTemplateRAWBlueprint {
	return &VMTemplateRAWBlueprint{Blueprint: *CreateBlueprint("RAW")}
}

// SetData sets DATA of a given VMTemplate.
func (vmtrb *VMTemplateRAWBlueprint) SetData(value string) {
	vmtrb.SetElement("DATA", value)
}

// SetType sets TYPE of a given VMTemplate.
func (vmtrb *VMTemplateRAWBlueprint) SetType(value string) {
	vmtrb.SetElement("TYPE", value)
}
