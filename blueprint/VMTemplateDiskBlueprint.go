package blueprint

// VMTemplateDiskBlueprint to set VMTemplate elements.
type VMTemplateDiskBlueprint struct {
	Blueprint
}

// CreateVMTemplateDiskBlueprint creates empty VMTemplateDiskBlueprint.
func CreateVMTemplateDiskBlueprint() *VMTemplateDiskBlueprint {
	return &VMTemplateDiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
}

// SetDevicePrefix sets device prefix in disk in VM template.
func (vmtdb *VMTemplateDiskBlueprint) SetDevicePrefix(value string) {
	vmtdb.SetElement("DEV_PREFIX", value)
}

// SetImage sets image in disk in VM template.
func (vmtdb *VMTemplateDiskBlueprint) SetImage(value string) {
	vmtdb.SetElement("IMAGE", value)
}
