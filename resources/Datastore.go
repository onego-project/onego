package resources

import "github.com/beevik/etree"

// Datastore structure represents OpenNebula datastore
type Datastore struct {
	Resource
}

// DatastoreTypeMap contains string representation of DatastoreType
var DatastoreTypeMap = map[DatastoreType]string{
	ImageDs:  "IMAGE_DS",
	SystemDs: "SYSTEM_DS",
	FileDs:   "FILE_DS",
}

// DatastoreType - type of datastore
type DatastoreType int

const (
	// ImageDs - Standard datastore for disk images
	ImageDs DatastoreType = iota
	// SystemDs - System datastore for disks of running VMs
	SystemDs
	// FileDs - File datastore for context, kernel, initrd files
	FileDs
)

// DatastoreDiskTypeMap contains string representation of DatastoreDiskType
var DatastoreDiskTypeMap = map[DatastoreDiskType]string{
	File:          "FILE",
	Block:         "BLOCK",
	Rbd:           "RBD",
	RbdCdRom:      "RBD_CDROM",
	Gluster:       "GLUSTER",
	GlusterCdRom:  "GLUSTER_CDROM",
	Sheepdog:      "SHEEPDOG",
	SheepdogCdRom: "SHEEPDOG_CDROM",
	Iscsi:         "ISCSI",
	None:          "",
}

// DatastoreDiskType - Type of Disks (used by the VMM_MAD). Values: BLOCK, CDROM or FILE
type DatastoreDiskType int

const (
	// File - File-based disk
	File DatastoreDiskType = iota
	// CdRom - An ISO9660 disk
	CdRom
	// Block - Block-device disk
	Block
	// Rbd - CEPH RBD disk
	Rbd
	// RbdCdRom - CEPH RBD CDROM disk
	RbdCdRom
	// Gluster - Gluster Block Device
	Gluster
	// GlusterCdRom - Gluster CDROM Device Device
	GlusterCdRom
	// Sheepdog - Sheepdog Block Device
	Sheepdog
	// SheepdogCdRom - Sheepdog CDROM Device Device
	SheepdogCdRom
	// Iscsi - iSCSI Volume (Devices Datastore)
	Iscsi
	// None - No disk type, error situation
	None = 255
)

// DatastoreState - state of datastore
type DatastoreState int

const (
	// Enabled state of datastore
	Enabled DatastoreState = iota
	// Disabled state of datastore
	Disabled
)

// CreateDatastoreWithID constructs Datastore
func CreateDatastoreWithID(id int) *Datastore {
	return &Datastore{*CreateResource("DATASTORE", id)}
}

// CreateDatastoreFromXML constructs datastore with full xml data
func CreateDatastoreFromXML(XMLdata *etree.Element) *Datastore {
	return &Datastore{Resource: Resource{XMLData: XMLdata}}
}

// User gets user ID of given datastore
func (d *Datastore) User() (int, error) {
	return d.intAttribute("UID")
}

// Group gets group ID of given datastore
func (d *Datastore) Group() (int, error) {
	return d.intAttribute("GID")
}

// Permissions gets datastore permissions
func (d *Datastore) Permissions() (*Permissions, error) {
	return d.permissions()
}

// DsMad gets DS_MAD of given datastore
func (d *Datastore) DsMad() (string, error) {
	return d.Attribute("DS_MAD")
}

// TmMad gets TM_MAD of given datastore
func (d *Datastore) TmMad() (string, error) {
	return d.Attribute("TM_MAD")
}

// BasePath gets base path of given datastore
func (d *Datastore) BasePath() (string, error) {
	return d.Attribute("BASE_PATH")
}

// Type gets type of given datastore
func (d *Datastore) Type() (DatastoreType, error) {
	i, err := d.intAttribute("TYPE")
	return DatastoreType(i), err
}

// DiskType gets disk type of given datastore
func (d *Datastore) DiskType() (DatastoreDiskType, error) {
	i, err := d.intAttribute("DISK_TYPE")
	return DatastoreDiskType(i), err
}

// State gets state of given datastore
func (d *Datastore) State() (DatastoreState, error) {
	i, err := d.intAttribute("STATE")
	return DatastoreState(i), err
}

// Clusters method returns an array of Cluster IDs for given datastore
func (d *Datastore) Clusters() ([]int, error) {
	return d.arrayOfIDs("CLUSTERS")
}

// TotalMB gets total size of memory of given datastore in MB
func (d *Datastore) TotalMB() (int, error) {
	return d.intAttribute("TOTAL_MB")
}

// FreeMB gets free size of memory of given datastore in MB
func (d *Datastore) FreeMB() (int, error) {
	return d.intAttribute("FREE_MB")
}

// UsedMB gets used size of memory of given datastore in MB
func (d *Datastore) UsedMB() (int, error) {
	return d.intAttribute("USED_MB")
}

// Images method returns an array of Image IDs for given datastore
func (d *Datastore) Images() ([]int, error) {
	return d.arrayOfIDs("IMAGES")
}
