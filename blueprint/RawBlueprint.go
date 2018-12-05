package blueprint

// RawBlueprint to set Raw elements.
type RawBlueprint struct {
	Blueprint
}

// CreateRawBlueprint creates empty RawBlueprint.
func CreateRawBlueprint() *RawBlueprint {
	return &RawBlueprint{Blueprint: *CreateBlueprint("RAW")}
}

// SetData sets DATA.
func (rb *RawBlueprint) SetData(value string) {
	rb.SetElement("DATA", value)
}

// SetType sets TYPE.
func (rb *RawBlueprint) SetType(value string) {
	rb.SetElement("TYPE", value)
}
