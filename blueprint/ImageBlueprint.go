package blueprint

import "github.com/onego-project/onego/resources"

// ImageBlueprint to set Image elements
type ImageBlueprint struct {
	Blueprint
}

// CreateAllocateImageBlueprint creates empty ImageBlueprint
func CreateAllocateImageBlueprint() *ImageBlueprint {
	return &ImageBlueprint{Blueprint: *CreateBlueprint("IMAGE")}
}

// CreateUpdateImageBlueprint creates empty ImageBlueprint
func CreateUpdateImageBlueprint() *ImageBlueprint {
	return &ImageBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
}

// SetDescription sets description of the given image
func (ds *ImageBlueprint) SetDescription(desc string) {
	ds.SetElement("DESCRIPTION", desc)
}

// SetDevPrefix sets dev. prefix of the given image
func (ds *ImageBlueprint) SetDevPrefix(devPrefix string) {
	ds.SetElement("DEV_PREFIX", devPrefix)
}

// SetDiskType sets disk type of the given image
func (ds *ImageBlueprint) SetDiskType(diskType resources.DiskType) {
	ds.SetElement("DISK_TYPE", resources.DiskTypeMap[diskType])
}

// SetDriver sets description of the given image
func (ds *ImageBlueprint) SetDriver(driver string) {
	ds.SetElement("DRIVER", driver)
}

// SetTarget sets target of the given image
func (ds *ImageBlueprint) SetTarget(target string) {
	ds.SetElement("TARGET", target)
}
