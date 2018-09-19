package blueprint

// GroupBlueprint to set Group elements
type GroupBlueprint struct {
	Blueprint
}

// CreateUpdateGroupBlueprint creates empty GroupBlueprint
func CreateUpdateGroupBlueprint() *GroupBlueprint {
	return &GroupBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
}
