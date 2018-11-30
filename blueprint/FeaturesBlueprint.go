package blueprint

// FeaturesBlueprint to set feature elements.
type FeaturesBlueprint struct {
	Blueprint
}

// CreateFeaturesBlueprint creates empty FeaturesBlueprint.
func CreateFeaturesBlueprint() *FeaturesBlueprint {
	return &FeaturesBlueprint{Blueprint: *CreateBlueprint("FEATURES")}
}

// SetACPI sets ACPI (YES/NO) to enable Advance Configuration and Power Interface.
func (fb *FeaturesBlueprint) SetACPI(value bool) {
	fb.SetElement("ACPI", boolToString(value))
}

// SetAPIC sets APIC (YES/NO) to enable Advanced programmable IRQ management.
func (fb *FeaturesBlueprint) SetAPIC(value bool) {
	fb.SetElement("APIC", boolToString(value))
}

// SetLocalTime sets LOCALTIME (YES/NO).
// The guest clock will be synchronized to the hosts configured timezone when booted.
func (fb *FeaturesBlueprint) SetLocalTime(value bool) {
	fb.SetElement("LOCALTIME", boolToString(value))
}

// SetPAE sets PAE (YES/NO) to enable Physical Address Extension.
func (fb *FeaturesBlueprint) SetPAE(value bool) {
	fb.SetElement("PAE", boolToString(value))
}

// SetHyperV sets HYPERV (YES/NO) to enable hyper-v feature.
func (fb *FeaturesBlueprint) SetHyperV(value bool) {
	fb.SetElement("HYPERV", boolToString(value))
}

// SetGuestAgent sets GUEST_AGENT (YES/NO) to enable QEMU Guest Agent communication.
func (fb *FeaturesBlueprint) SetGuestAgent(value bool) {
	fb.SetElement("GUEST_AGENT", boolToString(value))
}
