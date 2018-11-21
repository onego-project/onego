package blueprint

// GraphicsBlueprint to set  elements.
type GraphicsBlueprint struct {
	Blueprint
}

// CreateGraphicsBlueprint creates empty GraphicsBlueprint.
func CreateGraphicsBlueprint() *GraphicsBlueprint {
	return &GraphicsBlueprint{Blueprint: *CreateBlueprint("GRAPHICS")}
}

// SetListen sets LISTEN of a given .
func (gb *GraphicsBlueprint) SetListen(value string) {
	gb.SetElement("LISTEN", value)
}

// SetType sets TYPE of a given .
func (gb *GraphicsBlueprint) SetType(value string) {
	gb.SetElement("TYPE", value)
}
