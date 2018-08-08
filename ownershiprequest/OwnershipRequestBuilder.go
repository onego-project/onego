package ownershiprequest

import (
	"github.com/owlet123/onego/users"
)

// Builder struct
type Builder struct {
	OwnershipRequest OwnershipRequest
}

// CreateBuilder construct Builder
func CreateBuilder() Builder {
	return Builder{OwnershipRequest{*users.CreateGroup(-1), *users.CreateUser(-1)}}
}

// Group method
func (b Builder) Group(group users.Group) Builder {
	b.OwnershipRequest.Group = group
	return b
}

// User method
func (b Builder) User(user users.User) Builder {
	b.OwnershipRequest.User = user
	return b
}

// Build method
func (b Builder) Build() OwnershipRequest {
	return b.OwnershipRequest
}
