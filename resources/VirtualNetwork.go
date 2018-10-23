package resources

import "github.com/beevik/etree"

// VirtualNetwork struct
type VirtualNetwork struct {
	Resource
}

// CreateVirtualNetworkWithID constructs VirtualNetwork with id
func CreateVirtualNetworkWithID(id int) *VirtualNetwork {
	return &VirtualNetwork{*CreateResource("VNET", id)}
}

// CreateVirtualNetworkFromXML constructs VirtualNetwork with full xml data
func CreateVirtualNetworkFromXML(XMLdata *etree.Element) *VirtualNetwork {
	return &VirtualNetwork{Resource: Resource{XMLData: XMLdata}}
}
