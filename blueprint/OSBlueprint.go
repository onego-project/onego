package blueprint

// OSBlueprint to set  elements.
type OSBlueprint struct {
	Blueprint
}

// ArchitectureTypeMap to convert architecture type to string.
var ArchitectureTypeMap = map[ArchitectureType]string{
	ArchitectureTypeI686:   "i686",
	ArchitectureTypeX86_64: "x86_64",
}

// ArchitectureType - type of architecture.
type ArchitectureType int

const (
	// ArchitectureTypeI686 to set architecture type to i686.
	ArchitectureTypeI686 ArchitectureType = iota
	// ArchitectureTypeX86_64 to set architecture type to x86_64.
	ArchitectureTypeX86_64
)

// CreateOSBlueprint creates empty OSBlueprint.
func CreateOSBlueprint() *OSBlueprint {
	return &OSBlueprint{Blueprint: *CreateBlueprint("OS")}
}

// SetArchitecture sets ARCH of a given OS.
func (ob *OSBlueprint) SetArchitecture(value ArchitectureType) {
	ob.SetElement("ARCH", ArchitectureTypeMap[value])
}
