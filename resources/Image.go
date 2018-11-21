package resources

import (
	"time"

	"github.com/onego-project/onego/errors"

	"github.com/beevik/etree"
)

// Image structure represents OpenNebula Image
type Image struct {
	Resource
}

// ImageTypeMap contains string representation of ImageType
var ImageTypeMap = map[ImageType]string{
	ImageTypeOs:        "OS",
	ImageTypeCDRom:     "CDROM",
	ImageTypeDataBlock: "DATABLOCK",
	ImageTypeKernel:    "KERNEL",
	ImageTypeRAMDisk:   "RAMDISK",
	ImageTypeContext:   "CONTEXT",
}

// ImageType - type of Images
type ImageType int

const (
	// ImageTypeOs - Base OS image
	ImageTypeOs ImageType = iota
	// ImageTypeCDRom - An ISO9660 image
	ImageTypeCDRom
	// ImageTypeDataBlock - User persistent data device
	ImageTypeDataBlock
	// ImageTypeKernel - Kernel files
	ImageTypeKernel
	// ImageTypeRAMDisk - Initrd files
	ImageTypeRAMDisk
	// ImageTypeContext - Context files
	ImageTypeContext
)

// ImageStateMap contains string representation of ImageState
var ImageStateMap = map[ImageState]string{
	ImageStateInit:           "INIT",
	ImageStateReady:          "READY",
	ImageStateUsed:           "USED",
	ImageStateDisabled:       "DISABLED",
	ImageStateLocked:         "LOCKED",
	ImageStateError:          "ERROR",
	ImageStateClone:          "CLONE",
	ImageStateDelete:         "DELETE",
	ImageStateUsedPers:       "USED_PERS",
	ImageStateLockedUsed:     "LOCKED_USED",
	ImageStateLockedUsedPers: "LOCKED_USED_PERS",
}

// ImageState - state of Images
type ImageState int

const (
	// ImageStateInit - Initialization state
	ImageStateInit ImageState = iota
	// ImageStateReady - Image ready to use
	ImageStateReady
	// ImageStateUsed - Image in use
	ImageStateUsed
	// ImageStateDisabled - Image can not be instantiated by a VM
	ImageStateDisabled
	// ImageStateLocked - FS operation for the Image in process
	ImageStateLocked
	// ImageStateError - Error state the operation FAILED
	ImageStateError
	// ImageStateClone - Image is being cloned
	ImageStateClone
	// ImageStateDelete - DS is deleting the image
	ImageStateDelete
	// ImageStateUsedPers - Image is in use and persistent
	ImageStateUsedPers
	// ImageStateLockedUsed - FS operation in progress, VMs waiting
	ImageStateLockedUsed
	// ImageStateLockedUsedPers - FS operation in progress, VMs waiting. Persistent
	ImageStateLockedUsedPers
)

// ImageSnapshot represents snapshot created from Image
type ImageSnapshot struct {
	Active   string
	Children string
	Date     *time.Time
	ID       int
	Name     string
	Parent   int
	Size     int
}

// CreateImageWithID constructs Image with id
func CreateImageWithID(id int) *Image {
	return &Image{*CreateResource("IMAGE", id)}
}

// CreateImageFromXML constructs Image with full xml data
func CreateImageFromXML(XMLdata *etree.Element) *Image {
	return &Image{Resource: Resource{XMLData: XMLdata}}
}

// User gets User ID of given Image
func (i *Image) User() (int, error) {
	return i.intAttribute("UID")
}

// Group gets Group ID of given Image
func (i *Image) Group() (int, error) {
	return i.intAttribute("GID")
}

// Permissions gets image permissions
func (i *Image) Permissions() (*Permissions, error) {
	return i.permissions()
}

// Type gets image type.
// To get string representation of image type use ImageTypeMap.
func (i *Image) Type() (ImageType, error) {
	t, err := i.intAttribute("TYPE")
	return ImageType(t), err
}

// DiskType gets image disk type
func (i *Image) DiskType() (DiskType, error) {
	t, err := i.intAttribute("DISK_TYPE")
	return DiskType(t), err
}

// Persistent gets true if Image is persistent; gets false if Image is not persistent.
func (i *Image) Persistent() (bool, error) {
	ret, err := i.intAttribute("PERSISTENT")
	if err != nil {
		return false, err
	}
	return intToBool(ret), nil
}

// RegistrationTime gets time when Image was registered
func (i *Image) RegistrationTime() (*time.Time, error) {
	return i.registrationTime()
}

// Source gets Image source
func (i *Image) Source() (string, error) {
	return i.Attribute("SOURCE")
}

// Path gets path to Image
func (i *Image) Path() (string, error) {
	return i.Attribute("PATH")
}

// FileSystemType gets Image file system type
func (i *Image) FileSystemType() (string, error) {
	return i.Attribute("FSTYPE")
}

// Size gets Image size
func (i *Image) Size() (int, error) {
	return i.intAttribute("SIZE")
}

// State gets Image state.
// To get string representation of image type use ImageStateMap.
func (i *Image) State() (ImageState, error) {
	state, err := i.intAttribute("STATE")
	return ImageState(state), err
}

// RunningVMs gets number of running virtual machines
func (i *Image) RunningVMs() (int, error) {
	return i.intAttribute("RUNNING_VMS")
}

// Datastore gets Datastore ID of given Image
func (i *Image) Datastore() (int, error) {
	return i.intAttribute("DATASTORE_ID")
}

// VirtualMachines gets array of Vm IDs of given Image
func (i *Image) VirtualMachines() ([]int, error) {
	return i.arrayOfIDs("VMS")
}

// Clones gets array of clone IDs of given Image
func (i *Image) Clones() ([]int, error) {
	return i.arrayOfIDs("CLONES")
}

// AppClones gets array of App. clones of given Image
func (i *Image) AppClones() ([]int, error) {
	return i.arrayOfIDs("APP_CLONES")
}

// Snapshots gets array of Snapshot IDs of given Image
func (i *Image) Snapshots() ([]*ImageSnapshot, error) {
	elements := i.XMLData.FindElements("SNAPSHOTS/SNAPSHOT")
	if len(elements) == 0 {
		return make([]*ImageSnapshot, 0), nil
	}

	imageSnapshots := make([]*ImageSnapshot, len(elements))
	var err error

	for i, e := range elements {
		imageSnapshots[i], err = createImageSnapshotFromElement(e)
		if err != nil {
			return nil, err
		}
	}
	return imageSnapshots, nil
}

func createImageSnapshotFromElement(element *etree.Element) (*ImageSnapshot, error) {
	if element == nil {
		return nil, &errors.XMLElementError{Path: "snapshot"}
	}

	// occurrence 0 - 1 (we can ignore error)
	active, err := attributeFromElement(element, "ACTIVE")
	if err != nil {
		active = ""
	}

	children, err := attributeFromElement(element, "CHILDREN")
	if err != nil {
		children = ""
	}

	name, err := attributeFromElement(element, "NAME")
	if err != nil {
		name = ""
	}

	dateInt, err := intAttributeFromElement(element, "DATE")
	if err != nil {
		return nil, err
	}
	d := time.Unix(int64(dateInt), 0)
	date := &d

	id, err := intAttributeFromElement(element, "ID")
	if err != nil {
		return nil, err
	}

	parent, err := intAttributeFromElement(element, "PARENT")
	if err != nil {
		return nil, err
	}

	size, err := intAttributeFromElement(element, "SIZE")
	if err != nil {
		return nil, err
	}

	return &ImageSnapshot{Active: active, Children: children, Date: date,
		ID: id, Name: name, Parent: parent, Size: size}, nil
}
