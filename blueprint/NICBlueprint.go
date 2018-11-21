package blueprint

// NICBlueprint to set  network interface.
type NICBlueprint struct {
	Blueprint
}

// CreateNICBlueprint creates empty VMNICBlueprint.
func CreateNICBlueprint() *NICBlueprint {
	return &NICBlueprint{Blueprint: *CreateBlueprint("NIC")}
}

// SetNetwork sets NETWORK of a given VM.
func (nb *NICBlueprint) SetNetwork(value string) {
	nb.SetElement("NETWORK", value)
}
