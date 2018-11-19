package blueprint

// VMTemplateFeaturesBlueprint to set VMTemplate elements.
type VMTemplateFeaturesBlueprint struct {
	Blueprint
}

// CreateVMTemplateFeaturesBlueprint creates empty VMTemplateFeaturesBlueprint.
func CreateVMTemplateFeaturesBlueprint() *VMTemplateFeaturesBlueprint {
	return &VMTemplateFeaturesBlueprint{Blueprint: *CreateBlueprint("FEATURES")}
}

// SetGuestAgent sets GUEST_AGENT of a given VMTemplate.
func (vmtfb *VMTemplateFeaturesBlueprint) SetGuestAgent(value string) {
	vmtfb.SetElement("GUEST_AGENT", value)
}
