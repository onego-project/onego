package requests

import (
	"github.com/onego-project/onego/resources"
)

// OwnershipRequest struct
type OwnershipRequest struct {
	Group resources.Group
	User  resources.User
}
