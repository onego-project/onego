package requests

import (
	"github.com/onego-project/onego/resources"
)

// OwnershipRequest structure to change owners
type OwnershipRequest struct {
	Group resources.Group
	User  resources.User
}

// OwnershipRequestBuilder structure to create request to change owner
type OwnershipRequestBuilder struct {
	ownershipRequest OwnershipRequest
}

var noChangeUser = *resources.CreateUserWithID(-1)
var noChangeGroup = *resources.CreateGroupWithID(-1)

// CreateOwnershipRequestBuilder construct OwnershipRequestBuilder
func CreateOwnershipRequestBuilder() *OwnershipRequestBuilder {
	return &OwnershipRequestBuilder{OwnershipRequest{User: noChangeUser, Group: noChangeGroup}}
}

// Group to change group
func (orb *OwnershipRequestBuilder) Group(group resources.Group) *OwnershipRequestBuilder {
	orb.ownershipRequest.Group = group
	return orb
}

// User to change user
func (orb *OwnershipRequestBuilder) User(user resources.User) *OwnershipRequestBuilder {
	orb.ownershipRequest.User = user
	return orb
}

// Build create ownership request
func (orb *OwnershipRequestBuilder) Build() OwnershipRequest {
	return orb.ownershipRequest
}
