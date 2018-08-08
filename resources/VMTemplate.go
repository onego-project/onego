package resources

import (
	"github.com/beevik/etree"
	"strconv"
	"fmt"
	"github.com/onego-project/onego/requests"
)

// VMTemplate struct
type VMTemplate struct {
	XMLData *etree.Element
}

// CreateVMTemplate constructs VMTemplate
func CreateVMTemplate(id int) *Host {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0"`)

	el := doc.CreateElement("VMTEMPLATE").CreateElement("ID")
	el.SetText(fmt.Sprintf("%d", id))

	return &Host{doc.Root()}
}

// GetAttribute method
func (t VMTemplate) GetAttribute(path string) string {
	elements := t.XMLData.FindElements(path)
	if elements == nil {
		return ""
	}
	return elements[0].Text()
}

// GetID method
func (t VMTemplate) GetID() int {
	i, err := strconv.Atoi(t.GetAttribute("ID"))
	if err != nil {
		return -1
	}
	return i
}

// GetName method
func (t VMTemplate) GetName() string {
	return t.GetAttribute("NAME")
}

// GetUID method
func (t VMTemplate) GetUser() *User {
	i, err := strconv.Atoi(t.GetAttribute("UID"))
	if err != nil {
		return nil
	}
	return CreateUser(i)
}

// GetGID method
func (t VMTemplate) GetGroup() *Group {
	i, err := strconv.Atoi(t.GetAttribute("GID"))
	if err != nil {
		return nil
	}
	return CreateGroup(i)
}

// GetPermissionUserUse method
func (t VMTemplate) GetPermissionUserUse() int {
	i, err := strconv.Atoi(t.GetAttribute("PERMISSIONS/OWNER_U"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionUserManage method
func (t VMTemplate) GetPermissionUserManage() int {
	i, err := strconv.Atoi(t.GetAttribute("PERMISSIONS/OWNER_M"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionUserAdmin method
func (t VMTemplate) GetPermissionUserAdmin() int {
	i, err := strconv.Atoi(t.GetAttribute("PERMISSIONS/OWNER_A"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionGroupUse method
func (t VMTemplate) GetPermissionGroupUse() int {
	i, err := strconv.Atoi(t.GetAttribute("PERMISSIONS/GROUP_U"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionGroupManage method
func (t VMTemplate) GetPermissionGroupManage() int {
	i, err := strconv.Atoi(t.GetAttribute("PERMISSIONS/GROUP_M"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionGroupAdmin method
func (t VMTemplate) GetPermissionGroupAdmin() int {
	i, err := strconv.Atoi(t.GetAttribute("PERMISSIONS/GROUP_A"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionOtherUse method
func (t VMTemplate) GetPermissionOtherUse() int {
	i, err := strconv.Atoi(t.GetAttribute("PERMISSIONS/OTHER_U"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionOtherManage method
func (t VMTemplate) GetPermissionOtherManage() int {
	i, err := strconv.Atoi(t.GetAttribute("PERMISSIONS/OTHER_M"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionOtherAdmin method
func (t VMTemplate) GetPermissionOtherAdmin() int {
	i, err := strconv.Atoi(t.GetAttribute("PERMISSIONS/OTHER_A"))
	if err != nil {
		return -1
	}
	return i
}

// GetPermissionRequest method
func (t VMTemplate) GetPermissionRequest() *requests.PermissionRequest {
	return &requests.PermissionRequest{
		Permissions: [][]int{{t.GetPermissionUserUse(), t.GetPermissionUserManage(), t.GetPermissionUserAdmin()},
			{t.GetPermissionGroupUse(), t.GetPermissionGroupManage(), t.GetPermissionGroupAdmin()},
			{t.GetPermissionOtherUse(), t.GetPermissionOtherManage(), t.GetPermissionOtherAdmin()}}}
}
