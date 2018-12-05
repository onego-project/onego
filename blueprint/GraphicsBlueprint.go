package blueprint

import "net"

// GraphicsBlueprint to set  elements.
type GraphicsBlueprint struct {
	Blueprint
}

// GraphicsTypeMap to convert GraphicsType to string.
var GraphicsTypeMap = map[GraphicsType]string{
	VNC:   "VNC",
	Spice: "SPICE",
	SDL:   "SDL",
	None:  "NONE",
}

// GraphicsType to set graphics type.
type GraphicsType int

const (
	// VNC graphics type.
	VNC GraphicsType = iota
	// Spice graphics type.
	Spice
	// SDL graphics type.
	SDL
	// None graphics type.
	None
)

// CreateGraphicsBlueprint creates empty GraphicsBlueprint.
func CreateGraphicsBlueprint() *GraphicsBlueprint {
	return &GraphicsBlueprint{Blueprint: *CreateBlueprint("GRAPHICS")}
}

// SetListen sets listen on IP of a given Graphics.
func (gb *GraphicsBlueprint) SetListen(ip net.IP) { // nolint: interfacer
	gb.SetElement("LISTEN", ip.String())
}

// SetType sets TYPE of a given Graphics.
func (gb *GraphicsBlueprint) SetType(value GraphicsType) {
	gb.SetElement("TYPE", GraphicsTypeMap[value])
}

// SetPort sets server port for VNC/SPICE server of a given Graphics.
func (gb *GraphicsBlueprint) SetPort(port string) {
	gb.SetElement("PORT", port)
}

// SetKeyMap sets key map of a given Graphics.
func (gb *GraphicsBlueprint) SetKeyMap(value string) {
	gb.SetElement("KEYMAP", value)
}

// SetPassword sets password of a given Graphics.
func (gb *GraphicsBlueprint) SetPassword(value string) {
	gb.SetElement("PASSWD", value)
}
