package services

// OwnershipFilter - resources belonging to
type OwnershipFilter int

const (
	// OwnershipFilterPrimaryGroup - resources belonging to the user's primary group
	OwnershipFilterPrimaryGroup OwnershipFilter = iota - 4
	// OwnershipFilterUser - resources belonging to the user
	OwnershipFilterUser
	// OwnershipFilterAll - all resources
	OwnershipFilterAll
	// OwnershipFilterUserGroup  - resources belonging to the user and any of his groups
	OwnershipFilterUserGroup
)
