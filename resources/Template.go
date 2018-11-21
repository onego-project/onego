package resources

import (
	"time"

	"github.com/beevik/etree"
)

// Template structure represents OpenNebula Virtual Machine Template.
type Template struct {
	Resource
}

// CreateTemplateWithID constructs VM with given ID.
func CreateTemplateWithID(id int) *Template {
	return &Template{*CreateResource("VMTEMPLATE", id)}
}

// CreateTemplateFromXML constructs Template with full xml data.
func CreateTemplateFromXML(XMLdata *etree.Element) *Template {
	return &Template{Resource: Resource{XMLData: XMLdata}}
}

// User gets user ID of given Template.
func (t *Template) User() (int, error) {
	return t.intAttribute("UID")
}

// Group gets group ID of given Template.
func (t *Template) Group() (int, error) {
	return t.intAttribute("GID")
}

// Permissions gets Template permissions.
func (t *Template) Permissions() (*Permissions, error) {
	return t.permissions()
}

// RegistrationTime gets time when Template was registered.
func (t *Template) RegistrationTime() (*time.Time, error) {
	return t.registrationTime()
}
