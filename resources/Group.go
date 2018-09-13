package resources

import (
	"strconv"
)

// Group structure represents group resource.
type Group struct {
	Resource
}

// CreateGroup constructs OpenNebula Group.
func CreateGroup(id int) *Group {
	return &Group{*CreateResource("GROUP", id)}
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
