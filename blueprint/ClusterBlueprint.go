package blueprint

import "strconv"

// ClusterBlueprint to set Cluster elements
type ClusterBlueprint struct {
	Blueprint
}

// CreateUpdateClusterBlueprint creates empty ClusterBlueprint
func CreateUpdateClusterBlueprint() *ClusterBlueprint {
	return &ClusterBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
}

// SetReservedCPU sets size of reserved cpu of the given cluster
func (cb *ClusterBlueprint) SetReservedCPU(reservedCPU int) {
	cb.SetElement("RESERVED_CPU", strconv.Itoa(reservedCPU))
}

// SetReservedMemory sets size of reserved memory of the given cluster
func (cb *ClusterBlueprint) SetReservedMemory(reservedMemory int) {
	cb.SetElement("RESERVED_MEM", strconv.Itoa(reservedMemory))
}
