package blueprint

// VMTemplateNICBlueprint to set VMTemplate network interface.
type VMTemplateNICBlueprint struct {
	Blueprint
}

// CreateVMTemplateNICBlueprint creates empty VMTemplateNICBlueprint.
func CreateVMTemplateNICBlueprint() *VMTemplateNICBlueprint {
	return &VMTemplateNICBlueprint{Blueprint: *CreateBlueprint("NIC")}
}

// SetNetwork sets NETWORK of a given VMTemplate.
func (vmtnb *VMTemplateNICBlueprint) SetNetwork(value string) {
	vmtnb.SetElement("NETWORK", value)
}
