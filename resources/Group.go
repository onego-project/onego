package resources

import (
	"github.com/beevik/etree"
)

// Group structure represents group resource.
type Group struct {
	Resource
}

// CreateGroupWithID constructs OpenNebula Group.
func CreateGroupWithID(id int) *Group {
	return &Group{*CreateResource("GROUP", id)}
}

// CreateGroupFromXML constructs Group with full xml data
func CreateGroupFromXML(XMLdata *etree.Element) *Group {
	return &Group{Resource: Resource{XMLData: XMLdata}}
}

// Users gets list of users ids of the given group.
// It returns empty array when Group has no User.
func (g *Group) Users() ([]int, error) {
	return g.arrayOfIDs("USERS")
}

// Admins gets list of admins ids of the given group.
// It returns empty array when Group has no Admin.
func (g *Group) Admins() ([]int, error) {
	return g.arrayOfIDs("ADMINS")
}
