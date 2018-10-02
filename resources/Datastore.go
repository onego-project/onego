package resources

import "github.com/beevik/etree"

// Datastore structure represents OpenNebula datastore
type Datastore struct {
	Resource
}

// DatastoreTypeMap contains string representation of DatastoreType
var DatastoreTypeMap = map[DatastoreType]string{
	DatastoreTypeImage:  "IMAGE_DS",
	DatastoreTypeSystem: "SYSTEM_DS",
	DatastoreTypeFile:   "FILE_DS",
}

// DatastoreType - type of datastore
type DatastoreType int

const (
	// DatastoreTypeImage - Standard datastore for disk images
	DatastoreTypeImage DatastoreType = iota
	// DatastoreTypeSystem - System datastore for disks of running VMs
	DatastoreTypeSystem
	// DatastoreTypeFile - File datastore for context, kernel, initrd files
	DatastoreTypeFile
)

// DiskTypeMap contains string representation of DiskType
var DiskTypeMap = map[DiskType]string{
	DiskTypeFile:          "FILE",
	DiskTypeBlock:         "BLOCK",
	DiskTypeRbd:           "RBD",
	DiskTypeRbdCdRom:      "RBD_CDROM",
	DiskTypeGluster:       "GLUSTER",
	DiskTypeGlusterCdRom:  "GLUSTER_CDROM",
	DiskTypeSheepdog:      "SHEEPDOG",
	DiskTypeSheepdogCdRom: "SHEEPDOG_CDROM",
	DiskTypeIscsi:         "ISCSI",
	DiskTypeNone:          "",
}

// DiskType - Type of Disks (used by the VMM_MAD). Values: BLOCK, CDROM or FILE
type DiskType int

const (
	// DiskTypeFile - DiskTypeFile-based disk
	DiskTypeFile DiskType = iota
	// DiskTypeCdRom - An ISO9660 disk
	DiskTypeCdRom
	// DiskTypeBlock - Block-device disk
	DiskTypeBlock
	// DiskTypeRbd - CEPH RBD disk
	DiskTypeRbd
	// DiskTypeRbdCdRom - CEPH RBD CDROM disk
	DiskTypeRbdCdRom
	// DiskTypeGluster - Gluster Block Device
	DiskTypeGluster
	// DiskTypeGlusterCdRom - Gluster CDROM Device Device
	DiskTypeGlusterCdRom
	// DiskTypeSheepdog - Sheepdog Block Device
	DiskTypeSheepdog
	// DiskTypeSheepdogCdRom - Sheepdog CDROM Device Device
	DiskTypeSheepdogCdRom
	// DiskTypeIscsi - iSCSI Volume (Devices Datastore)
	DiskTypeIscsi
	// DiskTypeNone - No disk type, error situation
	DiskTypeNone = 255
)

// DatastoreStateMap contains string representation of DatastoreState
var DatastoreStateMap = map[DatastoreState]string{
	DatastoreStateReady:    "READY",
	DatastoreStateDisabled: "DISABLED",
}

// DatastoreState - state of datastore
type DatastoreState int

const (
	// DatastoreStateReady state of datastore
	DatastoreStateReady DatastoreState = iota
	// DatastoreStateDisabled state of datastore
	DatastoreStateDisabled
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
func (d *Datastore) DiskType() (DiskType, error) {
	i, err := d.intAttribute("DISK_TYPE")
	return DiskType(i), err
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
