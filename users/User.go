package users

import (
	"fmt"
	"github.com/beevik/etree"
	"strconv"
)

// User struct
type User struct {
	XMLData *etree.Element
}

// CreateUser constructs User
func CreateUser(id int) *User {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0"`)

	el := doc.CreateElement("USER").CreateElement("ID")
	el.SetText(fmt.Sprintf("%d", id))

	return &User{doc.Root()}
}

// GetAttribute method
func (u User) GetAttribute(path string) string {
	elements := u.XMLData.FindElements(path)
	if elements == nil {
		return ""
	}
	return elements[0].Text()
}

// GetID method
func (u User) GetID() int {
	i, err := strconv.Atoi(u.GetAttribute("ID"))
	if err != nil {
		return -1
	}
	return i
}

// GetName method
func (u User) GetName() string {
	return u.GetAttribute("NAME")
}

// GetPassword method
func (u User) GetPassword() string {
	return u.GetAttribute("PASSWORD")
}

// GetAuthDriver method
func (u User) GetAuthDriver() string {
	return u.GetAttribute("AUTH_DRIVER")
}

// GetGroups method
func (u User) GetGroups() []Group {
	elements := u.XMLData.FindElements("GROUPS/ID")
	groups := make([]Group, len(elements))
	for i, e := range elements {
		id, err := strconv.Atoi(e.Text())
		if err != nil {
			groups[i] = *CreateGroup(id)
		}
	}
	return groups
}
