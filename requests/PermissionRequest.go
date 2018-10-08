package requests

// PermissionRequest structure to change permission
type PermissionRequest struct {
	Permissions [][]permissionChangeType
}

// PermissionRequestBuilder structure to create permission permissionRequest
type PermissionRequestBuilder struct {
	permissionRequest PermissionRequest
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

type permissionChangeType int

const (
	allow    = 1
	deny     = 0
	noChange = -1
)

// CreatePermissionRequestBuilder constructs PermissionRequestBuilder of PermissionRequest
func CreatePermissionRequestBuilder() *PermissionRequestBuilder {
	return &PermissionRequestBuilder{
		PermissionRequest{[][]permissionChangeType{{noChange, noChange, noChange},
			{noChange, noChange, noChange}, {noChange, noChange, noChange}}}}
}

// Allow allows certain permission group to perform the specified action.
// For example I'm allowing a group of a resource to use said resource.
func (prb *PermissionRequestBuilder) Allow(pg PermissionGroup, pt PermissionType) *PermissionRequestBuilder {
	prb.permissionRequest.Permissions[pg][pt] = allow
	return prb
}

// Deny denies certain permission group to perform the specified action.
// For example I'm denying a user (owner) of a resource to manage said resource.
func (prb *PermissionRequestBuilder) Deny(pg PermissionGroup, pt PermissionType) *PermissionRequestBuilder {
	prb.permissionRequest.Permissions[pg][pt] = deny
	return prb
}

// Build to create permission permissionRequest
func (prb *PermissionRequestBuilder) Build() PermissionRequest {
	return prb.permissionRequest
}
