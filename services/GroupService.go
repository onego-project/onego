package services

import (
	"context"
	"fmt"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/resources"
)

// GroupService manages group in OpenNebula
type GroupService struct {
	Service
}

// Allocate allocates a new group in OpenNebula
func (gs *GroupService) Allocate(ctx context.Context, groupName string) (*resources.Group, error) {
	resArr, err := gs.call(ctx, "one.group.allocate", groupName)
	if err != nil {
		return nil, err
	}

	return gs.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Delete deletes the given group from the pool
func (gs *GroupService) Delete(ctx context.Context, group resources.Group) error {
	groupID, err := group.ID()
	if err != nil {
		return err
	}
	_, err = gs.call(ctx, "one.group.delete", groupID)

	return err
}

// Update replaces the group template contents
func (gs *GroupService) Update(ctx context.Context, group resources.Group, blueprint blueprint.Interface, updateType UpdateType) error {
	groupID, err := group.ID()
	if err != nil {
		return err
	}

	blueprintText, err := blueprint.Render()
	if err != nil {
		return err
	}

	_, err = gs.call(ctx, "one.group.update", groupID, blueprintText, updateType)
	return err
}

// RetrieveInfo retrieves group by ID. If group id is -1 then the connected user's group info is returned
func (gs *GroupService) RetrieveInfo(ctx context.Context, groupID int) (*resources.Group, error) {
	doc, err := gs.retrieveInfo(ctx, "one.group.info", groupID)
	if err != nil {
		return nil, err
	}

	return resources.CreateGroupFromXML(doc.Root()), nil
}

// List retrieves list of groups in the group pool
func (gs *GroupService) List(ctx context.Context) ([]*resources.Group, error) {
	doc, err := gs.list(ctx, "one.grouppool.info")
	if err != nil {
		return nil, err
	}

	elements := doc.FindElements("GROUP_POOL/GROUP")
	if len(elements) == 0 {
		return nil, fmt.Errorf("no group in group pool")
	}

	groups := make([]*resources.Group, len(elements))
	for i, e := range elements {
		groups[i] = resources.CreateGroupFromXML(e)
	}

	return groups, nil
}

// auxiliary method to add or remove admin to/from the Group administrators set
func (gs *GroupService) manageAdmin(ctx context.Context, methodName string, group resources.Group, admin resources.User) error {
	groupID, err := group.ID()
	if err != nil {
		return err
	}

	adminID, err := admin.ID()
	if err != nil {
		return err
	}

	_, err = gs.call(ctx, methodName, groupID, adminID)
	return err
}

// AddAdmin adds an Admin to the Group administrators set
func (gs *GroupService) AddAdmin(ctx context.Context, group resources.Group, admin resources.User) error {
	return gs.manageAdmin(ctx, "one.group.addadmin", group, admin)
}

// RemoveAdmin removes an Admin from the Group administrators set
func (gs *GroupService) RemoveAdmin(ctx context.Context, group resources.Group, admin resources.User) error {
	return gs.manageAdmin(ctx, "one.group.deladmin", group, admin)
}
