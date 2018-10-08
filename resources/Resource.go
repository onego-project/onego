package resources

import (
	"fmt"
	"strconv"

	"github.com/beevik/etree"
)

// Resource structure contains XML data and main methods for open nebula resources
type Resource struct {
	XMLData *etree.Element
}

// Permissions structure represents permissions
type Permissions struct {
	User  PermissionGroup
	Group PermissionGroup
	Other PermissionGroup
}

// PermissionGroup structure to create permission structure
type PermissionGroup struct {
	Use    bool
	Manage bool
	Admin  bool
}

const invalidCode = -1

// constants for value conversion in resources methods
const base10 = 10
const bitSize64 = 64

// CreateResource constructs Resource by tag and ID
func CreateResource(tag string, id int) *Resource {
	doc := etree.NewDocument()

	el := doc.CreateElement(tag).CreateElement("ID")
	el.SetText(fmt.Sprintf("%d", id))

	return &Resource{doc.Root()}
}

// Attribute gets resource attribute founded on the path
func (r *Resource) Attribute(path string) (string, error) {
	return attributeFromElement(r.XMLData, path)
}

func (r *Resource) intAttribute(path string) (int, error) {
	return intAttributeFromElement(r.XMLData, path)
}

func attributeFromElement(e *etree.Element, path string) (string, error) {
	if e == nil {
		return "", fmt.Errorf("no element, unable to get %s", path)
	}

	element := e.FindElement(path)
	if element == nil {
		return "", fmt.Errorf("unable to find %s", path)
	}
	return element.Text(), nil
}

func intAttributeFromElement(e *etree.Element, path string) (int, error) {
	attribute, err := attributeFromElement(e, path)
	if err != nil {
		return invalidCode, err
	}
	i, err := strconv.Atoi(attribute)
	if err != nil {
		return invalidCode, err
	}
	return i, nil
}

// ID gets resource id
func (r *Resource) ID() (int, error) {
	return r.intAttribute("ID")
}

// Name gets resource name
func (r *Resource) Name() (string, error) {
	return r.Attribute("NAME")
}

func (r *Resource) arrayOfIDs(tag string) ([]int, error) {
	elements := r.XMLData.FindElements(tag + "/ID")
	if len(elements) == 0 {
		return make([]int, 0), nil
	}

	object := make([]int, len(elements))

	for i, e := range elements {
		id, err := strconv.Atoi(e.Text())
		if err != nil {
			return nil, err
		}
		object[i] = id
	}
	return object, nil
}

func intToBool(i int) bool {
	return i == 1
}

func (r *Resource) createPermission(perm string) (*PermissionGroup, error) {
	permissionTypeArray := [3]string{"U", "M", "A"}

	var resPermTypeArray [len(permissionTypeArray)]int
	var err error

	for i := 0; i < len(permissionTypeArray); i++ {
		resPermTypeArray[i], err = r.intAttribute("PERMISSIONS/" + perm + "_" + permissionTypeArray[i])
		if err != nil {
			return nil, err
		}
	}

	return &PermissionGroup{Use: intToBool(resPermTypeArray[0]), Manage: intToBool(resPermTypeArray[1]),
		Admin: intToBool(resPermTypeArray[2])}, nil
}

func (r *Resource) permissions() (*Permissions, error) {
	permissionGroupArray := [3]string{"OWNER", "GROUP", "OTHER"}

	var resPermGroupArray [len(permissionGroupArray)]*PermissionGroup
	var err error

	for i := 0; i < len(permissionGroupArray); i++ {
		resPermGroupArray[i], err = r.createPermission(permissionGroupArray[i])
		if err != nil {
			return nil, err
		}
	}

	return &Permissions{User: *resPermGroupArray[0], Group: *resPermGroupArray[1], Other: *resPermGroupArray[2]}, nil
}
