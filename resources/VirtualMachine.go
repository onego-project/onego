package resources

// GraphicsTypeMap to convert GraphicsType to string.
var GraphicsTypeMap = map[GraphicsType]string{
	GraphicsTypeVNC:   "VNC",
	GraphicsTypeSpice: "SPICE",
	GraphicsTypeSDL:   "SDL",
	GraphicsTypeNone:  "NONE",
}

// GraphicsType to set graphics type.
type GraphicsType int

const (
	// GraphicsTypeVNC graphics type.
	GraphicsTypeVNC GraphicsType = iota
	// GraphicsTypeSpice graphics type.
	GraphicsTypeSpice
	// GraphicsTypeSDL graphics type.
	GraphicsTypeSDL
	// GraphicsTypeNone graphics type.
	GraphicsTypeNone
)

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

