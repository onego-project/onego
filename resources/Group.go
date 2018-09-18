package resources

import (
	"github.com/beevik/etree"
	"strconv"
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
	elements := g.XMLData.FindElements("USERS/ID")
	if len(elements) == 0 {
		return make([]int, 0), nil
	}

	users := make([]int, len(elements))

	for i, e := range elements {
		id, err := strconv.Atoi(e.Text())
		if err != nil {
			return nil, err
		}
		users[i] = id
	}
	return users, nil
}

// Admins gets list of admins ids of the given group.
// It returns empty array when Group has no Admin.
func (g *Group) Admins() ([]int, error) {
	elements := g.XMLData.FindElements("ADMINS/ID")
	if len(elements) == 0 {
		return make([]int, 0), nil
	}

	admins := make([]int, len(elements))

	for i, e := range elements {
		id, err := strconv.Atoi(e.Text())
		if err != nil {
			return nil, err
		}
		admins[i] = id
	}
	return admins, nil
}
