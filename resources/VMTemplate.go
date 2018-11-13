package resources

import (
	"time"

	"github.com/beevik/etree"
)

// VMTemplate structure represents OpenNebula Virtual Machine Template.
type VMTemplate struct {
	Resource
}

// CreateVMTemplateWithID constructs VM with given ID.
func CreateVMTemplateWithID(id int) *VMTemplate {
	return &VMTemplate{*CreateResource("VMTEMPLATE", id)}
}

// CreateVMTemplateFromXML constructs VMTemplate with full xml data.
func CreateVMTemplateFromXML(XMLdata *etree.Element) *VMTemplate {
	return &VMTemplate{Resource: Resource{XMLData: XMLdata}}
}

// User gets user ID of given VMTemplate.
func (t *VMTemplate) User() (int, error) {
	return t.intAttribute("UID")
}

// Group gets group ID of given VMTemplate.
func (t *VMTemplate) Group() (int, error) {
	return t.intAttribute("GID")
}

// Permissions gets VMTemplate permissions.
func (t *VMTemplate) Permissions() (*Permissions, error) {
	return t.permissions()
}

// RegistrationTime gets time when VMTemplate was registered.
func (t *VMTemplate) RegistrationTime() (*time.Time, error) {
	return t.registrationTime()
}
