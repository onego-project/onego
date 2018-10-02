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
