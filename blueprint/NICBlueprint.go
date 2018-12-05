package blueprint

import "strconv"

// NICBlueprint to set  network interface.
type NICBlueprint struct {
	Blueprint
}

// CreateNICBlueprint creates empty VMNICBlueprint.
func CreateNICBlueprint() *NICBlueprint {
	return &NICBlueprint{Blueprint: *CreateBlueprint("NIC")}
}

// SetNetworkName sets NETWORK of a given VM.
func (nb *NICBlueprint) SetNetworkName(value string) {
	nb.SetElement("NETWORK", value)
}

// SetNetworkOwnerName sets NETWORK_UNAME of a given VM.
func (nb *NICBlueprint) SetNetworkOwnerName(value string) {
	nb.SetElement("NETWORK_UNAME", value)
}

// SetNetworkID sets NETWORK_ID of a given VM.
func (nb *NICBlueprint) SetNetworkID(value int) {
	nb.SetElement("NETWORK_ID", strconv.Itoa(value))
}

// SetNetworkOwnerID sets NETWORK_UID of a given VM.
func (nb *NICBlueprint) SetNetworkOwnerID(value int) {
	nb.SetElement("NETWORK_UID", strconv.Itoa(value))
}
