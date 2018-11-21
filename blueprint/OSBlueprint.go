package blueprint

// OSBlueprint to set  elements.
type OSBlueprint struct {
	Blueprint
}

// CreateOSBlueprint creates empty OSBlueprint.
func CreateOSBlueprint() *OSBlueprint {
	return &OSBlueprint{Blueprint: *CreateBlueprint("OS")}
}

// SetArchitecture sets ARCH of a given OS.
func (ob *OSBlueprint) SetArchitecture(value string) {
	ob.SetElement("ARCH", value)
}
