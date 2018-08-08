package requests

// PermissionRequest struct
type PermissionRequest struct {
	Permissions [][]int
}

// Builder struct
type Builder struct {
	request PermissionRequest
}

// PermissionGroup type
type PermissionGroup int

const (
	// User of PermissionGroup
	User PermissionGroup = iota
	// Group of PermissionGroup
	Group
	// Other of PermissionGroup
	Other
)

// PermissionType type
type PermissionType int

const (
	// Use of PermissionType
	Use PermissionType = iota
	// Manage of PermissionGroup
	Manage
	// Admin of PermissionGroup
	Admin
)

// CreateBuilder constructs Builder of PermissionRequest
func CreateBuilder() Builder {
	return Builder{
		PermissionRequest{[][]int{{-1, -1, -1}, {-1, -1, -1}, {-1, -1, -1}}}}
}

// Allow method
func (b Builder) Allow(pg PermissionGroup, pt PermissionType) Builder {
	b.request.Permissions[pg][pt] = 1
	return b
}

// Deny method
func (b Builder) Deny(pg PermissionGroup, pt PermissionType) Builder {
	b.request.Permissions[pg][pt] = 0
	return b
}

// Build method
func (b Builder) Build() PermissionRequest {
	return b.request
}
