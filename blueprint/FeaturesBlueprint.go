package blueprint

// FeaturesBlueprint to set  elements.
type FeaturesBlueprint struct {
	Blueprint
}

// CreateFeaturesBlueprint creates empty FeaturesBlueprint.
func CreateFeaturesBlueprint() *FeaturesBlueprint {
	return &FeaturesBlueprint{Blueprint: *CreateBlueprint("FEATURES")}
}

// SetGuestAgent sets GUEST_AGENT of a given .
func (fb *FeaturesBlueprint) SetGuestAgent(value string) {
	fb.SetElement("GUEST_AGENT", value)
}
