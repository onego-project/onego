package blueprint

// ClusterBlueprint to set Cluster elements
type ClusterBlueprint struct {
	Blueprint
}

// CreateUpdateClusterBlueprint creates empty ClusterBlueprint
func CreateUpdateClusterBlueprint() *ClusterBlueprint {
	return &ClusterBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
}

// SetReservedCPU sets size of reserved cpu of the given cluster
func (cb *ClusterBlueprint) SetReservedCPU(reservedCPU string) {
	cb.SetElement("RESERVED_CPU", reservedCPU)
}

// SetReservedMemory sets size of reserved memory of the given cluster
func (cb *ClusterBlueprint) SetReservedMemory(reservedMemory string) {
	cb.SetElement("RESERVED_MEM", reservedMemory)
}
