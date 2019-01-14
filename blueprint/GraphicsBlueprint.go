package blueprint

import (
	"net"

	"github.com/onego-project/onego/resources"
)

// GraphicsBlueprint to set  elements.
type GraphicsBlueprint struct {
	Blueprint
}

// CreateGraphicsBlueprint creates empty GraphicsBlueprint.
func CreateGraphicsBlueprint() *GraphicsBlueprint {
	return &GraphicsBlueprint{Blueprint: *CreateBlueprint("GRAPHICS")}
}

// SetListen sets listen on IP of a given Graphics.
func (gb *GraphicsBlueprint) SetListen(ip net.IP) { // nolint: interfacer
	gb.SetElement("LISTEN", ip.String())
}

// SetType sets TYPE of a given Graphics.
func (gb *GraphicsBlueprint) SetType(value resources.GraphicsType) {
	gb.SetElement("TYPE", resources.GraphicsTypeMap[value])
}

// SetPort sets server port for GraphicsTypeVNC/SPICE server of a given Graphics.
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
