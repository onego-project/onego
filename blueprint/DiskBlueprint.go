package blueprint

// DiskBlueprint to set  elements.
type DiskBlueprint struct {
	Blueprint
}

// CreateDiskBlueprint creates empty DiskBlueprint.
func CreateDiskBlueprint() *DiskBlueprint {
	return &DiskBlueprint{Blueprint: *CreateBlueprint("DISK")}
}

// SetDevicePrefix sets device prefix in disk in VM template.
func (db *DiskBlueprint) SetDevicePrefix(value string) {
	db.SetElement("DEV_PREFIX", value)
}

// SetImage sets image in a disk.
func (db *DiskBlueprint) SetImage(value string) {
	db.SetElement("IMAGE", value)
}
