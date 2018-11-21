package blueprint

// RawBlueprint to set  elements.
type RawBlueprint struct {
	Blueprint
}

// CreateRAWBlueprint creates empty RawBlueprint.
func CreateRAWBlueprint() *RawBlueprint {
	return &RawBlueprint{Blueprint: *CreateBlueprint("RAW")}
}

// SetData sets DATA of a given .
func (rb *RawBlueprint) SetData(value string) {
	rb.SetElement("DATA", value)
}

// SetType sets TYPE of a given .
func (rb *RawBlueprint) SetType(value string) {
	rb.SetElement("TYPE", value)
}
