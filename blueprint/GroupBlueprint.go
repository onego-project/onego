package blueprint

// GroupBlueprint to set Group elements
type GroupBlueprint struct {
	Blueprint
}

// CreateGroupBlueprint creates empty GroupBlueprint
func CreateGroupBlueprint() *GroupBlueprint {
	return &GroupBlueprint{Blueprint: *CreateBlueprint()}
}
