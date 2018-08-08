package ownershiprequest

import (
	"github.com/owlet123/onego/users"
)

// OwnershipRequest struct
type OwnershipRequest struct {
	Group users.Group
	User  users.User
}
