package blueprint

// UserTemplateBlueprint to set UserTemplate elements.
type UserTemplateBlueprint struct {
	Blueprint
}

// CreateUpdateUserTemplateBlueprint creates empty UserTemplateBlueprint.
func CreateUpdateUserTemplateBlueprint() *UserTemplateBlueprint {
	return &UserTemplateBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
}

// SetDescription sets description of a given UserTemplate.
func (tb *UserTemplateBlueprint) SetDescription(description string) {
	tb.SetElement("DESCRIPTION", description)
}
