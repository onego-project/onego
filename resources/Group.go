package resources

import (
	"fmt"
	"github.com/beevik/etree"
	"strconv"
)

// Group struct
type Group struct {
	XMLData *etree.Element
}

// CreateGroup constructs Group
func CreateGroup(id int) *Group {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0"`)

	el := doc.CreateElement("GROUP").CreateElement("ID")
	el.SetText(fmt.Sprintf("%d", id))

	return &Group{doc.Root()}
}

// GetAttribute method
func (g Group) GetAttribute(path string) string {
	elements := g.XMLData.FindElements(path)
	if elements == nil {
		return ""
	}
	return elements[0].Text()
}

// GetID method
func (g Group) GetID() int {
	i, err := strconv.Atoi(g.GetAttribute("ID"))
	if err != nil {
		return -1
	}
	return i
}

// GetName method
func (g Group) GetName() string {
	return g.GetAttribute("NAME")
}

// GetUsers method
func (g Group) GetUsers() []User {
	elements := g.XMLData.FindElements("USERS/ID")
	users := make([]User, len(elements))
	for i, e := range elements {
		id, err := strconv.Atoi(e.Text())
		if err != nil {
			users[i] = *CreateUser(id)
		}
	}
	return users
}
