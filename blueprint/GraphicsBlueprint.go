package blueprint

// GraphicsBlueprint to set  elements.
type GraphicsBlueprint struct {
	Blueprint
}

// GraphicsMap to convert GraphicsType to string.
var GraphicsMap = map[GraphicsType]string{
	VNC:   "VNC",
	Spice: "SPICE",
}

// GraphicsType to set graphics type.
type GraphicsType int

const (
	// VNC graphics type.
	VNC GraphicsType = iota
	// Spice graphics type.
	Spice
)

// CreateGraphicsBlueprint creates empty GraphicsBlueprint.
func CreateGraphicsBlueprint() *GraphicsBlueprint {
	return &GraphicsBlueprint{Blueprint: *CreateBlueprint("GRAPHICS")}
}

// SetListen sets LISTEN of a given Graphics.
func (gb *GraphicsBlueprint) SetListen(value string) {
	gb.SetElement("LISTEN", value)
}

// SetType sets TYPE of a given Graphics.
func (gb *GraphicsBlueprint) SetType(value GraphicsType) {
	gb.SetElement("TYPE", GraphicsMap[value])
}
