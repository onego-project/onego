package blueprint

// HostBlueprint to set Host elements
type HostBlueprint struct {
	Blueprint
}

// CreateUpdateHostBlueprint creates empty HostBlueprint
func CreateUpdateHostBlueprint() *HostBlueprint {
	return &HostBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
}

// SetArchitecture sets ARCH of the given Host
func (hb *HostBlueprint) SetArchitecture(arch string) {
	hb.SetElement("ARCH", arch)
}

// SetClusterName sets CLUSTER_NAME of the given Host
func (hb *HostBlueprint) SetClusterName(clusterName string) {
	hb.SetElement("CLUSTER_NAME", clusterName)
}

// SetHostName sets HOSTNAME of the given Host
func (hb *HostBlueprint) SetHostName(hostName string) {
	hb.SetElement("HOSTNAME", hostName)
}

// SetImMad sets IM_MAD of the given Host
func (hb *HostBlueprint) SetImMad(imMad string) {
	hb.SetElement("IM_MAD", imMad)
}

// SetPriority sets PRIORITY of the given Host
func (hb *HostBlueprint) SetPriority(priority string) {
	hb.SetElement("PRIORITY", priority)
}

// SetReservedCPU sets RESERVED_CPU of the given Host
func (hb *HostBlueprint) SetReservedCPU(reservedCPU string) {
	hb.SetElement("RESERVED_CPU", reservedCPU)
}

// SetReservedMemory sets RESERVED_MEM of the given Host
func (hb *HostBlueprint) SetReservedMemory(reservedMemory string) {
	hb.SetElement("RESERVED_MEM", reservedMemory)
}

// SetStatus sets STATUS of the given Host
func (hb *HostBlueprint) SetStatus(status string) {
	hb.SetElement("STATUS", status)
}

// SetVMMad sets VM_MAD of the given Host
func (hb *HostBlueprint) SetVMMad(vmMad string) {
	hb.SetElement("VM_MAD", vmMad)
}
