package blueprint

import (
	"strconv"

	"github.com/onego-project/onego/resources"
)

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

// SetClone sets clone (YES/NO) in a disk.
func (db *DiskBlueprint) SetClone(value bool) {
	db.SetElement("CLONE", boolToString(value))
}

// SetCloneTarget sets clone target in a disk.
func (db *DiskBlueprint) SetCloneTarget(value string) {
	db.SetElement("CLONE_TARGET", value)
}

// SetClusterID sets cluster ID in a disk.
func (db *DiskBlueprint) SetClusterID(value int) {
	db.SetElement("CLUSTER_ID", strconv.Itoa(value))
}

// SetDatastoreID sets datastore ID in a disk.
func (db *DiskBlueprint) SetDatastoreID(value int) {
	db.SetElement("DATASTORE_ID", strconv.Itoa(value))
}

// SetDatastore sets datastore in a disk.
func (db *DiskBlueprint) SetDatastore(value string) {
	db.SetElement("DATASTORE", value)
}

// SetDiskID sets disk ID in a disk.
func (db *DiskBlueprint) SetDiskID(value int) {
	db.SetElement("DISK_ID", strconv.Itoa(value))
}

// SetDiskType sets disk type in a disk.
func (db *DiskBlueprint) SetDiskType(value resources.DiskType) {
	db.SetElement("DISK_TYPE", resources.DiskTypeMap[value])
}

// SetDriver sets driver in a disk.
func (db *DiskBlueprint) SetDriver(value string) {
	db.SetElement("DRIVER", value)
}

// SetImage sets image in a disk.
func (db *DiskBlueprint) SetImage(value string) {
	db.SetElement("IMAGE", value)
}

// SetImageID sets image id in a disk.
func (db *DiskBlueprint) SetImageID(value int) {
	db.SetElement("IMAGE_ID", strconv.Itoa(value))
}

// SetImageUserName sets image user name in a disk.
func (db *DiskBlueprint) SetImageUserName(value string) {
	db.SetElement("IMAGE_UNAME", value)
}

// SetPoolName sets pool name in a disk.
func (db *DiskBlueprint) SetPoolName(value string) {
	db.SetElement("POOL_NAME", value)
}

// SetReadOnly sets if the disk is read only (YES/NO).
func (db *DiskBlueprint) SetReadOnly(value bool) {
	db.SetElement("READONLY", boolToString(value))
}

// SetSave sets if the disk is saved (YES/NO).
func (db *DiskBlueprint) SetSave(value bool) {
	db.SetElement("SAVE", boolToString(value))
}

// SetSize sets size in a disk.
func (db *DiskBlueprint) SetSize(value int) {
	db.SetElement("SIZE", strconv.Itoa(value))
}

// SetSource sets source of a disk.
func (db *DiskBlueprint) SetSource(value string) {
	db.SetElement("SOURCE", value)
}

// SetTarget sets target of a disk.
func (db *DiskBlueprint) SetTarget(value string) {
	db.SetElement("TARGET", value)
}

// SetTmMad sets TM_MAD in a disk.
func (db *DiskBlueprint) SetTmMad(value string) {
	db.SetElement("TM_MAD", value)
}

// SetType sets type in a disk.
func (db *DiskBlueprint) SetType(value resources.DiskType) {
	db.SetElement("TYPE", resources.DiskTypeMap[value])
}
