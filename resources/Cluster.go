package resources

import "github.com/beevik/etree"

// Cluster structure represents OpenNebula cluster
type Cluster struct {
	Resource
}

// CreateClusterWithID constructs User with id
func CreateClusterWithID(id int) *Cluster {
	return &Cluster{*CreateResource("CLUSTER", id)}
}

// CreateClusterFromXML constructs Cluster with full xml data
func CreateClusterFromXML(XMLdata *etree.Element) *Cluster {
	return &Cluster{Resource: Resource{XMLData: XMLdata}}
}

// Hosts gets array of host ID of given Cluster
func (c *Cluster) Hosts() ([]int, error) {
	return c.arrayOfIDs("HOSTS")
}

// Datastores gets array of datastores ID of given Cluster
func (c *Cluster) Datastores() ([]int, error) {
	return c.arrayOfIDs("DATASTORES")
}

// VirtualNetworks gets array of virtual networks ID of given Cluster
func (c *Cluster) VirtualNetworks() ([]int, error) {
	return c.arrayOfIDs("VNETS")
}
