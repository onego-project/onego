package requests

import (
	"github.com/onego-project/onego/resources"
)

// OwnershipRequestBuilder struct
type OwnershipRequestBuilder struct {
	OwnershipRequest OwnershipRequest
}

// CreateOwnershipRequestBuilder construct OwnershipRequestBuilder
func CreateOwnershipRequestBuilder() OwnershipRequestBuilder {
	return OwnershipRequestBuilder{OwnershipRequest{*resources.CreateGroup(-1), *resources.CreateUser(-1)}}
}

// Group method
func (orb OwnershipRequestBuilder) Group(group resources.Group) OwnershipRequestBuilder {
	orb.OwnershipRequest.Group = group
	return orb
}

// User method
func (orb OwnershipRequestBuilder) User(user resources.User) OwnershipRequestBuilder {
	orb.OwnershipRequest.User = user
	return orb
}

// Build method
func (orb OwnershipRequestBuilder) Build() OwnershipRequest {
	return orb.OwnershipRequest
}
