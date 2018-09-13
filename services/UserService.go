package services

import (
	"context"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/resources"
)

// UserService structure to manage OpenNebula User
type UserService struct {
	Service
}

// Allocate allocates a new user in OpenNebula and returns User
func (us *UserService) Allocate(ctx context.Context, username, password, authDriver string, mainGroup resources.Group, secondaryGroups []resources.Group) (*resources.User, error) {
	groups := make([]int, len(secondaryGroups)+1)

	groupID, err := mainGroup.ID()
	if err != nil {
		return nil, err
	}
	groups[0] = groupID

	for i, group := range secondaryGroups {
		groupID, err = group.ID()
		if err != nil {
			return nil, err
		}
		groups[i+1] = groupID
	}

	resArr, err := us.call(ctx, "one.user.allocate", username, password, authDriver, groups)
	if err != nil {
		return nil, err
	}

	return us.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Delete deletes the given user from the pool
func (us *UserService) Delete(ctx context.Context, user resources.User) error {
	userID, err := user.ID()
	if err != nil {
		return err
	}
	_, err = us.call(ctx, "one.user.delete", userID)

	return err
}

// ChangePassword changes password for the given user
func (us *UserService) ChangePassword(ctx context.Context, user resources.User, password string) error {
	userID, err := user.ID()
	if err != nil {
		return err
	}

	_, err = us.call(ctx, "one.user.passwd", userID, password)
	return err
}

// Update replaces the user template contents
func (us *UserService) Update(ctx context.Context, user resources.User, blueprint blueprint.Interface, updateType UpdateType) error {
	userID, err := user.ID()
	if err != nil {
		return err
	}

	blueprintText, err := blueprint.Render()
	if err != nil {
		return err
	}

	_, err = us.call(ctx, "one.user.update", userID, blueprintText, updateType)
	return err
}

// ChangeAuthDriver changes the authentication driver for the given user; last argument of call avoids to change password via this method
func (us *UserService) ChangeAuthDriver(ctx context.Context, user resources.User, authDriver string) error {
	userID, err := user.ID()
	if err != nil {
		return err
	}

	_, err = us.call(ctx, "one.user.chauth", userID, authDriver, "")
	return err
}

// auxiliary method to change main group or add/remove secondary group
func (us *UserService) manageGroups(ctx context.Context, methodName string, user resources.User, group resources.Group) error {
	userID, err := user.ID()
	if err != nil {
		return err
	}

	groupID, err := group.ID()
	if err != nil {
		return err
	}

	_, err = us.call(ctx, methodName, userID, groupID)
	return err
}

// ChangeMainGroup changes the main group of the given user
func (us *UserService) ChangeMainGroup(ctx context.Context, user resources.User, group resources.Group) error {
	return us.manageGroups(ctx, "one.user.chgrp", user, group)
}

// AddSecondaryGroup adds the User to the secondary group
func (us *UserService) AddSecondaryGroup(ctx context.Context, user resources.User, group resources.Group) error {
	return us.manageGroups(ctx, "one.user.addgroup", user, group)
}

// RemoveSecondaryGroup removes the User from a secondary group
func (us *UserService) RemoveSecondaryGroup(ctx context.Context, user resources.User, group resources.Group) error {
	return us.manageGroups(ctx, "one.user.delgroup", user, group)
}

// RetrieveInfo retrieves information for the given user. If user ID is -1 the connected user's own info is returned
func (us *UserService) RetrieveInfo(ctx context.Context, userID int) (*resources.User, error) {
	doc, err := us.retrieveInfo(ctx, "one.user.info", userID)
	if err != nil {
		return nil, err
	}

	return resources.CreateUserFromXML(doc.Root()), nil
}

// RetrieveConnectedUserInfo retrieves information for the connected user
func (us *UserService) RetrieveConnectedUserInfo(ctx context.Context) (*resources.User, error) {
	return us.RetrieveInfo(ctx, -1)
}

// List retrieves information for all the users in the pool
func (us *UserService) List(ctx context.Context) ([]*resources.User, error) {
	doc, err := us.list(ctx, "one.userpool.info")
	if err != nil {
		return nil, err
	}

	elements := doc.FindElements("USER_POOL/USER")

	users := make([]*resources.User, len(elements))
	for i, e := range elements {
		users[i] = resources.CreateUserFromXML(e)
	}

	return users, nil
}
