package blueprint

// UserBlueprint to set User elements
type UserBlueprint struct {
	Blueprint
}

// CreateUpdateUserBlueprint creates empty UserBlueprint
func CreateUpdateUserBlueprint() *UserBlueprint {
	return &UserBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
}

// SetEmail sets email of the given user
func (u *UserBlueprint) SetEmail(email string) {
	u.SetElement("EMAIL", email)
}

// SetFullName sets full name of the given user
func (u *UserBlueprint) SetFullName(fullName string) {
	u.SetElement("NAME", fullName)
}

// SetSSHPublicKey sets ssh public key of the given user
func (u *UserBlueprint) SetSSHPublicKey(sshPublicKey string) {
	u.SetElement("SSH_PUBLIC_KEY", sshPublicKey)
}
